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

const (
	retvalName = "r"
	vaArgName  = "va"
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
	autovars   []string
	c          *ctx
	declInfos  declInfos
	flatScopes map[*cc.Scope]struct{}
	locals     map[*cc.Declarator]string // storage: static or automatic, linkage: none -> C renamed
	t          *cc.FunctionType
	tlsAllocs  int64

	maxValist int
	nextID    int
}

func (c *ctx) newFnCtx(t *cc.FunctionType, n *cc.CompoundStatement) (r *fnCtx) {
	fnScope := n.LexicalScope()
	// trc("%v: ==== fnScope %p, parent %p\n%s", n.Position(), fnScope, fnScope.Parent, dumpScope(fnScope))
	var flatScopes map[*cc.Scope]struct{}
next:
	for _, gotoStmt := range n.Gotos() {
		gotoScope := gotoStmt.LexicalScope()
		// trc("%v: '%s', gotoScope %p, parent %p\n%s", gotoStmt.Position(), cc.NodeSource(gotoStmt), gotoScope, gotoScope.Parent, dumpScope(gotoScope))
		var targetScope *cc.Scope
		switch x := gotoStmt.Label().(type) {
		case *cc.LabeledStatement:
			targetScope = x.LexicalScope()
			// trc("%v: '%s', targetScope %p, parent %p\n%s", x.Position(), cc.NodeSource(x), targetScope, targetScope.Parent, dumpScope(targetScope))
		default:
			c.err(errorf("TODO %T", x))
			continue next
		}

		m := map[*cc.Scope]struct{}{gotoScope: {}}
		// targetScope must be the same as gotoScope or any of its parent scopes.
		for sc := gotoScope; sc != nil && sc.Parent != nil; sc = sc.Parent {
			m[sc] = struct{}{}
			// trc("searching scope %p, parent %p\n%s", sc, sc.Parent, dumpScope(sc))
			if sc == targetScope {
				// trc("FOUND targetScope")
				continue next
			}
		}

		// Jumping into a block.
		if flatScopes == nil {
			flatScopes = map[*cc.Scope]struct{}{}
		}
		for sc := targetScope; sc != nil && sc != fnScope; sc = sc.Parent {
			// trc("FLAT[%p]", sc)
			flatScopes[sc] = struct{}{}
			if _, ok := m[sc]; ok {
				// trc("FOUND common scope")
				break
			}
		}
	}
	return &fnCtx{c: c, t: t, flatScopes: flatScopes}
}

func (f *fnCtx) newAutovarName() (nm string) {
	return fmt.Sprintf("%sv%d", tag(ccgoAutomatic), f.c.id())
}

func (f *fnCtx) newAutovar(n cc.Node, t cc.Type) (nm string) {
	nm = f.newAutovarName()
	f.registerAutoVar(fmt.Sprintf("var %s %s;", nm, f.c.typ(n, t)))
	return nm
}

func (f *fnCtx) registerAutoVar(s string) { f.autovars = append(f.autovars, s) }

func (f *fnCtx) registerLocal(d *cc.Declarator) {
	if f == nil {
		return
	}

	if f.locals == nil {
		f.locals = map[*cc.Declarator]string{}
	}
	f.locals[d] = ""
}

func (f *fnCtx) renameLocals() {
	var a []*cc.Declarator
	for k := range f.locals {
		a = append(a, k)
	}
	sort.Slice(a, func(i, j int) bool {
		x, y := a[i], a[j]
		if x.Name() < y.Name() {
			return true
		}

		if x.Name() > y.Name() {
			return false
		}

		return x.Visible() < y.Visible()
	})
	var r nameRegister
	for _, d := range a {
		f.locals[d] = r.put(f.c.declaratorTag(d) + d.Name())
	}
}

func (f *fnCtx) declareLocals() string {
	var a []string
	for k, v := range f.locals {
		if info := f.declInfos[k]; info != nil && info.pinned() {
			a = append(a, fmt.Sprintf("var %s_ /* %s at bp%+d */ %s;", tag(preserve), k.Name(), info.bpOff, f.c.typ(k, k.Type())))
			continue
		}

		if k.IsTypename() {
			continue
		}

		if k.StorageDuration() != cc.Static && v != "" {
			a = append(a, fmt.Sprintf("var %s %s;", v, f.c.typ(k, k.Type())))
		}
	}
	a = append(a, f.autovars...)
	sort.Strings(a)
	return strings.Join(a, "")
}

func (f *fnCtx) id() int { f.nextID++; return f.nextID }

func (c *ctx) externalDeclaration(w writer, n *cc.ExternalDeclaration) {
	switch n.Case {
	case cc.ExternalDeclarationFuncDef: // FunctionDefinition
		if d := n.FunctionDefinition.Declarator; d.Linkage() == cc.External {
			c.externsDefined[n.FunctionDefinition.Declarator.Name()] = struct{}{}
		}
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
	c.functionDefinition0(w, sep(n), n, n.Declarator, n.CompoundStatement, false)
}

func (c *ctx) functionDefinition0(w writer, sep string, pos cc.Node, d *cc.Declarator, cs *cc.CompoundStatement, local bool) {
	ft, ok := d.Type().(*cc.FunctionType)
	if !ok {
		c.err(errorf("%v: internal error %v", d.Position(), d.Type()))
		return
	}

	c.checkValidType(d, ft)
	f0, pass := c.f, c.pass
	c.f = c.newFnCtx(ft, cs)
	defer func() { c.f = f0; c.pass = pass }()
	c.pass = 1
	c.compoundStatement(discard{}, cs, true, "")
	c.f.renameLocals()
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
	c.f.nextID = 0
	isMain := d.Linkage() == cc.External && d.Name() == "main"
	if isMain {
		c.hasMain = true
	}
	s := strings.TrimRight(sep, "\n\r\t ")
	s += c.posComment(pos)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
		if s == "\n" {
			s = "\n\n"
		}
	}
	switch {
	case local:
		w.w("%s%s%s := func%s", s, c.declaratorTag(d), d.Name(), c.signature(ft, true, false, true))
	default:
		w.w("%sfunc %s%s%s ", s, c.declaratorTag(d), d.Name(), c.signature(ft, true, isMain, true))
	}
	c.compoundStatement(w, cs, true, "")
}

func (c *ctx) signature(f *cc.FunctionType, paramNames, isMain, useNames bool) string {
	var b strings.Builder
	switch {
	case paramNames:
		fmt.Fprintf(&b, "(%stls *%s%sTLS", tag(ccgo), c.task.tlsQualifier, tag(preserve))
	default:
		fmt.Fprintf(&b, "(*%s%sTLS", c.task.tlsQualifier, tag(preserve))
	}
	if f.MaxArgs() != 0 {
		for i, v := range f.Parameters() {
			if !c.checkValidParamType(v, v.Type()) {
				return ""
			}

			b.WriteString(", ")
			if paramNames {
				switch nm := v.Name(); {
				case nm == "":
					fmt.Fprintf(&b, "%sp%d ", tag(ccgo), i)
				default:
					switch info := c.f.declInfos.info(v.Declarator); {
					case info.pinned():
						fmt.Fprintf(&b, "%s_%s ", tag(ccgo), nm)
					default:
						fmt.Fprintf(&b, "%s%s ", tag(automatic), nm)
					}
				}
			}
			b.WriteString(c.typ2(v, v.Type().Decay(), useNames))
		}
	}
	switch {
	case isMain && len(f.Parameters()) == 0 || isMain && len(f.Parameters()) == 1 && f.Parameters()[0].Type().Kind() == cc.Void:
		fmt.Fprintf(&b, ", %sargc %sint32, %[1]sargv %suintptr", tag(ccgo), tag(preserve))
	case isMain && len(f.Parameters()) == 1:
		fmt.Fprintf(&b, ", %sargv %suintptr", tag(ccgo), tag(preserve))
	case f.IsVariadic():
		switch {
		case paramNames:
			fmt.Fprintf(&b, ", %s%s %suintptr", tag(ccgo), vaArgName, tag(preserve))
		default:
			fmt.Fprintf(&b, ", %suintptr", tag(preserve))
		}
	}
	b.WriteByte(')')
	if f.Result().Kind() != cc.Void {
		if paramNames {
			fmt.Fprintf(&b, "(%s%s ", tag(ccgo), retvalName)
		}
		b.WriteString(c.typ2(nil, f.Result(), useNames))
		if paramNames {
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
				break
			}

			sep := sep(n)
			switch x := n.DeclarationSpecifiers.Type().(type) {
			case *cc.EnumType:
				c.defineEnum(w, sep, n, x)
			case *cc.StructType:
				c.defineStruct(w, sep, n, x)
			case *cc.UnionType:
				c.defineUnion(w, sep, n, x)
			}
		default:
			w.w("%s", sep(n))
			for l := n.InitDeclaratorList; l != nil; l = l.InitDeclaratorList {
				c.initDeclarator(w, sep(l.InitDeclarator), l.InitDeclarator, external)
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
	if sc := d.LexicalScope(); sc.Parent == nil {
		hasInitializer := false
		for _, v := range sc.Nodes[d.Name()] {
			if x, ok := v.(*cc.Declarator); ok && x.HasInitializer() {
				hasInitializer = true
				break
			}
		}
		if hasInitializer && !d.HasInitializer() {
			return
		}
	}

	if attr := d.Type().Attributes(); attr != nil {
		if attr.Alias() != "" {
			c.err(errorf("TODO unsupported attribute(s)"))
			return
		}
	}

	if d.Type().Kind() == cc.Function && d.Linkage() == cc.External || d.IsExtern() && !d.Type().IsIncomplete() {
		c.externsDeclared[d.Name()] = d
	}

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

	nm := d.Name()
	linkName := c.declaratorTag(d) + nm
	var info *declInfo
	if c.f != nil {
		info = c.f.declInfos.info(d)
	}
	switch c.pass {
	case 1:
		if d.Linkage() == cc.None {
			c.f.registerLocal(d)
		}
	case 2:
		if nm := c.f.locals[d]; nm != "" {
			linkName = nm
		}
	}
	switch n.Case {
	case cc.InitDeclaratorDecl: // Declarator Asm
		switch {
		case d.IsTypename():
			if external && c.typenames.add(nm) && !d.Type().IsIncomplete() {
				w.w("\n\n%s%stype %s%s = %s;", sep, c.posComment(n), tag(typename), nm, c.typedef(d, d.Type()))
				c.defineEnumStructUnion(w, sep, n, d.Type())
			}
			if !external {
				return
			}
		default:
			if d.IsExtern() {
				return
			}

			c.defineEnumStructUnion(w, sep, n, d.Type())
			switch {
			case d.IsStatic():
				switch c.pass {
				case 1:
					// nop
				case 2:
					if nm := c.f.locals[d]; nm != "" {
						w.w("%s%svar %s %s;", sep, c.posComment(n), nm, c.typ(d, d.Type()))
						break
					}

					fallthrough
				default:
					w.w("%s%svar %s %s;", sep, c.posComment(n), linkName, c.typ(d, d.Type()))
				}
			default:
				switch c.pass {
				case 0:
					w.w("%s%svar %s %s;", sep, c.posComment(n), linkName, c.typ(d, d.Type()))
				}
			}
		}
	case cc.InitDeclaratorInit: // Declarator Asm '=' Initializer
		c.defineEnumStructUnion(w, sep, n, d.Type())
		switch {
		case d.Linkage() == cc.Internal:
			w.w("%s%svar %s = %s;", sep, c.posComment(n), linkName, c.initializerOuter(w, n.Initializer, d.Type()))
		case d.IsStatic():
			switch c.pass {
			case 1:
				// nop
			case 2:
				if nm := c.f.locals[d]; nm != "" {
					w.w("%s%svar %s = %s;", sep, c.posComment(n), nm, c.initializerOuter(w, n.Initializer, d.Type()))
					break
				}

				fallthrough
			default:
				w.w("%s%svar %s = %s;", sep, c.posComment(n), linkName, c.initializerOuter(w, n.Initializer, d.Type()))
			}
		default:
			switch {
			case info != nil && info.pinned():
				w.w("%s%s*(*%s)(%s) = %s;", sep, c.posComment(n), c.typ(d, d.Type()), unsafePointer(bpOff(info.bpOff)), c.initializerOuter(w, n.Initializer, d.Type()))
			default:
				switch {
				case d.LexicalScope().Parent == nil:
					w.w("%s%svar %s = %s;", sep, c.posComment(n), linkName, c.initializerOuter(w, n.Initializer, d.Type()))
				default:
					w.w("%s%s%s = %s;", sep, c.posComment(n), linkName, c.initializerOuter(w, n.Initializer, d.Type()))
				}
			}
		}

	default:
		c.err(errorf("internal error %T %v", n, n.Case))
	}
	if info != nil {
		// w.w("\n// read: %d, write: %d, address taken %v\n", d.ReadCount(), d.WriteCount(), d.AddressTaken()) //TODO-
		if d.StorageDuration() == cc.Automatic && d.ReadCount() == d.SizeofCount() && !info.pinned() {
			w.w("\n%s_ = %s;", tag(preserve), linkName)
		}
	}
}
