// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"sort"
	"strings"

	"modernc.org/cc/v4"
)

type declInfo struct {
	d     *cc.Declarator
	bpOff int64

	addressTaken bool
}

func (n *declInfo) pinned() bool { return n.d.StorageDuration() == cc.Automatic && n.addressTaken }

type declInfos map[*cc.Declarator]*declInfo

func (n *declInfos) info(d *cc.Declarator) (r *declInfo) {
	m := *n
	if m == nil {
		m = declInfos{}
		*n = m
	}
	if r = m[d]; r == nil {
		r = &declInfo{d: d}
		m[d] = r
	}
	return r
}

func (n *declInfos) takeAddress(d *cc.Declarator) { n.info(d).addressTaken = true }

type fnCtx struct {
	c         *ctx
	declInfos declInfos
	t         *cc.FunctionType
	tlsAllocs int64

	maxValist int
	nextID    int
}

func (c *ctx) newFnCtx(t *cc.FunctionType) (r *fnCtx) {
	return &fnCtx{c: c, t: t}
}

func (f *fnCtx) id() int { f.nextID++; return f.nextID }

func (c *ctx) externalDeclaration(w writer, n *cc.ExternalDeclaration) {
	switch n.Case {
	case cc.ExternalDeclarationFuncDef: // FunctionDefinition
		c.functionDefinition(w, n.FunctionDefinition)
	case cc.ExternalDeclarationDecl: // Declaration
		c.declaration(w, n.Declaration, true)
	case cc.ExternalDeclarationAsmStmt: // AsmStatement
		//TODO c.err(errorf("TODO %v", n.Case))
	case cc.ExternalDeclarationEmpty: // ';'
		//TODO c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) functionDefinition(w writer, n *cc.FunctionDefinition) {
	c.functionDefinition0(w, sep(n), n.Declarator, n.CompoundStatement, false)
}

func (c *ctx) functionDefinition0(w writer, sep string, d *cc.Declarator, cs *cc.CompoundStatement, local bool) {
	ft, ok := d.Type().(*cc.FunctionType)
	if !ok {
		c.err(errorf("%v: internal error %v", d.Position(), d.Type()))
		return
	}

	f0, pass := c.f, c.pass
	c.f = c.newFnCtx(ft)
	defer func() { c.f = f0; c.pass = pass }()
	c.pass = 1
	c.compoundStatement(discard{}, cs, true)
	var a []*cc.Declarator
	for d, n := range c.f.declInfos {
		if n.pinned() {
			a = append(a, d)
		}
	}
	sort.Slice(a, func(i, j int) bool {
		x := a[i].NameTok()
		y := a[j].NameTok()
		return x.Seq() < y.Seq()
	})
	for _, d := range a {
		info := c.f.declInfos[d]
		info.bpOff = roundup(c.f.tlsAllocs, int64(d.Type().Align()))
		c.f.tlsAllocs = info.bpOff + d.Type().Size()
	}
	c.pass = 2
	isMain := d.Linkage() == cc.External && d.Name() == "main"
	if isMain {
		c.hasMain = true
	}
	switch {
	case local:
		w.w("%s%s%s := func%s", sep, c.declaratorTag(d), d.Name(), c.signature(ft, true, false))
	default:
		w.w("%sfunc %s%s%s ", sep, c.declaratorTag(d), d.Name(), c.signature(ft, true, isMain))
	}
	c.compoundStatement(w, cs, true)
}

func (c *ctx) signature(f *cc.FunctionType, names, isMain bool) string {
	var b strings.Builder
	switch {
	case names:
		fmt.Fprintf(&b, "(%stls *%s%sTLS", tag(ccgo), c.task.tlsQualifier, tag(preserve))
	default:
		fmt.Fprintf(&b, "(*%s%sTLS", c.task.tlsQualifier, tag(preserve))
	}
	if f.MaxArgs() != 0 {
		for _, v := range f.Parameters() {
			b.WriteString(", ")
			if names {
				switch nm := v.Name(); {
				case nm == "":
					fmt.Fprintf(&b, "%sp ", tag(ccgo))
				default:
					switch info := c.f.declInfos.info(v.Declarator); {
					case info.pinned():
						fmt.Fprintf(&b, "%s_%s ", tag(ccgo), nm)
					default:
						fmt.Fprintf(&b, "%s%s ", tag(automatic), nm)
					}
				}
			}
			b.WriteString(c.typ(v.Type()))
		}
	}
	switch {
	case isMain && len(f.Parameters()) == 0 || isMain && len(f.Parameters()) == 1 && f.Parameters()[0].Type().Kind() == cc.Void:
		fmt.Fprintf(&b, ", %sargc int32, %[1]sargv uintptr", tag(ccgo))
	case isMain && len(f.Parameters()) == 1:
		fmt.Fprintf(&b, ", %sargv uintptr", tag(ccgo))
	case f.IsVariadic():
		switch {
		case names:
			fmt.Fprintf(&b, ", %sva uintptr", tag(ccgo))
		default:
			fmt.Fprintf(&b, ", uintptr")
		}
	}
	b.WriteByte(')')
	if f.Result().Kind() != cc.Void {
		if names {
			fmt.Fprintf(&b, "(%sr ", tag(ccgo))
		}
		b.WriteString(c.typ(f.Result()))
		if names {
			b.WriteByte(')')
		}
	}
	return b.String()
}

func (c *ctx) declaration(w writer, n *cc.Declaration, external bool) {
	switch n.Case {
	case cc.DeclarationDecl: // DeclarationSpecifiers InitDeclaratorList AttributeSpecifierList ';'
		switch {
		case n.InitDeclaratorList == nil:
			if !external {
				break
			}

			if n.DeclarationSpecifiers == nil {
				c.err(errorf("TODO %v:", n.Position()))
				break
			}

			switch x := n.DeclarationSpecifiers.Type().(type) {
			case *cc.EnumType:
				c.defineEnum(w, x)
			case *cc.StructType:
				c.defineStruct(w, x)
			case *cc.UnionType:
				c.defineUnion(w, x)
			}
		default:
			for l := n.InitDeclaratorList; l != nil; l = l.InitDeclaratorList {
				c.initDeclarator(w, sep(n), l.InitDeclarator, external)
			}
		}
	case cc.DeclarationAssert: // StaticAssertDeclaration
		c.err(errorf("TODO %v", n.Case))
	case cc.DeclarationAuto: // "__auto_type" Declarator '=' Initializer ';'
		c.err(errorf("TODO %v", n.Case))
	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
}

func (c *ctx) initDeclarator(w writer, sep string, n *cc.InitDeclarator, external bool) {
	d := n.Declarator
	if d.Type().Kind() == cc.Function || d.IsExtern() && d.Type().IsIncomplete() {
		return
	}

	if n.Asm != nil {
		w.w("//TODO %s %s // %v:", cc.NodeSource(d), cc.NodeSource(n.Asm), c.pos(n))
		if d.LexicalScope().Parent == nil {
			return
		}

		w.w("\n%spanic(0) // assembler statements not supported", tag(preserve))
	}

	var info *declInfo
	if c.f != nil {
		info = c.f.declInfos.info(d)
	}
	nm := d.Name()
	switch n.Case {
	case cc.InitDeclaratorDecl: // Declarator Asm
		switch {
		case d.IsTypename():
			if external && c.typenames.add(nm) {
				w.w("%stype %s%s = %s;", sep, tag(typename), nm, c.typedef(d.Type()))
				c.defineEnumStructUnion(w, d.Type())
			}
			if !external {
				return
			}
		default:
			if d.IsExtern() {
				return
			}

			c.defineEnumStructUnion(w, d.Type())
			switch {
			case info != nil && info.pinned():
				w.w("%svar _ /* %s */ %s;", sep, nm, c.typ(d.Type()))
			default:
				w.w("%svar %s%s %s;", sep, c.declaratorTag(d), nm, c.typ(d.Type()))
			}
		}
	case cc.InitDeclaratorInit: // Declarator Asm '=' Initializer
		c.defineEnumStructUnion(w, d.Type())
		switch {
		case d.Linkage() == cc.Internal:
			w.w("%svar %s%s = %s;", sep, c.declaratorTag(d), nm, c.initializerOuter(w, n.Initializer, d.Type()))
		case d.IsStatic():
			w.w("%svar %s%s = %s;", sep, c.declaratorTag(d), nm, c.initializerOuter(w, n.Initializer, d.Type()))
		default:
			switch {
			case info != nil && info.pinned():
				w.w("\n*(*%s)(unsafe.Pointer(%s)) = %s;", c.typ(d.Type()), bpOff(info.bpOff), c.initializerOuter(w, n.Initializer, d.Type()))
			default:
				switch {
				case d.LexicalScope().Parent == nil:
					w.w("%svar %s%s = %s;", sep, c.declaratorTag(d), nm, c.initializerOuter(w, n.Initializer, d.Type()))
				default:
					w.w("\n%s%s := %s;", c.declaratorTag(d), nm, c.initializerOuter(w, n.Initializer, d.Type()))
				}
			}
		}

	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	if info != nil {
		// w.w("\n// %p  %q read: %d, write: %d, address taken %v\n", d, d.Name(), d.ReadCount(), d.WriteCount(), d.AddressTaken()) //TODO-
		if d.StorageDuration() == cc.Automatic && d.ReadCount() == 0 && !info.pinned() {
			w.w("\n_ = %s%s;", c.declaratorTag(d), nm)
		}
	}
}
