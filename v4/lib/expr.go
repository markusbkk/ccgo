// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"modernc.org/cc/v4"
	"modernc.org/mathutil"
)

type mode int

const (
	_           mode = iota
	exprBool         // C scalar type, Go bool
	exprCall         // C func pointer, Go function value
	exprDefault      //
	exprIndex        // C pointer, Go array
	exprLvalue       //
	exprSelect       // C struct, Go struct
	exprUintptr      // C pointer, Go uintptr
	exprVoid         // C void, no Go equivalent
)

const (
	ccgoFP = "__ccgo_fp"
)

func (c *ctx) expr(w writer, n cc.ExpressionNode, to cc.Type, toMode mode) *buf {
	if toMode == 0 {
		c.err(errorf("internal error"))
		return &buf{}
	}

	if n == nil {
		if toMode != exprVoid {
			c.err(errorf("TODO"))
		}
		return &buf{}
	}

	if x, ok := n.(*cc.ExpressionList); ok && x == nil {
		if toMode != exprVoid {
			c.err(errorf("TODO"))
		}
		return &buf{}
	}

	if to == nil {
		to = n.Type()
	}
	// trc("%v: EXPR  pre call EXPR1 -> %s %s (%s) %T", c.pos(n), to, toMode, cc.NodeSource(n), n)
	r, from, fromMode := c.expr0(w, n, to, toMode)
	// trc("%v: EXPR post call EXPR0 %v %v -> %v %v (%s) %T", c.pos(n), from, fromMode, to, toMode, cc.NodeSource(n), n)
	if from == nil || fromMode == 0 {
		// trc("IN %v: from %v, %v to %v %v, src '%s', buf '%s'", c.pos(n), from, fromMode, to, toMode, cc.NodeSource(n), r.bytes())
		c.err(errorf("TODO %T %v %v -> %v %v", n, from, fromMode, to, toMode))
		return r
	}

	return c.convert(n, w, r, from, to, fromMode, toMode)
}

func (c *ctx) convert(n cc.ExpressionNode, w writer, s *buf, from, to cc.Type, fromMode, toMode mode) (r *buf) {
	// trc("IN %v: from %v, %v to %v %v, src '%s', buf '%s'", c.pos(n), from, fromMode, to, toMode, cc.NodeSource(n), s.bytes())
	// defer func() {
	// 	trc("OUT %v: from %v, %v to %v %v, src '%s', bufs '%s' -> '%s'", c.pos(n), from, fromMode, to, toMode, cc.NodeSource(n), s.bytes(), r.bytes())
	// }()
	if to == nil {
		// trc("ERR %v: from %v: %v, to <nil>: %v '%q' node: %T src: '%q' (%v: %v: %v:)", c.pos(n), from, fromMode, toMode, s.bytes(), s.n, cc.NodeSource(n), origin(4), origin(3), origin(2))
		c.err(errorf("TODO"))
		return s
	}

	if assert && fromMode == exprUintptr && from.Kind() != cc.Ptr && from.Kind() != cc.Function {
		trc("%v: %v %v -> %v %v", c.pos(n), from, fromMode, to, toMode)
		c.err(errorf("TODO assertion failed"))
	}
	if assert && toMode == exprUintptr && to.Kind() != cc.Ptr {
		trc("%v: %v %v -> %v %v", c.pos(n), from, fromMode, to, toMode)
		c.err(errorf("TODO assertion failed"))
	}
	if from != nil && from.Kind() == cc.Enum {
		from = from.(*cc.EnumType).UnderlyingType()
	}
	if to.Kind() == cc.Enum {
		to = to.(*cc.EnumType).UnderlyingType()
	}

	if cc.IsScalarType(from) && fromMode == exprDefault && toMode == exprBool {
		var b buf
		b.w("(%s != 0)", s)
		return &b
	}

	if from == to || from != nil && from.IsCompatible(to) || fromMode == exprBool && cc.IsIntegerType(to) {
		if fromMode == toMode {
			return s
		}

		if from == c.ast.SizeT || to == c.ast.SizeT {
			if toMode != exprVoid {
				return c.convertType(n, s, from, to, fromMode, toMode)
			}
		}

		if fromMode == toMode {
			return s
		}

		return c.convertMode(n, w, s, from, to, fromMode, toMode)
	}

	if fromMode == toMode {
		return c.convertType(n, s, from, to, fromMode, toMode)
	}

	if from != nil && from.Kind() == cc.Ptr {
		return c.convertFromPointer(n, s, from.(*cc.PointerType), to, fromMode, toMode)
	}

	if toMode == exprVoid || to.Kind() == cc.Void {
		return s
	}

	if to.Kind() == cc.Ptr {
		return c.convertToPointer(n, s, from, to.(*cc.PointerType), fromMode, toMode)
	}

	// trc("%v: %s", n.Position(), cc.NodeSource(n))
	// trc("TODO %q %s %s -> %s %s", s, from, fromMode, to, toMode)
	c.err(errorf("TODO %q %s %s -> %s %s", s, from, fromMode, to, toMode))
	return s //TODO
}

func (c *ctx) convertToPointer(n cc.ExpressionNode, s *buf, from cc.Type, to *cc.PointerType, fromMode, toMode mode) (r *buf) {
	var b buf
	switch fromMode {
	case exprDefault:
		switch toMode {
		case exprUintptr:
			b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n, s)))
			return &b
		}
	case exprBool:
		switch toMode {
		case exprDefault:
			b.w("%s%sBool%s(%s)", c.task.tlsQualifier, tag(preserve), c.helper(n, to), s)
			return &b
		}
	}

	// trc("%v: from %v, %v to %v %v, src '%s', buf '%s'", c.pos(n), from, fromMode, to, toMode, cc.NodeSource(n), s.bytes())
	c.err(errorf("TODO %q %s %s -> %s %s", s, from, fromMode, to, toMode))
	return s //TODO
}

func (c *ctx) pin(n cc.Node, b *buf) *buf {
	switch x := b.n.(type) {
	case *cc.Declarator:
		switch c.pass {
		case 0:
			// ok
		case 1:
			switch symKind(string(b.bytes())) {
			case automatic, ccgoAutomatic:
				c.f.declInfos.takeAddress(x)
				// trc("%v: PIN %v at %v (%v: %v: %v:)", c.pos(n), x.Name(), c.pos(x), origin(4), origin(3), origin(2))
				// trc("%s", debug.Stack())
			}
		case 2:
			// ok
		default:
			c.err(errorf("%v: internal error: %d", n.Position(), c.pass))
		}
	case *cc.PostfixExpression:
		if y := c.declaratorOf(x.PostfixExpression); y != nil {
			s := strings.Trim(string(b.bytes()), "()")
			s = s[:strings.IndexByte(s, '.')]
			switch symKind(s) {
			case automatic, ccgoAutomatic:
				c.f.declInfos.takeAddress(y)
				// trc("%v: PIN %v at %v (%v: %v:)", c.pos(n), y.Name(), c.pos(y), origin(3), origin(2))
				// trc("%s", debug.Stack())
			}
			return b
		}
	}
	return b
}

// type unchanged
func (c *ctx) convertMode(n cc.ExpressionNode, w writer, s *buf, from, to cc.Type, fromMode, toMode mode) (r *buf) {
	// defer func() { trc("%v: from %v: %v, to %v: %v %q -> %q", c.pos(n), from, fromMode, to, toMode, b, r) }()
	var b buf
	switch fromMode {
	case exprDefault:
		switch toMode {
		case exprLvalue:
			return s
		case exprCall:
			return s
		case exprVoid:
			return s
		case exprBool:
			b.w("(%s != 0)", s)
			return &b
		case exprIndex:
			switch x := from.(type) {
			case *cc.PointerType:
				switch y := from.Undecay().(type) {
				case *cc.ArrayType:
					b.w("((*%s)(%s))", c.typ(n, y), unsafePointer(s))
					return &b
				default:
					c.err(errorf("TODO %T", y))
				}
			case *cc.ArrayType:
				b.w("(*(*%s)(%s))", c.typ(n, x), unsafeAddr(s))
				return &b
			default:
				trc("%v:", n.Position())
				c.err(errorf("TODO %T", x))
			}
		}
	case exprUintptr:
		switch toMode {
		case exprDefault:
			return s
		case exprBool:
			b.w("(%s != 0)", s)
			return &b
		case exprCall:
			// v := fmt.Sprintf("%sf%d", tag(ccgo), c.id())
			// ft := from.(*cc.PointerType).Elem().(*cc.FunctionType)
			// vs := fmt.Sprintf("var %s func%s;", v, c.signature(ft, false, false, true))
			// switch {
			// case c.f != nil:
			// 	c.f.registerAutoVar(vs)
			// default:
			// 	w.w("%s", vs)
			// }
			// w.w("\n*(*%suintptr)(%s) = %s;", tag(preserve), unsafeAddr(v), s) // Free pass from .pin
			// var b buf
			// b.w("%s", v)
			// return &b

			var b buf
			ft := from.(*cc.PointerType).Elem().(*cc.FunctionType)
			b.w("(*(*func%s)(%sunsafe.%sPointer(&struct{%[3]suintptr}{%s})))", c.signature(ft, false, false, true), tag(importQualifier), tag(preserve), s)
			return &b
		}
	case exprBool:
		switch toMode {
		case exprDefault:
			switch {
			case cc.IsIntegerType(to):
				b.w("%s%sBool%s(%s)", c.task.tlsQualifier, tag(preserve), c.helper(n, to), s)
				return &b
			}
		case exprVoid:
			return s
		}
	case exprVoid:
		switch toMode {
		case exprDefault:
			return s
		}
	}
	//trc("%v: from %v, %v to %v %v, src '%s', buf '%s'", c.pos(n), from, fromMode, to, toMode, cc.NodeSource(n), s.bytes())
	c.err(errorf("TODO %q %s %s -> %s %s", s, from, fromMode, to, toMode))
	return s //TODO
}

//TODO- func (c *ctx) isIdent(s string) bool {
//TODO- 	for i, v := range s {
//TODO- 		switch {
//TODO- 		case i == 0:
//TODO- 			if !unicode.IsLetter(v) && v != '_' {
//TODO- 				return false
//TODO- 			}
//TODO- 		default:
//TODO- 			if !unicode.IsLetter(v) && v != '_' && !unicode.IsDigit(v) {
//TODO- 				return false
//TODO- 			}
//TODO- 		}
//TODO- 	}
//TODO- 	return len(s) != 0
//TODO- }

// mode unchanged
func (c *ctx) convertType(n cc.ExpressionNode, s *buf, from, to cc.Type, fromMode, toMode mode) (r *buf) {
	// defer func() { trc("%v: from %v: %v, to %v: %v %q -> %q", c.pos(n), from, fromMode, to, toMode, s, r) }()
	if from.Kind() == cc.Ptr && to.Kind() == cc.Ptr || to.Kind() == cc.Void {
		return s
	}

	var b buf
	if cc.IsScalarType(from) && cc.IsScalarType(to) {
		//b.w("(%s(%s))", c.typ(n, to), s)
		switch {
		case from.Kind() == cc.Int128:
			//TODO
		case from.Kind() == cc.UInt128:
			//TODO
		case to.Kind() == cc.Int128:
			//TODO
		case to.Kind() == cc.UInt128:
			//TODO
		default:
			b.w("(%s%s%sFrom%s(%s))", c.task.tlsQualifier, tag(preserve), c.helper(n, to), c.helper(n, from), s)
			return &b
		}
	}

	if from.Kind() == cc.Function && to.Kind() == cc.Ptr && to.(*cc.PointerType).Elem().Kind() == cc.Function {
		return s
	}

	c.err(errorf("TODO %q %s %s -> %s %s (%v:)", s, from, fromMode, to, toMode, c.pos(n)))
	// panic(errorf("TODO %q %s %s -> %s %s (%v:)", s, from, fromMode, to, toMode, c.pos(n)))
	return s //TODO
}

func (c *ctx) isCharType(t cc.Type) bool {
	switch t.Kind() {
	case cc.Char, cc.UChar, cc.SChar:
		return true
	}

	return false
}

func (c *ctx) convertFromPointer(n cc.ExpressionNode, s *buf, from *cc.PointerType, to cc.Type, fromMode, toMode mode) (r *buf) {
	var b buf
	if to.Kind() == cc.Ptr {
		if fromMode == exprUintptr && toMode == exprDefault {
			return s
		}

		if fromMode == exprDefault && toMode == exprUintptr {
			b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n, s)))
			return &b
		}
	}

	if cc.IsIntegerType(to) {
		if toMode == exprDefault {
			b.w("(%s(%s))", c.typ(n, to), s)
			return &b
		}
	}

	c.err(errorf("TODO %q %s %s, %s -> %s %s, %s", s, from, from.Kind(), fromMode, to, to.Kind(), toMode))
	// trc("%v: TODO %q %s %s, %s -> %s %s, %s", cpos(n), s, from, from.Kind(), fromMode, to, to.Kind(), toMode)
	return s //TODO
}

func (c *ctx) reduceBitFieldValue(expr *buf, f *cc.Field, t cc.Type, mode mode) (r *buf) {
	if mode != exprDefault || f == nil {
		return expr
	}

	var b buf
	bits := f.ValueBits()
	if bits >= t.Size()*8 {
		return expr
	}

	m := ^uint64(0) >> (64 - bits)
	switch {
	case cc.IsSignedInteger(t):
		w := t.Size() * 8
		b.w("(((%s)&%#0x)<<%d>>%[3]d)", expr, m, w-bits)
	default:
		b.w("((%s)&%#0x)", expr, m)
	}
	return &b
}

func (c *ctx) expr0(w writer, n cc.ExpressionNode, t cc.Type, mod mode) (r *buf, rt cc.Type, rmode mode) {
	// trc("%v: %T (%q), %v, %v (%v: %v: %v:) (IN)", n.Position(), n, cc.NodeSource(n), t, mod, origin(4), origin(3), origin(2))
	// defer func() {
	// 	trc("%v: %T (%q), %v, %v (RET)", n.Position(), n, cc.NodeSource(n), t, mod)
	// }()

	defer func(mod mode) {
		if r == nil || rt == nil || !cc.IsIntegerType(rt) {
			return
		}

		if x, ok := n.(*cc.PrimaryExpression); ok {
			switch x.Case {
			case cc.PrimaryExpressionIdent: // IDENTIFIER
				if x.Value() != nil {
					return
				}
			case
				cc.PrimaryExpressionInt,     // INTCONST
				cc.PrimaryExpressionFloat,   // FLOATCONST
				cc.PrimaryExpressionChar,    // CHARCONST
				cc.PrimaryExpressionLChar,   // LONGCHARCONST
				cc.PrimaryExpressionString,  // STRINGLITERAL
				cc.PrimaryExpressionLString: // LONGSTRINGLITERAL
				return
			case
				cc.PrimaryExpressionExpr,    // '(' ExpressionList ')'
				cc.PrimaryExpressionStmt,    // '(' CompoundStatement ')'
				cc.PrimaryExpressionGeneric: // GenericSelection
				// ok
			default:
				c.err(errorf("internal error %T %v", x, x.Case))
				return
			}
		}
		if bf := rt.BitField(); bf != nil && bf.ValueBits() > 32 {
			r = c.reduceBitFieldValue(r, bf, rt, mod)
		}
	}(mod)

	blank := false
out:
	switch {
	case mod == exprBool:
		mod = exprDefault
	case mod == exprDefault && n.Type().Undecay().Kind() == cc.Array:
		if d := c.declaratorOf(n); d == nil || !d.IsParam() {
			mod = exprUintptr
		}
	case mod == exprVoid:
		if _, ok := n.(*cc.ExpressionList); ok {
			break out
		}

		if n.Type().Kind() == cc.Void {
			switch x := n.(type) {
			case *cc.CastExpression:
				if x.Case == cc.CastExpressionCast && x.CastExpression.Type().Kind() != cc.Void {
					blank = true
				}
			}
			break out
		}

		switch x := n.(type) {
		case *cc.AssignmentExpression:
			break out
		case *cc.PostfixExpression:
			switch x.Case {
			case cc.PostfixExpressionCall, cc.PostfixExpressionDec, cc.PostfixExpressionInc:
				break out
			}
		case *cc.UnaryExpression:
			switch x.Case {
			case cc.UnaryExpressionDec, cc.UnaryExpressionInc:
				break out
			}
		case *cc.PrimaryExpression:
			switch x.Case {
			case cc.PrimaryExpressionExpr:
				break out
			}
		}

		blank = true

	}
	if blank {
		defer func() {
			if len(r.bytes()) != 0 {
				var b buf
				b.w("%s_ = %s", tag(preserve), r.bytes())
				r.b = b.b
			}
		}()
	}
	if t == nil {
		t = n.Type()
	}
	switch x := n.(type) {
	case *cc.AdditiveExpression:
		return c.additiveExpression(w, x, t, mod)
	case *cc.AndExpression:
		return c.andExpression(w, x, t, mod)
	case *cc.AssignmentExpression:
		return c.assignmentExpression(w, x, t, mod)
	case *cc.CastExpression:
		return c.castExpression(w, x, t, mod)
	case *cc.ConstantExpression:
		return c.expr0(w, x.ConditionalExpression, t, mod)
	case *cc.ConditionalExpression:
		return c.conditionalExpression(w, x, t, mod)
	case *cc.EqualityExpression:
		return c.equalityExpression(w, x, t, mod)
	case *cc.ExclusiveOrExpression:
		return c.exclusiveOrExpression(w, x, t, mod)
	case *cc.ExpressionList:
		return c.expressionList(w, x, t, mod)
	case *cc.InclusiveOrExpression:
		return c.inclusiveOrExpression(w, x, t, mod)
	case *cc.LogicalAndExpression:
		return c.logicalAndExpression(w, x, t, mod)
	case *cc.LogicalOrExpression:
		return c.logicalOrExpression(w, x, t, mod)
	case *cc.MultiplicativeExpression:
		return c.multiplicativeExpression(w, x, t, mod)
	case *cc.PostfixExpression:
		return c.postfixExpression(w, x, t, mod)
	case *cc.PrimaryExpression:
		return c.primaryExpression(w, x, t, mod)
	case *cc.RelationalExpression:
		return c.relationExpression(w, x, t, mod)
	case *cc.ShiftExpression:
		return c.shiftExpression(w, x, t, mod)
	case *cc.UnaryExpression:
		return c.unaryExpression(w, x, t, mod)
	default:
		c.err(errorf("TODO %T", x))
		return nil, nil, 0
	}
}

func (c *ctx) andExpression(w writer, n *cc.AndExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.AndExpressionEq: // EqualityExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.AndExpressionAnd: // AndExpression '&' EqualityExpression
		b.w("(%s & %s)", c.expr(w, n.AndExpression, n.Type(), exprDefault), c.expr(w, n.EqualityExpression, n.Type(), exprDefault))
		rt, rmode = n.Type(), exprDefault
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) exclusiveOrExpression(w writer, n *cc.ExclusiveOrExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.ExclusiveOrExpressionAnd: // AndExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.ExclusiveOrExpressionXor: // ExclusiveOrExpression '^' AndExpression
		b.w("(%s ^ %s)", c.expr(w, n.ExclusiveOrExpression, n.Type(), exprDefault), c.expr(w, n.AndExpression, n.Type(), exprDefault))
		rt, rmode = n.Type(), exprDefault
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) inclusiveOrExpression(w writer, n *cc.InclusiveOrExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.InclusiveOrExpressionXor: // ExclusiveOrExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.InclusiveOrExpressionOr: // InclusiveOrExpression '|' ExclusiveOrExpression
		b.w("(%s | %s)", c.expr(w, n.InclusiveOrExpression, n.Type(), exprDefault), c.expr(w, n.ExclusiveOrExpression, n.Type(), exprDefault))
		rt, rmode = n.Type(), exprDefault
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) shiftExpression(w writer, n *cc.ShiftExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.ShiftExpressionAdd: // AdditiveExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.ShiftExpressionLsh: // ShiftExpression "<<" AdditiveExpression
		b.w("(%s << %s)", c.expr(w, n.ShiftExpression, n.Type(), exprDefault), c.expr(w, n.AdditiveExpression, nil, exprDefault))
		rt, rmode = n.Type(), exprDefault
	case cc.ShiftExpressionRsh: // ShiftExpression ">>" AdditiveExpression
		b.w("(%s >> %s)", c.expr(w, n.ShiftExpression, n.Type(), exprDefault), c.expr(w, n.AdditiveExpression, nil, exprDefault))
		rt, rmode = n.Type(), exprDefault
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) logicalAndExpression(w writer, n *cc.LogicalAndExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.LogicalAndExpressionOr: // InclusiveOrExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.LogicalAndExpressionLAnd: // LogicalAndExpression "&&" InclusiveOrExpression
		rt, rmode = n.Type(), exprBool
		var al, ar buf
		bl := c.expr(&al, n.LogicalAndExpression, nil, exprBool)
		br := c.expr(&ar, n.InclusiveOrExpression, nil, exprBool)
		switch {
		default:
			// case al.len() == 0 || ar.len() == 0:
			b.w("((%s) && (%s))", bl, br)
		case al.len() == 0 && ar.len() != 0:
			// Sequence point
			// if v = bl; v { ar };
			// v && br
			v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
			vs := fmt.Sprintf("var %s %sbool;", v, tag(preserve))
			c.f.registerAutoVar(vs)
			w.w("\nif %s = %s; %s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) && (%s))", v, br)
		case al.len() != 0 && ar.len() == 0:
			// Sequence point
			// al;
			// bl && br
			w.w("%s;", al.bytes())
			b.w("((%s) && (%s))", bl, br)
		case al.len() != 0 && ar.len() != 0:
			// Sequence point
			// al; if v = bl; v { ar };
			// v && br
			v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
			vs := fmt.Sprintf("var %s %sbool;", v, tag(preserve))
			c.f.registerAutoVar(vs)
			w.w("%s;", al.bytes())
			w.w("\nif %s = %s; %s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) && (%s))", v, br)
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, c.ast.Int, rmode
}

func (c *ctx) logicalOrExpression(w writer, n *cc.LogicalOrExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.LogicalOrExpressionLAnd: // LogicalAndExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.LogicalOrExpressionLOr: // LogicalOrExpression "||" LogicalAndExpression
		rt, rmode = n.Type(), exprBool
		var al, ar buf
		bl := c.expr(&al, n.LogicalOrExpression, nil, exprBool)
		br := c.expr(&ar, n.LogicalAndExpression, nil, exprBool)
		switch {
		default:
			// case al.len() == 0 || ar.len() == 0:
			b.w("((%s) || (%s))", bl, br)
		case al.len() == 0 && ar.len() != 0:
			// Sequence point
			// if v = bl; !v { ar };
			// v || br
			v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
			vs := fmt.Sprintf("var %s %sbool;", v, tag(preserve))
			c.f.registerAutoVar(vs)
			w.w("\nif %s = %s; !%s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) || (%s))", v, br)
		case al.len() != 0 && ar.len() == 0:
			// Sequence point
			// al;
			// bl || br
			w.w("%s;", al.bytes())
			b.w("((%s) || (%s))", bl, br)
		case al.len() != 0 && ar.len() != 0:
			// Sequence point
			// al; if v = bl; !v { ar };
			// v || br
			v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
			vs := fmt.Sprintf("var %s %sbool;", v, tag(preserve))
			c.f.registerAutoVar(vs)
			w.w("%s;", al.bytes())
			w.w("\nif %s = %s; !%s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) || (%s))", v, br)
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, c.ast.Int, rmode
}

func (c *ctx) conditionalExpression(w writer, n *cc.ConditionalExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.ConditionalExpressionLOr: // LogicalOrExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.ConditionalExpressionCond: // LogicalOrExpression '?' ExpressionList ':' ConditionalExpression
		v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
		switch mode {
		case exprCall:
			rt, rmode = n.Type(), mode
			vs := fmt.Sprintf("var %s func%s;", v, c.signature(n.Type().(*cc.PointerType).Elem().(*cc.FunctionType), false, false, true))
			switch {
			case c.f != nil:
				c.f.registerAutoVar(vs)
			default:
				w.w("%s", vs)
			}
			w.w("if %s {", c.expr(w, n.LogicalOrExpression, nil, exprBool))
			w.w("%s = %s;", v, c.expr(w, n.ExpressionList, n.Type(), mode))
			w.w("} else {")
			w.w("%s = %s;", v, c.expr(w, n.ConditionalExpression, n.Type(), mode))
			w.w("};")
			b.w("%s", v)
		case exprIndex:
			rt, rmode = n.Type(), exprUintptr
			vs := fmt.Sprintf("var %s %s;", v, c.typ(n, n.Type()))
			switch {
			case c.f != nil:
				c.f.registerAutoVar(vs)
			default:
				w.w("%s", vs)
			}
			w.w("if %s {", c.expr(w, n.LogicalOrExpression, nil, exprBool))
			w.w("%s = %s;", v, c.pin(n, c.expr(w, n.ExpressionList, n.Type(), exprUintptr)))
			w.w("} else {")
			w.w("%s = %s;", v, c.pin(n, c.expr(w, n.ConditionalExpression, n.Type(), exprUintptr)))
			w.w("};")
			b.w("%s", v)
		case exprVoid:
			rt, rmode = n.Type(), mode
			w.w("if %s {", c.expr(w, n.LogicalOrExpression, nil, exprBool))
			w.w("%s;", c.expr(w, n.ExpressionList, n.Type(), exprVoid))
			w.w("} else {")
			w.w("%s;", c.expr(w, n.ConditionalExpression, n.Type(), exprVoid))
			w.w("};")
		default:
			rt, rmode = n.Type(), mode
			vs := fmt.Sprintf("var %s %s;", v, c.typ(n, n.Type()))
			switch {
			case c.f != nil:
				c.f.registerAutoVar(vs)
			default:
				w.w("%s", vs)
			}
			w.w("if %s {", c.expr(w, n.LogicalOrExpression, nil, exprBool))
			w.w("%s = %s;", v, c.expr(w, n.ExpressionList, n.Type(), exprDefault))
			w.w("} else {")
			w.w("%s = %s;", v, c.expr(w, n.ConditionalExpression, n.Type(), exprDefault))
			w.w("};")
			b.w("%s", v)
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) castExpression(w writer, n *cc.CastExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.CastExpressionUnary: // UnaryExpression
		return c.expr0(w, n.UnaryExpression, t, mode)
	case cc.CastExpressionCast: // '(' TypeName ')' CastExpression
		switch x := t.(type) {
		case *cc.PointerType:
			switch x.Elem().(type) {
			case *cc.FunctionType:
				if mode == exprCall {
					rt, rmode = n.Type(), exprUintptr
					b.w("%s", c.expr(w, n.CastExpression, n.Type(), exprDefault))
					return &b, rt, rmode
				}
			}
		}

		rt, rmode = n.Type(), mode
		b.w("%s", c.expr(w, n.CastExpression, n.Type(), exprDefault))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) multiplicativeExpression(w writer, n *cc.MultiplicativeExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	rt, rmode = n.Type(), exprDefault
	var b buf
	switch n.Case {
	case cc.MultiplicativeExpressionCast: // CastExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.MultiplicativeExpressionMul: // MultiplicativeExpression '*' CastExpression
		x, y := c.binopArgs(w, n.MultiplicativeExpression, n.CastExpression, n.Type())
		b.w("(%s * %s)", x, y)
	case cc.MultiplicativeExpressionDiv: // MultiplicativeExpression '/' CastExpression
		x, y := c.binopArgs(w, n.MultiplicativeExpression, n.CastExpression, n.Type())
		b.w("(%s / %s)", x, y)
	case cc.MultiplicativeExpressionMod: // MultiplicativeExpression '%' CastExpression
		x, y := c.binopArgs(w, n.MultiplicativeExpression, n.CastExpression, n.Type())
		b.w("(%s %% %s)", x, y)
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) additiveExpression(w writer, n *cc.AdditiveExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	rt, rmode = n.Type(), exprDefault
	var b buf
	switch n.Case {
	case cc.AdditiveExpressionMul: // MultiplicativeExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.AdditiveExpressionAdd: // AdditiveExpression '+' MultiplicativeExpression
		switch x, y := n.AdditiveExpression.Type(), n.MultiplicativeExpression.Type(); {
		case cc.IsArithmeticType(x) && cc.IsArithmeticType(y):
			x, y := c.binopArgs(w, n.AdditiveExpression, n.MultiplicativeExpression, n.Type())
			b.w("(%s + %s)", x, y)
		case x.Kind() == cc.Ptr && cc.IsIntegerType(y):
			s := ""
			if sz := x.(*cc.PointerType).Elem().Undecay().Size(); sz != 1 {
				s = fmt.Sprintf("*%d", sz)
			}
			b.w("(%s + ((%s)%s))", c.expr(w, n.AdditiveExpression, n.Type(), exprDefault), c.expr(w, n.MultiplicativeExpression, n.Type(), exprDefault), s)
		case cc.IsIntegerType(x) && y.Kind() == cc.Ptr:
			s := ""
			if sz := y.(*cc.PointerType).Elem().Undecay().Size(); sz != 1 {
				s = fmt.Sprintf("*%d", sz)
			}
			b.w("(((%s)%s)+%s)", c.expr(w, n.AdditiveExpression, n.Type(), exprDefault), s, c.expr(w, n.MultiplicativeExpression, n.Type(), exprDefault))
		default:
			c.err(errorf("TODO %v + %v", x, y)) // -
		}
	case cc.AdditiveExpressionSub: // AdditiveExpression '-' MultiplicativeExpression
		switch x, y := n.AdditiveExpression.Type(), n.MultiplicativeExpression.Type(); {
		case cc.IsArithmeticType(x) && cc.IsArithmeticType(y):
			x, y := c.binopArgs(w, n.AdditiveExpression, n.MultiplicativeExpression, n.Type())
			b.w("(%s - %s)", x, y)
		case x.Kind() == cc.Ptr && y.Kind() == cc.Ptr:
			b.w("((%s - %s)", c.expr(w, n.AdditiveExpression, n.Type(), exprDefault), c.expr(w, n.MultiplicativeExpression, n.Type(), exprDefault))
			if v := x.(*cc.PointerType).Elem().Undecay().Size(); v > 1 {
				b.w("/%d", v)
			}
			b.w(")")
		case x.Kind() == cc.Ptr && cc.IsIntegerType(y):
			s := ""
			if sz := x.(*cc.PointerType).Elem().Undecay().Size(); sz != 1 {
				s = fmt.Sprintf("*%d", sz)
			}
			b.w("(%s - ((%s)%s))", c.expr(w, n.AdditiveExpression, n.Type(), exprDefault), c.expr(w, n.MultiplicativeExpression, n.Type(), exprDefault), s)
		default:
			c.err(errorf("TODO %v - %v", x, y))
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) binopArgs(w writer, a, b cc.ExpressionNode, t cc.Type) (x, y *buf) {
	return c.expr(w, a, t, exprDefault), c.expr(w, b, t, exprDefault)
}

func (c *ctx) equalityExpression(w writer, n *cc.EqualityExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	if n.Case == cc.EqualityExpressionRel { // RelationalExpression
		c.err(errorf("TODO %v", n.Case))
		return &b, nil, exprBool
	}

	ct := c.usualArithmeticConversions(n.EqualityExpression.Type(), n.RelationalExpression.Type())
	switch n.Case {
	case cc.EqualityExpressionEq: // EqualityExpression "==" RelationalExpression
		b.w("(%s == %s)", c.expr(w, n.EqualityExpression, ct, exprDefault), c.expr(w, n.RelationalExpression, ct, exprDefault))
		rt, rmode = n.Type(), exprBool
	case cc.EqualityExpressionNeq: // EqualityExpression "!=" RelationalExpression
		b.w("(%s != %s)", c.expr(w, n.EqualityExpression, ct, exprDefault), c.expr(w, n.RelationalExpression, ct, exprDefault))
		rt, rmode = n.Type(), exprBool
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) relationExpression(w writer, n *cc.RelationalExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	if n.Case == cc.RelationalExpressionShift { // ShiftExpression
		c.err(errorf("TODO %v", n.Case))
		return &b, nil, exprBool
	}

	ct := c.usualArithmeticConversions(n.RelationalExpression.Type(), n.ShiftExpression.Type())
	rt, rmode = n.Type(), exprBool
	switch n.Case {
	case cc.RelationalExpressionLt: // RelationalExpression '<' ShiftExpression
		b.w("(%s < %s)", c.expr(w, n.RelationalExpression, ct, exprDefault), c.expr(w, n.ShiftExpression, ct, exprDefault))
	case cc.RelationalExpressionGt: // RelationalExpression '>' ShiftExpression
		b.w("(%s > %s)", c.expr(w, n.RelationalExpression, ct, exprDefault), c.expr(w, n.ShiftExpression, ct, exprDefault))
	case cc.RelationalExpressionLeq: // RelationalExpression "<=" ShiftExpression
		b.w("(%s <= %s)", c.expr(w, n.RelationalExpression, ct, exprDefault), c.expr(w, n.ShiftExpression, ct, exprDefault))
	case cc.RelationalExpressionGeq: // RelationalExpression ">=" ShiftExpression
		b.w("(%s >= %s)", c.expr(w, n.RelationalExpression, ct, exprDefault), c.expr(w, n.ShiftExpression, ct, exprDefault))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) usualArithmeticConversions(a, b cc.Type) (r cc.Type) {
	if a.Kind() == cc.Ptr && (cc.IsIntegerType(b) || b.Kind() == cc.Ptr) {
		return a
	}

	if b.Kind() == cc.Ptr && (cc.IsIntegerType(a) || a.Kind() == cc.Ptr) {
		return b
	}

	return cc.UsualArithmeticConversions(a, b)
}

func (c *ctx) isBitField(n cc.ExpressionNode) bool {
	for {
		switch x := n.(type) {
		case *cc.PostfixExpression:
			switch x.Case {
			case cc.PostfixExpressionSelect: // PostfixExpression '.' IDENTIFIER
				return x.Field().IsBitfield()
			case cc.PostfixExpressionPSelect: // PostfixExpression "->" IDENTIFIER
				return x.Field().IsBitfield()
			default:
				return false
			}
		case *cc.PrimaryExpression:
			switch x.Case {
			case cc.PrimaryExpressionExpr: // '(' ExpressionList ')'
				n = x.ExpressionList
			default:
				return false
			}
		case *cc.UnaryExpression:
			switch x.Case {
			case cc.UnaryExpressionPostfix: // PostfixExpression
				n = x.PostfixExpression
			default:
				return false
			}
		default:
			trc("TODO %T", x)
			return false
		}
	}
}

func (c *ctx) preIncDecBitField(op string, w writer, n cc.ExpressionNode, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	var p *buf
	var f *cc.Field
	switch x := n.(type) {
	case *cc.PostfixExpression:
		switch x.Case {
		case cc.PostfixExpressionSelect:
			p = c.pin(n, c.expr(w, x.PostfixExpression, x.PostfixExpression.Type().Pointer(), exprUintptr))
			f = x.Field()
		case cc.PostfixExpressionPSelect:
			p = c.expr(w, x.PostfixExpression, nil, exprDefault)
			f = x.Field()
		default:
			trc("%v: BITFIELD %v", n.Position(), x.Case)
			c.err(errorf("TODO %T", x))
			return &b, rt, rmode
		}
	default:
		trc("%v: BITFIELD %v", n.Position(), mode)
		c.err(errorf("TODO %T", x))
		return &b, rt, rmode
	}

	switch mode {
	case exprDefault:
		v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
		vs := fmt.Sprintf("var %s %s;", v, c.typ(n, f.Type()))
		switch {
		case c.f != nil:
			c.f.registerAutoVar(vs)
		default:
			w.w("%s", vs)
		}
		bf, _, _ := c.bitField(w, n, p, f, exprDefault)
		w.w("\n%v = %sAssignBitFieldPtr%d%s(%s+%d, (%s)%s1, %d, %d, %#0x);", v, c.task.tlsQualifier, f.AccessBytes()*8, c.helper(n, f.Type()), p, f.Offset(), bf, op, f.ValueBits(), f.OffsetBits(), f.Mask())
		b.w("%s", v)
		return &b, f.Type(), exprDefault
	case exprVoid:
		sop := "Inc"
		if op == "-" {
			sop = "Dec"
		}
		w.w("\n%sPost%sBitFieldPtr%d%s(%s+%d, 1, %d, %d, %#0x);", c.task.tlsQualifier, sop, f.AccessBytes()*8, c.helper(n, f.Type()), p, f.Offset(), f.ValueBits(), f.OffsetBits(), f.Mask())
		return &b, n.Type(), exprVoid
	default:
		trc("%v: BITFIELD %v", n.Position(), mode)
		c.err(errorf("TODO %v", mode))
	}
	return &b, rt, rmode
}

func (c *ctx) unaryExpression(w writer, n *cc.UnaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
out:
	switch n.Case {
	case cc.UnaryExpressionPostfix: // PostfixExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.UnaryExpressionInc: // "++" UnaryExpression
		if c.isBitField(n.UnaryExpression) {
			return c.preIncDecBitField("+", w, n.UnaryExpression, mode)
		}

		rt, rmode = n.Type(), mode
		switch ue := n.UnaryExpression.Type(); {
		case ue.Kind() == cc.Ptr && ue.(*cc.PointerType).Elem().Undecay().Size() != 1:
			sz := ue.(*cc.PointerType).Elem().Undecay().Size()
			switch mode {
			case exprVoid:
				b.w("%s += %d", c.expr(w, n.UnaryExpression, nil, exprDefault), sz)
			case exprDefault:
				switch d := c.declaratorOf(n.UnaryExpression); {
				case d != nil:
					v := c.f.newAutovar(n, n.UnaryExpression.Type())
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s += %d;", ds, sz)
					w.w("\n%s = %s;", v, ds)
					b.w("%s", v)
				default:
					c.err(errorf("TODO")) // -
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		default:
			switch mode {
			case exprVoid:
				b.w("%s++", c.expr(w, n.UnaryExpression, nil, exprDefault))
			case exprDefault:
				switch d := c.declaratorOf(n.UnaryExpression); {
				case d != nil:
					v := c.f.newAutovar(n, n.UnaryExpression.Type())
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s++;", ds)
					w.w("\n%s = %s;", v, ds)
					b.w("%s", v)
				default:
					v := c.f.newAutovar(n, n.UnaryExpression.Type())
					v2 := c.f.newAutovar(n, n.UnaryExpression.Type().Pointer())
					w.w("%s = %s;", v2, c.expr(w, n.UnaryExpression, n.UnaryExpression.Type().Pointer(), exprUintptr))
					w.w("(*(*%s)(%s))++;", c.typ(n, n.UnaryExpression.Type()), unsafePointer(v2))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n, n.UnaryExpression.Type()), unsafePointer(v2))
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		}
	case cc.UnaryExpressionDec: // "--" UnaryExpression
		if c.isBitField(n.UnaryExpression) {
			return c.preIncDecBitField("-", w, n.UnaryExpression, mode)
		}

		rt, rmode = n.Type(), mode
		switch ue := n.UnaryExpression.Type(); {
		case ue.Kind() == cc.Ptr && ue.(*cc.PointerType).Elem().Undecay().Size() != 1:
			sz := ue.(*cc.PointerType).Elem().Undecay().Size()
			switch mode {
			case exprVoid:
				b.w("%s -= %d", c.expr(w, n.UnaryExpression, nil, exprDefault), sz)
			case exprDefault:
				switch d := c.declaratorOf(n.UnaryExpression); {
				case d != nil:
					v := c.f.newAutovar(n, n.UnaryExpression.Type())
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s -= %d;", ds, sz)
					w.w("\n%s = %s;", v, ds)
					b.w("%s", v)
				default:
					c.err(errorf("TODO")) // -
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		default:
			switch mode {
			case exprVoid:
				b.w("%s--", c.expr(w, n.UnaryExpression, nil, exprDefault))
			case exprDefault:
				switch d := c.declaratorOf(n.UnaryExpression); {
				case d != nil:
					v := c.f.newAutovar(n, n.UnaryExpression.Type())
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s--;", ds)
					w.w("\n%s = %s;", v, ds)
					b.w("%s", v)
				default:
					v := c.f.newAutovar(n, n.UnaryExpression.Type())
					v2 := c.f.newAutovar(n, n.UnaryExpression.Type().Pointer())
					w.w("%s = %s;", v2, c.expr(w, n.UnaryExpression, n.UnaryExpression.Type().Pointer(), exprUintptr))
					w.w("(*(*%s)(%s))--;", c.typ(n, n.UnaryExpression.Type()), unsafePointer(v2))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n, n.UnaryExpression.Type()), unsafePointer(v2))
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		}
	case cc.UnaryExpressionAddrof: // '&' CastExpression
		// trc("%v: nt %v, ct %v, '%s' %v", n.Token.Position(), n.Type(), n.CastExpression.Type(), cc.NodeSource(n), mode)
		switch n.Type().Undecay().(type) {
		case *cc.FunctionType:
			rt, rmode = n.Type(), mode
			b.w("%s", c.expr(w, n.CastExpression, nil, mode))
			break out
		}

		rt, rmode = n.Type(), exprUintptr
		b.w("%s", c.expr(w, n.CastExpression, rt, exprUintptr))
	case cc.UnaryExpressionDeref: // '*' CastExpression
		// trc("%v: nt %v, ct %v, '%s' %v", n.Token.Position(), n.Type(), n.CastExpression.Type(), cc.NodeSource(n), mode)
		if ce, ok := n.CastExpression.(*cc.CastExpression); ok && ce.Case == cc.CastExpressionCast {
			if pfe, ok := ce.CastExpression.(*cc.PostfixExpression); ok && pfe.Case == cc.PostfixExpressionCall {
				if pe, ok := pfe.PostfixExpression.(*cc.PrimaryExpression); ok && pe.Case == cc.PrimaryExpressionIdent && pe.Token.SrcStr() == "__builtin_va_arg_impl" {
					if argumentExpressionListLen(pfe.ArgumentExpressionList) != 1 {
						c.err(errorf("internal error"))
						break out
					}

					p, ok := ce.Type().(*cc.PointerType)
					if !ok {
						c.err(errorf("internal error"))
						break out
					}

					rt, rmode = n.Type(), mode
					t := p.Elem()
					if !cc.IsScalarType(t) {
						c.err(errorf("unsupported va_arg type: %v", t.Kind()))
						t = p
					}
					b.w("%sVa%s(&%s)", c.task.tlsQualifier, c.helper(n, t), c.expr(w, pfe.ArgumentExpressionList.AssignmentExpression, nil, exprDefault))
					break out
				}
			}
		}
		switch n.Type().Undecay().(type) {
		case *cc.FunctionType:
			rt, rmode = n.Type(), mode
			b.w("%s", c.expr(w, n.CastExpression, nil, mode))
			break out
		}

		switch mode {
		case exprDefault, exprLvalue, exprVoid:
			rt, rmode = n.Type(), mode
			b.w("(*(*%s)(%s))", c.typ(n, n.CastExpression.Type().(*cc.PointerType).Elem()), unsafePointer(c.expr(w, n.CastExpression, nil, exprDefault)))
		case exprSelect:
			rt, rmode = n.Type(), mode
			b.w("((*%s)(%s))", c.typ(n, n.CastExpression.Type().(*cc.PointerType).Elem()), unsafePointer(c.expr(w, n.CastExpression, nil, exprDefault)))
		case exprUintptr:
			rt, rmode = n.CastExpression.Type(), mode
			b.w("%s", c.expr(w, n.CastExpression, nil, exprDefault))
		case exprCall:
			rt, rmode = n.CastExpression.Type().(*cc.PointerType).Elem(), exprUintptr
			b.w("(*(*%suintptr)(%s))", tag(preserve), unsafePointer(c.expr(w, n.CastExpression, nil, exprDefault)))
		default:
			// trc("%v: %s", n.Token.Position(), cc.NodeSource(n))
			c.err(errorf("TODO %v", mode))
		}
	case cc.UnaryExpressionPlus: // '+' CastExpression
		rt, rmode = n.Type(), exprDefault
		b.w("(+(%s))", c.expr(w, n.CastExpression, n.Type(), exprDefault))
	case cc.UnaryExpressionMinus: // '-' CastExpression
		rt, rmode = n.Type(), exprDefault
		b.w("(-(%s))", c.expr(w, n.CastExpression, n.Type(), exprDefault))
	case cc.UnaryExpressionCpl: // '~' CastExpression
		rt, rmode = n.Type(), exprDefault
		b.w("(^(%s))", c.expr(w, n.CastExpression, n.Type(), exprDefault))
	case cc.UnaryExpressionNot: // '!' CastExpression
		rt, rmode = n.Type(), exprBool
		b.w("(!(%s))", c.expr(w, n.CastExpression, nil, exprBool))
	case cc.UnaryExpressionSizeofExpr: // "sizeof" UnaryExpression
		if t.Kind() == cc.Void {
			t = n.Type()
		}
		rt, rmode = t, exprDefault
		if c.isValidType(n.UnaryExpression, n.UnaryExpression.Type(), true) {
			b.w("(%s%s%sFromInt64(%d))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), n.Value())
		}
	case cc.UnaryExpressionSizeofType: // "sizeof" '(' TypeName ')'
		if t.Kind() == cc.Void {
			t = n.Type()
		}
		rt, rmode = t, exprDefault
		if c.isValidType(n.TypeName, n.TypeName.Type(), true) {
			b.w("(%s%s%sFromInt64(%d))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), n.Value())
		}
	case cc.UnaryExpressionLabelAddr: // "&&" IDENTIFIER
		c.err(errorf("TODO %v", n.Case))
	case cc.UnaryExpressionAlignofExpr: // "_Alignof" UnaryExpression
		if t.Kind() == cc.Void {
			t = n.Type()
		}
		rt, rmode = t, exprDefault
		b.w("(%s%s%sFromInt32(%d))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), n.UnaryExpression.Type().Align())
	case cc.UnaryExpressionAlignofType: // "_Alignof" '(' TypeName ')'
		if t.Kind() == cc.Void {
			t = n.Type()
		}
		rt, rmode = t, exprDefault
		b.w("(%s%s%sFromInt32(%d))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), n.TypeName.Type().Align())
	case cc.UnaryExpressionImag: // "__imag__" UnaryExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.UnaryExpressionReal: // "__real__" UnaryExpression
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) postfixExpressionIndex(w writer, p, index cc.ExpressionNode, pt *cc.PointerType, nt, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	// trc("%v: %s[%s] %v", c.pos(p), cc.NodeSource(p), cc.NodeSource(index), mode)
	// defer func() { trc("%v: %s[%s] %v -> %q", c.pos(p), cc.NodeSource(p), cc.NodeSource(index), mode, r.bytes()) }()
	var b buf
	elem := pt.Elem()
	var mul string
	if v := elem.Size(); v != 1 {
		if v < 0 {
			c.err(errorf("TODO"))
		}
		mul = fmt.Sprintf("*%v", v)
	}

	rt, rmode = nt, mode
	if f := c.isLastStructOrUnionField(p); f != nil {
		switch f.Type().(type) {
		case *cc.ArrayType:
			// Flexible array member.
			//
			//  https://en.wikipedia.org/wiki/Flexible_array_member
			switch mode {
			case exprLvalue, exprDefault, exprSelect:
				b.w("(*(*%s)(%sunsafe.%sPointer(%s + %[3]suintptr(%[5]s)%s)))", c.typ(p, elem), tag(importQualifier), tag(preserve), c.expr(w, p, nil, exprDefault), c.expr(w, index, nil, exprDefault), mul)
				return &b, nt, mode
			case exprUintptr:
				b.w("(%s+%suintptr(%s)%s)", c.expr(w, p, nil, exprDefault), tag(preserve), c.expr(w, index, nil, exprDefault), mul)
				return &b, nt.Pointer(), mode
			default:
				c.err(errorf("TODO %v", mode))
				return &b, t, mode
			}
		}
	}

	if mode == exprVoid {
		mode = exprDefault
	}
	switch mode {
	case exprSelect, exprLvalue, exprDefault, exprIndex:
		switch x := pt.Undecay().(type) {
		case *cc.ArrayType:
			if d := c.declaratorOf(p); d != nil && !d.IsParam() {
				b.w("%s[%s]", c.expr(w, p, nil, exprIndex), c.expr(w, index, nil, exprDefault))
				break
			}

			b.w("(*(*%s)(%sunsafe.%sPointer(%s + %[3]suintptr(%[5]s)%s)))", c.typ(p, elem), tag(importQualifier), tag(preserve), c.expr(w, p, nil, exprDefault), c.expr(w, index, nil, exprDefault), mul)
		case *cc.PointerType:
			b.w("(*(*%s)(%sunsafe.%sPointer(%s + %[3]suintptr(%[5]s)%s)))", c.typ(p, elem), tag(importQualifier), tag(preserve), c.expr(w, p, nil, exprDefault), c.expr(w, index, nil, exprDefault), mul)
		default:
			// trc("%v: %s[%s] %v %T", c.pos(p), cc.NodeSource(p), cc.NodeSource(index), mode, x)
			c.err(errorf("TODO %T", x))
		}
	case exprCall:
		rt, rmode = t.(*cc.PointerType), exprUintptr
		switch x := pt.Undecay().(type) {
		case *cc.ArrayType:
			if d := c.declaratorOf(p); d != nil && !d.IsParam() {
				b.w("%s[%s]", c.expr(w, p, nil, exprIndex), c.expr(w, index, nil, exprDefault))
				break
			}

			b.w("(*(*%s)(%sunsafe.%sPointer(%s + %[3]suintptr(%[5]s)%s)))", c.typ(p, elem), tag(importQualifier), tag(preserve), c.expr(w, p, nil, exprDefault), c.expr(w, index, nil, exprDefault), mul)
		case *cc.PointerType:
			b.w("(*(*%s)(%sunsafe.%sPointer(%s + %[3]suintptr(%[5]s)%s)))", c.typ(p, elem), tag(importQualifier), tag(preserve), c.expr(w, p, nil, exprDefault), c.expr(w, index, nil, exprDefault), mul)
		default:
			// trc("%v: %s[%s] %v %T", c.pos(p), cc.NodeSource(p), cc.NodeSource(index), mode, x)
			c.err(errorf("TODO %T", x))
		}
	case exprUintptr:
		rt, rmode = nt.Pointer(), mode
		if elem.Kind() == cc.Array {
			if d := c.declaratorOf(p); d != nil && d.Type().Kind() == cc.Ptr {
				b.w("((%s)+(%s)%s)", c.expr(w, p, nil, exprDefault), c.expr(w, index, c.pvoid, exprDefault), mul)
				break
			}
		}

		b.w("(%s + %suintptr(%[2]suintptr(%s)%s))", c.expr(w, p, nil, exprDefault), tag(preserve), c.expr(w, index, nil, exprDefault), mul)
	default:
		// trc("%v: %s[%s] %v", c.pos(p), cc.NodeSource(p), cc.NodeSource(index), mode)
		c.err(errorf("TODO %v", mode))
	}
	return &b, rt, rmode
}

func (c *ctx) postIncDecBitField(op string, w writer, n cc.ExpressionNode, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	var p *buf
	var f *cc.Field
	switch x := n.(type) {
	case *cc.PostfixExpression:
		switch x.Case {
		case cc.PostfixExpressionSelect:
			p = c.pin(n, c.expr(w, x.PostfixExpression, x.PostfixExpression.Type().Pointer(), exprUintptr))
			f = x.Field()
		case cc.PostfixExpressionPSelect:
			p = c.expr(w, x.PostfixExpression, nil, exprDefault)
			f = x.Field()
		default:
			trc("%v: BITFIELD %v", n.Position(), x.Case)
			c.err(errorf("TODO %T", x))
			return &b, rt, rmode
		}
	default:
		trc("%v: BITFIELD %v", n.Position(), mode)
		c.err(errorf("TODO %T", x))
		return &b, rt, rmode
	}

	switch mode {
	case exprDefault, exprVoid:
		b.w("%sPost%sBitFieldPtr%d%s(%s+%d, 1, %d, %d, %#0x)", c.task.tlsQualifier, op, f.AccessBytes()*8, c.helper(n, f.Type()), p, f.Offset(), f.ValueBits(), f.OffsetBits(), f.Mask())
		return &b, f.Type(), exprDefault
	default:
		trc("%v: BITFIELD %v", n.Position(), mode)
		c.err(errorf("TODO %v", mode))
	}
	return &b, rt, rmode
}

func (c *ctx) postfixExpression(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
out:
	switch n.Case {
	case cc.PostfixExpressionPrimary: // PrimaryExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.PostfixExpressionIndex: // PostfixExpression '[' ExpressionList ']'
		if x, ok := n.PostfixExpression.Type().(*cc.PointerType); ok {
			return c.postfixExpressionIndex(w, n.PostfixExpression, n.ExpressionList, x, n.Type(), t, mode)
		}

		if x, ok := n.ExpressionList.Type().(*cc.PointerType); ok {
			return c.postfixExpressionIndex(w, n.ExpressionList, n.PostfixExpression, x, n.Type(), t, mode)
		}

		c.err(errorf("TODO %v", n.Case))
	case cc.PostfixExpressionCall: // PostfixExpression '(' ArgumentExpressionList ')'
		//TODO __builtin_object_size 28_strings.c on darwin/amd64
		switch c.declaratorOf(n.PostfixExpression).Name() {
		case "__builtin_constant_p":
			w.w("%s_ = %s;", tag(preserve), c.expr(w, n.ArgumentExpressionList.AssignmentExpression, nil, exprDefault))
			switch mode {
			case exprBool:
				rt, rmode = n.Type(), mode
				switch {
				case n.Value().(cc.Int64Value) == 0:
					b.w("(false)")
				default:
					b.w("(true)")
				}
			default:
				rt, rmode = n.Type(), exprDefault
				b.w("(%v)", n.Value())
			}
			break out
		case "__builtin_va_start":
			if argumentExpressionListLen(n.ArgumentExpressionList) != 2 || mode != exprVoid {
				c.err(errorf("internal error"))
				break out
			}

			rt, rmode = n.Type(), mode
			w.w("%s = %s%s", c.expr(w, n.ArgumentExpressionList.AssignmentExpression, nil, exprDefault), tag(ccgo), vaArgName)
			break out
		case "__builtin_va_end":
			if argumentExpressionListLen(n.ArgumentExpressionList) != 1 || mode != exprVoid {
				c.err(errorf("internal error"))
				break out
			}

			rt, rmode = n.Type(), mode
			w.w("%s_ = %s;", tag(preserve), c.expr(w, n.ArgumentExpressionList.AssignmentExpression, nil, exprDefault))
			break out
		case "__atomic_load_n":
			return c.atomicLoadN(w, n, t, mode)
		case "__atomic_store_n":
			return c.atomicStoreN(w, n, t, mode)
		case "__builtin_sub_overflow":
			return c.subOverflow(w, n, t, mode)
		case "__builtin_mul_overflow":
			return c.mulOverflow(w, n, t, mode)
		case "__builtin_add_overflow":
			return c.addOverflow(w, n, t, mode)
		}

		switch mode {
		case exprSelect:
			switch n.Type().(type) {
			case *cc.StructType:
				c.err(errorf("TODO"))
			case *cc.UnionType:
				v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
				e, _, _ := c.postfixExpressionCall(w, n)
				w.w("%s := %s;", v, e)
				b.w("%s", v)
				return &b, n.Type(), mode
			}
		default:
			return c.postfixExpressionCall(w, n)
		}
	case cc.PostfixExpressionSelect: // PostfixExpression '.' IDENTIFIER
		return c.postfixExpressionSelect(w, n, t, mode)
	case cc.PostfixExpressionPSelect: // PostfixExpression "->" IDENTIFIER
		return c.postfixExpressionPSelect(w, n, t, mode)
	case cc.PostfixExpressionInc: // PostfixExpression "++"
		if c.isBitField(n.PostfixExpression) {
			return c.postIncDecBitField("Inc", w, n.PostfixExpression, mode)
		}

		rt, rmode = n.Type(), mode
		switch pe := n.PostfixExpression.Type(); {
		case pe.Kind() == cc.Ptr && pe.(*cc.PointerType).Elem().Undecay().Size() != 1:
			sz := pe.(*cc.PointerType).Elem().Undecay().Size()
			switch mode {
			case exprVoid:
				b.w("%s += %d", c.expr(w, n.PostfixExpression, nil, exprDefault), sz)
			case exprDefault, exprUintptr:
				v := c.f.newAutovar(n, n.PostfixExpression.Type())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("%s = %s;", v, ds)
					w.w("%s += %d;", ds, sz)
					b.w("%s", v)
				default:
					v2 := c.f.newAutovar(n, n.PostfixExpression.Type().Pointer())
					w.w("%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintptr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s)) += %d;", c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2), sz)
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		default:
			switch mode {
			case exprVoid:
				b.w("%s++", c.expr(w, n.PostfixExpression, nil, exprDefault))
			case exprDefault:
				v := c.f.newAutovar(n, n.PostfixExpression.Type())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("%s = %s;", v, ds)
					w.w("%s++;", ds)
					b.w("%s", v)
				default:
					v2 := c.f.newAutovar(n, n.PostfixExpression.Type().Pointer())
					w.w("%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintptr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s))++;", c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2))
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		}
	case cc.PostfixExpressionDec: // PostfixExpression "--"
		if c.isBitField(n.PostfixExpression) {
			return c.postIncDecBitField("Dec", w, n.PostfixExpression, mode)
		}

		rt, rmode = n.Type(), mode
		switch pe := n.PostfixExpression.Type(); {
		case pe.Kind() == cc.Ptr && pe.(*cc.PointerType).Elem().Undecay().Size() != 1:
			sz := pe.(*cc.PointerType).Elem().Undecay().Size()
			switch mode {
			case exprVoid:
				b.w("%s -= %d", c.expr(w, n.PostfixExpression, nil, exprDefault), sz)
			case exprDefault:
				v := c.f.newAutovar(n, n.PostfixExpression.Type())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("%s = %s;", v, ds)
					w.w("%s -= %d;", ds, sz)
					b.w("%s", v)
				default:
					v2 := c.f.newAutovar(n, n.PostfixExpression.Type().Pointer())
					w.w("%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintptr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s)) -= %d;", c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2), sz)
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		default:
			switch mode {
			case exprVoid:
				b.w("%s--", c.expr(w, n.PostfixExpression, nil, exprDefault))
			case exprDefault:
				v := c.f.newAutovar(n, n.PostfixExpression.Type())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("%s = %s;", v, ds)
					w.w("%s--;", ds)
					b.w("%s", v)
				default:
					v2 := c.f.newAutovar(n, n.PostfixExpression.Type().Pointer())
					w.w("%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintptr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s))--;", c.typ(n, n.PostfixExpression.Type()), unsafePointer(v2))
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO")) // -
			}
		}
	case cc.PostfixExpressionComplit: // '(' TypeName ')' '{' InitializerList ',' '}'
		var a []*cc.Initializer
		for l := n.InitializerList; l != nil; l = l.InitializerList {
			a = append(a, c.initalizerFlatten(l.Initializer, nil)...)
		}
		t := n.TypeName.Type()
		return c.initializer(w, n, a, t, 0, t.Kind() == cc.Array), t, exprDefault
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) mulOverflow(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	args := argumentExpressionList(n.ArgumentExpressionList)
	if len(args) != 3 {
		c.err(errorf("%v: invalid number of arguments to __builtin_mul_overflow", n.ArgumentExpressionList.Position()))
		return &b, t, mode
	}

	switch {
	case cc.IsScalarType(args[0].Type()):
		// ok
	default:
		c.err(errorf("%v: invalid first argument to __builtin_mul_overflow: %s", n.ArgumentExpressionList.Position(), args[0].Type()))
		return &b, t, mode
	}

	switch {
	case cc.IsScalarType(args[1].Type()):
		// ok
	default:
		c.err(errorf("%v: invalid second argument to __builtin_mul_overflow: %s", n.ArgumentExpressionList.Position(), args[1].Type()))
		return &b, t, mode
	}

	if args[2].Type().Kind() != cc.Ptr {
		c.err(errorf("%v: invalid third argument to __builtin_mul_overflow: %s", n.ArgumentExpressionList.Position(), args[2].Type()))
		return &b, t, mode
	}

	b.w("%s__builtin_mul_overflow%s(%stls, %s, %s, %s)", tag(external), c.helper(n, args[0].Type()), tag(ccgo), c.expr(w, args[0], nil, exprDefault), c.expr(w, args[1], args[0].Type(), exprDefault), c.expr(w, args[2], nil, exprDefault))
	return &b, c.ast.Int, exprDefault
}

func (c *ctx) addOverflow(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	args := argumentExpressionList(n.ArgumentExpressionList)
	if len(args) != 3 {
		c.err(errorf("%v: invalid number of arguments to __builtin_add_overflow", n.ArgumentExpressionList.Position()))
		return &b, t, mode
	}

	switch {
	case cc.IsScalarType(args[0].Type()):
		// ok
	default:
		c.err(errorf("%v: invalid first argument to __builtin_add_overflow: %s", n.ArgumentExpressionList.Position(), args[0].Type()))
		return &b, t, mode
	}

	switch {
	case cc.IsScalarType(args[1].Type()):
		// ok
	default:
		c.err(errorf("%v: invalid second argument to __builtin_add_overflow: %s", n.ArgumentExpressionList.Position(), args[1].Type()))
		return &b, t, mode
	}

	if args[2].Type().Kind() != cc.Ptr {
		c.err(errorf("%v: invalid third argument to __builtin_add_overflow: %s", n.ArgumentExpressionList.Position(), args[2].Type()))
		return &b, t, mode
	}

	b.w("%s__builtin_add_overflow%s(%stls, %s, %s, %s)", tag(external), c.helper(n, args[0].Type()), tag(ccgo), c.expr(w, args[0], nil, exprDefault), c.expr(w, args[1], args[0].Type(), exprDefault), c.expr(w, args[2], nil, exprDefault))
	return &b, c.ast.Int, exprDefault
}

func (c *ctx) subOverflow(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	args := argumentExpressionList(n.ArgumentExpressionList)
	if len(args) != 3 {
		c.err(errorf("%v: invalid number of arguments to __builtin_sub_overflow", n.ArgumentExpressionList.Position()))
		return &b, t, mode
	}

	switch {
	case cc.IsScalarType(args[0].Type()):
		// ok
	default:
		c.err(errorf("%v: invalid first argument to __builtin_sub_overflow: %s", n.ArgumentExpressionList.Position(), args[0].Type()))
		return &b, t, mode
	}

	switch {
	case cc.IsScalarType(args[1].Type()):
		// ok
	default:
		c.err(errorf("%v: invalid second argument to __builtin_sub_overflow: %s", n.ArgumentExpressionList.Position(), args[1].Type()))
		return &b, t, mode
	}

	if args[2].Type().Kind() != cc.Ptr {
		c.err(errorf("%v: invalid third argument to __builtin_add_overflow: %s", n.ArgumentExpressionList.Position(), args[2].Type()))
		return &b, t, mode
	}

	b.w("%s__builtin_sub_overflow%s(%stls, %s, %s, %s)", tag(external), c.helper(n, args[0].Type()), tag(ccgo), c.expr(w, args[0], nil, exprDefault), c.expr(w, args[1], args[0].Type(), exprDefault), c.expr(w, args[2], nil, exprDefault))
	return &b, c.ast.Int, exprDefault
}

func (c *ctx) atomicLoadN(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf

	args := argumentExpressionList(n.ArgumentExpressionList)
	if len(args) != 2 {
		c.err(errorf("%v: invalid number of arguments to __atomic_store_n", n.ArgumentExpressionList.Position()))
		return &b, t, mode
	}

	pt, ok := args[0].Type().(*cc.PointerType)
	if !ok {
		c.err(errorf("%v: invalid first argument to __atomic_store_n: %s", n.ArgumentExpressionList.Position(), args[0].Type()))
		return &b, t, mode
	}

	rt = pt.Elem()
	b.w("%sAtomicLoadN%s(%s, %s)", c.task.tlsQualifier, c.helper(n, rt), c.expr(w, args[0], nil, exprDefault), c.expr(w, args[1], nil, exprDefault))
	return &b, rt, mode
}

func (c *ctx) atomicStoreN(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	if mode != exprVoid {
		c.err(errorf("%v: __atomic_store_n used as a value", n.Position()))
		return &b, t, mode
	}

	args := argumentExpressionList(n.ArgumentExpressionList)
	if len(args) != 3 {
		c.err(errorf("%v: invalid number of arguments to __atomic_store_n", n.ArgumentExpressionList.Position()))
		return &b, t, mode
	}

	if args[0].Type().Kind() != cc.Ptr {
		c.err(errorf("%v: invalid first argument to __atomic_store_n: %s", n.ArgumentExpressionList.Position(), args[0].Type()))
		return &b, t, mode
	}

	switch a1 := args[1]; {
	case cc.IsScalarType(a1.Type()):
		b.w("%sAtomicStoreN%s(%s, %s, %s)", c.task.tlsQualifier, c.helper(n, a1.Type()), c.expr(w, args[0], nil, exprDefault), c.expr(w, a1, nil, exprDefault), c.expr(w, args[2], nil, exprDefault))
	default:
		c.err(errorf("%v: invalid second argument to __atomic_store_n: %s", n.ArgumentExpressionList.Position(), a1.Type()))
	}
	return &b, t, mode
}

func (c *ctx) bitField(w writer, n cc.Node, p *buf, f *cc.Field, mode mode) (r *buf, rt cc.Type, rmode mode) {
	//TODO do not pin expr.fld
	rt = f.Type()
	if f.ValueBits() < c.ast.Int.Size()*8 {
		rt = c.ast.Int
	}
	var b buf
	switch mode {
	case exprDefault, exprVoid:
		rmode = exprDefault
		b.w("((%s(%s((*(*uint%d)(%sunsafe.%sPointer(%s +%d))&%#0x)>>%d)", c.typ(n, rt), c.typ(n, f.Type()), f.AccessBytes()*8, tag(importQualifier), tag(preserve), p, f.Offset(), f.Mask(), f.OffsetBits())
		if cc.IsSignedInteger(f.Type()) && !c.isPositiveEnum(f.Type()) {
			w := f.Type().Size() * 8
			b.w("<<%d>>%[1]d", w-f.ValueBits())
		}
		b.w(")))")
	case exprUintptr:
		rt, rmode = rt.Pointer(), mode
		b.w("(uintptr)(%sunsafe.%sPointer(%s +%d))", tag(importQualifier), tag(preserve), p, f.Offset())
	default:
		c.err(errorf("TODO %v", mode))
	}
	return &b, rt, rmode
}

// t is enum type and all its enum consts are >= 0.
func (c *ctx) isPositiveEnum(t cc.Type) bool {
	switch x := t.(type) {
	case *cc.EnumType:
		return x.Min() >= 0
	}

	return false
}

// PostfixExpression "->" IDENTIFIER
func (c *ctx) postfixExpressionPSelect(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	f := n.Field()
	if f.IsBitfield() {
		return c.bitField(w, n, c.expr(w, n.PostfixExpression, nil, exprDefault), n.Field(), mode)
	}

	pe, ok := n.PostfixExpression.Type().(*cc.PointerType)
	if !ok {
		c.err(errorf("TODO %T", n.PostfixExpression.Type()))
		return &b, rt, rmode
	}

	if u, ok := pe.Elem().(*cc.UnionType); ok && f != firstPositiveSizedField(u) {
		switch mode {
		case exprSelect, exprLvalue, exprDefault:
			rt, rmode = n.Type(), mode
			switch {
			case f.Offset() != 0:
				c.err(errorf("TODO %v", mode))
				//b.w("(*(*%s)(%sunsafe.%sAdd(%[2]sunsafe.%sPointer((%s)), %d)))", c.typ(n, f.Type()), tag(importQualifier), tag(preserve), c.expr(w, n.PostfixExpression, nil, exprDefault), f.Offset())
			default:
				b.w("(*(*%s)(%s))", c.typ(n, f.Type()), unsafePointer(c.expr(w, n.PostfixExpression, nil, exprDefault)))
			}
		default:
			c.err(errorf("TODO %v", mode))
		}
		return &b, rt, rmode
	}

	switch mode {
	case exprDefault, exprLvalue, exprIndex, exprSelect:
		rt, rmode = n.Type(), mode
		b.w("((*%s)(%s).", c.typ(n, pe.Elem()), unsafePointer(c.expr(w, n.PostfixExpression, nil, exprDefault)))
		switch {
		case f.Parent() != nil:
			c.parentFields(&b, n.Token, f)
		default:
			b.w("%s%s", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
		}
		b.w(")")
	case exprUintptr:
		rt, rmode = n.Type().Pointer(), mode
		b.w("((%s)%s)", c.expr(w, n.PostfixExpression, nil, exprDefault), fldOff(f.Offset()))
	case exprCall:
		rt, rmode = n.Type().(*cc.PointerType), exprUintptr
		b.w("((*%s)(%s).", c.typ(n, pe.Elem()), unsafePointer(c.expr(w, n.PostfixExpression, nil, exprDefault)))
		switch {
		case f.Parent() != nil:
			c.parentFields(&b, n.Token, f)
		default:
			b.w("%s%s", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
		}
		b.w(")")
	default:
		c.err(errorf("TODO %v", mode))
	}
	return &b, rt, rmode
}

// PostfixExpression '.' IDENTIFIER
func (c *ctx) postfixExpressionSelect(w writer, n *cc.PostfixExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	b.n = n
	f := n.Field()
	if f.IsBitfield() {
		return c.bitField(w, n, c.pin(n, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintptr)), f, mode)
	}

	if mode == exprVoid {
		mode = exprDefault
	}
	if u, ok := n.PostfixExpression.Type().(*cc.UnionType); ok && f != firstPositiveSizedField(u) {
		switch mode {
		case exprLvalue, exprDefault, exprSelect:
			rt, rmode = n.Type(), mode
			switch {
			case f.Offset() != 0:
				b.w("(*(*%s)(%sunsafe.%sAdd(%[2]sunsafe.%sPointer(&(%s)), %d)))", c.typ(n, f.Type()), tag(importQualifier), tag(preserve), c.expr(w, n.PostfixExpression, nil, exprSelect), f.Offset())
			default:
				b.w("(*(*%s)(%s))", c.typ(n, f.Type()), unsafeAddr(c.expr(w, n.PostfixExpression, nil, exprSelect)))
			}
		case exprCall:
			rt, rmode = n.Type().(*cc.PointerType), exprUintptr
			switch {
			case f.Offset() != 0:
				b.w("(*(*%s)(%sunsafe.%sAdd(%[2]sunsafe.%sPointer(&(%s)), %d)))", c.typ(n, f.Type()), tag(importQualifier), tag(preserve), c.expr(w, n.PostfixExpression, nil, exprSelect), f.Offset())
			default:
				b.w("(*(*%s)(%s))", c.typ(n, f.Type()), unsafeAddr(c.expr(w, n.PostfixExpression, nil, exprSelect)))
			}
		case exprUintptr:
			rt, rmode = n.Type().Pointer(), mode
			switch {
			case f.Offset() != 0:
				// b.w("%suintptr(%sunsafe.%s[1]Add(%[2]sunsafe.%[1]sPointer(&(%[3]s)), %d))", tag(preserve), tag(importQualifier), c.expr(w, n.PostfixExpression, nil, exprSelect), f.Offset())
				b.w("%suintptr(0)", tag(preserve))
			default:
				b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.expr(w, n.PostfixExpression, nil, exprSelect)))
			}
		case exprIndex:
			switch x := n.Type().Undecay().(type) {
			case *cc.ArrayType:
				rt, rmode = n.Type(), mode
				switch {
				case f.Offset() != 0:
					//TODO XXX panic(todo("", n.Position()))
					// b.w("((*%s)(%sunsafe.%sAdd(%[2]sunsafe.%sPointer(&(%s)), %d)))", c.typ(n, f.Type()), tag(importQualifier), tag(preserve), c.pin(n.PostfixExpression, c.expr(w, n.PostfixExpression, nil, exprSelect)), f.Offset())
					c.err(errorf("TODO"))
				default:
					b.w("((*%s)(%s))", c.typ(n, f.Type()), unsafeAddr(c.expr(w, n.PostfixExpression, nil, exprSelect)))
				}
			default:
				c.err(errorf("TODO %T", x))
			}
		default:
			c.err(errorf("TODO %v", mode))
		}
		return &b, rt, rmode
	}

	if mode == exprVoid {
		mode = exprDefault
	}
	switch mode {
	case exprLvalue, exprDefault, exprSelect, exprIndex:
		rt, rmode = n.Type(), mode
		b.w("(%s.", c.expr(w, n.PostfixExpression, nil, exprSelect))
		switch {
		case f.Parent() != nil:
			c.parentFields(&b, n.Token, f)
		default:
			b.w("%s%s", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
		}
		b.w(")")
	case exprUintptr:
		rt, rmode = n.Type().Pointer(), mode
		b.w("%suintptr(%sunsafe.%[1]sPointer(&(%[3]s.", tag(preserve), tag(importQualifier), c.pin(n, c.expr(w, n.PostfixExpression, nil, exprLvalue)))
		switch {
		case f.Parent() != nil:
			c.parentFields(&b, n.Token, f)
		default:
			b.w("%s%s", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
		}
		b.w(")))")
	case exprCall:
		rt, rmode = n.Type().(*cc.PointerType), exprUintptr
		b.w("(%s.", c.expr(w, n.PostfixExpression, nil, exprSelect))
		switch {
		case f.Parent() != nil:
			//TODO XXX panic(todo("", n.Position()))
			c.err(errorf("TODO %v", n.Case))
		default:
			b.w("%s%s)", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
		}
	default:
		c.err(errorf("TODO %v", mode))
	}
	return &b, rt, rmode
}

func (c *ctx) parentFields(b *buf, n cc.Node, f *cc.Field) {
	if p := f.Parent(); p != nil {
		c.parentFields(b, n, p)
		b.w(".")
	}
	switch {
	case f.Name() == "":
		b.w("%s__ccgo%d", tag(field), f.Offset())
	default:
		b.w("%s%s", tag(field), f.Name())
	}
}

func (c *ctx) isLastStructOrUnionField(n cc.ExpressionNode) *cc.Field {
	for {
		switch x := n.(type) {
		case *cc.PostfixExpression:
			var f *cc.Field
			var t cc.Type
			switch x.Case {
			case cc.PostfixExpressionSelect: // PostfixExpression '.' IDENTIFIER
				f = x.Field()
				t = x.PostfixExpression.Type()
			case cc.PostfixExpressionPSelect: // PostfixExpression "->" IDENTIFIER
				f = x.Field()
				t = x.PostfixExpression.Type().(*cc.PointerType).Elem()
			}
			switch x := t.(type) {
			case *cc.StructType:
				if f.Index() == x.NumFields()-1 {
					return f
				}
			case *cc.UnionType:
				return f
			}
		}

		return nil
	}
}

func (c *ctx) declaratorOf(n cc.ExpressionNode) *cc.Declarator {
	for n != nil {
		switch x := n.(type) {
		case *cc.PrimaryExpression:
			switch x.Case {
			case cc.PrimaryExpressionIdent: // IDENTIFIER
				switch y := x.ResolvedTo().(type) {
				case *cc.Declarator:
					return y
				case *cc.Parameter:
					return y.Declarator
				case nil:
					return nil
				default:
					c.err(errorf("TODO %T", y))
					return nil
				}
			case cc.PrimaryExpressionExpr: // '(' ExpressionList ')'
				n = x.ExpressionList
			default:
				return nil
			}
		case *cc.PostfixExpression:
			switch x.Case {
			case cc.PostfixExpressionPrimary: // PrimaryExpression
				n = x.PrimaryExpression
			default:
				return nil
			}
		case *cc.ExpressionList:
			if x == nil {
				return nil
			}

			for l := x; l != nil; l = l.ExpressionList {
				n = l.AssignmentExpression
			}
		case *cc.CastExpression:
			switch x.Case {
			case cc.CastExpressionUnary: // UnaryExpression
				n = x.UnaryExpression
			default:
				return nil
			}
		case *cc.UnaryExpression:
			switch x.Case {
			case
				cc.UnaryExpressionInc,
				cc.UnaryExpressionDec,
				cc.UnaryExpressionPostfix: // PostfixExpression

				n = x.PostfixExpression
			default:
				return nil
			}
		case *cc.ConditionalExpression:
			switch x.Case {
			case cc.ConditionalExpressionLOr: // LogicalOrExpression
				n = x.LogicalOrExpression
			default:
				return nil
			}
		case *cc.AdditiveExpression:
			switch x.Case {
			case cc.AdditiveExpressionMul: // MultiplicativeExpression
				n = x.MultiplicativeExpression
			default:
				return nil
			}
		default:
			panic(todo("%T", n))
		}
	}
	return nil
}

func (c *ctx) postfixExpressionCall(w writer, n *cc.PostfixExpression) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	var ft *cc.FunctionType
	var d *cc.Declarator
	switch d = c.declaratorOf(n.PostfixExpression); {
	case d != nil:
		switch x := d.Type().(type) {
		case *cc.PointerType:
			var ok bool
			if ft, ok = x.Elem().(*cc.FunctionType); !ok {
				c.err(errorf("TODO %T", x.Elem()))
				return
			}
		case *cc.FunctionType:
			ft = x
		default:
			c.err(errorf("TODO %T", d.Type()))
			return
		}
	default:
		pt, ok := n.PostfixExpression.Type().(*cc.PointerType)
		if !ok {
			c.err(errorf("TODO %T", n.PostfixExpression.Type()))
			return
		}

		if ft, ok = pt.Elem().(*cc.FunctionType); !ok {
			c.err(errorf("TODO %T", pt.Elem()))
			return
		}
	}

	var args []cc.ExpressionNode
	for l := n.ArgumentExpressionList; l != nil; l = l.ArgumentExpressionList {
		args = append(args, l.AssignmentExpression)
	}
	if len(args) < ft.MinArgs() {
		c.err(errorf("%v: too few arguments to function '%s', type '%v' in '%v'", c.pos(n.PostfixExpression), cc.NodeSource(n.PostfixExpression), ft, cc.NodeSource(n)))
		return &b, nil, 0
	}

	if len(args) > ft.MaxArgs() && ft.MaxArgs() >= 0 {
		c.err(errorf("%v: too many arguments to function '%s', type '%v' in '%v'", c.pos(n.PostfixExpression), cc.NodeSource(n.PostfixExpression), ft, cc.NodeSource(n)))
		return &b, nil, 0
	}

	// trc("%v: len(args) %v, ft.MaxArgs %v, ft.IsVariadic() %v, d != nil %v, d.IsSynthetic() %v, d.IsFuncDef() %v", n.Position(), len(args), ft.MaxArgs(), ft.IsVariadic(), d != nil, d.IsSynthetic(), d.IsFuncDef())
	if len(args) > ft.MaxArgs() && !ft.IsVariadic() && d != nil && !d.IsSynthetic() && d.IsFuncDef() {
		max := mathutil.Max(ft.MaxArgs(), 0)
		for _, v := range args[max:] {
			w.w("%s_ = %s;", tag(preserve), c.expr(w, v, nil, exprDefault))
		}
		args = args[:max]
	}

	params := ft.Parameters()
	var xargs []*buf
	for i, v := range args {
		mode := exprDefault
		var t cc.Type
		switch {
		case i < len(params):
			t = params[i].Type()
		default:
			switch t = v.Type(); {
			case cc.IsIntegerType(t):
				t = cc.IntegerPromotion(t)
			case t.Kind() == cc.Float:
				t = c.ast.Double
			}
		}
		switch v.Type().Undecay().Kind() {
		case cc.Function:
			mode = exprUintptr
		}
		xargs = append(xargs, c.expr(w, v, t, mode))
	}
	switch {
	case c.f == nil:
		b.w("%s(%snil", c.expr(w, n.PostfixExpression, nil, exprCall), tag(preserve))
	default:
		b.w("%s(%stls", c.expr(w, n.PostfixExpression, nil, exprCall), tag(ccgo))
	}
	switch {
	case ft.IsVariadic():
		for _, v := range xargs[:ft.MinArgs()] {
			b.w(", %s", v)
		}
		switch {
		case len(xargs) == ft.MinArgs():
			b.w(", 0")
		default:
			b.w(", %s%sVaList(%s", c.task.tlsQualifier, tag(preserve), bpOff(c.f.tlsAllocs+8))
			if n := len(xargs[ft.MinArgs():]); n > c.f.maxValist {
				c.f.maxValist = n
			}
			for _, v := range xargs[ft.MinArgs():] {
				b.w(", %s", v)
			}
			b.w(")")
		}
	default:
		for _, v := range xargs {
			b.w(", %s", v)
		}
	}
	b.w(")")
	rt, rmode = ft.Result(), exprDefault
	if rt.Kind() == cc.Void {
		rmode = exprVoid
	}
	return &b, rt, rmode
}

func (c *ctx) assignmentExpression(w writer, n *cc.AssignmentExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.AssignmentExpressionCond: // ConditionalExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.AssignmentExpressionAssign: // UnaryExpression '=' AssignmentExpression
		switch x := n.UnaryExpression.(type) {
		case *cc.PostfixExpression:
			switch x.Case {
			case cc.PostfixExpressionSelect:
				f := x.Field()
				if !f.IsBitfield() {
					break
				}

				//TODO do not pin/use pointer
				p := c.pin(n, c.expr(w, x.PostfixExpression, x.PostfixExpression.Type().Pointer(), exprUintptr))
				switch mode {
				case exprDefault:
					b.w("%sAssignBitFieldPtr%d%s(%s+%d, %s, %d, %d, %#0x)", c.task.tlsQualifier, f.AccessBytes()*8, c.helper(n, f.Type()), p, f.Offset(), c.expr(w, n.AssignmentExpression, f.Type(), exprDefault), f.ValueBits(), f.OffsetBits(), f.Mask())
					return &b, f.Type(), exprDefault
				case exprVoid:
					b.w("%sSetBitFieldPtr%d%s(%s+%d, %s, %d, %#0x)", c.task.tlsQualifier, f.AccessBytes()*8, c.helper(n, f.Type()), p, f.Offset(), c.expr(w, n.AssignmentExpression, f.Type(), exprDefault), f.OffsetBits(), f.Mask())
					return &b, n.Type(), exprVoid
				default:
					trc("%v: BITFIELD", n.Position())
					c.err(errorf("TODO %v", mode))
					return &b, rt, rmode
				}
			case cc.PostfixExpressionPSelect:
				f := x.Field()
				if !f.IsBitfield() {
					break
				}

				switch mode {
				case exprDefault:
					b.w("%sAssignBitFieldPtr%d%s(%s+%d, %s, %d, %d, %#0x)", c.task.tlsQualifier, f.AccessBytes()*8, c.helper(n, f.Type()), c.expr(w, x.PostfixExpression, nil, exprDefault), f.Offset(), c.expr(w, n.AssignmentExpression, f.Type(), exprDefault), f.ValueBits(), f.OffsetBits(), f.Mask())
					return &b, f.Type(), exprDefault
				case exprVoid:
					b.w("%sSetBitFieldPtr%d%s(%s+%d, %s, %d, %#0x)", c.task.tlsQualifier, f.AccessBytes()*8, c.helper(n, f.Type()), c.expr(w, x.PostfixExpression, nil, exprDefault), f.Offset(), c.expr(w, n.AssignmentExpression, f.Type(), exprDefault), f.OffsetBits(), f.Mask())
					return &b, n.Type(), exprVoid
				default:
					trc("%v: BITFIELD", n.Position())
					c.err(errorf("TODO %v", mode))
					return &b, rt, rmode
				}
			}
		}
		switch mode {
		case exprDefault:
			rt, rmode = n.Type(), exprDefault
			v := c.f.newAutovar(n, n.UnaryExpression.Type())
			w.w("%s = %s;", v, c.expr(w, n.AssignmentExpression, n.UnaryExpression.Type(), exprDefault))
			w.w("%s = %s;", c.expr(w, n.UnaryExpression, nil, exprDefault), v)
			b.w("%s", v)
		case exprVoid:
			b.w("%s = %s", c.expr(w, n.UnaryExpression, nil, exprLvalue), c.expr(w, n.AssignmentExpression, n.UnaryExpression.Type(), exprDefault))
			rt, rmode = n.Type(), exprVoid
		default:
			c.err(errorf("TODO %v", mode))
		}
	case cc.AssignmentExpressionMul, // UnaryExpression "*=" AssignmentExpression
		cc.AssignmentExpressionDiv, // UnaryExpression "/=" AssignmentExpression
		cc.AssignmentExpressionMod, // UnaryExpression "%=" AssignmentExpression
		cc.AssignmentExpressionAdd, // UnaryExpression "+=" AssignmentExpression
		cc.AssignmentExpressionSub, // UnaryExpression "-=" AssignmentExpression
		cc.AssignmentExpressionLsh, // UnaryExpression "<<=" AssignmentExpression
		cc.AssignmentExpressionRsh, // UnaryExpression ">>=" AssignmentExpression
		cc.AssignmentExpressionAnd, // UnaryExpression "&=" AssignmentExpression
		cc.AssignmentExpressionXor, // UnaryExpression "^=" AssignmentExpression
		cc.AssignmentExpressionOr:  // UnaryExpression "|=" AssignmentExpression

		rt, rmode = n.Type(), mode
		op := n.Token.SrcStr()
		op = op[:len(op)-1]
		x, y := n.UnaryExpression.Type(), n.AssignmentExpression.Type()
		ct := c.usualArithmeticConversions(x, y)
		ut := n.UnaryExpression.Type()
		switch x := n.UnaryExpression.(type) {
		case *cc.PostfixExpression:
			switch x.Case {
			case cc.PostfixExpressionSelect:
				f := x.Field()
				if !f.IsBitfield() {
					break
				}

				p := c.pin(n, c.expr(w, x.PostfixExpression, x.PostfixExpression.Type().Pointer(), exprUintptr))
				bf, _, _ := c.bitField(w, n, p, f, exprDefault)
				switch mode {
				case exprDefault, exprVoid:
					b.w("%sAssignBitFieldPtr%d%s(%s+%d, %s(%s(%s)%s%[7]s(%[10]s)), %d, %d, %#0x)",
						c.task.tlsQualifier, f.AccessBytes()*8, c.helper(n, ut), p, f.Offset(),
						c.typ(n, ut), c.typ(n, ct), bf,
						op,
						c.expr(w, n.AssignmentExpression, ut, exprDefault),
						f.ValueBits(), f.OffsetBits(), f.Mask(),
					)
					return c.reduceBitFieldValue(&b, f, f.Type(), rmode), rt, exprDefault
				default:
					trc("%v: BITFIELD %v", n.Position(), mode)
					c.err(errorf("TODO %v", mode))
					return &b, rt, rmode
				}
			case cc.PostfixExpressionPSelect:
				f := x.Field()
				if !f.IsBitfield() {
					break
				}

				p := c.expr(w, x.PostfixExpression, nil, exprDefault)
				bf, _, _ := c.bitField(w, n, p, f, exprDefault)
				switch mode {
				case exprDefault, exprVoid:
					b.w("%sAssignBitFieldPtr%d%s(%s+%d, %s(%s(%s)%s%[7]s(%[10]s)), %d, %d, %#0x)",
						c.task.tlsQualifier, f.AccessBytes()*8, c.helper(n, ut), p, f.Offset(),
						c.typ(n, ut), c.typ(n, ct), bf,
						op,
						c.expr(w, n.AssignmentExpression, ut, exprDefault),
						f.ValueBits(), f.OffsetBits(), f.Mask(),
					)
					return c.reduceBitFieldValue(&b, f, f.Type(), rmode), rt, exprDefault
				default:
					trc("%v: BITFIELD", n.Position())
					c.err(errorf("TODO %v", mode))
					return &b, rt, rmode
				}
			}
		}

		var k, v string
		switch n.Case {
		case cc.AssignmentExpressionAdd: // UnaryExpression "+=" AssignmentExpression
			switch {
			case x.Kind() == cc.Ptr && cc.IsIntegerType(y):
				if sz := x.(*cc.PointerType).Elem().Undecay().Size(); sz != 1 {
					k = fmt.Sprintf("*%d", sz)
				}
			case cc.IsIntegerType(x) && y.Kind() == cc.Ptr:
				c.err(errorf("TODO")) // -
			}
		case cc.AssignmentExpressionSub: // UnaryExpression "-=" AssignmentExpression
			switch {
			case x.Kind() == cc.Ptr && cc.IsIntegerType(y):
				if sz := x.(*cc.PointerType).Elem().Undecay().Size(); sz != 1 {
					k = fmt.Sprintf("*%d", sz)
				}
			case x.Kind() == cc.Ptr && y.Kind() == cc.Ptr:
				if sz := x.(*cc.PointerType).Elem().Undecay().Size(); sz != 1 {
					k = fmt.Sprintf("/%d", sz)
				}
			}
		}
		switch mode {
		case exprDefault, exprVoid:
			switch d := c.declaratorOf(n.UnaryExpression); {
			case d != nil:
				v = fmt.Sprintf("%s", c.expr(w, n.UnaryExpression, nil, exprDefault))
			default:
				p := fmt.Sprintf("%sp%d", tag(ccgo), c.id())
				ps := fmt.Sprintf("var %s %s;", p, c.typ(n, ut.Pointer()))
				c.f.registerAutoVar(ps)
				w.w("\n%s = %s;", p, c.expr(w, n.UnaryExpression, ut.Pointer(), exprUintptr))
				v = fmt.Sprintf("(*(*%s)(%s))", c.typ(n, ut), unsafePointer(p))
			}
			w.w("\n%s = %s((%s(%s)) %s ((%s)%s));", v, c.typ(n, ut), c.typ(n, ct), v, op, c.expr(w, n.AssignmentExpression, ct, exprDefault), k)
			if mode == exprDefault {
				b.w("%s", v)
			}
		default:
			c.err(errorf("TODO %v", mode))
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) expressionList(w writer, n *cc.ExpressionList, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	for ; n != nil; n = n.ExpressionList {
		switch {
		case n.ExpressionList == nil:
			return c.expr0(w, n.AssignmentExpression, t, mode)
		default:
			w.w("%s%s;", sep(n.AssignmentExpression), c.expr(w, n.AssignmentExpression, nil, exprVoid))
		}
	}
	c.err(errorf("TODO internal error", n))
	return r, rt, rmode
}

func (c *ctx) expressionListLast(n cc.ExpressionNode) cc.ExpressionNode {
	for {
		switch x := n.(type) {
		case *cc.ExpressionList:
			for n := x; n != nil; n = n.ExpressionList {
				switch {
				case n.ExpressionList == nil:
					return n.AssignmentExpression
				}
			}
			return nil
		case *cc.PrimaryExpression:
			if x.Case == cc.PrimaryExpressionExpr {
				n = x.ExpressionList
				break
			}

			return n
		default:
			return n
		}
	}
}

func (c *ctx) primaryExpression(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
out:
	switch n.Case {
	case cc.PrimaryExpressionIdent: // IDENTIFIER
		rt, rmode = n.Type(), mode
		switch x := n.ResolvedTo().(type) {
		case *cc.Declarator:
			nm := x.Name()
			linkName := c.declaratorTag(x) + nm
			if c.pass == 2 {
				if nm := c.f.locals[x]; nm != "" {
					linkName = nm
				}
			}
			c.externsMentioned[nm] = struct{}{}
			b.n = x
			var info *declInfo
			if c.f != nil {
				info = c.f.declInfos.info(x)
			}
			switch {
			case info != nil && info.pinned():
				switch mode {
				case exprLvalue, exprSelect, exprIndex:
					b.w("(*(*%s)(%s))", c.typ(n, x.Type()), unsafePointer(bpOff(info.bpOff)))
				case exprUintptr:
					rt = x.Type().Pointer()
					b.w("%s", bpOff(info.bpOff))
				case exprDefault, exprVoid:
					rmode = exprDefault
					switch _, ok := n.Type().Undecay().(*cc.ArrayType); {
					case ok && !x.IsParam():
						b.w("%s", bpOff(info.bpOff))
					default:
						b.w("(*(*%s)(%s))", c.typ(n, x.Type()), unsafePointer(bpOff(info.bpOff)))
					}
				case exprCall:
					switch y := x.Type().Undecay().(type) {
					case *cc.PointerType:
						if ft, ok := y.Elem().(*cc.FunctionType); ok {
							b.w("(*(*func%s)(%s))", c.signature(ft, false, false, true), unsafePointer(bpOff(info.bpOff)))
							break
						}

						c.err(errorf("TODO %T:", y.Elem()))
					default:
						//b.w("(*(*func%s)(%s))", c.signature(x.Type().(*cc.FunctionType), false, false), unsafePointer(bpOff(info.bpOff)))
						c.err(errorf("TODO %T", y))
					}
				default:
					c.err(errorf("TODO %v %v:", mode, n.Position()))
				}
			default:
				switch mode {
				case exprVoid:
					r, rt, _ = c.primaryExpression(w, n, t, exprDefault)
					return r, rt, exprDefault
				case exprDefault:
					switch x.Type().Kind() {
					case cc.Array:
						p := &buf{n: x}
						p.w("%s", linkName)
						b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n, p)))
					case cc.Function:
						// v := fmt.Sprintf("%sf%d", tag(ccgo), c.id())
						// switch {
						// case c.f != nil:
						// 	w.w("%s := %s;", v, linkName)
						// default:
						// 	w.w("var %s = %s;", v, linkName)
						// }
						// b.w("(*(*%suintptr)(%s))", tag(preserve), unsafeAddr(v))

						// b.w("(*(*%suintptr)(%sunsafe.%[2]sPointer(&struct{f func%[3]s}{%s})))", tag(preserve), tag(importQualifier), c.signature(x.Type().(*cc.FunctionType), false, false, true), linkName)
						b.w("%s%s(%s)", tag(preserve), ccgoFP, linkName)
					default:
						b.w("%s", linkName)
					}
				case exprLvalue, exprSelect:
					b.w("%s", linkName)
				case exprCall:
					switch y := x.Type().(type) {
					case *cc.FunctionType:
						if !c.task.strictISOMode {
							if _, ok := forcedBuiltins[nm]; ok {
								nm = "__builtin_" + nm
								linkName = c.declaratorTag(x) + nm
							}
						}
						b.w("%s", linkName)
					case *cc.PointerType:
						switch z := y.Elem().(type) {
						case *cc.FunctionType:
							rmode = exprUintptr
							b.w("%s", linkName)
						default:
							// trc("%v: %s", x.Position(), cc.NodeSource(n))
							c.err(errorf("TODO %T", z))
						}
					default:
						c.err(errorf("TODO %T", y))
					}
				case exprIndex:
					switch x.Type().Undecay().Kind() {
					case cc.Array:
						b.w("%s", linkName)
					default:
						panic(todo(""))
						c.err(errorf("TODO %v", mode))
					}
				case exprUintptr:
					rt = x.Type().Pointer()
					switch {
					case x.Type().Kind() == cc.Function:
						// v := fmt.Sprintf("%sf%d", tag(ccgo), c.id())
						// switch {
						// case c.f != nil:
						// 	w.w("%s := %s;", v, linkName)
						// default:
						// 	w.w("var %s = %s;", v, linkName)
						// }
						// b.w("(*(*%suintptr)(%s))", tag(preserve), unsafeAddr(v)) // Free pass from .pin

						// b.w("(*(*%suintptr)(%sunsafe.%[1]sPointer(&struct{f func%[3]s}{%s})))", tag(preserve), tag(importQualifier), c.signature(x.Type().(*cc.FunctionType), false, false, true), linkName)
						b.w("%s%s(%s)", tag(preserve), ccgoFP, linkName)
					default:
						p := &buf{n: x}
						p.w("%s", linkName)
						b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n, p)))
					}
				default:
					c.err(errorf("TODO %v", mode))
				}
			}
		case *cc.Enumerator:
			switch {
			case x.ResolvedIn().Parent == nil:
				rt, rmode = n.Type(), exprDefault
				b.w("(%s%s%sFrom%s(%s%s))", c.task.tlsQualifier, tag(preserve), c.helper(n, n.Type()), c.helper(n, n.Type()), tag(enumConst), x.Token.Src())
			default:
				rt, rmode = n.Type(), exprDefault
				b.w("%v", n.Value())
			}
		case nil:
			switch mode {
			case exprCall:
				b.w("%s%s", tag(external), n.Token.Src())
				break out
			default:
				c.err(errorf("TODO %v: %v", n.Position(), mode))
				break out
			}
		default:
			c.err(errorf("TODO %T", x))
		}
	case cc.PrimaryExpressionInt: // INTCONST
		return c.primaryExpressionIntConst(w, n, t, mode)
	case cc.PrimaryExpressionFloat: // FLOATCONST
		return c.primaryExpressionFloatConst(w, n, t, mode)
	case cc.PrimaryExpressionChar: // CHARCONST
		return c.primaryExpressionCharConst(w, n, t, mode)
	case cc.PrimaryExpressionLChar: // LONGCHARCONST
		return c.primaryExpressionLCharConst(w, n, t, mode)
	case cc.PrimaryExpressionString: // STRINGLITERAL
		return c.primaryExpressionStringConst(w, n, t, mode)
	case cc.PrimaryExpressionLString: // LONGSTRINGLITERAL
		return c.primaryExpressionLStringConst(w, n, t, mode)
	case cc.PrimaryExpressionExpr: // '(' ExpressionList ')'
		return c.expr0(w, n.ExpressionList, nil, mode)
	case cc.PrimaryExpressionStmt: // '(' CompoundStatement ')'
		// trc("%v: %v %s", n.Position(), n.Type(), cc.NodeSource(n))
		switch n.Type().Kind() {
		case cc.Void:
			rt, rmode = n.Type(), exprVoid
			c.compoundStatement(w, n.CompoundStatement, false, "")
		default:
			rt, rmode = n.Type(), exprDefault
			v := fmt.Sprintf("%sv%d", tag(ccgo), c.id())
			w.w("var %s %s;", v, c.typ(n, n.Type()))
			c.compoundStatement(w, n.CompoundStatement, false, v)
			b.w("%s", v)
		}
	case cc.PrimaryExpressionGeneric: // GenericSelection
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) utf16(n cc.Node, s cc.UTF16StringValue) string {
	b := bytes.NewBuffer(make([]byte, 0, 2*len(s)))
	bo := c.ast.ABI.ByteOrder
	for _, v := range s {
		if err := binary.Write(b, bo, v); err != nil {
			c.err(errorf("%v: %v", n.Position(), err))
			return ""
		}
	}
	return b.String()
}

func (c *ctx) utf32(n cc.Node, s cc.UTF32StringValue) string {
	b := bytes.NewBuffer(make([]byte, 0, 4*len(s)))
	bo := c.ast.ABI.ByteOrder
	for _, v := range s {
		if err := binary.Write(b, bo, v); err != nil {
			c.err(errorf("%v: %v", n.Position(), err))
			return ""
		}
	}
	return b.String()
}

func (c *ctx) primaryExpressionLStringConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch x := n.Type().Undecay().(type) {
	case *cc.ArrayType:
		switch y := t.(type) {
		case *cc.ArrayType:
			switch z := n.Value().(type) {
			case cc.UTF16StringValue:
				c.err(errorf("TODO %T", z))
			case cc.UTF32StringValue:
				for len(z) != 0 && z[len(z)-1] == 0 {
					z = z[:len(z)-1]
				}
				b.w("%s{", c.typ(n, y))
				for _, c := range z {
					b.w("%s, ", strconv.QuoteRune(c))
				}
				b.w("}")
			default:
				c.err(errorf("TODO %T", z))
			}
		case *cc.PointerType:
			switch z := n.Value().(type) {
			case cc.UTF16StringValue:
				b.w("%q", c.utf16(n, z))
			case cc.UTF32StringValue:
				b.w("%q", c.utf32(n, z))
			default:
				c.err(errorf("TODO %T", z))
			}
		default:
			c.err(errorf("TODO %T", y))
		}
	default:
		c.err(errorf("TODO %T", x))
	}
	return &b, t, exprDefault
}

func (c *ctx) primaryExpressionStringConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	s := n.Value()
	var b buf
	switch x := n.Type().Undecay().(type) {
	case *cc.ArrayType:
		switch {
		case c.isCharType(x.Elem()):
			s := string(s.(cc.StringValue))
			switch t.Kind() {
			case cc.Array:
				to := t.(*cc.ArrayType)
				max := to.Len()
				a := []byte(s)
				for len(a) != 0 && a[len(a)-1] == 0 {
					a = a[:len(a)-1]
				}
				b.w("%s{", c.typ(n, to))
				for i := 0; i < len(a) && int64(i) < max; i++ {
					b.w("%s, ", c.stringCharConst(a[i], to.Elem()))
				}
				b.w("}")
			case cc.Ptr:
				t := t.(*cc.PointerType)
				if c.isCharType(t.Elem()) || t.Elem().Kind() == cc.Void {
					b.w("%q", s)
					break
				}

				c.err(errorf("TODO"))
			default:
				if cc.IsIntegerType(t) {
					b.w("(%s(%q))", c.typ(n, t), s)
					break
				}

				trc("%v: %s <- %q, convert to %s", n.Position(), x, s, t)
				c.err(errorf("TODO %s", t))
			}
		default:
			c.err(errorf("TODO"))
		}
	default:
		c.err(errorf("TODO %T", x))
	}
	return &b, t, exprDefault
}

func (c *ctx) primaryExpressionLCharConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	rt, rmode = t, exprDefault
	var b buf
	src := n.Token.SrcStr()
	lit := src[1:] // L
	val, err := strconv.Unquote(lit)
	if err != nil {
		switch {
		case strings.HasPrefix(lit, `'\`) && isOctalString(lit[2:len(lit)-1]):
			lit = fmt.Sprintf(`'\%03o'`, n.Value())
			if val, err = strconv.Unquote(lit); err != nil {
				lit = fmt.Sprintf(`'\u%04x'`, n.Value())
			}
		case src == `'\"'`:
			lit = `'"'`
		}
		if val, err = strconv.Unquote(lit); err != nil {
			c.err(errorf("TODO `%s` -> %s", lit, err))
			return &b, rt, rmode
		}
	}

	ch := []rune(val)[0]
	switch x := n.Value().(type) {
	case cc.Int64Value:
		if rune(x) != ch {
			c.err(errorf("TODO `%s` -> |% x|, exp %#0x", lit, val, x))
			return &b, rt, rmode
		}
	case cc.UInt64Value:
		if rune(x) != ch {
			c.err(errorf("TODO `%s` -> |% x|, exp %#0x", lit, val, x))
			return &b, rt, rmode
		}
	}

	b.w("(%s%s%sFromInt32(%s))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), lit)
	return &b, rt, rmode
}

func (c *ctx) primaryExpressionCharConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	rt, rmode = t, exprDefault
	var b buf
	src := n.Token.SrcStr()
	lit := src
	val, err := strconv.Unquote(lit)
	if err != nil {
		switch {
		case strings.HasPrefix(src, `'\`) && isOctalString(src[2:len(src)-1]):
			lit = fmt.Sprintf(`'\%03o'`, n.Value())
		case src == `'\"'`:
			lit = `'"'`
		}
		if val, err = strconv.Unquote(lit); err != nil {
			c.err(errorf("TODO `%s` -> %s", lit, err))
			return &b, rt, rmode
		}
	}
	if len(val) != 1 {
		c.err(errorf("TODO `%s` -> |% x|", lit, val))
		return &b, rt, rmode
	}

	ch := val[0]
	switch x := n.Value().(type) {
	case cc.Int64Value:
		if byte(x) != ch {
			c.err(errorf("TODO `%s` -> |% x|, exp %#0x", lit, val, x))
			return &b, rt, rmode
		}
	case cc.UInt64Value:
		if byte(x) != ch {
			c.err(errorf("TODO `%s` -> |% x|, exp %#0x", lit, val, x))
			return &b, rt, rmode
		}
	}

	b.w("(%s%s%sFromUint8(%s))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), lit)
	return &b, rt, rmode
}

func (c *ctx) primaryExpressionIntConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	rt, rmode = t, exprDefault
	var b buf
	src := n.Token.SrcStr()
	lit := strings.TrimRight(src, "uUlL")
	var want uint64
	switch x := n.Value().(type) {
	case cc.Int64Value:
		want = uint64(x)
	case cc.UInt64Value:
		want = uint64(x)
	default:
		c.err(errorf("TODO %T", x))
		return &b, rt, rmode
	}

	val, err := strconv.ParseUint(lit, 0, 64)
	if err != nil {
		c.err(errorf("TODO `%s` -> %s", lit, err))
		return &b, rt, rmode
	}

	if val != want {
		c.err(errorf("TODO `%s` -> got %v, want %v", lit, val, want))
		return &b, rt, rmode
	}

	if t.Kind() == cc.Void {
		b.w("(%s)", lit)
		return &b, rt, rmode
	}

	b.w("(%s%s%sFrom%s(%s))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), c.helper(n, n.Type()), lit)
	return &b, rt, rmode
}

func (c *ctx) primaryExpressionFloatConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	rt, rmode = t, exprDefault
	switch x := n.Value().(type) {
	case *cc.LongDoubleValue:
		// b.w("(%s%s%sFrom%s(%v))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), c.helper(n, n.Type()), (*big.Float)(x))
		b.w("(%s%s%sFrom%s(%v))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), c.helper(n, c.ast.Double), (*big.Float)(x))
	case cc.Float64Value:
		b.w("(%s%s%sFrom%s(%v))", c.task.tlsQualifier, tag(preserve), c.helper(n, t), c.helper(n, n.Type()), x)
	default:
		c.err(errorf("TODO %T", x))
	}
	return &b, rt, rmode
}

func (c *ctx) stringCharConst(b byte, t cc.Type) string {
	switch {
	case b >= ' ' && b < 0x7f:
		return strconv.QuoteRuneToASCII(rune(b))
	case cc.IsSignedInteger(t):
		return fmt.Sprint(int8(b))
	default:
		return fmt.Sprint(b)
	}
}
