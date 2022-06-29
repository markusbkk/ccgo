// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"sort"

	"modernc.org/cc/v4"
)

func (c *ctx) initializerOuter(w writer, n *cc.Initializer, t cc.Type) (r *buf) {
	return c.initializer(w, n, c.initalizerFlatten(n, nil), t, 0)
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

func (c *ctx) initializer(w writer, n cc.Node, a []*cc.Initializer, t cc.Type, off0 int64) (r *buf) {
	if len(a) == 0 {
		c.err(errorf("TODO"))
		return nil
	}

	// trc("==== (init A) typ %s off0 %#0x (%v:)", t, off0, a[0].Position())
	// dumpInitializer(a, "")
	// trc("---- (init Z)")
	if cc.IsScalarType(t) {
		if len(a) != 1 {
			// trc("%v: FAIL", a[0].Position())
			c.err(errorf("TODO scalar %s, len(initializers) %v", t, len(a)))
			return nil
		}

		if a[0].Offset()-off0 != 0 {
			c.err(errorf("TODO"))
			return nil
		}

		return c.expr(w, a[0].AssignmentExpression, t, exprDefault)
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

		return c.initializerUnion(w, n, a, x, off0)
	default:
		// trc("%v: in type %v, in expr type %v, t %v", a[0].Position(), a[0].Type(), a[0].AssignmentExpression.Type(), t)
		c.err(errorf("TODO %T", x))
		return nil
	}
}

func (c *ctx) initializerUnion(w writer, n cc.Node, a []*cc.Initializer, t *cc.UnionType, off0 int64) (r *buf) {
	// trc("==== %v: (struct A) %s off0 %#0x", n.Position(), t, off0)
	// dumpInitializer(a, "")
	// trc("---- (struct Z)")
	var b buf
	b.w("%s{", c.typ(n, t))
	f := t.FieldByIndex(0)
	flds := []*cc.Field{f}
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
		b.w("%s%s: %s, ", tag(field), c.fieldName(t, f), c.initializer(w, n, v, f.Type(), off0+f.Offset()))
	}
	b.w("}")
	return &b
}

func (c *ctx) initializerStruct(w writer, n cc.Node, a []*cc.Initializer, t *cc.StructType, off0 int64) (r *buf) {
	// trc("==== %v: (struct A) %s off0 %#0x", n.Position(), t, off0)
	// dumpInitializer(a, "")
	// trc("---- (struct Z)")
	var b buf
	b.w("%s{", c.typ(n, t))
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
		b.w("%s%s: %s, ", tag(field), c.fieldName(t, f), c.initializer(w, n, v, f.Type(), off0+f.Offset()))
	}
	b.w("}")
	return &b
}

func (c *ctx) initializerArray(w writer, n cc.Node, a []*cc.Initializer, t *cc.ArrayType, off0 int64) (r *buf) {
	// trc("==== (array A) %s off0 %#0x", t, off0)
	// dumpInitializer(a, "")
	// trc("---- (array Z)")
	var b buf
	b.w("%s{", c.typ(n, t))
	et := t.Elem()
	esz := et.Size()
	s := sortInitializers(a, func(n int64) int64 { n -= off0; return n - n%esz })
	for _, v := range s {
		off := v[0].Offset() - off0
		off -= off % esz
		b.w("%d: %s, ", off/esz, c.initializer(w, n, v, et, off0+off))
	}
	b.w("}")
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
		switch v.Case {
		case cc.InitializerExpr:
			fmt.Printf("%s %v: off %#05x '%s' %s\n", pref, pos(v.AssignmentExpression), v.Offset(), cc.NodeSource(v.AssignmentExpression), v.AssignmentExpression.Type())
		case cc.InitializerInitList:
			s := pref + "· "
			for l := v.InitializerList; l != nil; l = l.InitializerList {
				dumpInitializer([]*cc.Initializer{l.Initializer}, s)
			}
		}
	}
}
