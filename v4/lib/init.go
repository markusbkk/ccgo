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
	// p := n.Position()
	// if len(a) != 0 {
	// 	p = a[0].Position()
	// }
	// trc("==== (init A) typ %s off0 %#0x (%v:) (from %v: %v: %v:)", t, off0, p, origin(4), origin(3), origin(2))
	// dumpInitializer(a, "")
	// defer trc("---- (init Z) typ %s off0 %#0x (%v:)", t, off0, p)
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
	// trc("==== (array A, size %v) %s off0 %#0x (%v:)", t.Size(), t, off0, pos(n))
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
		b.w("\n%d: %s, ", off/esz, c.initializer(w, n, v, et, off0+off, true))
	}
	b.w("}")
	return &b
}

func (c *ctx) initializerStruct(w writer, n cc.Node, a []*cc.Initializer, t *cc.StructType, off0 int64) (r *buf) {
	// trc("==== %v: (struct A, size %v) %s off0 %#0x", n.Position(), t.Size(), t, off0)
	// dumpInitializer(a, "")
	// defer trc("---- %v: (struct Z, size %v) %s off0 %#0x", n.Position(), t.Size(), t, off0)
	var b buf
	b.w("%s{", c.initTyp(n, t))
	if c.isZeroInitializerSlice(a) {
		b.w("}")
		return &b
	}

	var flds []*cc.Field
	for i := 0; ; i++ {
		if f := t.FieldByIndex(i); f != nil {
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
						c.err(errorf("TODO %T", x))
						return nil
					}
				default:
					c.err(errorf("TODO %T", x))
					return nil
				}
				continue
			}

			if f.IsBitfield() && f.ValueBits() == 0 {
				continue
			}

			flds = append(flds, f)
			// trc("appended: flds[%d] %q %s off %#0x ogo %#0x sz %#0x", len(flds)-1, f.Name(), f.Type(), f.Offset(), f.OuterGroupOffset(), f.Type().Size())
			continue
		}

		break
	}
	s := sortInitializers(a, func(off int64) int64 {
		off -= off0
		i := sort.Search(len(flds), func(i int) bool {
			return flds[i].OuterGroupOffset() >= off
		})
		if i < len(flds) && flds[i].OuterGroupOffset() == off {
			return off
		}

		return flds[i-1].OuterGroupOffset()
	})
	// trc("==== initializers (A)")
	// for i, v := range s {
	// 	for j, w := range v {
	// 		if w.Field() == nil {
	// 			trc("%d.%d: %q off %v, %s", i, j, "", w.Offset(), cc.NodeSource(w.AssignmentExpression))
	// 			continue
	// 		}

	// 		trc("%d.%d: %q off %v, ogo %v, %s", i, j, w.Field().Name(), w.Field().Offset(), w.Field().OuterGroupOffset(), cc.NodeSource(w.AssignmentExpression))
	// 	}
	// }
	// trc("==== initializers (Z)")
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
		if f.IsBitfield() {
			// trc("==== %v: TODO bitfield", cpos(n))
			// trc("%q %s off %#0x ogo %#0x sz %#0x", f.Name(), f.Type(), f.Offset(), f.OuterGroupOffset(), f.Type().Size())
			// for i, v := range v {
			// 	trc("%d: %q %s", i, v.Field().Name(), cc.NodeSource(v.AssignmentExpression))
			// }
			// trc("----")
			for len(flds) != 0 && flds[0].OuterGroupOffset() == f.OuterGroupOffset() {
				flds = flds[1:]
			}
			b.w("%s__ccgo%d: ", tag(field), f.OuterGroupOffset())
			sort.Slice(v, func(i, j int) bool {
				a, b := v[i].Field(), v[j].Field()
				return a.Offset()*8+int64(a.OffsetBits()) < b.Offset()*8+int64(b.OffsetBits())
			})
			ogo := f.OuterGroupOffset()
			gsz := 8 * (int64(f.GroupSize()) + f.Offset() - ogo)
			for i, in := range v {
				if i != 0 {
					b.w("|")
				}
				f = in.Field()
				sh := f.OffsetBits() + 8*int(f.Offset()-ogo)
				b.w("(((%suint%d(%s))&%#0x)<<%d)", tag(preserve), gsz, c.expr(w, in.AssignmentExpression, nil, exprDefault), uint(1)<<f.ValueBits()-1, sh)
			}
			b.w(", ")
			continue
		}

		for isEmpty(v[0].Type()) {
			v = v[1:]
		}
		// trc("f %q %s off %#0x v[0].Type() %v", f.Name(), f.Type(), f.Offset(), v[0].Type())
		flds = flds[1:]
		b.w("%s%s: %s, ", tag(field), c.fieldName(t, f), c.initializer(w, n, v, f.Type(), off0+f.Offset(), false))
	}
	b.w("}")
	return &b
}

func (c *ctx) initializerUnion(w writer, n cc.Node, a []*cc.Initializer, t *cc.UnionType, off0 int64, arrayElem bool) (r *buf) {
	// trc("==== %v: (union A, size %v) %s off0 %#0x, arrayElem %v", n.Position(), t.Size(), t, off0, arrayElem)
	// dumpInitializer(a, "")
	// trc("---- (union Z)")
	var b buf
	if c.isZeroInitializerSlice(a) {
		b.w("%s{}", c.typ(n, t))
		return &b
	}

	if t.NumFields() == 1 {
		b.w("%s{%s%s: %s}", c.typ(n, t), tag(field), c.fieldName(t, t.FieldByIndex(0)), c.initializer(w, n, a, t.FieldByIndex(0).Type(), off0, false))
		return &b
	}

	b.w("*(*%s)(%sunsafe.%sPointer(&struct{ ", c.typ(n, t), tag(importQualifier), tag(preserve))
	switch len(a) {
	case 1:
		b.w("%s", c.initializerUnionOne(w, n, a, t, off0))
	default:
		b.w("%s", c.initializerUnionMany(w, n, a, t, off0, arrayElem))
	}
	b.w("))")
	return &b
}

func (c *ctx) initializerUnionOne(w writer, n cc.Node, a []*cc.Initializer, t *cc.UnionType, off0 int64) (r *buf) {
	var b buf
	in := a[0]
	pre := in.Offset() - off0
	if pre != 0 {
		b.w("%s_ [%d]byte;", tag(preserve), pre)
	}
	b.w("%sf ", tag(preserve))
	f := in.Field()
	switch {
	case f != nil && f.IsBitfield():
		b.w("%suint%d", tag(preserve), f.AccessBytes()*8)
	default:
		b.w("%s ", c.typ(n, in.Type()))
	}
	if post := t.Size() - (pre + in.Type().Size()); post != 0 {
		b.w("; %s_ [%d]byte", tag(preserve), post)
	}
	b.w("}{%sf: ", tag(preserve))
	switch f := in.Field(); {
	case f != nil && f.IsBitfield():
		b.w("(((%suint%d(%s))&%#0x)<<%d)", tag(preserve), f.AccessBytes()*8, c.expr(w, in.AssignmentExpression, nil, exprDefault), uint(1)<<f.ValueBits()-1, f.OffsetBits())
	default:
		switch x := in.Type().(type) {
		case *cc.PredefinedType, *cc.PointerType:
			b.w("%s", c.expr(w, in.AssignmentExpression, in.Type(), exprDefault))
		default:
			c.err(errorf("TODO %T", x))
		}
	}
	b.w("}")
	return &b
}

func (c *ctx) initializerUnionMany(w writer, n cc.Node, a []*cc.Initializer, t *cc.UnionType, off0 int64, arrayElem bool) (r *buf) {
	var b buf
	path, x := c.initlializerLCA(a)
	lca := path[x]
	ft := lca.Type()
	fOff := lca.Offset()
out:
	switch {
	case ft == t:
		f := lca.InitializerList.UnionField()
		ft = f.Type()
		fOff = off0 + f.Offset()
	case arrayElem:
		// trc("arrayElem, x %v", x)
		for _, v := range path {
			if v.InitializerList != nil && v.InitializerList.UnionField() != nil {
				if f := v.InitializerList.UnionField(); f == t.FieldByIndex(f.Index()) {
					ft = f.Type()
					fOff = off0 + f.Offset()
					break out
				}
			}
		}

		// for i, v := range path {
		// 	var s string
		// 	if v.InitializerList != nil && v.InitializerList.UnionField() != nil {
		// 		f := v.InitializerList.UnionField()
		// 		s = fmt.Sprintf("%q %v", f.Name(), f.Type())
		// 	}
		// 	trc("%d/%d: %v: %s", i, len(path), v.Position(), s)
		// }
		c.err(errorf("TODO"))
		return &b
	}
	pre := fOff - off0
	if pre != 0 {
		b.w("%s_ [%d]byte;", tag(preserve), pre)
	}
	b.w("%sf ", tag(preserve))
	b.w("%s ", c.typ(n, ft))
	if post := t.Size() - (pre + ft.Size()); post != 0 {
		b.w("; %s_ [%d]byte", tag(preserve), post)
	}
	b.w("}{%sf: ", tag(preserve))
	switch x := ft.(type) {
	case *cc.ArrayType:
		b.w("%s", c.initializerArray(w, n, a, x, off0))
	case *cc.StructType:
		b.w("%s", c.initializerStruct(w, n, a, x, off0))
	case *cc.UnionType:
		b.w("%s", c.initializerUnion(w, n, a, x, off0, false))
	default:
		c.err(errorf("TODO %T", x))
	}
	b.w("}")
	return &b
}

// https://en.wikipedia.org/wiki/Lowest_common_ancestor
func (c ctx) initlializerLCA(a []*cc.Initializer) (r []*cc.Initializer, ri int) {
	if len(a) < 2 {
		panic(todo("internal error"))
	}

	nodes := map[*cc.Initializer]struct{}{}
	var path []*cc.Initializer
	for p := a[0].Parent(); p != nil; p = p.Parent() {
		path = append(path, p)
		r = append(r, p)
		nodes[p] = struct{}{}
	}
	for _, v := range a[1:] {
		for p := v.Parent(); p != nil; p = p.Parent() {
			if _, ok := nodes[p]; ok {
				for path[0] != p {
					delete(nodes, p)
					path = path[1:]
					ri++
				}
				break
			}
		}
	}
	return r, ri
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
		for p := v.Parent(); p != nil; p = p.Parent() {
			switch d := p.Type().Typedef(); {
			case d != nil:
				t = fmt.Sprintf("[%s].", d.Name()) + t
			default:
				switch x, ok := p.Type().(interface{ Tag() cc.Token }); {
				case ok:
					tag := x.Tag()
					t = fmt.Sprintf("[%s].", tag.SrcStr()) + t
				default:
					t = fmt.Sprintf("[%s].", p.Type()) + t
				}
			}
		}
		var fs string
		if f := v.Field(); f != nil {
			var ps string
			for p := f.Parent(); p != nil; p = p.Parent() {
				ps = ps + fmt.Sprintf("{%q %v}", p.Name(), p.Type())
			}
			fs = fmt.Sprintf(
				" %s(field %q, IsBitfield %v, Offset %v, OffsetBits %v, OuterGroupOffset %v, InOverlapGroup %v, Mask %#0x, ValueBits %v)",
				ps, f.Name(), f.IsBitfield(), f.Offset(), f.OffsetBits(), f.OuterGroupOffset(), f.InOverlapGroup(), f.Mask(), f.ValueBits(),
			)
		}
		switch v.Case {
		case cc.InitializerExpr:
			fmt.Printf("%s %v: off %#05x '%s' %s type %q <- %s%s\n", pref, pos(v.AssignmentExpression), v.Offset(), cc.NodeSource(v.AssignmentExpression), t, v.Type(), v.AssignmentExpression.Type(), fs)
		case cc.InitializerInitList:
			s := pref + "· " + fs
			for l := v.InitializerList; l != nil; l = l.InitializerList {
				dumpInitializer([]*cc.Initializer{l.Initializer}, s)
			}
		}
	}
}
