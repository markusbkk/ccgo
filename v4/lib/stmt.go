// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
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
		c.compoundStatement(w, n.CompoundStatement, false)
	case cc.StatementExpr: // ExpressionStatement
		var a buf
		b := c.expr(&a, n.ExpressionStatement.ExpressionList, nil, exprVoid)
		if a.len() != 0 || b.len() != 0 {
			w.w("%s%s", sep, c.posComment(n))
			if a.len() != 0 {
				w.w("%s;", a.bytes())
			}
			w.w("%s;", b.bytes())
		}
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
		// c.err(errorf("%v: assembler statements not supported", c.pos(n)))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) labeledStatement(w writer, n *cc.LabeledStatement) {
	switch n.Case {
	case cc.LabeledStatementLabel: // IDENTIFIER ':' Statement
		c.err(errorf("TODO %v", n.Case))
	case cc.LabeledStatementCaseLabel: // "case" ConstantExpression ':' Statement
		if n.CaseOrdinal() != 0 {
			w.w("fallthrough;")
		}
		w.w("case %s:", c.expr(nil, n.ConstantExpression, c.switchExpr, exprDefault))
		c.statement(w, n.Statement)
	case cc.LabeledStatementRange: // "case" ConstantExpression "..." ConstantExpression ':' Statement
		if n.CaseOrdinal() != 0 {
			w.w("fallthrough;")
		}
		c.err(errorf("TODO %v", n.Case))
	case cc.LabeledStatementDefault: // "default" ':' Statement
		if n.CaseOrdinal() != 0 {
			w.w("fallthrough;")
		}
		w.w("default:")
		c.statement(w, n.Statement)
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) compoundStatement(w writer, n *cc.CompoundStatement, fnBlock bool) {
	switch {
	case fnBlock:
		w.w(" {")
		switch {
		case c.f.tlsAllocs+int64(c.f.maxValist) != 0:
			c.f.tlsAllocs = roundup(c.f.tlsAllocs, 8)
			v := c.f.tlsAllocs
			if c.f.maxValist != 0 {
				v += 8 * int64((c.f.maxValist + 1))
			}
			w.w("%sbp := %[1]stls.Alloc(%d); /* tlsAllocs %v maxValist %v */", tag(ccgo), v, c.f.tlsAllocs, c.f.maxValist)
			w.w("defer %stls.Free(%d);", tag(ccgo), v)
			for _, v := range c.f.t.Parameters() {
				if d := v.Declarator; d != nil && c.f.declInfos.info(d).pinned() {
					w.w("*(*%s)(%s) = %s_%s;", c.typ(d.Type()), unsafePointer(bpOff(c.f.declInfos.info(d).bpOff)), tag(ccgo), d.Name())
				}
			}
			w.w("%s%s", sep(n.Token), c.posComment(n))
		default:
			w.w("%s%s", strings.TrimSpace(sep(n.Token)), c.posComment(n))
		}
	default:
		w.w(" { %s%s", sep(n.Token), c.posComment(n))
	}
	var bi *cc.BlockItem
	for l := n.BlockItemList; l != nil; l = l.BlockItemList {
		bi = l.BlockItem
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
		se := c.switchExpr
		defer func() { c.switchExpr = se }()
		c.switchExpr = cc.IntegerPromotion(n.ExpressionList.Type())
		w.w("switch %s", c.expr(w, n.ExpressionList, c.switchExpr, exprDefault))
		c.statement(w, n.Statement)
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
		var a buf
		switch b := c.expr(&a, n.ExpressionList, nil, exprBool); {
		case a.len() != 0:
			w.w("for {")
			c.unbracedStatement(w, n.Statement)
			w.w("%s", a.bytes())
			w.w("\nif !(%s) { break };", b)
			w.w("\n};")
		default:
			w.w("for %scond := true; %[1]scond; %[1]scond = %s", tag(ccgo), b)
			c.bracedStatement(w, n.Statement)
		}
	case cc.IterationStatementFor: // "for" '(' ExpressionList ';' ExpressionList ';' ExpressionList ')' Statement
		var a, a2, a3 buf
		var b2 []byte
		if n.ExpressionList2 != nil {
			b2 = c.expr(&a2, n.ExpressionList2, nil, exprBool).bytes()
		}
		switch b, b3 := c.expr(&a, n.ExpressionList, nil, exprVoid), c.expr(&a3, n.ExpressionList3, nil, exprVoid); {
		case a.len() == 0 && a2.len() == 0 && a3.len() == 0:
			w.w("for %s; %s; %s {", b, b2, b3)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() == 0 && a2.len() == 0 && a3.len() != 0:
			w.w("for %s; %s;  {", b, b2)
			c.unbracedStatement(w, n.Statement)
			w.w("\n%s%s;}", a3.bytes(), b3.bytes())
			w.w("};")
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
			w.w("\n%s%s;}", a3.bytes(), b3.bytes())
		case a.len() == 0 && a2.len() != 0 && a3.len() == 0:
			w.w("for %s; ; %s {", b, b3)
			w.w("\n%s", a2.bytes())
			w.w("\nif !(%s) { break };", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() != 0 && a2.len() == 0 && a3.len() == 0:
			w.w("%s%s", a.bytes(), b.bytes())
			w.w("\nfor ;%s; %s{", b2, b3)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		case a.len() != 0 && a2.len() == 0 && a3.len() != 0:
			w.w("%s%s", a.bytes(), b.bytes())
			w.w("\nfor %s {", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("\n%s%s;}", a3.bytes(), b3.bytes())
		case a.len() != 0 && a2.len() != 0 && a3.len() == 0:
			w.w("%s%s", a.bytes(), b.bytes())
			w.w("\nfor ; ; %s {", b3)
			w.w("\n%s", a2.bytes())
			w.w("\nif !(%s) { break };", b2)
			c.unbracedStatement(w, n.Statement)
			w.w("};")
		default:
			trc("", c.pos(n))
			c.err(errorf("TODO"))
		}
	case cc.IterationStatementForDecl: // "for" '(' Declaration ExpressionList ';' ExpressionList ')' Statement
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) jumpStatement(w writer, n *cc.JumpStatement) {
	switch n.Case {
	case cc.JumpStatementGoto: // "goto" IDENTIFIER ';'
		c.err(errorf("TODO %v", n.Case))
	case cc.JumpStatementGotoExpr: // "goto" '*' ExpressionList ';'
		c.err(errorf("TODO %v", n.Case))
	case cc.JumpStatementContinue: // "continue" ';'
		w.w("continue;")
	case cc.JumpStatementBreak: // "break" ';'
		w.w("break;")
	case cc.JumpStatementReturn: // "return" ExpressionList ';'
		switch {
		case n.ExpressionList != nil:
			switch {
			case c.f.t.Result().Kind() == cc.Void:
				if n.ExpressionList.Type().Kind() != cc.Void {
					w.w("%s_ = ", tag(preserve))
				}
				w.w("%s; return;", c.expr(w, n.ExpressionList, c.f.t.Result(), exprDefault))
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
