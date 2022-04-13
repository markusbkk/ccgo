// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"modernc.org/cc/v4"
)

func (c *ctx) typedef(t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, t, false, false)
	return b.String()
}

func (c *ctx) typ(t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, t, true, true)
	return b.String()
}

func (c *ctx) typ0(b *strings.Builder, t cc.Type, useTypename, useStructUnionTag bool) {
	if tn := t.Typedef(); tn != nil && useTypename && tn.LexicalScope().Parent == nil {
		fmt.Fprintf(b, "%s%s", tag(typename), tn.Name())
		return
	}

	switch x := t.(type) {
	case *cc.PointerType:
		b.WriteString("uintptr")
	case *cc.PredefinedType:
		if t.VectorSize() > 0 {
			c.err(errorf("TODO"))
		}
		switch {
		case cc.IsIntegerType(t):
			if t.Size() <= 8 {
				if !cc.IsSignedInteger(t) {
					b.WriteByte('u')
				}
				fmt.Fprintf(b, "int%d", 8*t.Size())
				break
			}

			b.WriteString("int")
			c.err(errorf("TODO %T", x))
		case t.Kind() == cc.Void:
			b.WriteString("struct{}")
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
			c.typ0(b, x.UnderlyingType(), false, false)
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
					c.err(errorf("TODO"))
				default:
					fmt.Fprintf(b, "%s%s", tag(field), nm)
				}
				b.WriteByte(' ')
				c.typ0(b, f.Type(), true, true)
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
			switch t.Align() {
			case 1:
				// ok
			case 2:
				fmt.Fprintf(b, "\n%s0 [0]uint16", tag(field))
			case 3, 4:
				fmt.Fprintf(b, "\n%s0 [0]uint32", tag(field))
			default:
				fmt.Fprintf(b, "\n%s0 [0]uint64", tag(field))
			}
			fmt.Fprintf(b, "\n%s1 [%d]byte\n}", tag(field), t.Size())
		}
	case *cc.ArrayType:
		fmt.Fprintf(b, "[%d]", x.Len())
		c.typ0(b, x.Elem(), true, true)
	default:
		b.WriteString("int")
		c.err(errorf("TODO %T", x))
	}
}

func (c *ctx) unionLiteral(t *cc.UnionType) string {
	var b strings.Builder
	c.typ0(&b, t, true, false)
	return b.String()
}

func (c *ctx) structLiteral(t *cc.StructType) string {
	var b strings.Builder
	c.typ0(&b, t, true, false)
	return b.String()
}

func (c *ctx) defineUnion(w writer, t *cc.UnionType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedUnions.add(nm) {
			return
		}

		w.w("\ntype %s%s = %s // %v:\n", tag(taggedUnion), nm, c.unionLiteral(t), c.pos(nmt))
	}
}

func (c *ctx) defineStruct(w writer, t *cc.StructType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedStructs.add(nm) {
			return
		}

		w.w("\ntype %s%s = %s // %v:\n", tag(taggedStruct), nm, c.structLiteral(t), c.pos(nmt))
	}
}

func (c *ctx) defineEnum(w writer, t *cc.EnumType) {
	if t.IsIncomplete() {
		return
	}

	nmt := t.Tag()
	if nm := nmt.SrcStr(); nm != "" && t.LexicalScope().Parent == nil {
		if !c.taggedEnums.add(nm) {
			return
		}

		w.w("\ntype %s%s = %s // %v:", tag(taggedEum), nm, c.typ(t.UnderlyingType()), c.pos(nmt))
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
		w.w("\n\t%s%s = %v // %v: ", tag(enumConst), nm, v.Value(), c.pos(v))
	}
	w.w("\n)\n")
}

func (c *ctx) defineEnumStructUnion(w writer, t cc.Type) {
	switch x := t.(type) {
	case *cc.EnumType:
		c.defineEnum(w, x)
	case *cc.StructType:
		c.defineStruct(w, x)
	case *cc.UnionType:
		c.defineUnion(w, x)
	}
}

func typeID(fset *token.FileSet, in map[string]ast.Expr, out map[string]string, typ ast.Expr) (r string, err error) {
	var b strings.Builder
	if err = typeID0(&b, fset, in, out, typ, map[string]struct{}{}); err != nil {
		return "", err
	}

	return b.String(), nil
}

func typeID0(b *strings.Builder, fset *token.FileSet, in map[string]ast.Expr, out map[string]string, typ ast.Expr, m map[string]struct{}) (err error) {
	switch x := typ.(type) {
	case *ast.Ident:
		switch symKind(x.Name) {
		case -1:
			b.WriteString(x.Name)
		case typename:
			if id, ok := out[x.Name]; ok {
				b.WriteString(id)
				break
			}

			t2, ok := in[x.Name]
			if !ok {
				return errorf("undefined type %s", x.Name)
			}

			if _, ok := m[x.Name]; ok {
				return errorf("invalid recursive type %s", x.Name)
			}

			m[x.Name] = struct{}{}
			id, err := typeID(fset, in, out, t2)
			if err != nil {
				return err
			}

			out[x.Name] = id
			b.WriteString(id)
		default:
			panic(todo("", x.Name, symKind(x.Name)))
		}
	case *ast.StructType:
		b.WriteString("struct{")
		for _, f := range x.Fields.List {
			ft, err := typeID(fset, in, out, f.Type)
			if err != nil {
				return err
			}

			for _, nm := range f.Names {
				fmt.Fprintf(b, "%s %s;", nm, ft)
			}
		}
		b.WriteByte('}')
	case *ast.ArrayType:
		fmt.Fprintf(b, "[")
		printer.Fprint(b, fset, x.Len)
		b.WriteByte(']')
		if err = typeID0(b, fset, in, out, x.Elt, m); err != nil {
			return err
		}
	default:
		panic(todo("%T", x))
	}
	return nil
}
