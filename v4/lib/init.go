// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"math/big"
	"sort"

	"modernc.org/cc/v4"
)

type initPatch struct {
	d   *cc.Declarator
	off int64
	b   *buf
}

func (c *ctx) initializerOuter(w writer, n *cc.Initializer, t cc.Type) (r *buf) {
	a := c.initalizerFlatten(n, nil)
	// dumpInitializer(a, "")
	return c.initializer(w, n, a, t, 0, false)
}

func (c *ctx) initalizerFlatten(n *cc.Initializer, a []*cc.Initializer) (r []*cc.Initializer) {
	r = a
	switch n.Case {
	case cc.InitializerExpr: // AssignmentExpression
		return append(r, n)
	case cc.InitializerInitList: // '{' InitializerList ',' '}'
		for l := n.InitializerList; l != nil; l = l.InitializerList {
			r = append(r, c.initalizerFlatten(l.Initializer, nil)...)
		}
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	return r
}

func (c *ctx) initializer(w writer, n cc.Node, a []*cc.Initializer, t cc.Type, off0 int64, arrayElem bool) (r *buf) {
	// trc("==== (init A) typ %s off0 %#0x (%v:)", t, off0, a[0].Position())
	// dumpInitializer(a, "")
	// trc("---- (init Z)")
	if cc.IsScalarType(t) {
		if len(a) == 0 {
			c.err(errorf("TODO"))
			return nil
		}

		if a[0].Offset()-off0 != 0 {
			c.err(errorf("TODO"))
			return nil
		}

		r = c.expr(w, a[0].AssignmentExpression, t, exprDefault)
		if t.Kind() == cc.Ptr && t.(*cc.PointerType).Elem().Kind() == cc.Function && c.initPatch != nil {
			c.initPatch(off0, r)
			var b buf
			b.w("(%suintptr(0))", tag(preserve))
			return &b
		}

		return r
	}

	switch x := t.(type) {
	case *cc.ArrayType:
		if len(a) == 1 && a[0].Type().Kind() == cc.Array && a[0].Value() != cc.Unknown {
			return c.expr(w, a[0].AssignmentExpression, t, exprDefault)
		}

		return c.initializerArray(w, n, a, x, off0)
	case *cc.StructType:
		if len(a) == 1 && a[0].Type().Kind() == cc.Struct {
			return c.expr(w, a[0].AssignmentExpression, t, exprDefault)
		}

		return c.initializerStruct(w, n, a, x, off0)
	case *cc.UnionType:
		if len(a) == 1 && a[0].Type().Kind() == cc.Union {
			return c.expr(w, a[0].AssignmentExpression, t, exprDefault)
		}

		return c.initializerUnion(w, n, a, x, off0, arrayElem)
	default:
		// trc("%v: in type %v, in expr type %v, t %v", a[0].Position(), a[0].Type(), a[0].AssignmentExpression.Type(), t)
		c.err(errorf("TODO %T", x))
		return nil
	}
}

func (c *ctx) isZeroInitializerSlice(s []*cc.Initializer) bool {
	for _, v := range s {
		if !c.isZero(v.AssignmentExpression.Value()) {
			return false
		}
	}

	return true
}

func (c *ctx) isZero(v cc.Value) bool {
	switch x := v.(type) {
	case cc.Int64Value:
		return x == 0
	case cc.UInt64Value:
		return x == 0
	case cc.Float64Value:
		return x == 0
	case *cc.ZeroValue:
		return true
	case cc.Complex128Value:
		return x == 0
	case cc.Complex64Value:
		return x == 0
	case *cc.ComplexLongDoubleValue:
		return c.isZero(x.Re) && c.isZero(x.Im)
	case *cc.LongDoubleValue:
		return !(*big.Float)(x).IsInf() && (*big.Float)(x).Sign() == 0
	default:
		return false
	}
}

func (c *ctx) initializerArray(w writer, n cc.Node, a []*cc.Initializer, t *cc.ArrayType, off0 int64) (r *buf) {
	// trc("==== (array A) %s off0 %#0x", t, off0)
	// dumpInitializer(a, "")
	// trc("---- (array Z)")
	var b buf
	b.w("%s{", c.typ(n, t))
	if c.isZeroInitializerSlice(a) {
		b.w("}")
		return &b
	}

	et := t.Elem()
	esz := et.Size()
	s := sortInitializers(a, func(n int64) int64 { n -= off0; return n - n%esz })
	for _, v := range s {
		off := v[0].Offset() - off0
		off -= off % esz
		b.w("%d: %s, ", off/esz, c.initializer(w, n, v, et, off0+off, true))
	}
	b.w("}")
	return &b
}

func (c *ctx) initializerStruct(w writer, n cc.Node, a []*cc.Initializer, t *cc.StructType, off0 int64) (r *buf) {
	// trc("==== %v: (struct A) %s off0 %#0x", n.Position(), t, off0)
	// dumpInitializer(a, "")
	// trc("---- (struct Z)")
	var b buf
	b.w("%s{", c.initTyp(n, t))
	if c.isZeroInitializerSlice(a) {
		b.w("}")
		return &b
	}

	var flds []*cc.Field
	for i := 0; ; i++ {
		if f := t.FieldByIndex(i); f != nil {
			if f.IsBitfield() {
				c.err(errorf("TODO bitfield"))
				return nil
			}

			if f.Type().Size() <= 0 {
				switch x := f.Type().(type) {
				case *cc.StructType:
					if x.NumFields() != 0 {
						c.err(errorf("TODO %T", x))
						return nil
					}
				case *cc.UnionType:
					if x.NumFields() != 0 {
						c.err(errorf("TODO %T", x))
						return nil
					}
				case *cc.ArrayType:
					if x.Len() > 0 {
						trc("", x.Len())
						c.err(errorf("TODO %T", x))
						return nil
					}
				default:
					c.err(errorf("TODO %T", x))
					return nil
				}
				continue
			}

			flds = append(flds, f)
			// trc("appended: flds[%d] %q %s off %#0x sz %#0x", len(flds)-1, f.Name(), f.Type(), f.Offset(), f.Type().Size())
			continue
		}

		break
	}
	s := sortInitializers(a, func(off int64) int64 {
		off -= off0
		i := sort.Search(len(flds), func(i int) bool {
			return flds[i].Offset() >= off
		})
		if i < len(flds) && flds[i].Offset() == off {
			return off
		}

		return flds[i-1].Offset()
	})
	for _, v := range s {
		first := v[0]
		off := first.Offset() - off0
		// trc("first.Offset() %#0x, off %#0x", first.Offset(), off)
		for off > flds[0].Offset()+flds[0].Type().Size()-1 {
			// trc("skip %q off %#0x", flds[0].Name(), flds[0].Offset())
			flds = flds[1:]
			if len(flds) == 0 {
				panic(todo("", n.Position()))
			}
		}
		f := flds[0]
		// trc("f %q %s off %#0x", f.Name(), f.Type(), f.Offset())
		b.w("%s%s: %s, ", tag(field), c.fieldName(t, f), c.initializer(w, n, v, f.Type(), off0+f.Offset(), false))
	}
	b.w("}")
	return &b
}

func (c *ctx) initializerUnion(w writer, n cc.Node, a []*cc.Initializer, t *cc.UnionType, off0 int64, arrayElem bool) (r *buf) {
	// trc("==== %v: (union A) %s off0 %#0x", n.Position(), t, off0)
	// dumpInitializer(a, "")
	// trc("---- (union Z)")
	var b buf
	if t.NumFields() == 1 {
		b.w("%s{%s%s: %s}", c.typ(n, t), tag(field), c.fieldName(t, t.FieldByIndex(0)), c.initializer(w, n, a, t.FieldByIndex(0).Type(), off0, false))
		return &b
	}

	if c.isZeroInitializerSlice(a) {
		b.w("%s{}", c.typ(n, t))
		return &b
	}

	p := a[0].Parent()
	if assert {
		for _, v := range a[1:] {
			if v.Parent() != p {
				c.err(errorf("TODO"))
				return &b
			}
		}
	}
	if p == nil {
		c.err(errorf("TODO"))
		return &b
	}

	ts := t.Size()
	pt := p.Type()
	if _, ok := pt.(*cc.ArrayType); ok && arrayElem {
		pt = t
	}
	fs := pt.Size()
	switch x := pt.(type) {
	case *cc.ArrayType:
		switch {
		case fs < ts:
			b.w("*(*%s)(%sunsafe.%sPointer(&struct{ f %s; _ [%d]byte}{f: %s}))", c.typ(n, t), tag(importQualifier), tag(preserve), c.typ(n, x), ts-fs, c.initializerArray(w, n, a, x, off0))
		case fs == ts:
			b.w("*(*%s)(%sunsafe.%sPointer(&%s))", c.typ(n, t), tag(importQualifier), tag(preserve), c.initializerArray(w, n, a, x, off0))
		default:
			c.err(errorf("TODO %s %d, ft %s %d", t, t.Size(), x, x.Size()))
		}
	case *cc.StructType:
		switch {
		case fs < ts:
			b.w("*(*%s)(%sunsafe.%sPointer(&struct{ f %s; _ [%d]byte}{f: %s}))", c.typ(n, t), tag(importQualifier), tag(preserve), c.typ(n, x), ts-fs, c.initializerStruct(w, n, a, x, off0))
		case fs == ts:
			b.w("*(*%s)(%sunsafe.%sPointer(&%s))", c.typ(n, t), tag(importQualifier), tag(preserve), c.initializerStruct(w, n, a, x, off0))
		default:
			c.err(errorf("TODO %s %d, ft %s %d", t, t.Size(), x, x.Size()))
		}
	case *cc.UnionType:
		if t.IsCompatible(x) {
			if len(a) != 1 {
				c.err(errorf("TODO %T", x))
				break
			}

			v0 := a[0]
			switch y := v0.Type().(type) {
			case
				*cc.EnumType,
				*cc.PointerType,
				*cc.PredefinedType:

				switch fs := y.Size(); {
				case fs < ts:
					b.w("*(*%s)(%sunsafe.%sPointer(&struct{ f %s; _ [%d]byte}{f: %s}))", c.typ(n, t), tag(importQualifier), tag(preserve), c.typ(n, y), ts-fs, c.expr(w, v0.AssignmentExpression, y, exprDefault))
				case fs == ts:
					b.w("*(*%s)(%sunsafe.%sPointer(&struct{ f %s}{%s}))", c.typ(n, t), tag(importQualifier), tag(preserve), c.typ(n, y), c.expr(w, v0.AssignmentExpression, y, exprDefault))
				default:
					c.err(errorf("TODO %s %d, ft %s %d", t, t.Size(), y, y.Size()))
				}
			default:
				c.err(errorf("TODO %T", y))
			}
			break
		}

		c.err(errorf("TODO %T", x))
	default:
		c.err(errorf("TODO %T %v", x, arrayElem))
	}
	return &b
}

func sortInitializers(a []*cc.Initializer, group func(int64) int64) (r [][]*cc.Initializer) {
	// [0]6.7.8/23: The order in which any side effects occur among the
	// initialization list expressions is unspecified.
	m := map[int64][]*cc.Initializer{}
	for _, v := range a {
		off := group(v.Offset())
		m[off] = append(m[off], v)
	}
	for _, v := range m {
		r = append(r, v)
	}
	sort.Slice(r, func(i, j int) bool { return r[i][0].Offset() < r[j][0].Offset() })
	return r
}

func dumpInitializer(a []*cc.Initializer, pref string) {
	for _, v := range a {
		var t string
		if p := v.Parent(); p != nil {
			switch d := p.Type().Typedef(); {
			case d != nil:
				t = fmt.Sprintf("[%s].", d.Name())
			default:
				t = fmt.Sprintf("[%s].", p.Type())
			}
		}
		switch v.Case {
		case cc.InitializerExpr:
			fmt.Printf("%s %v: off %#05x '%s' %s%s <- %s\n", pref, pos(v.AssignmentExpression), v.Offset(), cc.NodeSource(v.AssignmentExpression), t, v.Type(), v.AssignmentExpression.Type())
		case cc.InitializerInitList:
			s := pref + "· "
			for l := v.InitializerList; l != nil; l = l.InitializerList {
				dumpInitializer([]*cc.Initializer{l.Initializer}, s)
			}
		}
	}
}
