// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"

	"modernc.org/cc/v4"
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
	exprUintpr       // C pointer, Go uintptr
	exprUntyped      // C primary expr literal, Go typed literal
	exprVoid         // C void, no Go equivalent
)

func (c *ctx) expr(w writer, n cc.ExpressionNode, to cc.Type, toMode mode) *buf {
	if toMode == 0 {
		c.err(errorf("internal error"))
		return nil
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
	r, from, fromMode := c.expr0(w, n, to, toMode)
	if from == nil || fromMode == 0 {
		// trc("", cpos(n))
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

	if from != nil && from.Kind() == cc.Enum {
		from = from.(*cc.EnumType).UnderlyingType()
	}
	if to.Kind() == cc.Enum {
		to = to.(*cc.EnumType).UnderlyingType()
	}
	if fromMode == exprUntyped {
		return c.convertUntyped(n, s, from, to, fromMode, toMode)
	}

	if from == to || from != nil && from.IsCompatible(to) {
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

	if toMode == exprVoid {
		return s
	}

	if to.Kind() == cc.Ptr {
		return c.convertToPointer(n, s, from, to.(*cc.PointerType), fromMode, toMode)
	}

	c.err(errorf("TODO %q %s %s -> %s %s", s, from, fromMode, to, toMode))
	return s //TODO
}

func (c *ctx) convertToPointer(n cc.ExpressionNode, s *buf, from cc.Type, to *cc.PointerType, fromMode, toMode mode) (r *buf) {
	var b buf
	switch fromMode {
	case exprDefault:
		switch toMode {
		case exprUintpr:
			b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n, s)))
			return &b
		}
	}

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
				//trc("%v: PIN %v at %v (%v: %v: %v:)", c.pos(n), x.Name(), c.pos(x), origin(4), origin(3), origin(2))
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
				//trc("%v: PIN %v at %v (%v: %v:)", c.pos(n), y.Name(), c.pos(y), origin(3), origin(2))
			}
			return b
		}
	}
	return b
}

func (c *ctx) convertUntyped(n cc.ExpressionNode, s *buf, from, to cc.Type, fromMode, toMode mode) (r *buf) {
	// defer func() { trc("%v: from %v: %v, to %v: %v %q -> %q", c.pos(n), from, fromMode, to, toMode, s, r) }()
	var b buf
	switch toMode {
	case exprDefault:
		// ok
	case exprBool:
		toMode = exprDefault
		var b2 buf
		defer func() { b2.w("(%s != 0)", b.bytes()); r = &b2 }()
	default:
		c.err(errorf("TODO %v: %v", n.Position(), toMode))
	}
	var val cc.Value
	switch x := s.n.(type) {
	case *cc.PrimaryExpression:
		val = x.Value()
	case nil:
		// ok
	default:
		c.err(errorf("TODO %T", x))
	}

	switch x := val.(type) {
	case cc.Int64Value:
		switch {
		case to.Kind() == cc.Ptr:
			b.w("%s(%d)", c.typ(to), uint64(x))
			return &b
		case cc.IsIntegerType(to):
			switch {
			case cc.IsSignedInteger(to):
				if to.Size() < 8 {
					sbit := uint64(1) << (to.Size()*8 - 1)
					switch {
					case x&cc.Int64Value(sbit) != 0:
						x |= ^cc.Int64Value(sbit<<1 - 1)
					default:
						x &= cc.Int64Value(sbit - 1)
					}
				}
				b.w("%s(%d)", c.typ(to), int64(x))
				return &b
			default:
				if to.Size() < 8 {
					m := uint64(1)<<(to.Size()*8) - 1
					x &= cc.Int64Value(m)
				}
				b.w("%s(%d)", c.typ(to), uint64(x))
				return &b
			}
		case cc.IsArithmeticType(to):
			b.w("%s(%v)", c.typ(to), x)
			return &b
		default:
			//TODO c.err(errorf("TODO %v", to))
		}
	case cc.UInt64Value:
		switch {
		case to.Kind() == cc.Ptr:
			b.w("%s(%d)", c.typ(to), x)
			return &b
		case cc.IsIntegerType(to):
			switch {
			case cc.IsSignedInteger(to):
				if to.Size() < 8 {
					sbit := uint64(1) << (to.Size()*8 - 1)
					switch {
					case x&cc.UInt64Value(sbit) != 0:
						x |= ^cc.UInt64Value(sbit<<1 - 1)
					default:
						x &= cc.UInt64Value(sbit - 1)
					}
				}
				b.w("%s(%d)", c.typ(to), int64(x))
				return &b
			default:
				if to.Size() < 8 {
					m := uint64(1)<<(to.Size()*8) - 1
					x &= cc.UInt64Value(m)
				}
				b.w("%s(%d)", c.typ(to), x)
				return &b
			}
		case cc.IsArithmeticType(to):
			b.w("%s(%v)", c.typ(to), x)
			return &b
		default:
			//TODO c.err(errorf("TODO %v", to))
		}
	case cc.Float64Value:
		switch {
		case x == 0:
			b.w("(%s%s%s(%v))", c.task.tlsQualifier, tag(preserve), c.helper(to), x)
		default:
			b.w("%s(%v)", c.typ(to), x)
		}
		return &b
	case *cc.LongDoubleValue:
		bf := (*big.Float)(x)
		switch {
		case bf.Sign() == 0:
			b.w("(%s%s%s(%v))", c.task.tlsQualifier, tag(preserve), c.helper(to), bf)
		default:
			b.w("%s(%v)", c.typ(to), bf)
		}
		return &b
	case nil:
		// ok
	default:
		//c.err(errorf("TODO %T", x))
	}

	from = from.Undecay()
	switch {
	case from.Kind() == cc.Array:
		switch from := from.(*cc.ArrayType); {
		case c.isCharType(from.Elem()):
			switch to.Kind() {
			case cc.Array:
				to := to.(*cc.ArrayType)
				max := to.Len()
				a := s.bytes()
				for len(a) != 0 && a[len(a)-1] == 0 {
					a = a[:len(a)-1]
				}
				b.w("%s{", c.typ(to))
				for i := 0; i < len(a) && int64(i) < max; i++ {
					b.w("%s, ", c.stringCharConst(a[i], to.Elem()))
				}
				b.w("}")
				return &b
			case cc.Ptr:
				to := to.(*cc.PointerType)
				if c.isCharType(to.Elem()) || to.Elem().Kind() == cc.Void {
					b.w("%q", s)
					return &b
				}
			default:
				c.err(errorf("TODO"))
			}
		default:
			c.err(errorf("TODO"))
		}
	case cc.IsIntegerType(from):
		switch {
		case cc.IsIntegerType(to):
			var val uint64
			switch {
			case cc.IsSignedInteger(from):
				v, err := strconv.ParseInt(string(s.bytes()), 10, 64)
				if err != nil {
					c.err(errorf("internal error: %v", err))
					break
				}

				val = uint64(v)
			default:
				v, err := strconv.ParseUint(string(s.bytes()), 10, 64)
				if err != nil {
					c.err(errorf("internal error: %v", err))
					break
				}

				val = uint64(v)
			}
			switch {
			case cc.IsSignedInteger(to):
				if to.Size() < 8 {
					m := uint64(1)<<(to.Size()*8) - 1
					switch {
					case int64(val) < 0:
						val |= m
					default:
						val &= m
					}
				}
				b.w("int%d(%d)", 8*to.Size(), int64(val))
			default:
				if to.Size() < 8 {
					m := uint64(1)<<(to.Size()*8) - 1
					val &= m
				}
				b.w("uint%d(%d)", 8*to.Size(), uint64(val))
			}
		case to.Kind() == cc.Ptr:
			b.w("%suintptr(%s)", tag(preserve), s)
		case cc.IsArithmeticType(to):
			b.w("%s(%s)", c.typ(to), s)
		case to.Kind() == cc.Struct:
			c.err(errorf("TODO"))
		case to.Kind() == cc.Union:
			c.err(errorf("TODO"))
		case to.Kind() == cc.Array:
			if !isZeroString(string(s.bytes())) {
				c.err(errorf("TODO"))
			}
			return s
		case to.Kind() == cc.Void:
			return nil
		default:
			// trc("%v: %q %s %s -> %s %s", cpos(n), s, from, fromMode, to, toMode)
			c.err(errorf("TODO"))
		}
		return &b
	case cc.IsArithmeticType(from):
		switch {
		case cc.IsArithmeticType(to):
			b.w("%s(%s)", c.typ(to), s)
		default:
			c.err(errorf("TODO"))
		}
		return &b
	}

	c.err(errorf("TODO %q %s %s -> %s %s", s, from, fromMode, to, toMode))
	return s //TODO
}

func isZeroString(s string) bool { return s == "0" }

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
					b.w("((*%s)(%s))", c.typ(y), unsafeAddr(s))
					return &b
				default:
					trc("%T", y)
				}
			default:
				trc("%T", x)
			}
		}
	case exprUintpr:
		switch toMode {
		case exprDefault:
			return s
		case exprCall:
			v := fmt.Sprintf("%sf%d", tag(ccgo), c.id())
			ft := from.(*cc.PointerType).Elem().(*cc.FunctionType)
			w.w("var %s func%s;/**/", v, c.signature(ft, false, false))
			w.w("\n*(*%suintptr)(%s) = %s;", tag(preserve), unsafeAddr(v), s) // Free pass from .pin
			var b buf
			b.w("%s", v)
			return &b
		}
	case exprBool:
		switch toMode {
		case exprDefault:
			switch {
			case cc.IsIntegerType(to):
				b.w("%s%sBool%s(%s)", c.task.tlsQualifier, tag(preserve), c.helper(to), s)
				return &b
			}
		}
	case exprVoid:
		switch toMode {
		case exprDefault:
			return s
		}
	}
	c.err(errorf("TODO %q %s %s -> %s %s", s, from, fromMode, to, toMode))
	return s //TODO
}

func (c *ctx) isIdent(s string) bool {
	for i, v := range s {
		switch {
		case i == 0:
			if !unicode.IsLetter(v) && v != '_' {
				return false
			}
		default:
			if !unicode.IsLetter(v) && v != '_' && !unicode.IsDigit(v) {
				return false
			}
		}
	}
	return len(s) != 0
}

// mode unchanged
func (c *ctx) convertType(n cc.ExpressionNode, s *buf, from, to cc.Type, fromMode, toMode mode) (r *buf) {
	// defer func() { trc("%v: from %v: %v, to %v: %v %q -> %q", c.pos(n), from, fromMode, to, toMode, s, r) }()
	var b buf
	if from.Kind() == cc.Ptr && to.Kind() == cc.Ptr {
		return s
	}

	if to.Kind() == cc.Void {
		return s
	}

	if cc.IsScalarType(from) && cc.IsScalarType(to) {
		switch {
		case from.Kind() == cc.Ptr && to.Kind() == cc.Ptr:
			return s
		case
			cc.IsFloatingPointType(from) && cc.IsFloatingPointType(to),
			cc.IsIntegerType(from) == cc.IsIntegerType(to) && cc.IsSignedInteger(from) == cc.IsSignedInteger(to) && from.Size() <= to.Size(),
			from.Kind() == cc.Ptr && cc.IsIntegerType(to),
			to.Kind() == cc.Ptr && cc.IsIntegerType(from):

			b.w("(%s(%s))", c.typ(to), s)
		default:
			b.w("(%s%sFrom%s(%s))", c.task.tlsQualifier, c.helper(to), c.helper(from), s)
		}
		return &b
	}

	c.err(errorf("TODO %q %s %s -> %s %s (%v:)", s, from, fromMode, to, toMode, c.pos(n)))
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
		if fromMode == exprUintpr && toMode == exprDefault {
			return s
		}

		if fromMode == exprDefault && toMode == exprUintpr {
			b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n, s)))
			return &b
		}
	}

	if cc.IsIntegerType(to) {
		if toMode == exprDefault {
			b.w("(%s(%s))", c.typ(to), s)
			return &b
		}
	}

	c.err(errorf("TODO %q %s %s, %s -> %s %s, %s", s, from, from.Kind(), fromMode, to, to.Kind(), toMode))
	// trc("%v: TODO %q %s %s, %s -> %s %s, %s", cpos(n), s, from, from.Kind(), fromMode, to, to.Kind(), toMode)
	return s //TODO
}

func (c *ctx) expr0(w writer, n cc.ExpressionNode, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	// trc("%v: %T (%q), %v, %v (%v: %v: %v:)", n.Position(), n, cc.NodeSource(n), t, mode, origin(4), origin(3), origin(2))
	// defer func() {
	// 	trc("%v: %T (%q), %v, %v (RET)", n.Position(), n, cc.NodeSource(n), t, mode)
	// }()
	switch {
	case mode == exprBool:
		mode = exprDefault
	case mode == exprDefault && n.Type().Undecay().Kind() == cc.Array:
		if d := c.declaratorOf(n); d == nil || !d.IsParam() {
			mode = exprUintpr
		}
	}
	if t == nil {
		t = n.Type()
	}
	switch x := n.(type) {
	case *cc.AdditiveExpression:
		return c.additiveExpression(w, x, t, mode)
	case *cc.AndExpression:
		return c.andExpression(w, x, t, mode)
	case *cc.AssignmentExpression:
		return c.assignmentExpression(w, x, t, mode)
	case *cc.CastExpression:
		return c.castExpression(w, x, t, mode)
	case *cc.ConstantExpression:
		return c.expr0(w, x.ConditionalExpression, t, mode)
	case *cc.ConditionalExpression:
		return c.conditionalExpression(w, x, t, mode)
	case *cc.EqualityExpression:
		return c.equalityExpression(w, x, t, mode)
	case *cc.ExclusiveOrExpression:
		return c.exclusiveOrExpression(w, x, t, mode)
	case *cc.ExpressionList:
		return c.expressionList(w, x, t, mode)
	case *cc.InclusiveOrExpression:
		return c.inclusiveOrExpression(w, x, t, mode)
	case *cc.LogicalAndExpression:
		return c.logicalAndExpression(w, x, t, mode)
	case *cc.LogicalOrExpression:
		return c.logicalOrExpression(w, x, t, mode)
	case *cc.MultiplicativeExpression:
		return c.multiplicativeExpression(w, x, t, mode)
	case *cc.PostfixExpression:
		return c.postfixExpression(w, x, t, mode)
	case *cc.PrimaryExpression:
		return c.primaryExpression(w, x, t, mode)
	case *cc.RelationalExpression:
		return c.relationExpression(w, x, t, mode)
	case *cc.ShiftExpression:
		return c.shiftExpression(w, x, t, mode)
	case *cc.UnaryExpression:
		return c.unaryExpression(w, x, t, mode)
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
		switch {
		case n.ShiftExpression.Value() != nil && n.AdditiveExpression.Value() != nil:
			b.w("(%s%s(%s) << %s)", c.task.tlsQualifier, c.helper(n.Type()), c.expr(w, n.ShiftExpression, n.Type(), exprDefault), c.expr(w, n.AdditiveExpression, nil, exprDefault))
		default:
			b.w("(%s << %s)", c.expr(w, n.ShiftExpression, n.Type(), exprDefault), c.expr(w, n.AdditiveExpression, nil, exprDefault))
		}
		rt, rmode = n.Type(), exprDefault
	case cc.ShiftExpressionRsh: // ShiftExpression ">>" AdditiveExpression
		switch {
		case n.ShiftExpression.Value() != nil && n.AdditiveExpression.Value() != nil:
			b.w("(%s%s(%s) >> %s)", c.task.tlsQualifier, c.helper(n.Type()), c.expr(w, n.ShiftExpression, n.Type(), exprDefault), c.expr(w, n.AdditiveExpression, nil, exprDefault))
		default:
			b.w("(%s >> %s)", c.expr(w, n.ShiftExpression, n.Type(), exprDefault), c.expr(w, n.AdditiveExpression, nil, exprDefault))
		}
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
			w.w("var %s %sbool;/**/", v, tag(preserve))
			w.w("%s", ar.vars("\n"))
			w.w("\nif %s = %s; %s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) && (%s))", v, br)
		case al.len() != 0 && ar.len() == 0:
			// Sequence point
			// al;
			// bl && br
			w.w("%s", al.vars(""))
			w.w("%s;", al.bytes())
			b.w("((%s) && (%s))", bl, br)
		case al.len() != 0 && ar.len() != 0:
			c.err(errorf("TODO %v", n.Case))
			// Sequence point
			// al; if v = bl; v { ar };
			// v && br
			v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
			w.w("var %s %sbool;/**/", v, tag(preserve))
			w.w("%s", al.vars("\n"))
			w.w("%s", ar.vars("\n"))
			w.w("\nif %s = %s; %s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) && (%s))", v, br)
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
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
			w.w("var %s %sbool;/**/", v, tag(preserve))
			w.w("%s", ar.vars("\n"))
			w.w("\nif %s = %s; !%s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) || (%s))", v, br)
		case al.len() != 0 && ar.len() == 0:
			// Sequence point
			// al;
			// bl || br
			w.w("%s", al.vars(""))
			w.w("%s;", al.bytes())
			b.w("((%s) || (%s))", bl, br)
		case al.len() != 0 && ar.len() != 0:
			c.err(errorf("TODO %v", n.Case))
			// Sequence point
			// al; if v = bl; !v { ar };
			// v || br
			v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
			w.w("var %s %sbool;/**/", v, tag(preserve))
			w.w("%s", al.vars("\n"))
			w.w("%s", ar.vars("\n"))
			w.w("\nif %s = %s; !%s { %s };", v, bl, v, ar.bytes())
			b.w("((%s) || (%s))", v, br)
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) conditionalExpression(w writer, n *cc.ConditionalExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	switch n.Case {
	case cc.ConditionalExpressionLOr: // LogicalOrExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.ConditionalExpressionCond: // LogicalOrExpression '?' ExpressionList ':' ConditionalExpression
		rt, rmode = n.Type(), exprDefault
		v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
		w.w("var %s %s;/**/", v, c.typ(n.Type()))
		w.w("\nif %s {", c.expr(w, n.LogicalOrExpression, nil, exprBool))
		w.w("%s = %s;", v, c.expr(w, n.ExpressionList, n.Type(), exprDefault))
		w.w("} else {")
		w.w("%s = %s;", v, c.expr(w, n.ConditionalExpression, n.Type(), exprDefault))
		w.w("};")
		b.w("%s", v)
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
	rt, rmode = n.Type(), mode
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
			if sz := x.(*cc.PointerType).Elem().Size(); sz != 1 {
				s = fmt.Sprintf("*%d", sz)
			}
			b.w("(%s + ((%s)%s))", c.expr(w, n.AdditiveExpression, n.Type(), exprDefault), c.expr(w, n.MultiplicativeExpression, n.Type(), exprDefault), s)
		case cc.IsIntegerType(x) && y.Kind() == cc.Ptr:
			s := ""
			if sz := y.(*cc.PointerType).Elem().Size(); sz != 1 {
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
			if v := x.(*cc.PointerType).Elem().Size(); v > 1 {
				b.w("/%d", v)
			}
			b.w(")")
		case x.Kind() == cc.Ptr && cc.IsIntegerType(y):
			s := ""
			if sz := x.(*cc.PointerType).Elem().Size(); sz != 1 {
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
	x = c.expr(w, a, t, exprDefault)
	y = c.expr(w, b, t, exprDefault)
	if !cc.IsIntegerType(t) || a.Value() == cc.Unknown || b.Value() == cc.Unknown {
		return x, y
	}

	var p, q buf
	p.w("(%s%s(%s))", c.task.tlsQualifier, c.helper(t), x)
	q.w("(%s%s(%s))", c.task.tlsQualifier, c.helper(t), y)
	return &p, &q
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

func (c *ctx) unaryExpression(w writer, n *cc.UnaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
out:
	switch n.Case {
	case cc.UnaryExpressionPostfix: // PostfixExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.UnaryExpressionInc: // "++" UnaryExpression
		rt, rmode = n.Type(), mode
		switch ue := n.UnaryExpression.Type(); {
		case ue.Kind() == cc.Ptr && ue.(*cc.PointerType).Elem().Size() != 1:
			sz := ue.(*cc.PointerType).Elem().Size()
			switch mode {
			case exprVoid:
				b.w("%s += %d", c.expr(w, n.UnaryExpression, nil, exprDefault), sz)
			case exprDefault:
				switch d := c.declaratorOf(n.UnaryExpression); {
				case d != nil:
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s += %d;", ds, sz)
					b.w("%s", ds)
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
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s++;", ds)
					b.w("%s", ds)
				default:
					c.err(errorf("TODO")) // 1: bit field
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		}
	case cc.UnaryExpressionDec: // "--" UnaryExpression
		rt, rmode = n.Type(), mode
		switch ue := n.UnaryExpression.Type(); {
		case ue.Kind() == cc.Ptr && ue.(*cc.PointerType).Elem().Size() != 1:
			sz := ue.(*cc.PointerType).Elem().Size()
			switch mode {
			case exprVoid:
				c.err(errorf("TODO"))
			case exprDefault:
				switch d := c.declaratorOf(n.UnaryExpression); {
				case d != nil:
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s -= %d;", ds, sz)
					b.w("%s", ds)
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
					ds := c.expr(w, n.UnaryExpression, nil, exprDefault)
					w.w("%s--;", ds)
					b.w("%s", ds)
				default:
					v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
					v2 := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
					w.w("var %s %s;/**/", v, c.typ(n.UnaryExpression.Type()))
					w.w("\nvar %s %s;/**/", v2, c.typ(n.UnaryExpression.Type().Pointer()))
					w.w("\n%s = %s;", v2, c.expr(w, n.UnaryExpression, n.UnaryExpression.Type().Pointer(), exprUintpr))
					w.w("(*(*%s)(%s))--;", c.typ(n.UnaryExpression.Type()), unsafePointer(v2))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n.UnaryExpression.Type()), unsafePointer(v2))
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		}
	case cc.UnaryExpressionAddrof: // '&' CastExpression
		rt, rmode = n.Type(), exprUintpr
		b.w("%s", c.expr(w, n.CastExpression, nil, exprUintpr))
	case cc.UnaryExpressionDeref: // '*' CastExpression
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
					t = p.Elem()
					if !cc.IsScalarType(t) {
						c.err(errorf("unsupported va_arg type: %v", t.Kind()))
						t = p
					}
					b.w("%sVa%s(&%s)", c.task.tlsQualifier, c.helper(t), c.expr(w, pfe.ArgumentExpressionList.AssignmentExpression, nil, exprDefault))
					break out
				}
			}
		}
		switch n.Type().Undecay().(type) {
		case *cc.FunctionType:
			if mode == exprCall {
				rt, rmode = n.Type(), mode
				b.w("%s", c.expr(w, n.CastExpression, nil, mode))
				break out
			}
		}

		switch mode {
		case exprDefault, exprLvalue:
			rt, rmode = n.Type(), exprDefault
			b.w("(*(*%s)(%s))", c.typ(n.CastExpression.Type().(*cc.PointerType).Elem()), unsafePointer(c.expr(w, n.CastExpression, nil, exprDefault)))
		case exprVoid:
			rt, rmode = n.Type(), mode
			b.w("_ = (*(*%s)(%s))", c.typ(n.CastExpression.Type().(*cc.PointerType).Elem()), unsafePointer(c.expr(w, n.CastExpression, nil, exprDefault)))
		case exprUintpr:
			rt, rmode = n.CastExpression.Type(), mode
			b.w("%s", c.expr(w, n.CastExpression, nil, exprDefault))
		default:
			c.err(errorf("TODO %v", mode))
		}
	case cc.UnaryExpressionPlus: // '+' CastExpression
		rt, rmode = n.Type(), exprDefault
		switch ct := n.CastExpression.Type(); {
		case (cc.IsIntegerType(ct) || ct.Kind() == cc.Ptr) && n.CastExpression.Value() != cc.Unknown:
			b.w("(+%s%s%s(%s))", c.task.tlsQualifier, tag(preserve), c.helper(ct), c.expr(w, n.CastExpression, n.Type(), exprDefault))
		default:
			b.w("(+%s)", c.expr(w, n.CastExpression, n.Type(), exprDefault))
		}
	case cc.UnaryExpressionMinus: // '-' CastExpression
		rt, rmode = n.Type(), exprDefault
		switch ct := n.CastExpression.Type(); {
		case (cc.IsIntegerType(ct) || ct.Kind() == cc.Ptr) && n.CastExpression.Value() != cc.Unknown:
			b.w("(-%s%s%s(%s))", c.task.tlsQualifier, tag(preserve), c.helper(ct), c.expr(w, n.CastExpression, n.Type(), exprDefault))
		default:
			b.w("(-%s)", c.expr(w, n.CastExpression, n.Type(), exprDefault))
		}
	case cc.UnaryExpressionCpl: // '~' CastExpression
		rt, rmode = n.Type(), exprDefault
		switch ct := n.CastExpression.Type(); {
		case n.CastExpression.Value() != cc.Unknown:
			b.w("(^%s%s%s(%s))", c.task.tlsQualifier, tag(preserve), c.helper(ct), c.expr(w, n.CastExpression, n.Type(), exprDefault))
		default:
			b.w("(^%s)", c.expr(w, n.CastExpression, n.Type(), exprDefault))
		}
	case cc.UnaryExpressionNot: // '!' CastExpression
		rt, rmode = n.Type(), exprBool
		b.w("(!%s)", c.expr(w, n.CastExpression, nil, exprBool))
	case cc.UnaryExpressionSizeofExpr: // "sizeof" UnaryExpression
		rt, rmode = n.Type(), exprDefault
		switch n.UnaryExpression.Type().Undecay().Kind() {
		case cc.Array:
			b.w("%s(%s)", c.typ(n.Type()), unsafe("Sizeof", c.expr(w, n.UnaryExpression, nil, exprIndex)))
		default:
			b.w("%s(%s)", c.typ(n.Type()), unsafe("Sizeof", c.expr(w, n.UnaryExpression, nil, exprDefault)))
		}
	case cc.UnaryExpressionSizeofType: // "sizeof" '(' TypeName ')'
		rt, rmode = n.Type(), exprDefault
		switch tn := n.TypeName.Type(); {
		case cc.IsScalarType(tn):
			b.w("%s(%sunsafe.%[2]sSizeof(%s(0)))", c.typ(n.Type()), tag(preserve), c.typ(tn))
		case tn.Kind() == cc.Array && tn.(*cc.ArrayType).IsVLA():
			c.err(errorf("TODO %v", n.Case))
		default:
			b.w("%s(%sunsafe.%[2]sSizeof(%s{}))", c.typ(n.Type()), tag(preserve), c.typ(tn.Undecay()))
		}
	case cc.UnaryExpressionLabelAddr: // "&&" IDENTIFIER
		c.err(errorf("TODO %v", n.Case))
	case cc.UnaryExpressionAlignofExpr: // "_Alignof" UnaryExpression
		rt, rmode = n.Type(), exprDefault
		switch n.UnaryExpression.Type().Undecay().Kind() {
		case cc.Array:
			b.w("%s(%s)", c.typ(n.Type()), unsafe("Alignof", c.expr(w, n.UnaryExpression, nil, exprIndex)))
		default:
			b.w("%s(%s)", c.typ(n.Type()), unsafe("Alignof", c.expr(w, n.UnaryExpression, nil, exprDefault)))
		}
	case cc.UnaryExpressionAlignofType: // "_Alignof" '(' TypeName ')'
		rt, rmode = n.Type(), exprDefault
		switch tn := n.TypeName.Type(); {
		case cc.IsScalarType(tn):
			b.w("%s(%sunsafe.%[2]sAlignof(%s(0)))", c.typ(n.Type()), tag(preserve), c.typ(tn))
		default:
			b.w("%s(%sunsafe.%[2]sAlignof(%s{}))", c.typ(n.Type()), tag(preserve), c.typ(tn.Undecay()))
		}
	case cc.UnaryExpressionImag: // "__imag__" UnaryExpression
		c.err(errorf("TODO %v", n.Case))
	case cc.UnaryExpressionReal: // "__real__" UnaryExpression
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
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
		switch x := n.PostfixExpression.Type().(type) {
		case *cc.PointerType:
			switch mode {
			case exprSelect, exprLvalue, exprDefault, exprIndex:
				switch y := x.Undecay().(type) {
				case *cc.ArrayType:
					rt, rmode = n.Type(), mode
					if d := c.declaratorOf(n.PostfixExpression); d != nil && d.IsParam() {
						var s string
						if v := x.Elem().Size(); v != 1 {
							s = fmt.Sprintf("*%v", v)
						}
						b.w("(*(*%s)(%sunsafe.%[2]sAdd(%[2]sunsafe.%[2]sPointer(%s), (%s)%s)))", c.typ(x.Elem()), tag(preserve), c.expr(w, n.PostfixExpression, nil, exprDefault), c.expr(w, n.ExpressionList, nil, exprDefault), s)
						break
					}

					b.w("%s[%s]", c.expr(w, n.PostfixExpression, nil, exprIndex), c.expr(w, n.ExpressionList, nil, exprDefault))
				case *cc.PointerType:
					rt, rmode = n.Type(), mode
					var s string
					if v := y.Elem().Size(); v != 1 {
						s = fmt.Sprintf("*%v", v)
					}
					b.w("(*(*%s)(%sunsafe.%[2]sAdd(%[2]sunsafe.%[2]sPointer(%s), (%s)%s)))", c.typ(y.Elem()), tag(preserve), c.expr(w, n.PostfixExpression, nil, exprDefault), c.expr(w, n.ExpressionList, nil, exprDefault), s)
				default:
					c.err(errorf("TODO %T", x))
				}
			case exprUintpr:
				rt, rmode = n.Type(), mode
				s := ""
				if sz := x.Elem().Size(); sz != 1 {
					s = fmt.Sprintf("*%v", sz)
				}
				if x.Elem().Kind() == cc.Array {
					if d := c.declaratorOf(n.PostfixExpression); d != nil && d.Type().Kind() == cc.Ptr {
						b.w("((%s)+(%s)%s)", c.expr(w, n.PostfixExpression, nil, exprDefault), c.expr(w, n.ExpressionList, c.pvoid, exprDefault), s)
						break
					}
				}

				b.w("%suintptr(%[1]sunsafe.%[1]sAdd(%[1]sunsafe.%[1]sPointer(%s), (%s%s)))", tag(preserve), c.expr(w, n.PostfixExpression, nil, exprUintpr), c.expr(w, n.ExpressionList, nil, exprDefault), s)
			default:
				//trc("", c.pos(n), mode, cc.NodeSource(n))
				c.err(errorf("TODO %v", mode))
			}
		case *cc.PredefinedType:
			switch {
			case x.VectorSize() < 0:
				switch mode {
				case exprDefault:
					switch y := n.ExpressionList.Type().Undecay().(type) {
					case *cc.PointerType:
						rt, rmode = n.Type(), mode
						var s string
						if v := y.Elem().Size(); v != 1 {
							s = fmt.Sprintf("*%v", v)
						}
						b.w("(*(*%s)(%sunsafe.%[2]sAdd(%[2]sunsafe.%[2]sPointer(%s), (%s)%s)))", c.typ(y.Elem()), tag(preserve), c.expr(w, n.ExpressionList, nil, exprDefault), c.expr(w, n.PostfixExpression, nil, exprDefault), s)
					default:
						c.err(errorf("TODO %T", y))
					}
				default:
					c.err(errorf("TODO %v %v", x, mode))
				}
			default:
				c.err(errorf("TODO %v", x))
			}
		default:
			c.err(errorf("TODO %T", x))
		}
	case cc.PostfixExpressionCall: // PostfixExpression '(' ArgumentExpressionList ')'
		//TODO __builtin_object_size 28_strings.c on darwin/amd64
		switch c.declaratorOf(n.PostfixExpression).Name() {
		case "__builtin_constant_p":
			w.w("_ = %s;", c.expr(w, n.ArgumentExpressionList.AssignmentExpression, nil, exprDefault))
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
			break out
		}

		return c.postfixExpressionCall(w, n)
	case cc.PostfixExpressionSelect: // PostfixExpression '.' IDENTIFIER
		b.n = n
		f := n.Field()
		if _, ok := n.PostfixExpression.Type().(*cc.UnionType); ok && f.Index() != 0 {
			switch mode {
			case exprLvalue, exprDefault, exprSelect:
				rt, rmode = n.Type(), mode
				switch {
				case f.Offset() != 0:
					b.w("(*(*%s)(%sunsafe.%[2]sAdd(%[2]sunsafe.%[2]sPointer(&(%s)), %d)))", c.typ(f.Type()), tag(preserve), c.expr(w, n.PostfixExpression, nil, exprSelect), f.Offset())
				default:
					b.w("(*(*%s)(%s))", c.typ(f.Type()), unsafeAddr(c.expr(w, n.PostfixExpression, nil, exprSelect)))
				}
			case exprUintpr:
				rt, rmode = n.Type(), mode
				switch {
				case f.Offset() != 0:
					b.w("%suintptr(%[1]sunsafe.%s[1]Add(%[1]sunsafe.%[1]sPointer(&(%s)), %d))", tag(preserve), c.pin(n.PostfixExpression, c.expr(w, n.PostfixExpression, nil, exprSelect)), f.Offset())
				default:
					b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n.PostfixExpression, c.expr(w, n.PostfixExpression, nil, exprSelect))))
				}
			case exprIndex:
				switch x := n.Type().Undecay().(type) {
				case *cc.ArrayType:
					rt, rmode = n.Type(), mode
					switch {
					case f.Offset() != 0:
						b.w("((*%s)(%sunsafe.%[2]sAdd(unsafe.%[2]sPointer(&(%s)), %d)))", c.typ(f.Type()), tag(preserve), c.pin(n.PostfixExpression, c.expr(w, n.PostfixExpression, nil, exprSelect)), f.Offset())
					default:
						b.w("((*%s)(%s))", c.typ(f.Type()), unsafeAddr(c.pin(n.PostfixExpression, c.expr(w, n.PostfixExpression, nil, exprSelect))))
					}
				default:
					c.err(errorf("TODO %T", x))
				}
			default:
				c.err(errorf("TODO %v", mode))
			}
			break
		}

		switch mode {
		case exprLvalue, exprDefault, exprSelect, exprIndex:
			rt, rmode = n.Type(), mode
			b.w("(%s.", c.expr(w, n.PostfixExpression, nil, exprSelect))
			switch {
			case f.Parent() != nil:
				c.err(errorf("TODO %v", n.Case))
			default:
				b.w("%s%s)", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
			}
		case exprUintpr:
			rt, rmode = n.Type(), mode
			b.w("%suintptr(%[1]sunsafe.%[1]sPointer(&(%s.", tag(preserve), c.pin(n, c.expr(w, n.PostfixExpression, nil, exprLvalue)))
			switch {
			case f.Parent() != nil:
				c.err(errorf("TODO %v", n.Case))
			default:
				b.w("%s%s)))", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
			}
		default:
			c.err(errorf("TODO %v", mode))
		}
	case cc.PostfixExpressionPSelect: // PostfixExpression "->" IDENTIFIER
		f := n.Field()
		switch mode {
		case exprDefault, exprLvalue, exprIndex:
			rt, rmode = n.Type(), exprDefault
			b.w("((*%s)(%s).", c.typ(n.PostfixExpression.Type().(*cc.PointerType).Elem()), unsafePointer(c.expr(w, n.PostfixExpression, nil, exprDefault)))
			switch {
			case f.Parent() != nil:
				c.err(errorf("TODO %v", n.Case))
			default:
				b.w("%s%s)", tag(field), c.fieldName(n.PostfixExpression.Type(), f))
			}
		case exprUintpr:
			rt, rmode = n.Type(), mode
			b.w("((%s)%s)", c.expr(w, n.PostfixExpression, nil, exprDefault), fldOff(f.Offset()))
		default:
			c.err(errorf("TODO %v", mode))
		}
	case cc.PostfixExpressionInc: // PostfixExpression "++"
		rt, rmode = n.Type(), mode
		switch pe := n.PostfixExpression.Type(); {
		case pe.Kind() == cc.Ptr && pe.(*cc.PointerType).Elem().Size() != 1:
			sz := pe.(*cc.PointerType).Elem().Size()
			switch mode {
			case exprVoid:
				b.w("%s += %d", c.expr(w, n.PostfixExpression, nil, exprDefault), sz)
			case exprDefault:
				v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\n%s = %s;", v, ds)
					w.w("%s += %d;", ds, sz)
					b.w("%s", v)
				default:
					v2 := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\nvar %v %s;/**/", v2, c.typ(n.PostfixExpression.Type().Pointer()))
					w.w("\n%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintpr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s)) += %d;", c.typ(n.PostfixExpression.Type()), unsafePointer(v2), sz)
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
				v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\n%s = %s;", v, ds)
					w.w("%s++;", ds)
					b.w("%s", v)
				default:
					v2 := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\nvar %v %s;/**/", v2, c.typ(n.PostfixExpression.Type().Pointer()))
					w.w("\n%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintpr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s))++;", c.typ(n.PostfixExpression.Type()), unsafePointer(v2))
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		}
	case cc.PostfixExpressionDec: // PostfixExpression "--"
		switch pe := n.PostfixExpression.Type(); {
		case pe.Kind() == cc.Ptr && pe.(*cc.PointerType).Elem().Size() != 1:
			sz := pe.(*cc.PointerType).Elem().Size()
			switch mode {
			case exprVoid:
				b.w("%s = %d", c.expr(w, n.PostfixExpression, nil, exprDefault), sz)
			case exprDefault:
				v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\n%s = %s;/**/", v, ds)
					w.w("%s -= %d;", ds, sz)
					b.w("%s", v)
				default:
					v2 := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\nvar %v %s;/**/", v2, c.typ(n.PostfixExpression.Type().Pointer()))
					w.w("\n%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintpr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s)) -= %d;", c.typ(n.PostfixExpression.Type()), unsafePointer(v2), sz)
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO %v", mode)) // -
			}
		default:
			switch mode {
			case exprVoid:
				rt, rmode = n.Type(), exprVoid
				b.w("%s--", c.expr(w, n.PostfixExpression, nil, exprDefault))
			case exprDefault:
				rt, rmode = n.Type(), exprDefault
				v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
				switch d := c.declaratorOf(n.PostfixExpression); {
				case d != nil:
					ds := c.expr(w, n.PostfixExpression, nil, exprDefault)
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\n%s = %s;", v, ds)
					w.w("%s--;", ds)
					b.w("%s", v)
				default:
					v2 := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
					w.w("var %v %s;/**/", v, c.typ(n.PostfixExpression.Type()))
					w.w("\nvar %v %s;/**/", v2, c.typ(n.PostfixExpression.Type().Pointer()))
					w.w("\n%s = %s;", v2, c.expr(w, n.PostfixExpression, n.PostfixExpression.Type().Pointer(), exprUintpr))
					w.w("%s = (*(*%s)(%s));", v, c.typ(n.PostfixExpression.Type()), unsafePointer(v2))
					w.w("(*(*%s)(%s))--;", c.typ(n.PostfixExpression.Type()), unsafePointer(v2))
					b.w("%s", v)
				}
			default:
				c.err(errorf("TODO")) // -
			}
		}
	case cc.PostfixExpressionComplit: // '(' TypeName ')' '{' InitializerList ',' '}'
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
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
	pt, ok := n.PostfixExpression.Type().(*cc.PointerType)
	if !ok {
		c.err(errorf("TODO %T", n.PostfixExpression.Type()))
		return
	}

	ft, ok := pt.Elem().(*cc.FunctionType)
	if !ok {
		c.err(errorf("TODO %T", pt.Elem()))
		return
	}

	var args []cc.ExpressionNode
	for l := n.ArgumentExpressionList; l != nil; l = l.ArgumentExpressionList {
		args = append(args, l.AssignmentExpression)
	}
	if len(args) < ft.MinArgs() {
		c.err(errorf("%v: too few arguments to function '%s', type '%v' in '%v'", c.pos(n.PostfixExpression), cc.NodeSource(n.PostfixExpression), ft, cc.NodeSource(n)))
		return nil, nil, 0
	}

	if len(args) > ft.MaxArgs() && ft.MaxArgs() >= 0 {
		c.err(errorf("%v: too many arguments to function '%s', type '%v' in '%v'", c.pos(n.PostfixExpression), cc.NodeSource(n.PostfixExpression), ft, cc.NodeSource(n)))
		return nil, nil, 0
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
			mode = exprUintpr
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
		switch mode {
		case exprDefault:
			rt, rmode = n.Type(), exprDefault
			v := fmt.Sprintf("%sv%d", tag(ccgoAutomatic), c.id())
			w.w("var %v %s;/**/", v, c.typ(n.UnaryExpression.Type()))
			w.w("\n%s = %s;", v, c.expr(w, n.AssignmentExpression, n.UnaryExpression.Type(), exprDefault))
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
		var mul, div string
		switch n.Case {
		case cc.AssignmentExpressionAdd: // UnaryExpression "+=" AssignmentExpression
			switch {
			case x.Kind() == cc.Ptr && cc.IsIntegerType(y):
				if sz := x.(*cc.PointerType).Elem().Size(); sz != 1 {
					mul = fmt.Sprintf("*%d", sz)
				}
			case cc.IsIntegerType(x) && y.Kind() == cc.Ptr:
				c.err(errorf("TODO")) // -
			}
		case cc.AssignmentExpressionSub: // UnaryExpression "-=" AssignmentExpression
			switch {
			case x.Kind() == cc.Ptr && cc.IsIntegerType(y):
				if sz := x.(*cc.PointerType).Elem().Size(); sz != 1 {
					mul = fmt.Sprintf("*%d", sz)
				}
			case x.Kind() == cc.Ptr && y.Kind() == cc.Ptr:
				if sz := x.(*cc.PointerType).Elem().Size(); sz != 1 {
					div = fmt.Sprintf("/%d", sz)
				}
			}
		}
		ct := c.usualArithmeticConversions(x, y)
		switch mode {
		case exprVoid:
			switch d := c.declaratorOf(n.UnaryExpression); {
			case d != nil:
				b.w("%s = ", c.expr(w, n.UnaryExpression, nil, exprDefault))
				var b2 buf
				b2.w("(%s %s (%s%s))", c.expr(w, n.UnaryExpression, ct, exprDefault), op, c.expr(w, n.AssignmentExpression, ct, exprDefault), mul)
				b.w("%s", c.convert(n, w, &b2, ct, t, mode, mode))
			default:
				p := fmt.Sprintf("%sp%d", tag(ccgo), c.id())
				ut := n.UnaryExpression.Type()
				w.w("var %s %s;/**/", p, c.typ(ut.Pointer()))
				w.w("\n%s = %s;", p, c.expr(w, n.UnaryExpression, ut.Pointer(), exprUintpr))
				var b2 buf
				p2 := newBufFromtring(fmt.Sprintf("(*(*%s)(%s))", c.typ(ut), unsafePointer(p)))
				b2.w("((%s %s (%s%s))%s)", c.convert(n, w, p2, n.UnaryExpression.Type(), ct, exprDefault, exprDefault), op, c.expr(w, n.AssignmentExpression, ct, exprDefault), mul, div)
				b.w("*(*%s)(%s) = %s", c.typ(ut), unsafePointer(p), c.convert(n, w, &b2, ct, t, mode, mode))
			}
		case exprDefault:
			switch d := c.declaratorOf(n.UnaryExpression); {
			case d != nil:
				w.w("%s = ", c.expr(w, n.UnaryExpression, nil, exprDefault))
				var b2 buf
				b2.w("(%s %s (%s%s))", c.expr(w, n.UnaryExpression, ct, exprDefault), op, c.expr(w, n.AssignmentExpression, ct, exprDefault), mul)
				w.w("%s;", c.convert(n, w, &b2, ct, t, mode, mode))
				b.w("%s", c.expr(w, n.UnaryExpression, nil, exprDefault))
			default:
				c.err(errorf("TODO"))
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

func (c *ctx) primaryExpression(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
out:
	switch n.Case {
	case cc.PrimaryExpressionIdent: // IDENTIFIER
		rt, rmode = n.Type(), mode
		switch x := n.ResolvedTo().(type) {
		case *cc.Declarator:
			b.n = x
			var info *declInfo
			if c.f != nil {
				info = c.f.declInfos.info(x)
			}
			switch {
			case info != nil && info.pinned():
				switch mode {
				case exprLvalue, exprSelect, exprIndex:
					b.w("(*(*%s)(%s))", c.typ(x.Type()), unsafePointer(bpOff(info.bpOff)))
				case exprUintpr:
					b.w("%s", bpOff(info.bpOff))
				case exprDefault:
					switch n.Type().Undecay().(type) {
					case *cc.ArrayType:
						b.w("%s", bpOff(info.bpOff))
					default:
						b.w("(*(*%s)(%s))", c.typ(x.Type()), unsafePointer(bpOff(info.bpOff)))
					}
				default:
					c.err(errorf("TODO %v %v:", mode, n.Position()))
				}
			default:
				switch mode {
				case exprDefault:
					switch x.Type().Kind() {
					case cc.Array:
						p := &buf{n: x}
						p.w("%s%s", c.declaratorTag(x), x.Name())
						b.w("%suintptr(%s)", tag(preserve), unsafeAddr(c.pin(n, p)))
					case cc.Function:
						v := fmt.Sprintf("%sf%d", tag(ccgo), c.id())
						switch {
						case c.f == nil:
							w.w("var %s = %s%s;", v, c.declaratorTag(x), x.Name())
						default:
							w.w("%s := %s%s;", v, c.declaratorTag(x), x.Name())
						}
						b.w("(*(*%suintptr)(%s))", tag(preserve), unsafeAddr(v))
					default:
						b.w("%s%s", c.declaratorTag(x), x.Name())
					}
				case exprLvalue, exprSelect:
					b.w("%s%s", c.declaratorTag(x), x.Name())
				case exprCall:
					switch y := x.Type().(type) {
					case *cc.FunctionType:
						b.w("%s%s", c.declaratorTag(x), x.Name())
					case *cc.PointerType:
						switch z := y.Elem().(type) {
						case *cc.FunctionType:
							rmode = exprUintpr
							b.w("%s%s", c.declaratorTag(x), x.Name())
						default:
							c.err(errorf("TODO %T", z))
						}
					default:
						c.err(errorf("TODO %T", y))
					}
				case exprIndex:
					switch x.Type().Kind() {
					case cc.Array:
						b.w("%s%s", c.declaratorTag(x), x.Name())
					default:
						panic(todo(""))
						c.err(errorf("TODO %v", mode))
					}
				case exprUintpr:
					switch {
					case x.Type().Kind() == cc.Function:
						v := fmt.Sprintf("%sf%d", tag(ccgo), c.id())
						switch {
						case c.f == nil:
							w.w("var %s = %s%s;", v, c.declaratorTag(x), x.Name())
						default:
							w.w("%s := %s%s;", v, c.declaratorTag(x), x.Name())
						}
						b.w("(*(*%suintptr)(%s))", tag(preserve), unsafeAddr(v)) // Free pass from .pin
					default:
						p := &buf{n: x}
						p.w("%s%s", c.declaratorTag(x), x.Name())
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
				b.w("%s%s", tag(enumConst), x.Token.Src())
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
		rt, rmode = n.Type(), exprUntyped
		b.w("%d", n.Value()) //TODO-
	case cc.PrimaryExpressionLChar: // LONGCHARCONST
		c.err(errorf("TODO %v", n.Case))
		b.n = n
	case cc.PrimaryExpressionString: // STRINGLITERAL
		rt, rmode = n.Type(), exprUntyped
		b.n = n
		b.w("%s", n.Value()) //TODO-
	case cc.PrimaryExpressionLString: // LONGSTRINGLITERAL
		c.err(errorf("TODO %v", n.Case))
	case cc.PrimaryExpressionExpr: // '(' ExpressionList ')'
		return c.expr0(w, n.ExpressionList, nil, mode)
	case cc.PrimaryExpressionStmt: // '(' CompoundStatement ')'
		c.err(errorf("TODO %v", n.Case))
	case cc.PrimaryExpressionGeneric: // GenericSelection
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return &b, rt, rmode
}

func (c *ctx) primaryExpressionIntConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	rt, rmode = t, exprDefault
	switch {
	case n.Type().Kind() == t.Kind():
		b.w("(%s(%v))", c.typ(t), n.Value())
	default:
		b.w("(%s%sFrom%s(%v))", c.task.tlsQualifier, c.helper(t), c.helper(n.Type()), n.Value())
	}
	return &b, rt, rmode
}

func (c *ctx) primaryExpressionFloatConst(w writer, n *cc.PrimaryExpression, t cc.Type, mode mode) (r *buf, rt cc.Type, rmode mode) {
	var b buf
	rt, rmode = t, exprDefault
	v := n.Value()
	f := interface{}(v)
	zero := false
	switch x := v.(type) {
	case *cc.LongDoubleValue:
		g := (*big.Float)(x)
		zero = g.Sign() == 0
		f = g
	case cc.Float64Value:
		zero = x == 0
	}
	switch {
	case !zero && n.Type().Kind() == t.Kind():
		b.w("(%s(%v))", c.typ(t), f)
	default:
		b.w("(%s%sFrom%s(%v))", c.task.tlsQualifier, c.helper(t), c.helper(n.Type()), f)
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
