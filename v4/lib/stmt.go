// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"modernc.org/cc/v4"
)

func (c *ctx) statement(w writer, n *cc.Statement) {
	switch n.Case {
	case cc.StatementLabeled: // LabeledStatement
		c.labeledStatement(w, n.LabeledStatement)
	case cc.StatementCompound: // CompoundStatement
		c.compoundStatement(w, n.CompoundStatement, false)
	case cc.StatementExpr: // ExpressionStatement
		c.expressionStatement(w, n.ExpressionStatement)
	case cc.StatementSelection: // SelectionStatement
		c.selectionStatement(w, n.SelectionStatement)
	case cc.StatementIteration: // IterationStatement
		c.iterationStatement(w, n.IterationStatement)
	case cc.StatementJump: // JumpStatement
		c.jumpStatement(w, n.JumpStatement)
	case cc.StatementAsm: // AsmStatement
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	w.w(" // %v:", c.pos(n))
}

func (c *ctx) labeledStatement(w writer, n *cc.LabeledStatement) {
	switch n.Case {
	case cc.LabeledStatementLabel: // IDENTIFIER ':' Statement
		c.err(errorf("TODO %v", n.Case))
	case cc.LabeledStatementCaseLabel: // "case" ConstantExpression ':' Statement
		if n.CaseOrdinal() != 0 {
			w.w("\nfallthrough")
		}
		var a1 buf
		b1 := c.expr(&a1, n.ConstantExpression, nil, expr)
		switch {
		case a1.len() == 0:
			w.w("\ncase %s:", b1)
			c.statement(w, n.Statement)
		default:
			c.err(errorf("TODO %v", n.Case))
		}
	case cc.LabeledStatementRange: // "case" ConstantExpression "..." ConstantExpression ':' Statement
		if n.CaseOrdinal() != 0 {
			w.w("\nfallthrough")
		}
		c.err(errorf("TODO %v", n.Case))
	case cc.LabeledStatementDefault: // "default" ':' Statement
		if n.CaseOrdinal() != 0 {
			w.w("\nfallthrough")
		}
		w.w("\ndefault:")
		c.statement(w, n.Statement)
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) compoundStatement(w writer, n *cc.CompoundStatement, fnBlock bool) {
	w.w(" { // %v:", c.pos(n))
	if fnBlock && c.f.tlsAllocs+int64(c.f.maxValist) != 0 {
		c.f.tlsAllocs = roundup(c.f.tlsAllocs, 8)
		v := c.f.tlsAllocs
		if c.f.maxValist != 0 {
			v += 8 * int64((c.f.maxValist + 1))
		}
		w.w("\n%sbp := %[1]stls.Alloc(%d)", tag(ccgoAutomatic), v)
		w.w("\ndefer %stls.Free(%d)", tag(ccgoAutomatic), v)
	}
	for l := n.BlockItemList; l != nil; l = l.BlockItemList {
		c.blockItem(w, l.BlockItem)
	}
	w.w("\n}")
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
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) selectionStatement(w writer, n *cc.SelectionStatement) {
	switch n.Case {
	case cc.SelectionStatementIf: // "if" '(' ExpressionList ')' Statement
		var a1 buf
		b1 := c.expr(&a1, n.ExpressionList, nil, exprBool)
		switch {
		case a1.len() == 0:
			w.w("\nif %s {", b1)
			c.statement(w, n.Statement)
			w.w("\n}")
		default:
			c.err(errorf("TODO %v", n.Case))
		}
	case cc.SelectionStatementIfElse: // "if" '(' ExpressionList ')' Statement "else" Statement
		var a1 buf
		b1 := c.expr(&a1, n.ExpressionList, nil, exprBool)
		switch {
		case a1.len() == 0:
			w.w("\nif %s {", b1)
			c.statement(w, n.Statement)
			w.w("\n} else {")
			c.statement(w, n.Statement2)
			w.w("\n}")
		default:
			c.err(errorf("TODO %v", n.Case))
		}
	case cc.SelectionStatementSwitch: // "switch" '(' ExpressionList ')' Statement
		var a1 buf
		b1 := c.expr(&a1, n.ExpressionList, nil, expr)
		switch {
		case a1.len() == 0:
			w.w("\nswitch %s", b1)
			c.statement(w, n.Statement)
		default:
			c.err(errorf("TODO %v", n.Case))
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) iterationStatement(w writer, n *cc.IterationStatement) {
	switch n.Case {
	case cc.IterationStatementWhile: // "while" '(' ExpressionList ')' Statement
		var a1 buf
		b1 := c.expr(&a1, n.ExpressionList, nil, exprBool)
		switch {
		case a1.len() == 0:
			w.w("\nfor %s {", b1)
			switch n.Statement.Case {
			case cc.StatementCompound:
				for l := n.Statement.CompoundStatement.BlockItemList; l != nil; l = l.BlockItemList {
					c.blockItem(w, l.BlockItem)
				}
			default:
				c.statement(w, n.Statement)
			}
			w.w("\n}")
		default:
			c.err(errorf("TODO %v", n.Case))
		}
	case cc.IterationStatementDo: // "do" Statement "while" '(' ExpressionList ')' ';'
		var a1 buf
		b1 := c.expr(&a1, n.ExpressionList, nil, exprBool)
		switch {
		case a1.len() == 0:
			w.w("\nfor %scond := true; %[1]scond; %[1]scond = %s {", tag(ccgoAutomatic), b1)
			switch n.Statement.Case {
			case cc.StatementCompound:
				for l := n.Statement.CompoundStatement.BlockItemList; l != nil; l = l.BlockItemList {
					c.blockItem(w, l.BlockItem)
				}
			default:
				c.statement(w, n.Statement)
			}
			w.w("\n}")
		default:
			c.err(errorf("TODO %v", n.Case))
		}
	case cc.IterationStatementFor: // "for" '(' ExpressionList ';' ExpressionList ';' ExpressionList ')' Statement
		var a1, a2, a3 buf
		var b1, b2, b3 []byte
		if n.ExpressionList != nil {
			b1 = c.expr(&a1, n.ExpressionList, nil, exprVoid)
		}
		switch {
		case a1.len() == 0:
			if n.ExpressionList2 != nil {
				b2 = c.expr(&a2, n.ExpressionList2, nil, exprBool)
			}
			switch {
			case a2.len() == 0:
				if n.ExpressionList3 != nil {
					b3 = c.expr(&a3, n.ExpressionList3, nil, exprVoid)
				}
				switch {
				case a3.len() == 0:
					w.w("\nfor %s; %s; %s {", b1, b2, b3)
					switch n.Statement.Case {
					case cc.StatementCompound:
						for l := n.Statement.CompoundStatement.BlockItemList; l != nil; l = l.BlockItemList {
							c.blockItem(w, l.BlockItem)
						}
					default:
						c.statement(w, n.Statement)
					}
					w.w("\n}")
				default:
					c.err(errorf("TODO %v", n.Case))
				}
			default:
				c.err(errorf("TODO %v", n.Case))
			}
		default:
			c.err(errorf("TODO %v", n.Case))
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
		c.err(errorf("TODO %v", n.Case))
	case cc.JumpStatementBreak: // "break" ';'
		w.w("\nbreak")
	case cc.JumpStatementReturn: // "return" ExpressionList ';'
		w.w("\nreturn %s", c.expr(w, n.ExpressionList, c.f.t.Result(), expr))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) expressionStatement(w writer, n *cc.ExpressionStatement) {
	w.w("\n%s", c.expr(w, n.ExpressionList, nil, exprVoid))
}
