// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"strings"

	"modernc.org/cc/v4"
	"modernc.org/gc/v2"
)

func (c *ctx) typedef(t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, t, false, false, false)
	return b.String()
}

func (c *ctx) helper(t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, t, false, false, false)
	return c.export(b.String())
}

func (c *ctx) typ(t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, t, true, true, false)
	return b.String()
}

func (c *ctx) typeSuffix(t cc.Type) string {
	if !cc.IsIntegerType(t) {
		c.err(errorf("TODO"))
	}
	return c.export(c.typ(t))
}

func (c *ctx) typ0(b *strings.Builder, t cc.Type, useTypename, useStructUnionTag, isField bool) {
	if tn := t.Typedef(); tn != nil && useTypename && tn.LexicalScope().Parent == nil {
		fmt.Fprintf(b, "%s%s", tag(typename), tn.Name())
		return
	}

	switch x := t.(type) {
	case *cc.PointerType, *cc.FunctionType:
		b.WriteString("uintptr")
	case *cc.PredefinedType:
		if t.VectorSize() > 0 {
			c.err(errorf("TODO"))
		}
		switch {
		case cc.IsIntegerType(t):
			switch {
			case t.Size() <= 8:
				if !cc.IsSignedInteger(t) {
					b.WriteByte('u')
				}
				fmt.Fprintf(b, "int%d", 8*t.Size())
			case t.Size() == 16:
				b.WriteString("[2]uint64")
			default:
				b.WriteString("int")
				c.err(errorf("TODO %T %v", x, t))
			}
		case t.Kind() == cc.Void:
			b.WriteString("struct{}")
		case t.Kind() == cc.Float:
			if t.Size() != 4 {
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
			b.WriteString("float32")
		case t.Kind() == cc.Double:
			if t.Size() != 8 {
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
			b.WriteString("float64")
		case t.Kind() == cc.LongDouble:
			// if t.Size() != 8 {
			// 	c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			// }
			switch t.Size() {
			case 8:
				b.WriteString("float64")
			case 16:
				switch {
				case isField:
					b.WriteString("float128")
				default:
					b.WriteString("float64")
				}
			default:
				c.err(errorf("C %v of unexpected size %d", x.Kind(), t.Size()))
			}
		default:
			b.WriteString("int")
			c.err(errorf("TODO %T %v %v", x, x, x.Kind()))
		}
	case *cc.EnumType:
		nmTag := x.Tag()
		switch nm := nmTag.SrcStr(); {
		case nm != "" && x.LexicalScope().Parent == nil:
			fmt.Fprintf(b, "%s%s", tag(taggedEum), nm)
		default:
			c.typ0(b, x.UnderlyingType(), false, false, false)
		}
	case *cc.StructType:
		nmTag := x.Tag()
		switch nm := nmTag.SrcStr(); {
		case nm != "" && x.LexicalScope().Parent == nil && useStructUnionTag:
			fmt.Fprintf(b, "%s%s", tag(taggedStruct), nm)
		default:
			b.WriteString("struct {")
			for i := 0; ; i++ {
				f := x.FieldByIndex(i)
				if f == nil {
					break
				}

				b.WriteByte('\n')
				switch nm := f.Name(); {
				case nm == "":
					fmt.Fprintf(b, "%s__ccgo%d", tag(field), i)
				default:
					fmt.Fprintf(b, "%s%s", tag(field), c.fieldName(x, f))
				}
				b.WriteByte(' ')
				c.typ0(b, f.Type(), true, true, true)
			}
			b.WriteString("\n}")
		}
	case *cc.UnionType:
		nmTag := x.Tag()
		switch nm := nmTag.SrcStr(); {
		case nm != "" && x.LexicalScope().Parent == nil && useStructUnionTag:
			fmt.Fprintf(b, "%s%s", tag(taggedUnion), nm)
		default:
			fmt.Fprintf(b, "struct {")
			var sz1 int64
			for i := 0; ; i++ {
				f := x.FieldByIndex(i)
				if f == nil {
					break
				}

				if i == 0 {
					sz1 = f.Type().Size()
				}
				b.WriteByte('\n')
				switch nm := f.Name(); {
				case nm == "":
					c.err(errorf("TODO"))
				default:
					fmt.Fprintf(b, "%s%s", tag(field), c.fieldName(x, f))
				}
				b.WriteByte(' ')
				if i != 0 {
					b.WriteString("[0]")
				}
				c.typ0(b, f.Type(), true, true, true)
			}
			if n := t.Size() - sz1; n != 0 {
				fmt.Fprintf(b, "\n%s__ccgo [%d]byte", tag(field), t.Size()-sz1)
			}
			b.WriteString("\n}")
		}
	case *cc.ArrayType:
		fmt.Fprintf(b, "[%d]", x.Len())
		c.typ0(b, x.Elem(), true, true, true)
	default:
		b.WriteString("int")
		c.err(errorf("TODO %T", x))
	}
}

func (c *ctx) unionLiteral(t *cc.UnionType) string {
	var b strings.Builder
	c.typ0(&b, t, true, false, false)
	return b.String()
}

func (c *ctx) structLiteral(t *cc.StructType) string {
	var b strings.Builder
	c.typ0(&b, t, true, false, false)
	return b.String()
}

func (c *ctx) defineUnion(w writer, sep string, n cc.Node, t *cc.UnionType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedUnions.add(nm) {
			return
		}

		w.w("%s%stype %s%s = %s;", sep, c.posComment(n), tag(taggedUnion), nm, c.unionLiteral(t))
	}
}

type fielder interface {
	FieldByIndex(int) *cc.Field
}

func (c *ctx) fieldName(t cc.Type, f *cc.Field) string {
	if ft := c.registerFields(t); ft != nil {
		return c.fields[ft].dict[f.Name()]
	}

	return tag(field) + f.Name()
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

func (c *ctx) defineStruct(w writer, sep string, n cc.Node, t *cc.StructType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedStructs.add(nm) {
			return
		}

		w.w("%s%stype %s%s = %s;", sep, c.posComment(n), tag(taggedStruct), nm, c.structLiteral(t))
	}
}

func (c *ctx) defineEnum(w writer, sepStr string, n cc.Node, t *cc.EnumType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedEnums.add(nm) {
			return
		}

		w.w("%s%stype %s%s = %s;", sepStr, c.posComment(n), tag(taggedEum), nm, c.typ(t.UnderlyingType()))
	}
	enums := t.Enumerators()
	if len(enums) == 0 {
		return
	}

	if !c.enumerators.add(enums[0].Token.SrcStr()) {
		return
	}

	w.w("\n\nconst (")
	for _, v := range enums {
		nm := v.Token.SrcStr()
		c.enumerators.add(nm)
		w.w("%s%s%s%s = %v;", sep(v), c.posComment(v), tag(enumConst), nm, v.Value())
	}
	w.w("\n)\n")
}

func (c *ctx) defineEnumStructUnion(w writer, sep string, n cc.Node, t cc.Type) {
	switch x := t.(type) {
	case *cc.EnumType:
		c.defineEnum(w, sep, n, x)
	case *cc.StructType:
		c.defineStruct(w, sep, n, x)
	case *cc.UnionType:
		c.defineUnion(w, sep, n, x)
	}
}

func typeID(in map[string]gc.Node, out map[string]string, typ gc.Node) (r string, err error) {
	var b strings.Builder
	if err = typeID0(&b, in, out, typ, map[string]struct{}{}); err != nil {
		return "", err
	}

	return "T" + b.String(), nil
}

func typeID0(b *strings.Builder, in map[string]gc.Node, out map[string]string, typ gc.Node, m map[string]struct{}) (err error) {
	switch x := typ.(type) {
	case *gc.StructType:
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
	case *gc.ArrayType:
		fmt.Fprintf(b, "[%s]", x.ArrayLength.Source(true))
		if err = typeID0(b, in, out, x.ElementType, m); err != nil {
			return err
		}
	case *gc.TypeName:
		if x.TypeArgs != nil || x.Name.PackageName.IsValid() {
			panic(todo("%T %s", x, x.Source(true)))
		}

		nm := x.Name.Ident.Src()
		switch symKind(nm) {
		case -1:
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
