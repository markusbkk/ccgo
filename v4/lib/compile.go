// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"bufio"
	"fmt"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
	"modernc.org/cc/v4"
)

type name int

const (
	defaultLibc = "modernc.org/libc"

	generatedFilePrefix = "Code generated for "
	generatedFileSuffix = ", DO NOT EDIT."
	//  package __ccgo_object_file_v1
	objectFilePackageName       = objectFilePackageNamePrefix + objectFileSemver
	objectFilePackageNamePrefix = "__ccgo_object_file_"
	objectFileSemver            = "v1"
)

const (
	// Lower number has higher priority in name allocation.
	external name = iota // storage class static, linkage external
	//TODO externalUnpinned

	typename
	taggedStruct
	taggedUnion
	taggedEum
	enumConst
	importQualifier

	macro
	define

	staticInternal // storage class static, linkage internal
	staticNone     // storage class static, linkage none
	automatic      // storage class automatic, linkage none, must be pinned if address taken
	ccgoAutomatic  // storage class automatic, linkage none, must be pinned if address taken
	ccgo           // not visible to transpiled C code, taking address is ok
	field          // field name
	meta           // linker metadata

	//TODO unpinned
	preserve
)

var (
	_ writer = (*buf)(nil)

	// Don't change the association once established, otherwise the major
	// objectFileSemver must be incremented.
	//
	// The concatenation of a tag and a valid C identifier must not create a Go
	// keyword neither it can be a prefix of a Go predefined/protected identifier,
	// see reservedNames.
	tags = [...]string{
		meta:            "_",
		ccgo:            "aa",
		ccgoAutomatic:   "cc",
		define:          "df", // #define
		enumConst:       "ec", // enumerator constant
		external:        "X",  // external linkage
		field:           "fd", // struct field
		importQualifier: "iq",
		macro:           "mv", // macro value
		automatic:       "an", // storage class automatic, linkage none
		staticInternal:  "si", // storage class static, linkage internal
		staticNone:      "sn", // storage class static, linkage none
		preserve:        "pp", // eg. TLS in iqlibc.ppTLS -> libc.TLS
		taggedEum:       "te", // tagged enum
		taggedStruct:    "ts", // tagged struct
		taggedUnion:     "tu", // tagged union
		typename:        "tn", // type name
		//TODO unpinned:        "un", // unpinned
	}
)

func init() {
	if !semver.IsValid(objectFileSemver) {
		panic(todo("internal error: invalid objectFileSemver: %q", objectFileSemver))
	}
}

type writer interface {
	w(s string, args ...interface{})
}

type discard struct{}

func (discard) w(s string, args ...interface{}) {}

type buf struct {
	b []byte
	n cc.Node
}

func newBufFromtring(s string) *buf { return &buf{b: []byte(s)} }

func (b *buf) Write(p []byte) (int, error) { b.b = append(b.b, p...); return len(p), nil }
func (b *buf) len() int                    { return len(b.b) }

func (b *buf) w(s string, args ...interface{}) {
	//trc("%v: %q %s", origin(2), s, args)
	fmt.Fprintf(b, s, args...)
}

func (b *buf) bytes() []byte {
	if b == nil {
		return nil
	}

	return b.b
}

func (b *buf) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		f.Write(b.bytes())
	case 'q':
		fmt.Fprintf(f, "%q", b.bytes())
	default:
		panic(todo("%q", string(verb)))
	}
}

func tag(nm name) string {
	if nm >= 0 {
		return tags[nm]
	}

	return ""
}

// errHandler is a function called on error.
type errHandler func(msg string, args ...interface{})

type ctx struct {
	ast                 *cc.AST
	breakCtx            string
	cfg                 *cc.Config
	compoundStmtValue   string
	continueCtx         string
	defineTaggedStructs map[string]*cc.StructType
	defineTaggedUnions  map[string]*cc.UnionType
	eh                  errHandler
	enumerators         nameSet
	externsDeclared     map[string]*cc.Declarator
	externsDefined      map[string]struct{}
	externsMentioned    map[string]struct{}
	f                   *fnCtx
	fields              map[fielder]*nameSpace
	ifn                 string
	imports             map[string]string // import path: qualifier
	initPatch           func(int64, *buf)
	out                 io.Writer
	pvoid               cc.Type
	switchCtx           []string
	taggedEnums         nameSet
	taggedStructs       nameSet
	taggedUnions        nameSet
	task                *Task
	typenames           nameSet
	verify              map[cc.Type]struct{}
	void                cc.Type

	nextID int
	pass   int // 0: out of function, 1: func 1st pass, 2: func 2nd pass.

	closed    bool
	hasMain   bool
	hasErrors bool
}

func newCtx(task *Task, eh errHandler) *ctx {
	return &ctx{
		cfg:                 task.cfg,
		defineTaggedStructs: map[string]*cc.StructType{},
		defineTaggedUnions:  map[string]*cc.UnionType{},
		eh:                  eh,
		externsDeclared:     map[string]*cc.Declarator{},
		externsDefined:      map[string]struct{}{},
		externsMentioned:    map[string]struct{}{},
		fields:              map[fielder]*nameSpace{},
		imports:             map[string]string{},
		task:                task,
		verify:              map[cc.Type]struct{}{},
	}
}

func (c *ctx) setBreakCtx(s string) func() {
	save := c.breakCtx
	c.breakCtx = s
	return func() { c.breakCtx = save }
}

func (c *ctx) setContinueCtx(s string) func() {
	save := c.continueCtx
	c.continueCtx = s
	return func() { c.continueCtx = save }
}

func (c *ctx) setSwitchCtx(s []string) func() {
	save := c.switchCtx
	c.switchCtx = s
	return func() { c.switchCtx = save }
}

func (c *ctx) baseName(n cc.Node) string {
	p := c.pos(n)
	return filepath.Base(p.Filename)
}

func (c *ctx) id() int {
	if c.f != nil {
		return c.f.id()
	}

	c.nextID++
	return c.nextID
}

func (c *ctx) err(err error) {
	c.hasErrors = true
	c.eh("%s", err.Error())
}

func (c *ctx) w(s string, args ...interface{}) {
	if c.closed {
		return
	}

	if _, err := fmt.Fprintf(c.out, s, args...); err != nil {
		c.err(err)
		c.closed = true
	}
}

func (c *ctx) compile(ifn, ofn string) (err error) {
	f, err := os.Create(ofn)
	if err != nil {
		return err
	}

	defer func() {
		if err2 := f.Close(); err2 != nil {
			c.err(errorf("%v", err2))
			if err == nil {
				err = err2
			}
			return
		}

		if c.hasErrors {
			return
		}

		if err2 := exec.Command("gofmt", "-s", "-w", "-r", "(x) -> x", ofn).Run(); err2 != nil {
			c.err(errorf("%s: gofmt: %v", ifn, err2))
			if err == nil {
				err = err2
			}
		}
		if *oTraceL {
			b, _ := os.ReadFile(ofn)
			fmt.Fprintf(os.Stderr, "%s\n", b)
		}
	}()

	w := bufio.NewWriter(f)
	c.out = w

	defer func() {
		if err := w.Flush(); err != nil {
			c.err(errorf("%v", err))
		}
	}()

	if c.ast, err = cc.Translate(c.cfg, sourcesFor(c.cfg, ifn, c.task.defs)); err != nil {
		return err
	}

	c.void = c.ast.Void
	c.pvoid = c.ast.PVoid
	c.ifn = ifn
	c.prologue(c)
	c.defines(c)
	for n := c.ast.TranslationUnit; n != nil; n = n.TranslationUnit {
		c.w("\n\n")
		c.externalDeclaration(c, n.ExternalDeclaration)
	}
	for len(c.defineTaggedStructs) != 0 {
		var a []string
		for k := range c.defineTaggedStructs {
			a = append(a, k)
		}
		sort.Strings(a)
		for k, t := range c.defineTaggedStructs {
			c.defineStructType(c, "\n\n", nil, t)
			delete(c.defineTaggedStructs, k)
		}
	}
	for len(c.defineTaggedUnions) != 0 {
		var a []string
		for k := range c.defineTaggedUnions {
			a = append(a, k)
		}
		sort.Strings(a)
		for k, t := range c.defineTaggedUnions {
			c.defineUnionType(c, "\n\n", nil, t)
			delete(c.defineTaggedUnions, k)
		}
	}
	c.verifyTypes()
	c.w("%s", sep(c.ast.EOF))
	if c.hasMain && c.task.tlsQualifier != "" {
		c.w("\n\nfunc %smain() { %s%[1]sStart(%[3]smain) }\n", tag(preserve), c.task.tlsQualifier, tag(external))
	}
	var a []string
	for k := range c.externsDefined {
		// trc("externsDefined %s", k)
		delete(c.externsDeclared, k)
	}
	for k := range c.externsDeclared {
		// trc("externsDeclared %s", k)
		if _, ok := c.externsMentioned[k]; ok {
			// trc("externsMentioned %s", k)
			a = append(a, k)
		}
	}
	sort.Strings(a)
	for _, k := range a {
		switch d := c.externsDeclared[k]; t := d.Type().(type) {
		case *cc.FunctionType:
			c.w("\n\nfunc %s%s%s", tag(meta), k, c.signature(t, false, false, false))
		default:
			c.w("\n\nvar %s%s %s", tag(meta), k, c.typ2(d, t, false))
		}
	}
	return nil
}

func (c *ctx) typeID(t cc.Type) string {
	var b strings.Builder
	c.typ0(&b, nil, t, false, false, false)
	return b.String()
}

func (c *ctx) verifyTypes() {
	if len(c.verify) == 0 {
		return
	}

	m := map[string]cc.Type{}
	for k := range c.verify {
		m[c.typeID(k)] = k
	}
	var a []string
	for k := range m {
		a = append(a, k)
	}
	sort.Strings(a)
	c.w("\n\nfunc init() {")
	for i, k := range a {
		t := m[k]
		v := fmt.Sprintf("%sv%d", tag(preserve), i)
		c.w("\n\tvar %s %s", v, c.initTyp(nil, t))
		if x, ok := t.(*cc.StructType); ok {
			t := x.Tag()
			if s := t.SrcStr(); s != "" {
				c.w("\n// struct %q", s)
			}
		}
		c.w("\nif g, e := %sunsafe.%sSizeof(%s), %[2]suintptr(%[4]d); g != e { panic(%[2]sg) }", tag(importQualifier), tag(preserve), v, t.Size())
		switch x := t.(type) {
		case *cc.StructType:
			groups := map[int64]struct{}{}
			for i := 0; i < x.NumFields(); i++ {
				f := x.FieldByIndex(i)
				switch {
				case f.IsBitfield():
					if f.InOverlapGroup() {
						continue
					}

					off := f.Offset()
					if _, ok := groups[off]; ok {
						break
					}

					groups[off] = struct{}{}
					sz := int64(f.GroupSize())
					nm := fmt.Sprintf("%s__ccgo%d", tag(field), off)
					c.w("\nif g, e := %sunsafe.%sSizeof(%s.%s), %[2]suintptr(%[5]d); g != e { panic(%[2]sg) }", tag(importQualifier), tag(preserve), v, nm, sz)
					c.w("\nif g, e := %sunsafe.%sOffsetof(%s.%s), %[2]suintptr(%[5]d); g != e { panic(%[2]sg) }", tag(importQualifier), tag(preserve), v, nm, off)
				default:
					if f.IsFlexibleArrayMember() {
						continue
					}

					if f.Type().Kind() == cc.Union {
						continue
					}

					off := f.Offset()
					sz := f.Type().Size()
					al := f.Type().FieldAlign()
					nm := tag(field) + c.fieldName(x, f)
					c.w("\nif g, e := %sunsafe.%sSizeof(%s.%s), %[2]suintptr(%[5]d); g != e { panic(%[2]sg) }", tag(importQualifier), tag(preserve), v, nm, sz)
					c.w("\nif g, e := %sunsafe.%sOffsetof(%s.%s), %[2]suintptr(%[5]d); g != e { panic(%[2]sg) }", tag(importQualifier), tag(preserve), v, nm, off)
					c.w("\nif g, e := %sunsafe.%sOffsetof(%s.%s) %% %[5]d, %[2]suintptr(0); g != e { panic(%[2]sg) }", tag(importQualifier), tag(preserve), v, nm, al)
				}
			}
		}
	}
	c.w("\n}")
}

func (c *ctx) defines(w writer) {
	var a []*cc.Macro
	for _, v := range c.ast.Macros {
		if v.IsConst && len(v.ReplacementList()) == 1 {
			a = append(a, v)
		}
	}
	if len(a) == 0 {
		return
	}

	sort.Slice(a, func(i, j int) bool { return a[i].Name.SrcStr() < a[j].Name.SrcStr() })
	for _, m := range a {
		r := m.ReplacementList()[0].SrcStr()
		w.w("%s%sconst %s%s = %q;", sep(m.Name), c.posComment(m), tag(define), m.Name.Src(), r)
		switch x := m.Value().(type) {
		case cc.Int64Value:
			w.w("%s%sconst %s%s = %v;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), x)
		case cc.UInt64Value:
			w.w("%s%sconst %s%s = %v;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), x)
		case cc.Float64Value:
			if s := fmt.Sprint(x); s == r {
				w.w("%s%sconst %s%s = %s;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), s)
			}
		case cc.StringValue:
			w.w("%s%sconst %s%s = %q;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), x[:len(x)-1])
		}
	}
}

var home = os.Getenv("HOME")

func (c *ctx) posComment(n cc.Node) string {
	if !c.task.positions {
		return ""
	}

	return fmt.Sprintf("\n//  %s:\n", c.pos(n))
}

func (c *ctx) pos(n cc.Node) (r token.Position) {
	if n == nil {
		return r
	}

	if r = token.Position(n.Position()); r.IsValid() {
		switch {
		case c.task.fullPaths:
			if strings.HasPrefix(r.Filename, home) {
				r.Filename = "$HOME" + r.Filename[len(home):]
			}
		default:
			r.Filename = filepath.Base(r.Filename)
		}
	}
	return r
}

func (c *ctx) prologue(w writer) {
	w.w(`// %s%s/%s by '%s %s'%s

//go:build ignore
// +build ignore

package %s
`,
		generatedFilePrefix,
		c.task.goos, c.task.goarch,
		filepath.Base(c.task.args[0]),
		strings.Join(c.task.args[1:], " "),
		generatedFileSuffix,
		objectFilePackageName,
	)
}

func (c *ctx) declaratorTag(d *cc.Declarator) string { return tag(c.declaratorKind(d)) }

func (c *ctx) declaratorKind(d *cc.Declarator) name {
	switch d.Linkage() {
	case cc.External:
		return external
	case cc.Internal:
		return staticInternal
	case cc.None:
		switch {
		case d.IsStatic():
			return staticNone
		default:
			return automatic
		}
	default:
		c.err(errorf("%v: internal error: %v", d.Position(), d.Linkage()))
		return -1
	}
}
