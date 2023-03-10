// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"strings"

	"modernc.org/cc/v4"
	"modernc.org/gc/v2"
	"modernc.org/mathutil"
)

func (c *ctx) typedef(n cc.Node, t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, n, t, false, false, false)
	return b.String()
}

func (c *ctx) helper(n cc.Node, t cc.Type) string {
	var b strings.Builder
	if t.Kind() == cc.Enum {
		t = t.(*cc.EnumType).UnderlyingType()
	}
	if !cc.IsScalarType(t) {
		c.err(errorf("%v: internal error: %s", n.Position(), t))
	}
	c.typ0(&b, n, t, false, false, false)
	s := b.String()
	if len(s) == 0 {
		c.err(errorf("%v: internal error: %s", n.Position(), t))
		return "__undefined__"
	}

	return export(s[len(tag(preserve)):])
}

func (c *ctx) typ(n cc.Node, t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, n, t, true, true, false)
	return b.String()
}

func (c *ctx) typ2(n cc.Node, t cc.Type, useNames bool) string {
	var b strings.Builder
	c.typ0(&b, n, t, useNames, useNames, false)
	return b.String()
}

func (c *ctx) initTyp(n cc.Node, t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, n, t, false, false, false)
	return b.String()
}

func (c *ctx) typ0(b *strings.Builder, n cc.Node, t cc.Type, useTypenames, useTags, isField bool) {
	if t.Kind() != cc.Array && !c.isValidType1(n, t, true) {
		b.WriteString(tag(preserve))
		b.WriteString("int32")
		return
	}

	if c.task.verifyTypes && !t.IsIncomplete() {
		switch t.Kind() {
		case cc.Struct, cc.Union, cc.Array:
			c.verify[t] = struct{}{}
		}
	}

	if tn := t.Typedef(); tn != nil && useTypenames && tn.LexicalScope().Parent == nil {
		fmt.Fprintf(b, "%s%s", tag(typename), tn.Name())
		return
	}

	switch x := t.(type) {
	case *cc.PointerType, *cc.FunctionType:
		b.WriteString(tag(preserve))
		b.WriteString("uintptr")
	case *cc.PredefinedType:
		if t.VectorSize() > 0 {
			c.err(errorf("TODO vector"))
			return
		}

		switch {
		case cc.IsIntegerType(t):
			switch {
			case t.Size() <= 8:
				b.WriteString(tag(preserve))
				if !cc.IsSignedInteger(t) {
					b.WriteByte('u')
				}
				fmt.Fprintf(b, "int%d", 8*t.Size())
			case t.Size() == 16:
				fmt.Fprintf(b, "[2]%suint64", tag(preserve))
			default:
				b.WriteString(tag(preserve))
				b.WriteString("int")
				c.err(errorf("TODO %T %v", x, t))
			}
		case t.Kind() == cc.Void:
			b.WriteString("struct{}")
		case t.Kind() == cc.Float:
			if t.Size() != 4 {
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
			b.WriteString(tag(preserve))
			b.WriteString("float32")
		case t.Kind() == cc.Double:
			if t.Size() != 8 {
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
			b.WriteString(tag(preserve))
			b.WriteString("float64")
		case t.Kind() == cc.LongDouble:
			// if t.Size() != 8 {
			// 	c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			// }
			switch t.Size() {
			case 8:
				b.WriteString(tag(preserve))
				b.WriteString("float64")
			case 16:
				fmt.Fprintf(b, "[2]%suint64", tag(preserve))
			default:
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
		case t.Kind() == cc.ComplexFloat:
			if t.Size() != 8 {
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
			b.WriteString(tag(preserve))
			b.WriteString("complex64")
		case t.Kind() == cc.ComplexDouble:
			if t.Size() != 16 {
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
			b.WriteString(tag(preserve))
			b.WriteString("complex128")
		case t.Kind() == cc.Float128:
			fmt.Fprintf(b, "[2]%suint64", tag(preserve))
		default:
			b.WriteString(tag(preserve))
			b.WriteString("int")
			c.err(errorf("TODO %T %v %v", x, x, x.Kind()))
		}
	case *cc.EnumType:
		nmTag := x.Tag()
		switch nm := nmTag.SrcStr(); {
		case nm != "" && x.LexicalScope().Parent == nil && useTags:
			fmt.Fprintf(b, "%s%s", tag(taggedEum), nm)
		default:
			c.typ0(b, n, x.UnderlyingType(), false, false, false)
		}
	case *cc.StructType:
		nmTag := x.Tag()
		switch nm := nmTag.SrcStr(); {
		case nm != "" && x.LexicalScope().Parent == nil && useTags:
			fmt.Fprintf(b, "%s%s", tag(taggedStruct), nm)
			c.defineTaggedStructs[nm] = x
		default:
			groups := map[int64]struct{}{}
			b.WriteString("struct {")
			b.WriteByte('\n')
			nf := x.NumFields()
			var al int
			for i := 0; i < nf; i++ {
				f := x.FieldByIndex(i)
				if f.IsFlexibleArrayMember() {
					continue
				}

				switch {
				case f.IsBitfield():
					al = mathutil.Max(al, mathutil.Min(f.GroupSize(), c.maxAlign))
				default:
					al = mathutil.Max(al, c.goFieldAlign(x.FieldByIndex(i).Type()))
				}
			}
			if al < x.Align() {
				c.alignPseudoField(b, mathutil.Min(x.Align(), c.maxAlign))
			}
			var off int64
			// trc("===== %s", x)
			for i := 0; i < nf; i++ {
				f := x.FieldByIndex(i)
				// trc("%v: %q, .off %v, .bitoff %v, .ab %v, .vbits %v, fam %v", i, f.Name(), f.Offset(), f.OffsetBits(), f.AccessBytes(), f.ValueBits(), f.IsFlexibleArrayMember())
				if f.InOverlapGroup() {
					continue
				}

				// trc("%q.%d", f.Name(), f.Index())
				// trc("off %v", off)
				switch {
				case f.IsBitfield():
					// trc("BF")
					var gsz int64
					foff := f.Offset()
					// trc("foff %v", foff)
					if _, ok := groups[foff]; !ok {
						groups[foff] = struct{}{}
						gsz = int64(f.GroupSize())
						// trc("gsz %v", gsz)
						off = roundup(off, mathutil.MinInt64(gsz, int64(c.maxAlign)))
						// trc("off %v", off)
						if p := foff - off; p > 0 {
							// trc("pad %v", p)
							b.WriteByte('\n')
							fmt.Fprintf(b, "%s__ccgo_align%d [%d]byte", tag(field), i, p)
							off += p
						}
						fmt.Fprintf(b, "\n%s__ccgo%d uint%d", tag(field), foff, gsz*8)
						off += gsz
					}
					// trc("post bf off %v", off)
				default:
					// trc("FLD")
					ft := f.Type()
					if f.IsFlexibleArrayMember() && ft.Size() <= 0 {
						break
					}

					cAlign := ft.Align()
					goAlign := c.goFieldAlign(ft)
					off0 := off
					off = roundup(off, int64(goAlign))
					// trc("cAlign %v, goAlign %v, off %v", cAlign, goAlign, off)
					switch {
					case off0%int64(goAlign) != 0 && goAlign < ft.Align():
						b.WriteByte('\n')
						fmt.Fprintf(b, "%s__ccgo_align%d [%d]byte", tag(field), i, f.Offset()-off0)
						off += int64(f.Offset() - off)
					case cAlign > goAlign && off%int64(cAlign) != 0:
						b.WriteByte('\n')
						fmt.Fprintf(b, "%s__ccgo_align%d [%d]byte", tag(field), i, cAlign-goAlign)
						off += int64(cAlign - goAlign)
					case off < f.Offset() && off%int64(goAlign) == 0:
						b.WriteByte('\n')
						fmt.Fprintf(b, "%s__ccgo_align%d [%d]byte", tag(field), i, f.Offset()-off)
						off += int64(f.Offset() - off)
					}
					if ft.Size() == 0 && i == x.NumFields()-1 {
						break
					}

					b.WriteByte('\n')
					fmt.Fprintf(b, "%s%s", tag(field), c.fieldName(x, f))
					b.WriteByte(' ')
					c.typ0(b, n, ft, true, true, true)
					off += ft.Size()
					// trc("post fld off %v", off)
				}
			}
			switch {
			case c.maxAlign < x.Align() && off < x.Size():
				b.WriteByte('\n')
				fmt.Fprintf(b, "%s__ccgo_pad%d [%d]byte", tag(field), nf, x.Size()-off)
			default:
				if p := x.Padding(); p != 0 {
					b.WriteByte('\n')
					fmt.Fprintf(b, "%s__ccgo_pad%d [%d]byte", tag(field), nf, p)
				}
			}
			b.WriteString("\n}")
		}
	case *cc.UnionType:
		nmTag := x.Tag()
		switch nm := nmTag.SrcStr(); {
		case nm != "" && x.LexicalScope().Parent == nil && useTags:
			fmt.Fprintf(b, "%s%s", tag(taggedUnion), nm)
			c.defineTaggedUnions[nm] = x
		default:
			fmt.Fprintf(b, "struct {")
			nf := x.NumFields()
			var al int
			for i := 0; i < nf; i++ {
				f := x.FieldByIndex(i)
				if f.IsFlexibleArrayMember() {
					continue
				}

				switch {
				case f.IsBitfield():
					al = mathutil.Max(al, mathutil.Min(f.GroupSize(), c.maxAlign))
				default:
					al = mathutil.Max(al, c.goFieldAlign(x.FieldByIndex(i).Type()))
				}
			}
			if al < x.Align() {
				c.alignPseudoField(b, mathutil.Min(x.Align(), c.maxAlign))
			}
			ff := firstPositiveSizedField(x)
			for i := 0; i < x.NumFields(); i++ {
				f := x.FieldByIndex(i)
				if f == ff || f.Type().Size() == 0 || f.IsBitfield() {
					continue
				}

				b.WriteByte('\n')
				switch nm := f.Name(); {
				case nm == "":
					fmt.Fprintf(b, "%s__ccgo%d", tag(field), i)
				default:
					fmt.Fprintf(b, "%s%s", tag(field), c.fieldName(x, f))
				}
				b.WriteByte(' ')
				b.WriteString("[0]")
				c.typ0(b, n, f.Type(), true, true, true)
			}
			if ff == nil {
				c.err(errorf("TODO"))
				return
			}

			sz1 := ff.Type().Size()
			b.WriteByte('\n')
			switch nm := ff.Name(); {
			case nm == "":
				fmt.Fprintf(b, "%s__ccgo%d", tag(field), ff.Index())
			default:
				fmt.Fprintf(b, "%s%s", tag(field), c.fieldName(x, ff))
			}
			b.WriteByte(' ')
			c.typ0(b, n, ff.Type(), true, true, true)
			if n := t.Size() - sz1; n != 0 {
				fmt.Fprintf(b, "\n%s__ccgo_pad%d [%d]byte", tag(field), nf, t.Size()-sz1)
			}
			b.WriteString("\n}")
		}
	case *cc.ArrayType:
		fmt.Fprintf(b, "[%d]", x.Len())
		c.typ0(b, n, x.Elem(), true, true, true)
	default:
		b.WriteString("int")
		c.err(errorf("TODO %T", x))
		return
	}
}

func (c *ctx) alignPseudoField(b *strings.Builder, align int) {
	var s string
	switch align {
	case 1:
		return
	case 2:
		s = "uint16"
	case 4:
		s = "uint32"
	case 8:
		s = "uint64" //TODO Won't work on [many] 32 archs.
	default:
		c.err(errorf("TODO %d", align))
		return
	}

	b.WriteByte('\n')
	fmt.Fprintf(b, "%s__ccgo_align [0]%s%s", tag(field), tag(preserve), s)
}

func (c *ctx) goFieldAlign(t cc.Type) (r int) {
	for t.Kind() == cc.Array {
		t = t.(*cc.ArrayType).Elem()
	}
	gk := gcKind(t.Kind(), c.ast.ABI)
	if r = c.task.goABI.Types[gk].FieldAlign; r != 0 {
		return r
	}

	return mathutil.Min(c.maxAlign, t.Align())
}

func (c *ctx) isValidParamType(n cc.Node, t cc.Type) (ok bool) {
	t = t.Undecay()
	if x, ok := t.(*cc.ArrayType); ok && x.IsIncomplete() && !x.IsVLA() {
		return true
	}

	return c.isValidType1(n, t, true)
}

func (c *ctx) isValidType(n cc.Node, t cc.Type, report bool) bool {
	switch x := t.Undecay().(type) {
	case *cc.ArrayType:
		if !c.isValidType(n, x.Elem(), report) {
			return false
		}
	case *cc.StructType:
		for i := 0; i < x.NumFields(); i++ {
			f := x.FieldByIndex(i)
			if !f.IsFlexibleArrayMember() && !c.isValidType(n, f.Type(), report) {
				return false
			}
		}
	case *cc.UnionType:
		for i := 0; i < x.NumFields(); i++ {
			f := x.FieldByIndex(i)
			if !f.IsFlexibleArrayMember() && !c.isValidType(n, f.Type(), report) {
				return false
			}
		}
	}
	return c.isValidType1(n, t, report)
}

func (c *ctx) isValidType1(n cc.Node, t cc.Type, report bool) bool {
	//trc("", pos(n), t, t.Attributes() != nil)
	if t == nil || t == cc.Invalid {
		if report {
			c.err(errorf("%v: invalid type", pos(n)))
		}
		return false
	}

	switch attr := t.Attributes(); {
	case t.Align() > 8 || (t.Size() > 0 && int64(t.Align()) > t.Size()):
		if c.task.ignoreUnsupportedAligment {
			break
		}

		if report {
			c.err(errorf("%v: unsupported alignment %d of %s", pos(n), t.Align(), t))
		}
		return false
	case attr != nil && (attr.Aligned() > 8 || (t.Size() > 0 && attr.Aligned() > t.Size())):
		if c.task.ignoreUnsupportedAligment {
			break
		}

		if report {
			c.err(errorf("%v: unsupported alignment %d of %s", pos(n), attr.Aligned(), t))
		}
		return false
	case attr != nil && (attr.VectorSize() > 0):
		if report {
			c.err(errorf("%v: unsupported vector type: %s", pos(n), t))
		}
		return false
	}

	switch x := t.Undecay().(type) {
	case *cc.ArrayType:
		if x.IsVLA() {
			if report {
				c.err(errorf("%v: variable length arrays are not supported", pos(n)))
			}
			return false
		}
	case *cc.FunctionType:
		if !c.isValidType(n, x.Result(), report) {
			return false
		}

		for _, v := range x.Parameters() {
			if !c.isValidParamType(n, v.Type()) {
				if report {
					c.err(errorf("%v: invalid parameter type: %s", pos(n), v.Type()))
				}
				return false
			}
		}
		return true
	}

	if t.IsIncomplete() {
		if report {
			c.err(errorf("%v: incomplete type: %s", pos(n), t))
		}
		return false
	}

	if t.Size() < 0 {
		if report {
			c.err(errorf("%v: invalid type size: %d", pos(n), t.Size()))
		}
		return false
	}

	return true
}

func (c *ctx) unionLiteral(n cc.Node, t *cc.UnionType) string {
	var b strings.Builder
	c.typ0(&b, n, t, true, false, false)
	return b.String()
}

func (c *ctx) structLiteral(n cc.Node, t *cc.StructType) string {
	var b strings.Builder
	c.typ0(&b, n, t, true, false, false)
	return b.String()
}

type fielder interface {
	NumFields() int
	FieldByIndex(int) *cc.Field
}

func (c *ctx) fieldName(t cc.Type, f *cc.Field) string {
	if f.Name() == "" {
		return fmt.Sprintf("__ccgo%d_%d", f.Index(), f.Offset())
	}

	if t != nil {
		if ft := c.registerFields(t); ft != nil {
			return c.fields[ft].dict[f.Name()]
		}
	}

	return f.Name()
}

func (c *ctx) registerFields(t cc.Type) (ft fielder) {
	if p, ok := t.(*cc.PointerType); ok {
		t = p.Elem()
	}
	ft, ok := t.(fielder)
	if !ok {
		c.err(errorf("internal error: %T", t))
		return ft
	}

	if _, ok := c.fields[ft]; ok {
		return ft
	}

	ns := &nameSpace{}
	c.fields[ft] = ns
	for i := 0; ; i++ {
		f := ft.FieldByIndex(i)
		if f == nil {
			break
		}

		nm := f.Name()
		if nm == "" {
			continue
		}

		ns.dict.put(nm, ns.reg.put(nm))
		if _, ok := f.Type().(fielder); ok {
			c.registerFields(f.Type())
		}
	}
	return ft
}

func (c *ctx) defineStructType(w writer, sep string, n cc.Node, t *cc.StructType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if c.pass != 0 || !c.taggedStructs.add(nm) {
			return
		}

		w.w("\n\n%s%stype %s%s = %s;", sep, c.posComment(n), tag(taggedStruct), nm, c.structLiteral(n, t))
	}
	for _, v := range c.structEnums(t) {
		c.defineEnumType(w, "\n", n, v)
	}
}

func (c *ctx) defineUnionType(w writer, sep string, n cc.Node, t *cc.UnionType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedUnions.add(nm) {
			return
		}

		w.w("\n\n%s%stype %s%s = %s;", sep, c.posComment(n), tag(taggedUnion), nm, c.unionLiteral(n, t))
	}
	for _, v := range c.unionEnums(t) {
		c.defineEnumType(w, "\n", n, v)
	}
}

func (c *ctx) structEnums(t *cc.StructType) (r []*cc.EnumType) {
	for i := 0; i < t.NumFields(); i++ {
		switch f := t.FieldByIndex(i); x := f.Type().(type) {
		case *cc.EnumType:
			r = append(r, x)
		}
	}
	return r
}

func (c *ctx) unionEnums(t *cc.UnionType) (r []*cc.EnumType) {
	for i := 0; i < t.NumFields(); i++ {
		switch f := t.FieldByIndex(i); x := f.Type().(type) {
		case *cc.EnumType:
			r = append(r, x)
		}
	}
	return r
}

func (c *ctx) defineEnumType(w writer, sepStr string, n cc.Node, t *cc.EnumType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedEnums.add(nm) {
			return
		}

		w.w("\n\n%s%stype %s%s = %s;", sepStr, c.posComment(n), tag(taggedEum), nm, c.typ(n, t.UnderlyingType()))
	}
	enums := t.Enumerators()
	if len(enums) == 0 {
		return
	}

	if !c.enumerators.add(enums[0].Token.SrcStr()) {
		return
	}

	for _, v := range enums {
		nm := v.Token.SrcStr()
		c.enumerators.add(nm)
		w.w("const %s%s%s%s = %v;", sep(v), c.posComment(v), tag(enumConst), nm, v.Value())
	}
}

func (c *ctx) defineType(w writer, sep string, n cc.Node, t cc.Type) {
	c.defineType0(w, sep, n, t)
}

func (c *ctx) defineType0(w writer, sep string, n cc.Node, t cc.Type) {
	if !c.isValidType1(n, t, false) {
		return
	}

	switch x := t.(type) {
	case *cc.EnumType:
		c.defineEnumType(w, sep, n, x)
	case *cc.StructType:
		c.defineStructType(w, sep, n, x)
	case *cc.UnionType:
		c.defineUnionType(w, sep, n, x)
	case *cc.PointerType:
		c.defineType0(w, sep, n, x.Elem())
	case *cc.ArrayType:
		c.defineType0(w, sep, n, x.Elem())
	}
}

func typeID(in map[string]gc.Node, out map[string]string, typ gc.Node) (r string, err error) {
	var b strings.Builder
	if err = typeID0(&b, in, out, typ, map[string]struct{}{}); err != nil {
		return "", err
	}

	// trc("`%s` -> type ID: `%s`", typ.Source(false), b.String())
	return b.String(), nil
}

func typeID0(b *strings.Builder, in map[string]gc.Node, out map[string]string, typ gc.Node, m map[string]struct{}) (err error) {
	switch x := typ.(type) {
	case *gc.StructTypeNode:
		b.WriteString("struct{")
		for _, f := range x.FieldDecls {
			switch y := f.(type) {
			case *gc.FieldDecl:
				ft, err := typeID(in, out, y.Type)
				if err != nil {
					return err
				}

				for _, nm := range y.IdentifierList {
					fmt.Fprintf(b, "%s %s;", nm.Ident.Src(), ft)
				}
			default:
				panic(todo("%T", y))
			}
		}
		b.WriteByte('}')
	case *gc.ArrayTypeNode:
		fmt.Fprintf(b, "[%s]", x.ArrayLength.Source(true))
		if err = typeID0(b, in, out, x.ElementType, m); err != nil {
			return err
		}
	case *gc.TypeNameNode:
		if x.TypeArgs != nil || x.Name.PackageName.IsValid() {
			panic(todo("%T %s", x, x.Source(true)))
		}

		nm := x.Name.Ident.Src()
		switch symKind(nm) {
		case -1, preserve:
			b.WriteString(nm)
		case typename, taggedStruct, taggedUnion, taggedEum:
			if id, ok := out[nm]; ok {
				b.WriteString(id)
				break
			}

			t2, ok := in[nm]
			if !ok {
				return errorf("undefined type %s", nm)
			}

			if _, ok := m[nm]; ok {
				return errorf("invalid recursive type %s", nm)
			}

			m[nm] = struct{}{}
			id, err := typeID(in, out, t2)
			if err != nil {
				return err
			}

			out[nm] = id
			b.WriteString(id)
		default:
			panic(todo("%T %s", x, x.Source(true)))
		}
	default:
		panic(todo("%T %s", x, x.Source(false)))
	}
	return nil
}

func firstPositiveSizedField(n *cc.UnionType) *cc.Field {
	for i := 0; i < n.NumFields(); i++ {
		if f := n.FieldByIndex(i); f.Type().Size() > 0 {
			return f
		}
	}
	return nil
}

func isEmpty(t cc.Type) bool {
	switch x := t.(type) {
	case *cc.StructType:
		return x.NumFields() == 0
	case *cc.UnionType:
		return x.NumFields() == 0
	}

	return false
}
