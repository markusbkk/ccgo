// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"strings"

	"modernc.org/cc/v4"
)

func (c *ctx) statement(w writer, n *cc.Statement) {
	sep := sep(n)
	if c.task.positions {
		sep = strings.TrimRight(sep, "\n\r\t ")
	}
	switch n.Case {
	case cc.StatementLabeled: // LabeledStatement
		w.w("%s%s", sep, c.posComment(n))
		c.labeledStatement(w, n.LabeledStatement)
	case cc.StatementCompound: // CompoundStatement
		c.compoundStatement(w, n.CompoundStatement, false, "")
	case cc.StatementExpr: // ExpressionStatement
		var a buf
		b := c.expr(&a, n.ExpressionStatement.ExpressionList, nil, exprVoid)
		if a.len() == 0 && b.len() == 0 {
			return
		}

		w.w("%s%s", sep, c.posComment(n))
		if a.len() != 0 {
			w.w("%s;", a.bytes())
		}
		w.w("%s;", b.bytes())
	case cc.StatementSelection: // SelectionStatement
		w.w("%s%s", sep, c.posComment(n))
		c.selectionStatement(w, n.SelectionStatement)
	case cc.StatementIteration: // IterationStatement
		w.w("%s%s", sep, c.posComment(n))
		c.iterationStatement(w, n.IterationStatement)
	case cc.StatementJump: // JumpStatement
		w.w("%s%s", sep, c.posComment(n))
		c.jumpStatement(w, n.JumpStatement)
	case cc.StatementAsm: // AsmStatement
		w.w("%s%s", sep, c.posComment(n))
		a := strings.Split(nodeSource(n.AsmStatement), "\n")
		w.w("\n// %s", strings.Join(a, "\n// "))
		w.w("\n%spanic(0) // assembler statements not supported", tag(preserve))
		c.err(errorf("%v: assembler statements not supported", c.pos(n)))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) labeledStatement(w writer, n *cc.LabeledStatement) {
	switch n.Case {
	case cc.LabeledStatementLabel: // IDENTIFIER ':' Statement
		w.w("%s%s:", tag(preserve), n.Token.Src()) //TODO use nameSpace
		c.statement(w, n.Statement)
	case cc.LabeledStatementCaseLabel: // "case" ConstantExpression ':' Statement
		switch {
		case len(c.switchCtx) != 0:
			w.w("%s:", c.switchCtx[0])
			c.switchCtx = c.switchCtx[1:]
		default:
			if n.CaseOrdinal() != 0 {
				w.w("fallthrough;")
			}
			w.w("case %s:", c.expr(nil, n.ConstantExpression, cc.IntegerPromotion(n.Switch().ExpressionList.Type()), exprDefault))
		}
		c.unbracedStatement(w, n.Statement)
	case cc.LabeledStatementRange: // "case" ConstantExpression "..." ConstantExpression ':' Statement
		if n.CaseOrdinal() != 0 {
			w.w("fallthrough;")
		}
		c.err(errorf("TODO %v", n.Case))
	case cc.LabeledStatementDefault: // "default" ':' Statement
		switch {
		case len(c.switchCtx) != 0:
			w.w("%s:", c.switchCtx[0])
			c.switchCtx = c.switchCtx[1:]
		default:
			if n.CaseOrdinal() != 0 {
				w.w("fallthrough;")
			}
			w.w("default:")
		}
		c.unbracedStatement(w, n.Statement)
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) compoundStatement(w writer, n *cc.CompoundStatement, fnBlock bool, value string) {
	defer func() { c.compoundStmtValue = value }()

	if _, ok := c.f.flatScopes[n.LexicalScope()]; ok {
		if fnBlock {
			trc("TODO")
			c.err(errorf("%v: internal error", n.Position()))
			return
		}

		trc("TODO")
		c.err(errorf("%v: TODO %v", n.Position(), fnBlock))
		return
	}

	switch {
	case fnBlock:
		if c.pass != 2 {
			break
		}

		w.w(" {")
		switch {
		case c.f.tlsAllocs+int64(c.f.maxValist) != 0:
			c.f.tlsAllocs = roundup(c.f.tlsAllocs, 8)
			v := c.f.tlsAllocs
			if c.f.maxValist != 0 {
				v += 8 * int64((c.f.maxValist + 1))
			}
			w.w("%sbp := %[1]stls.%sAlloc(%d); /* tlsAllocs %v maxValist %v */", tag(ccgo), tag(preserve), v, c.f.tlsAllocs, c.f.maxValist)
			w.w("defer %stls.%sFree(%d);", tag(ccgo), tag(preserve), v)
			for _, v := range c.f.t.Parameters() {
				if d := v.Declarator; d != nil && c.f.declInfos.info(d).pinned() {
					w.w("*(*%s)(%s) = %s_%s;", c.typ(n, d.Type()), unsafePointer(bpOff(c.f.declInfos.info(d).bpOff)), tag(ccgo), d.Name())
				}
			}
			w.w("%s%s", sep(n.Token), c.posComment(n))
		default:
			w.w("%s%s", strings.TrimSpace(sep(n.Token)), c.posComment(n))
		}
		w.w("%s", c.f.declareLocals())
	default:
		w.w(" { %s%s", sep(n.Token), c.posComment(n))
	}
	var bi *cc.BlockItem
	for l := n.BlockItemList; l != nil; l = l.BlockItemList {
		bi = l.BlockItem
		if l.BlockItemList == nil && value != "" {
			switch bi.Case {
			case cc.BlockItemStmt:
				s := bi.Statement
				for s.Case == cc.StatementLabeled {
					s = s.LabeledStatement.Statement
				}
				switch s.Case {
				case cc.StatementExpr:
					switch e := s.ExpressionStatement.ExpressionList; x := e.(type) {
					case *cc.PrimaryExpression:
						switch x.Case {
						case cc.PrimaryExpressionStmt:
							c.blockItem(w, bi)
							w.w("%s = %s;", value, c.compoundStmtValue)
						default:
							w.w("%s = ", value)
							c.blockItem(w, bi)
						}
					default:
						w.w("%s = ", value)
						c.blockItem(w, bi)
					}
				default:
					// trc("%v: %s", bi.Position(), cc.NodeSource(bi))
					c.err(errorf("TODO %v", s.Case))
				}
			default:
				c.err(errorf("TODO %v", bi.Case))
			}
			break
		}

		c.blockItem(w, bi)
	}
	switch {
	case fnBlock && c.f.t.Result().Kind() != cc.Void && !c.isReturn(bi):
		s := sep(n.Token2)
		if strings.Contains(s, "\n") {
			w.w("%s", s)
			s = ""
		}
		w.w("return %s%s;%s", tag(ccgo), retvalName, s)
	default:
		w.w("%s", sep(n.Token2))
	}
	w.w("};")
}

func (c *ctx) isReturn(n *cc.BlockItem) bool {
	if n == nil || n.Case != cc.BlockItemStmt {
		return false
	}

	return n.Statement.Case == cc.StatementJump && n.Statement.JumpStatement.Case == cc.JumpStatementReturn
}

func (c *ctx) blockItem(w writer, n *cc.BlockItem) {
	switch n.Case {
	case cc.BlockItemDecl: // Declaration
		c.declaration(w, n.Declaration, false)
	case cc.BlockItemLabel: // LabelDeclaration
		c.err(errorf("TODO %v", n.Case))
	case cc.BlockItemStmt: // Statement
		c.statement(w, n.Statement)
	case cc.BlockItemFuncDef: // DeclarationSpecifiers Declarator CompoundStatement
		if c.pass == 2 {
			c.err(errorf("%v: nested functions not supported", c.pos(n.Declarator)))
			// c.functionDefinition0(w, sep(n), n, n.Declarator, n.CompoundStatement, true) //TODO does not really work yet
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) selectionStatement(w writer, n *cc.SelectionStatement) {
	if _, ok := c.f.flatScopes[n.LexicalScope()]; ok {
		c.selectionStatementFlat(w, n)
		return
	}

	switch n.Statement.Case {
	case cc.StatementCompound:
		if _, ok := c.f.flatScopes[n.Statement.CompoundStatement.LexicalScope()]; ok {
			c.selectionStatementFlat(w, n)
			return
		}
	}

	switch n.Case {
	case cc.SelectionStatementIf: // "if" '(' ExpressionList ')' Statement
		w.w("if %s", c.expr(w, n.ExpressionList, nil, exprBool))
		c.bracedStatement(w, n.Statement)
	case cc.SelectionStatementIfElse: // "if" '(' ExpressionList ')' Statement "else" Statement
		w.w("if %s {", c.expr(w, n.ExpressionList, nil, exprBool))
		c.unbracedStatement(w, n.Statement)
		w.w("} else {")
		c.unbracedStatement(w, n.Statement2)
		w.w("};")
	case cc.SelectionStatementSwitch: // "switch" '(' ExpressionList ')' Statement
		for _, v := range n.LabeledStatements() {
			if v.Case == cc.LabeledStatementLabel {
				c.selectionStatementFlat(w, n)
				return
			}
		}

		ok := false
		switch s := n.Statement; s.Case {
		case cc.StatementCompound:
			if l := n.Statement.CompoundStatement.BlockItemList; l != nil {
				switch bi := l.BlockItem; bi.Case {
				case cc.BlockItemStmt:
					switch s := bi.Statement; s.Case {
					case cc.StatementLabeled:
						switch ls := s.LabeledStatement; ls.Case {
						case cc.LabeledStatementCaseLabel, cc.LabeledStatementDefault:
							// ok
							ok = true
						}
					}
				}
			}

		}
		if !ok {
			c.selectionStatementFlat(w, n)
			return
		}

		defer c.setSwitchCtx(nil)()
		defer c.setBreakCtx("")()

		w.w("switch %s", c.expr(w, n.ExpressionList, cc.IntegerPromotion(n.ExpressionList.Type()), exprDefault))
		c.statement(w, n.Statement)
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) selectionStatementFlat(w writer, n *cc.SelectionStatement) {
	switch n.Case {
	case cc.SelectionStatementIf: // "if" '(' ExpressionList ')' Statement
		//	if !expr goto a
		//	stmt
		// a:
		a := c.label()
		w.w("if !(%s) { goto %s };", c.expr(w, n.ExpressionList, nil, exprBool), a)
		c.unbracedStatement(w, n.Statement)
		w.w("%s:", a)
	case cc.SelectionStatementIfElse: // "if" '(' ExpressionList ')' Statement "else" Statement
		//	if !expr goto a
		//	stmt
		//	goto b
		// a:	stmt2
		// b:
		a := c.label()
		b := c.label()
		w.w("if !(%s) { goto %s };", c.expr(w, n.ExpressionList, nil, exprBool), a)
		c.unbracedStatement(w, n.Statement)
		w.w("goto %s; %s:", b, a)
		c.unbracedStatement(w, n.Statement2)
		w.w("%s:", b)
	case cc.SelectionStatementSwitch: // "switch" '(' ExpressionList ')' Statement
		//	switch expr {
		//	case 1:
		//		goto label1
		//	case 2:
		//		goto label2
		//	...
		//	default:
		//		goto labelN
		//	}
		//	goto brk
		//	label1:
		//		statements in case 1
		//	label2:
		//		statements in case 2
		//	...
		//	labelN:
		//		statements in default
		//	brk:
		t := cc.IntegerPromotion(n.ExpressionList.Type())
		w.w("switch %s {", c.expr(w, n.ExpressionList, t, exprDefault))
		var labels []string
		for _, v := range n.LabeledStatements() {
			switch v.Case {
			case cc.LabeledStatementLabel: // IDENTIFIER ':' Statement
				// nop
			case cc.LabeledStatementCaseLabel: // "case" ConstantExpression ':' Statement
				label := c.label()
				labels = append(labels, label)
				w.w("case %s: goto %s;", c.expr(nil, v.ConstantExpression, t, exprDefault), label)
			case cc.LabeledStatementRange: // "case" ConstantExpression "..." ConstantExpression ':' Statement
				c.err(errorf("TODO %v", n.Case))
			case cc.LabeledStatementDefault: // "default" ':' Statement
				label := c.label()
				labels = append(labels, label)
				w.w("default: goto %s;", label)
			default:
				c.err(errorf("internal error %T %v", n, n.Case))
			}
		}
		brk := c.label()
		w.w("\n}; goto %s;", brk)

		defer c.setSwitchCtx(labels)()
		defer c.setBreakCtx(brk)()

		c.unbracedStatement(w, n.Statement)
		w.w("%s:", brk)
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) bracedStatement(w writer, n *cc.Statement) {
	switch n.Case {
	case cc.StatementCompound:
		c.statement(w, n)
	default:
		w.w("{")
		c.statement(w, n)
		w.w("};")
	}
}

func (c *ctx) unbracedStatement(w writer, n *cc.Statement) {
	switch n.Case {
	case cc.StatementCompound:
		for l := n.CompoundStatement.BlockItemList; l != nil; l = l.BlockItemList {
			c.blockItem(w, l.BlockItem)
		}
	default:
		c.statement(w, n)
	}
}

func (c *ctx) iterationStatement(w writer, n *cc.IterationStatement) {
	defer c.setBreakCtx("")()
	defer c.setContinueCtx("")()

	if _, ok := c.f.flatScopes[n.LexicalScope()]; ok {
		c.iterationStatementFlat(w, n)
		return
	}

	switch n.Statement.Case {
	case cc.StatementCompound:
		if _, ok := c.f.flatScopes[n.Statement.CompoundStatement.LexicalScope()]; ok {
			c.iterationStatementFlat(w, n)
			return
		}
	}

	var cont string
	switch n.Case {
	case cc.IterationStatementWhile: // "while" '(' ExpressionList ')' Statement
		var a buf
		switch b := c.expr(&a, n.ExpressionList, nil, exprBool); {
		case a.len() != 0:
			w.w("for {")
			w.w("%s", a.bytes())
			w.w("\nif !(%s) { break };", b)
			c.unbracedStatement(w, n.Statement)
			w.w("\n};")
		default:
			w.w("for %s", b)
			c.bracedStatement(w, n.Statement)
		}
	case cc.IterationStatementDo: // "do" Statement "while" '(' ExpressionList ')' ';'
		if isZero(n.ExpressionList.Value()) {
			c.statement(w, n.Statement)
			break
		}

		var a buf
		switch b := c.expr(&a, n.ExpressionList, nil, exprBool); {
		case a.len() != 0:
			w.w("for %sfirst := true; ; %[1]sfirst = false {", tag(ccgo))
			w.w("\nif !first {")
			w.w("%s", a.bytes())
			w.w("\nif !(%s) { break };", b)
			w.w("\n};")
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		default:
			w.w("for %scond := true; %[1]scond; %[1]scond = %s", tag(ccgo), b)
			c.bracedStatement(w, n.Statement)
		}
	case cc.IterationStatementFor: // "for" '(' ExpressionList ';' ExpressionList ';' ExpressionList ')' Statement
		var a, a2, a3 buf
		var b2 []byte
		b := c.expr(&a, n.ExpressionList, nil, exprVoid)
		if n.ExpressionList2 != nil {
			b2 = c.expr(&a2, n.ExpressionList2, nil, exprBool).bytes()
		}
		b3 := c.expr(&a3, n.ExpressionList3, nil, exprVoid)
		if a3.len() != 0 {
			cont = c.label()
			defer c.setContinueCtx(cont)()
		}
		switch {
		case a.len() == 0 && a2.len() == 0 && a3.len() == 0:
			w.w("for %s; %s; %s {", b, b2, b3)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() == 0 && a2.len() == 0 && a3.len() != 0:
			w.w("for %s; %s;  {", b, b2)
			c.unbracedStatement(w, n.Statement)
			w.w("\ngoto %s; %[1]s:", cont)
			w.w("\n%s%s};", a3.bytes(), b3.bytes())
		case a.len() == 0 && a2.len() != 0 && a3.len() == 0:
			w.w("for %s; ; %s {", b, b3)
			w.w("\n%s", a2.bytes())
			w.w("\nif !(%s) { break };", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() == 0 && a2.len() != 0 && a3.len() != 0:
			w.w("for %s; ; {", b)
			w.w("\n%s", a2.bytes())
			w.w("\nif !(%s) { break };", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("\ngoto %s; %[1]s:", cont)
			w.w("\n%s%s};", a3.bytes(), b3.bytes())
		case a.len() != 0 && a2.len() == 0 && a3.len() == 0:
			w.w("%s%s", a.bytes(), b.bytes())
			w.w("\nfor ;%s; %s{", b2, b3)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() != 0 && a2.len() == 0 && a3.len() != 0:
			w.w("%s%s", a.bytes(), b.bytes())
			w.w("\nfor %s {", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("\ngoto %s; %[1]s:", cont)
			w.w("\n%s%s};", a3.bytes(), b3.bytes())
		case a.len() != 0 && a2.len() != 0 && a3.len() == 0:
			w.w("%s%s", a.bytes(), b.bytes())
			w.w("\nfor ; ; %s {", b3)
			w.w("\n%s", a2.bytes())
			w.w("\nif !(%s) { break };", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() != 0 && a2.len() != 0 && a3.len() != 0:
			w.w("%s%s", a.bytes(), b.bytes())
			w.w("\nfor ; ; %s {", b3)
			w.w("\n%s", a2.bytes())
			w.w("\nif !(%s) { break };", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("\ngoto %s; %[1]s:", cont)
			w.w("\n%s%s};", a3.bytes(), b3.bytes())
		}
	case cc.IterationStatementForDecl: // "for" '(' Declaration ExpressionList ';' ExpressionList ')' Statement
		c.declaration(w, n.Declaration, false)
		var a, a2 buf
		var b []byte
		if n.ExpressionList != nil {
			b = c.expr(&a, n.ExpressionList, nil, exprBool).bytes()
		}
		b2 := c.expr(&a2, n.ExpressionList2, nil, exprVoid)
		if a2.len() != 0 {
			cont = c.label()
			defer c.setContinueCtx(cont)()
		}
		switch {
		case a.len() == 0 && a2.len() == 0:
			w.w("for ; %s; %s {", b, b2)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() == 0 && a2.len() != 0:
			w.w("for %s  {", b)
			c.unbracedStatement(w, n.Statement)
			w.w("\ngoto %s; %[1]s:", cont)
			w.w("\n%s%s};", a2.bytes(), b2.bytes())
		case a.len() != 0 && a2.len() == 0:
			w.w("for ; ; %s {", b2)
			w.w("\n%s", a.bytes())
			w.w("\nif !(%s) { break };", b)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		default: // case a.len() != 0 && a2.len() != 0:
			w.w("for {")
			w.w("\n%s", a.bytes())
			w.w("\nif !(%s) { break };", b)
			c.unbracedStatement(w, n.Statement)
			w.w("\ngoto %s; %[1]s:", cont)
			w.w("\n%s%s};", a2.bytes(), b2.bytes())
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) iterationStatementFlat(w writer, n *cc.IterationStatement) {
	brk := c.label()
	cont := c.label()

	defer c.setBreakCtx(brk)()
	defer c.setContinueCtx(cont)()

	switch n.Case {
	case cc.IterationStatementWhile: // "while" '(' ExpressionList ')' Statement
		// cont:
		//	if !expr goto brk
		//	stmt
		//	goto cont
		// brk:
		w.w("%s: if !(%s) { goto %s };", cont, c.expr(w, n.ExpressionList, nil, exprBool), brk)
		c.unbracedStatement(w, n.Statement)
		w.w("goto %s; %s:", cont, brk)
	case cc.IterationStatementDo: // "do" Statement "while" '(' ExpressionList ')' ';'
		// cont:
		//	stmt
		//	if expr goto cont
		// brk:
		w.w("%s:", cont)
		c.unbracedStatement(w, n.Statement)
		w.w("if (%s) { goto %s }; goto %s; %[3]s:", c.expr(w, n.ExpressionList, nil, exprBool), cont, brk)
	case cc.IterationStatementFor: // "for" '(' ExpressionList ';' ExpressionList ';' ExpressionList ')' Statement
		//	expr1
		// a:	if !expr2 goto brk
		//	stmt
		// cont:
		//	expr3
		//	goto a
		// brk:
		a := c.label()
		c.expr(w, n.ExpressionList, nil, exprVoid)
		w.w("%s: if !(%s) { goto %s };", a, c.expr(w, n.ExpressionList2, nil, exprBool), brk)
		c.unbracedStatement(w, n.Statement)
		w.w("goto %s; %[1]s: %s; goto %s; %s:", cont, c.expr(w, n.ExpressionList3, nil, exprVoid), a, brk)
	case cc.IterationStatementForDecl: // "for" '(' Declaration ExpressionList ';' ExpressionList ')' Statement
		//	decl
		// a:	if !expr goto brk
		//	stmt
		// cont:
		//	expr2
		//	goto a
		// brk:
		a := c.label()
		c.declaration(w, n.Declaration, false)
		w.w("%s: if !(%s) { goto %s };", a, c.expr(w, n.ExpressionList, nil, exprBool), brk)
		c.unbracedStatement(w, n.Statement)
		w.w("goto %s; %[1]s: %s; goto %s; %s:", cont, c.expr(w, n.ExpressionList2, nil, exprVoid), a, brk)
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) label() string { return fmt.Sprintf("%s_%d", tag(ccgo), c.id()) }

func (c *ctx) jumpStatement(w writer, n *cc.JumpStatement) {
	switch n.Case {
	case cc.JumpStatementGoto: // "goto" IDENTIFIER ';'
		w.w("goto %s%s;", tag(preserve), n.Token2.Src())
	case cc.JumpStatementGotoExpr: // "goto" '*' ExpressionList ';'
		c.err(errorf("TODO %v", n.Case))
	case cc.JumpStatementContinue: // "continue" ';'
		if c.continueCtx != "" {
			w.w("goto %s;", c.continueCtx)
			break
		}

		w.w("continue;")
	case cc.JumpStatementBreak: // "break" ';'
		if c.breakCtx != "" {
			w.w("goto %s;", c.breakCtx)
			break
		}

		w.w("break;")
	case cc.JumpStatementReturn: // "return" ExpressionList ';'
		switch {
		case n.ExpressionList != nil:
			switch {
			case c.f.t.Result().Kind() == cc.Void:
				w.w("%s; return;", c.expr(w, n.ExpressionList, nil, exprVoid))
			default:
				w.w("return %s;", c.expr(w, n.ExpressionList, c.f.t.Result(), exprDefault))
			}
		default:
			w.w("return;")
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}
