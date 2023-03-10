// Copyright 2017 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v2"

import (
	"modernc.org/cc/v2"
	crtp "modernc.org/crt"
	"modernc.org/mathutil"
)

func (g *gen) compoundStmt(n *cc.CompoundStmt, vars []*cc.Declarator, cases map[*cc.LabeledStmt]int, sentinel bool, brk, cont *int, params, escParams []*cc.Declarator, main, value bool) {
	if vars != nil {
		g.w(" {")
	}
	vars = append([]*cc.Declarator(nil), vars...)
	w := 0
	for _, v := range vars {
		if v != allocaDeclarator {
			if v.Referenced == 0 && v.Initializer != nil && v.Linkage == cc.LinkageNone && v.DeclarationSpecifier.IsStatic() && v.Name() == idFuncName {
				continue
			}

			if v.Referenced == 0 && v.Initializer == nil && !v.AddressTaken {
				continue
			}

			if v.DeclarationSpecifier.IsStatic() || v.DeclarationSpecifier.IsExtern() {
				g.enqueueNumbered(v)
				continue
			}
		}

		vars[w] = v
		w++
	}
	vars = vars[:w]
	alloca := false
	var malloc int64
	var offp, offv []int64
	for _, v := range escParams {
		malloc = roundup(malloc, crtp.StackAlign)
		offp = append(offp, malloc)
		malloc += g.model.Sizeof(v.Type)
	}
	for _, v := range vars {
		if v == allocaDeclarator {
			continue
		}

		if isVLA(v) {
			continue
		}

		if g.escaped(v) {
			malloc = roundup(malloc, 16)
			offv = append(offv, malloc)
			//fmt.Printf("%v:\n", g.position(v)) //TODO- DBG
			malloc += g.model.Sizeof(v.Type)
		}
	}
	if malloc != 0 {
		malloc = roundup(malloc, crtp.StackAlign)
		if malloc > mathutil.MaxInt {
			todo("", g.position(n))
		}

		g.allocatedStack = int(malloc)
		g.w("\nesc := %sMallocStack(tls, %d)", g.crtPrefix, malloc)
	}
	if len(vars)+len(escParams) != 0 {
		localNames := map[int]struct{}{}
		num := 0
		for _, v := range append(params, vars...) {
			if v == nil || v == allocaDeclarator {
				continue
			}

			nm := v.Name()
			if _, ok := localNames[nm]; ok {
				num++
				g.nums[v] = num
			}
			localNames[nm] = struct{}{}
		}
		switch {
		case len(vars)+len(escParams) == 1:
			g.w("\nvar ")
		default:
			g.w("\nvar (\n")
		}
		for i, v := range escParams {
			g.w("\n\t%s = esc+%d // *%s", g.mangleDeclarator(v), offp[i], g.ptyp(v.Type, false, 1))
		}
		for _, v := range vars {
			switch {
			case v == allocaDeclarator:
				g.w("\nallocs []uintptr")
				g.needAlloca = true
				alloca = true
			case g.escaped(v):
				if isVLA(v) {
					g.w("\n\t%s uintptr // %s", g.mangleDeclarator(v), g.typeComment(v.Type))
					g.w("\n\t_ = %s", g.mangleDeclarator(v))
					break
				}

				g.w("\n\t%s = esc+%d // *%s", g.mangleDeclarator(v), offv[0], g.typeComment(v.Type))
				g.w("\n\t_ = %s", g.mangleDeclarator(v))
				offv = offv[1:]
			default:
				switch {
				case v.Type.Kind() == cc.Ptr:
					g.w("\n\t%s %s\t// %s", g.mangleDeclarator(v), g.typ(v.Type), g.typeComment(v.Type))
				default:
					g.w("\n\t%s %s", g.mangleDeclarator(v), g.typ(v.Type))
				}
				g.w("\n\t_ = %s", g.mangleDeclarator(v))
			}
		}
		if len(vars)+len(escParams) != 1 {
			g.w("\n)")
		}
	}
	if alloca {
		g.w("\ndefer func() {")
		g.w(`
for _, v := range allocs {
	%sFree(v)
}`, g.crtPrefix)
		g.w("\n}()")
	}
	for _, v := range escParams {
		g.w("\n*(*%s)(unsafe.Pointer(%s)) = a%s", g.typ(v.Type), g.mangleDeclarator(v), dict.S(v.Name()))
	}
	returned := g.blockItemListOpt(n.BlockItemListOpt, cases, brk, cont, main, value)
	if !returned && malloc != 0 {
		g.w(";%sFreeStack(tls, %d)", g.crtPrefix, malloc)
	}
	if vars != nil {
		if sentinel && !returned {
			g.w(";return r")
		}
		g.w("\n}")
	}
}

func (g *gen) blockItemListOpt(n *cc.BlockItemListOpt, cases map[*cc.LabeledStmt]int, brk, cont *int, main, value bool) (returned bool) {
	if n == nil {
		return false
	}

	return g.blockItemList(n.BlockItemList, cases, brk, cont, main, value)
}

func (g *gen) blockItemList(n *cc.BlockItemList, cases map[*cc.LabeledStmt]int, brk, cont *int, main, value bool) (returned bool) {
	for ; n != nil; n = n.BlockItemList {
		returned = g.blockItem(n.BlockItem, cases, brk, cont, main, value && n.BlockItemList == nil)
	}
	return returned
}

func (g *gen) blockItem(n *cc.BlockItem, cases map[*cc.LabeledStmt]int, brk, cont *int, main, value bool) (returned bool) {
	switch n.Case {
	case cc.BlockItemDecl: // Declaration
		g.declaration(n.Declaration, false)
	case cc.BlockItemStmt: // Stmt
		if g.tweaks.Watch {
			g.w("\n%sWatch(tls)", g.crtPrefix)
		}
		returned = g.stmt(n.Stmt, cases, brk, cont, main, value)
	default:
		todo("", g.position(n), n.Case)
	}
	return returned
}

func (g *gen) stmt(n *cc.Stmt, cases map[*cc.LabeledStmt]int, brk, cont *int, main, value bool) (returned bool) {
	//g.w("\n/* %v: %v %v */\n", g.position(n), n.Case, n.UseGotos)
	switch n.Case {
	case cc.StmtExpr: // ExprStmt
		g.exprStmt(n.ExprStmt, value)
	case cc.StmtJump: // JumpStmt
		returned = g.jumpStmt(n.JumpStmt, brk, cont, main)
	case cc.StmtIter: // IterationStmt
		g.iterationStmt(n.IterationStmt, cases, brk, cont, main, value)
	case cc.StmtBlock: // CompoundStmt
		g.compoundStmt(n.CompoundStmt, nil, cases, false, brk, cont, nil, nil, main, false)
	case cc.StmtSelect: // SelectionStmt
		g.selectionStmt(n.SelectionStmt, cases, brk, cont, main, value)
	case cc.StmtLabeled: // LabeledStmt
		returned = g.labeledStmt(n.LabeledStmt, cases, brk, cont, main, value)
	default:
		todo("", g.position(n), n.Case)
	}
	return returned
}

func (g *gen) labeledStmt(n *cc.LabeledStmt, cases map[*cc.LabeledStmt]int, brk, cont *int, main, value bool) (returned bool) {
	switch n.Case {
	case
		cc.LabeledStmtSwitchCase, // "case" ConstExpr ':' Stmt
		cc.LabeledStmtDefault:    // "default" ':' Stmt

		l, ok := cases[n]
		if !ok {
			todo("", g.position(n))
		}
		g.w("\n_%d:", l)
		g.stmt(n.Stmt, cases, brk, cont, main, value)
	case cc.LabeledStmtLabel: // IDENTIFIER ':' Stmt
		g.w("\ngoto %[1]s;%[1]s:\n", mangleLabel(n.Token.Val))
		returned = g.stmt(n.Stmt, cases, brk, cont, main, value)
	case cc.LabeledStmtLabel2: // TYPEDEF_NAME ':' Stmt
		g.w("\ngoto %[1]s;%[1]s:\n", mangleLabel(n.Token.Val))
		returned = g.stmt(n.Stmt, cases, brk, cont, main, value)
	default:
		todo("", g.position(n), n.Case)
	}
	return returned
}

func (g *gen) selectionStmt(n *cc.SelectionStmt, cases map[*cc.LabeledStmt]int, brk, cont *int, main, value bool) {
	switch n.Case {
	case cc.SelectionStmtSwitch: // "switch" '(' ExprList ')' Stmt
		if n.ExprList.Operand.Value != nil && g.voidCanIgnoreExprList(n.ExprList) {
			//TODO optimize
		}
		g.w("\nswitch ")
		switch el := n.ExprList; {
		case isSingleExpression(el):
			g.convert(n.ExprList.Expr, n.SwitchOp.Type)
		default:
			todo("", g.position(n))
		}
		g.w("{")
		after := -g.local()
		cases := map[*cc.LabeledStmt]int{}
		var deflt *cc.LabeledStmt
		for _, v := range n.Cases {
			l := g.local()
			cases[v] = l
			switch ce := v.ConstExpr; {
			case ce != nil:
				g.w("\ncase ")
				g.convert(ce.Expr, n.SwitchOp.Type)
				g.w(": goto _%d", l)
			default:
				deflt = v
				g.w("\ndefault: goto _%d\n", l)
			}
		}
		g.w("}")
		if deflt == nil {
			after = -after
			g.w("\ngoto _%d\n", after)
		}
		g.stmt(n.Stmt, cases, &after, cont, main, value)
		if after > 0 {
			g.w("\n_%d:", after)
		}
	case cc.SelectionStmtIf: // "if" '(' ExprList ')' Stmt
		g.w("\n")
		if e := n.ExprList; g.voidCanIgnoreExprList(e) {
			if e.IsZero() {
				if !n.UseGotos {
					g.exprList(e, true, false)
					//TODO- g.w("if false {\n")
					//TODO- g.stmt(n.Stmt, cases, brk, cont, main, value)
					//TODO- g.w("}")
					break
				}

				a := g.local()
				g.exprList(e, true, false)
				g.w("\ngoto _%d\n", a)
				g.stmt(n.Stmt, cases, brk, cont, main, value)
				g.w("\n_%d:", a)
				break
			}

			if e.IsNonZero() {
				g.exprList(e, true, false)
				g.stmt(n.Stmt, cases, brk, cont, main, value)
				break
			}
		}

		if !n.UseGotos {
			g.w("if ")
			g.exprList(n.ExprList, false, true)
			g.w(" != 0 {\n")
			g.stmt(n.Stmt, cases, brk, cont, main, value)
			g.w("}")
			break
		}

		// if exprList == 0 { goto A }
		// stmt
		// A:
		a := g.local()
		g.w("if ")
		g.exprList(n.ExprList, false, true)
		g.w(" == 0 { goto _%d }\n", a)
		g.stmt(n.Stmt, cases, brk, cont, main, value)
		g.w("\n_%d:", a)
	case cc.SelectionStmtIfElse: // "if" '(' ExprList ')' Stmt "else" Stmt
		g.w("\n")
		if e := n.ExprList; g.voidCanIgnoreExprList(e) {
			if e.IsZero() {
				if !n.UseGotos {
					g.exprList(n.ExprList, true, false)
					//TODO- g.w("if false {")
					//TODO- g.stmt(n.Stmt, cases, brk, cont, main, value)
					//TODO- g.w("} else {")
					g.stmt(n.Stmt2, cases, brk, cont, main, value)
					//TODO- g.w("}")
					break
				}

				a := g.local()
				b := g.local()
				g.exprList(n.ExprList, true, false)
				g.w("\ngoto _%d\n", a)
				g.stmt(n.Stmt, cases, brk, cont, main, value)
				g.w("\ngoto _%d\n", b)
				g.w("\n_%d:", a)
				g.stmt(n.Stmt2, cases, brk, cont, main, value)
				g.w("\n_%d:", b)
				break
			}

			if e.IsNonZero() {
				if !n.UseGotos {
					g.exprList(n.ExprList, true, false)
					//TODO- g.w("if true {")
					g.stmt(n.Stmt, cases, brk, cont, main, value)
					//TODO- g.w("} else {")
					//TODO- g.stmt(n.Stmt2, cases, brk, cont, main, value)
					//TODO- g.w("}")
					break
				}

				a := g.local()
				g.exprList(n.ExprList, true, false)
				g.stmt(n.Stmt, cases, brk, cont, main, value)
				g.w("\ngoto _%d\n", a)
				g.stmt(n.Stmt2, cases, brk, cont, main, value)
				g.w("\n_%d:", a)
				break
			}
		}

		if !n.UseGotos {
			g.w("if ")
			g.exprList(n.ExprList, false, true)
			g.w(" != 0 {\n")
			g.stmt(n.Stmt, cases, brk, cont, main, value)
			g.w("} else {\n")
			g.stmt(n.Stmt2, cases, brk, cont, main, value)
			g.w("}")
			break
		}

		// if exprList == 0 { goto A }
		// stmt
		// goto B
		// A:
		// stmt2
		// B:
		a := g.local()
		b := g.local()
		g.w("if ")
		g.exprList(n.ExprList, false, true)
		g.w(" == 0 { goto _%d }\n", a)
		g.stmt(n.Stmt, cases, brk, cont, main, value)
		g.w("\ngoto _%d\n", b)
		g.w("\n_%d:", a)
		g.stmt(n.Stmt2, cases, brk, cont, main, value)
		g.w("\n_%d:", b)
	default:
		todo("", g.position(n), n.Case)
	}
}

func (g *gen) iterationStmt(n *cc.IterationStmt, cases map[*cc.LabeledStmt]int, brk, cont *int, main, value bool) {
	switch n.Case {
	case cc.IterationStmtDo: // "do" Stmt "while" '(' ExprList ')' ';'
		if e := n.ExprList; g.voidCanIgnoreExprList(e) {
			if e.IsZero() {
				// stmt
				// A: <- continue, break
				a := -g.local()
				g.stmt(n.Stmt, cases, &a, &a, main, value)
				g.exprList(e, true, false)
				if a > 0 {
					g.w("\n_%d:", a)
				}
				return
			}

			if e.IsNonZero() {
				if !n.UseGotos {
					g.w("\nfor {")
					g.stmt(n.Stmt, cases, nil, nil, main, value)
					g.exprList(e, true, false)
					g.w("}")
					break
				}

				// A: <-continue
				// stmt
				// goto A
				// B: <- break
				a := g.local()
				b := -g.local()
				g.w("\n_%d:", a)
				g.stmt(n.Stmt, cases, &b, &a, main, value)
				g.exprList(e, true, false)
				g.w("\ngoto _%d\n", a)
				if b > 0 {
					g.w("\n_%d:", b)
				}
				return
			}
		}

		if !n.UseGotos {
			g.w("\nfor c := true; c ; c = ")
			g.exprList(n.ExprList, false, true)
			g.w(" != 0 {")
			g.stmt(n.Stmt, cases, nil, nil, main, value)
			g.w("}")
			break
		}

		// A:
		// stmt
		// B: <- continue
		// if exprList != 0 { goto A }
		// C: <- break
		a := g.local()
		b := -g.local()
		c := -g.local()
		g.w("\n_%d:", a)
		g.stmt(n.Stmt, cases, &c, &b, main, value)
		if b > 0 {
			g.w("\n_%d:", b)
		}
		g.w("\nif ")
		g.exprList(n.ExprList, false, true)
		g.w(" != 0 { goto _%d }\n", a)
		if c > 0 {
			g.w("\n_%d:", c)
		}
	case cc.IterationStmtForDecl: // "for" '(' Declaration ExprListOpt ';' ExprListOpt ')' Stmt
		if n.ExprListOpt == nil || n.ExprListOpt.ExprList.IsNonZero() && g.voidCanIgnoreExprList(n.ExprListOpt.ExprList) {
			if !n.UseGotos {
				g.w("\nfor ")
				g.declaration(n.Declaration, false)
				g.w(";")
				g.w(";")
				if n.ExprListOpt2 != nil {
					g.exprListOpt(n.ExprListOpt2, isSingleExpression(n.ExprListOpt2.ExprList), true)
				}
				g.w("{")
				g.stmt(n.Stmt, cases, nil, nil, main, value)
				g.w("}")
				break
			}

			// Declaration
			// A:
			// Stmt
			// B: <- continue
			// ExprListOpt2
			// goto A
			// C: <- break
			g.w("\n")
			g.declaration(n.Declaration, false)
			a := g.local()
			b := -g.local()
			c := -g.local()
			g.w("\n_%d:", a)
			g.stmt(n.Stmt, cases, &c, &b, main, value)
			if n.ExprListOpt2 != nil {
				g.w("\n")
			}
			if b > 0 {
				g.w("\n_%d:", b)
			}
			g.exprListOpt(n.ExprListOpt2, true, true)
			g.w("\ngoto _%d\n", a)
			if c > 0 {
				g.w("\n_%d:", c)
			}
			return
		}

		if n.ExprListOpt != nil && n.ExprListOpt.ExprList.IsZero() && g.voidCanIgnoreExprList(n.ExprListOpt.ExprList) {
			if !n.UseGotos {
				g.w("\nfor ")
				g.declaration(n.Declaration, true)
				g.w("; false ;")
				if n.ExprListOpt2 != nil {
					g.exprListOpt(n.ExprListOpt2, isSingleExpression(n.ExprListOpt2.ExprList), true)
				}
				g.w("{")
				g.stmt(n.Stmt, cases, nil, nil, main, value)
				g.w("}")
				break
			}

			// Declaration
			// goto A
			// Stmt
			// B: <- continue
			// ExprListOpt2
			// A: <- break
			g.w("\n")
			g.declaration(n.Declaration, false)
			a := g.local()
			b := -g.local()
			g.w("\ngoto _%d:", a)
			g.stmt(n.Stmt, cases, &a, &b, main, value)
			if n.ExprListOpt2 != nil {
				g.w("\n")
			}
			if b > 0 {
				g.w("\n_%d:", b)
			}
			g.exprListOpt(n.ExprListOpt2, true, true)
			g.w("\n_%d:", a)
			return
		}

		if !n.UseGotos {
			g.w("\nfor ")
			g.declaration(n.Declaration, true)
			g.w(";")
			if n.ExprListOpt != nil {
				g.exprList(n.ExprListOpt.ExprList, false, false)
				g.w(" != 0")
			}
			g.w(";")
			if n.ExprListOpt2 != nil {
				g.exprListOpt(n.ExprListOpt2, isSingleExpression(n.ExprListOpt2.ExprList), true)
			}
			g.w("{")
			g.stmt(n.Stmt, cases, nil, nil, main, value)
			g.w("}")
			break
		}

		// Declaration
		// A:
		// if ExprListOpt == 0 { goto C }
		// Stmt
		// B: <- continue
		// ExprListOpt2
		// goto A
		// C: <- break
		g.w("\n")
		g.declaration(n.Declaration, false)
		a := g.local()
		b := -g.local()
		c := -g.local()
		g.w("\n_%d:", a)
		if n.ExprListOpt != nil {
			g.w("if ")
			g.exprList(n.ExprListOpt.ExprList, false, false)
			c = -c
			g.w(" == 0 { goto _%d }\n", c)
		}
		g.stmt(n.Stmt, cases, &c, &b, main, value)
		if n.ExprListOpt2 != nil {
			g.w("\n")
		}
		if b > 0 {
			g.w("\n_%d:", b)
		}
		g.exprListOpt(n.ExprListOpt2, true, false)
		g.w("\ngoto _%d\n", a)
		if c > 0 {
			g.w("\n_%d:", c)
		}
	case cc.IterationStmtFor: // "for" '(' ExprListOpt ';' ExprListOpt ';' ExprListOpt ')' Stmt
		if n.ExprListOpt2 == nil || n.ExprListOpt2.ExprList.IsNonZero() && g.voidCanIgnoreExprList(n.ExprListOpt2.ExprList) {
			if !n.UseGotos {
				g.w("\nfor ")
				if n.ExprListOpt != nil {
					g.exprListOpt(n.ExprListOpt, isSingleExpression(n.ExprListOpt.ExprList), true)
				}
				g.w(";;")
				if n.ExprListOpt3 != nil {
					g.exprListOpt(n.ExprListOpt3, isSingleExpression(n.ExprListOpt3.ExprList), true)
				}
				g.w(" {")
				g.stmt(n.Stmt, cases, nil, nil, main, value)
				g.w("}")
				break
			}

			// ExprListOpt
			// A:
			// Stmt
			// B: <- continue
			// ExprListOpt3
			// goto A
			// C: <- break
			g.w("\n")
			g.exprListOpt(n.ExprListOpt, true, false)
			a := g.local()
			b := -g.local()
			c := -g.local()
			g.w("\n_%d:", a)
			g.stmt(n.Stmt, cases, &c, &b, main, value)
			if n.ExprListOpt3 != nil {
				g.w("\n")
			}
			if b > 0 {
				g.w("\n_%d:", b)
			}
			g.exprListOpt(n.ExprListOpt3, true, false)
			g.w("\ngoto _%d\n", a)
			if c > 0 {
				g.w("\n_%d:", c)
			}
			return
		}

		if n.ExprListOpt2 != nil && n.ExprListOpt2.ExprList.IsZero() && g.voidCanIgnoreExprList(n.ExprListOpt2.ExprList) {
			if !n.UseGotos {
				g.w("\nfor ")
				if n.ExprListOpt != nil {
					g.exprListOpt(n.ExprListOpt, isSingleExpression(n.ExprListOpt.ExprList), true)
				}
				g.w("; false ;")
				if n.ExprListOpt3 != nil {
					g.exprListOpt(n.ExprListOpt3, isSingleExpression(n.ExprListOpt3.ExprList), true)
				}
				g.w(" {")
				g.stmt(n.Stmt, cases, nil, nil, main, value)
				g.w("}")
				break
			}

			// ExprListOpt
			// A:
			// goto C
			// Stmt
			// B: <- continue
			// ExprListOpt3
			// goto A
			// C: <- break
			g.w("\n")
			g.exprListOpt(n.ExprListOpt, true, false)
			a := g.local()
			b := -g.local()
			c := -g.local()
			g.w("\n_%d:", a)
			g.w("\ngoto _%d }\n", c)
			g.stmt(n.Stmt, cases, &c, &b, main, value)
			if n.ExprListOpt3 != nil {
				g.w("\n")
			}
			if b > 0 {
				g.w("\n_%d:", b)
			}
			g.exprListOpt(n.ExprListOpt3, true, false)
			g.w("\ngoto _%d\n", a)
			if c > 0 {
				g.w("\n_%d:", c)
			}
			return
		}

		if !n.UseGotos {
			g.w("\nfor ")
			if n.ExprListOpt != nil {
				g.exprListOpt(n.ExprListOpt, isSingleExpression(n.ExprListOpt.ExprList), true)
			}
			g.w(";")
			if n.ExprListOpt2 != nil {
				g.exprList(n.ExprListOpt2.ExprList, false, true)
				g.w(" != 0")
			}
			g.w(";")
			if n.ExprListOpt3 != nil {
				g.exprListOpt(n.ExprListOpt3, isSingleExpression(n.ExprListOpt3.ExprList), true)
			}
			g.w(" {")
			g.stmt(n.Stmt, cases, nil, nil, main, value)
			g.w("}")
			break
		}

		// ExprListOpt
		// A:
		// if ExprListOpt2 == 0 { goto C }
		// Stmt
		// B: <- continue
		// ExprListOpt3
		// goto A
		// C: <- break
		g.w("\n")
		g.exprListOpt(n.ExprListOpt, true, false)
		a := g.local()
		b := -g.local()
		c := -g.local()
		g.w("\n_%d:", a)
		if n.ExprListOpt2 != nil {
			g.w("if ")
			g.exprList(n.ExprListOpt2.ExprList, false, true)
			c = -c
			g.w(" == 0 { goto _%d }\n", c)
		}
		g.stmt(n.Stmt, cases, &c, &b, main, value)
		if n.ExprListOpt3 != nil {
			g.w("\n")
		}
		if b > 0 {
			g.w("\n_%d:", b)
		}
		g.exprListOpt(n.ExprListOpt3, true, false)
		g.w("\ngoto _%d\n", a)
		if c > 0 {
			g.w("\n_%d:", c)
		}
	case cc.IterationStmtWhile: // "while" '(' ExprList ')' Stmt
		if e := n.ExprList; g.voidCanIgnoreExprList(e) {
			if e.IsZero() {
				if !n.UseGotos {
					g.w("\nfor false {")
					g.stmt(n.Stmt, cases, nil, nil, main, value)
					g.w("}")
					break
				}

				//TODO todo("", g.position(n))
				//TODO break
			}

			if e.IsNonZero() {
				if !n.UseGotos {
					g.w("\nfor {")
					g.stmt(n.Stmt, cases, nil, nil, main, value)
					g.w("}")
					break
				}

				// A:
				// exprList
				// stmt
				// goto A
				// B:
				a := g.local()
				b := -g.local()
				g.w("\n_%d:", a)
				g.exprList(n.ExprList, true, false)
				g.stmt(n.Stmt, cases, &b, &a, main, value)
				g.w("\ngoto _%d\n", a)
				if b > 0 {
					g.w("\n_%d:", b)
				}
				return
			}
		}

		if !n.UseGotos {
			g.w("\nfor ")
			g.exprList(n.ExprList, false, true)
			g.w(" != 0 {")
			g.stmt(n.Stmt, cases, nil, nil, main, value)
			g.w("}")
			break
		}

		// A:
		// if exprList == 0 { goto B }
		// stmt
		// goto A
		// B:
		a := g.local()
		b := g.local()
		g.w("\n_%d:\nif ", a)
		g.exprList(n.ExprList, false, true)
		g.w(" == 0 { goto _%d }\n", b)
		g.stmt(n.Stmt, cases, &b, &a, main, value)
		g.w("\ngoto _%d\n\n_%d:", a, b)
	default:
		todo("", g.position(n), n.Case)
	}
}

func (g *gen) local() int {
	r := g.nextLabel
	g.nextLabel++
	return r
}

func (g *gen) jumpStmt(n *cc.JumpStmt, brk, cont *int, main bool) (returned bool) {
	if main {
		n.ReturnOperand.Type = cc.Int
	}
	switch n.Case {
	case cc.JumpStmtReturn: // "return" ExprListOpt ';'
		switch o := n.ExprListOpt; {
		case o != nil:
			switch rt := n.ReturnOperand.Type; {
			case rt == nil:
				switch {
				case isSingleExpression(o.ExprList) && o.ExprList.Expr.Case == cc.ExprCond:
					todo("", g.position(n))
				default:
					g.exprList(o.ExprList, true, false)
					if g.allocatedStack != 0 {
						g.w("\n%sFreeStack(tls, %d)", g.crtPrefix, g.allocatedStack)
					}
					g.w("\nreturn")
				}
			default:
				switch {
				case g.allocatedStack != 0:
					switch {
					case isSingleExpression(o.ExprList) && o.ExprList.Expr.Case == cc.ExprCond:
						n := o.ExprList.Expr // Expr '?' ExprList ':' Expr
						switch {
						case n.Expr.IsZero() && g.voidCanIgnore(n.Expr):
							g.w("\nr = ")
							g.convert(n.Expr2, rt)
							g.w("\n%sFreeStack(tls, %d)", g.crtPrefix, g.allocatedStack)
							g.w("\nreturn r")
						case n.Expr.IsNonZero() && g.voidCanIgnore(n.Expr):
							g.w("\nr = ")
							g.exprList2(n.ExprList, rt)
							g.w("\n%sFreeStack(tls, %d)", g.crtPrefix, g.allocatedStack)
							g.w("\nreturn r;")
						default:
							g.w("\nif ")
							g.value(n.Expr, false)
							g.w(" != 0 { r = ")
							g.exprList2(n.ExprList, rt)
							g.w(";%sFreeStack(tls, %d)", g.crtPrefix, g.allocatedStack)
							g.w("; return r }")
							g.w("\nr = ")
							g.convert(n.Expr2, rt)
							g.w("\n%sFreeStack(tls, %d)", g.crtPrefix, g.allocatedStack)
							g.w("\nreturn r")
						}
					default:
						g.w("\nr = ")
						switch {
						case isSingleExpression(o.ExprList) && o.ExprList.Operand.Value != nil && o.ExprList.Operand.Type.Equal(rt) && g.voidCanIgnoreExprList(o.ExprList):
							g.constant(o.ExprList.Expr)
						default:
							g.exprList2(o.ExprList, rt)
						}
						g.w("\n%sFreeStack(tls, %d)", g.crtPrefix, g.allocatedStack)
						g.w("\nreturn r")
					}
				default:
					switch {
					case isSingleExpression(o.ExprList) && o.ExprList.Expr.Case == cc.ExprCond:
						n := o.ExprList.Expr // Expr '?' ExprList ':' Expr
						switch {
						case n.Expr.IsZero() && g.voidCanIgnore(n.Expr):
							g.w("\nreturn ")
							g.convert(n.Expr2, rt)
						case n.Expr.IsNonZero() && g.voidCanIgnore(n.Expr):
							g.w("\nreturn ")
							g.exprList2(n.ExprList, rt)
							g.w("\n")
						default:
							g.w("\nif ")
							g.value(n.Expr, false)
							g.w(" != 0 { return ")
							g.exprList2(n.ExprList, rt)
							g.w("}\nreturn ")
							g.convert(n.Expr2, rt)
						}
					default:
						g.w("\nreturn ")
						switch {
						case isSingleExpression(o.ExprList) && o.ExprList.Operand.Value != nil && o.ExprList.Operand.Type.Equal(rt) && g.voidCanIgnoreExprList(o.ExprList):
							g.constant(o.ExprList.Expr)
						default:
							g.exprList2(o.ExprList, rt)
						}
					}
				}
			}
		default:
			if g.allocatedStack != 0 {
				g.w("\n%sFreeStack(tls, %d)", g.crtPrefix, g.allocatedStack)
			}
			g.w("\nreturn ")
		}
		returned = true
	case cc.JumpStmtBreak: // "break" ';'
		if brk == nil {
			g.w("\nbreak")
			break
		}

		if *brk < 0 {
			*brk = -*brk // Signal used.
		}
		g.w("\ngoto _%d\n", *brk)
	case cc.JumpStmtGoto: // "goto" IDENTIFIER ';'
		g.w("\ngoto %s\n", mangleLabel(n.Token2.Val))
	case cc.JumpStmtContinue: // "continue" ';'
		if cont == nil {
			g.w("\ncontinue")
			break
		}

		if *cont < 0 {
			*cont = -*cont // Signal used.
		}
		g.w("\ngoto _%d\n", *cont)
	default:
		todo("", g.position(n), n.Case)
	}
	return returned
}

func (g *gen) exprStmt(n *cc.ExprStmt, value bool) {
	if o := n.ExprListOpt; o != nil {
		g.w("\n")
		if value {
			g.w("return ")
		}
		g.exprList(o.ExprList, !value, false)
	}
}
