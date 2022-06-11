// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"bufio"
	"bytes"
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
	// keyword neither it can be a prefix of a Go predefined identifier.
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

func (b *buf) Write(p []byte) (int, error)     { b.b = append(b.b, p...); return len(p), nil }
func (b *buf) len() int                        { return len(b.b) }
func (b *buf) w(s string, args ...interface{}) { fmt.Fprintf(b, s, args...) }

func (b *buf) bytes() []byte {
	if b == nil {
		return nil
	}

	return b.b
}

var varTag = []byte(";/**/\n")

func (b *buf) vars(prefix string) (r string) {
	if !bytes.Contains(b.b, varTag) {
		return ""
	}

	a := strings.Split(string(b.bytes()), "\n")
	for i, v := range a {
		if !strings.HasPrefix(v, "var") || !strings.HasSuffix(v, ";/**/") {
			b.b = []byte(strings.Join(a[i:], "\n"))
			return prefix + strings.Join(a[:i], "\n")
		}
	}
	panic(todo("unrechable"))
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
	cfg                 *cc.Config
	defineTaggedStructs map[string]*cc.StructType
	eh                  errHandler
	enumerators         nameSet
	externsDeclared     map[string]*cc.Declarator
	externsDefined      map[string]struct{}
	externsMentioned    map[string]struct{}
	f                   *fnCtx
	fields              map[fielder]*nameSpace
	ifn                 string
	imports             map[string]string // import path: qualifier
	out                 io.Writer
	pvoid               cc.Type
	switchExpr          cc.Type
	taggedEnums         nameSet
	taggedStructs       nameSet
	taggedUnions        nameSet
	task                *Task
	typenames           nameSet
	void                cc.Type

	nextID int
	pass   int // 0: out of function, 1: func 1st pass, 2: func 2nd pass.

	closed  bool
	hasMain bool
}

func newCtx(task *Task, eh errHandler) *ctx {
	return &ctx{
		cfg:                 task.cfg,
		defineTaggedStructs: map[string]*cc.StructType{},
		eh:                  eh,
		externsDeclared:     map[string]*cc.Declarator{},
		externsDefined:      map[string]struct{}{},
		externsMentioned:    map[string]struct{}{},
		fields:              map[fielder]*nameSpace{},
		imports:             map[string]string{},
		task:                task,
	}
}

func (c *ctx) id() int {
	if c.f != nil {
		return c.f.id()
	}

	c.nextID++
	return c.nextID
}

func (c *ctx) err(err error) { c.eh("%s", err.Error()) }

func (c ctx) w(s string, args ...interface{}) {
	if c.closed {
		return
	}

	if _, err := fmt.Fprintf(c.out, s, args...); err != nil {
		c.err(err)
		c.closed = true
	}
}

func (c *ctx) compile(ifn, ofn string) error {
	f, err := os.Create(ofn)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			c.err(errorf("%v", err))
			return
		}

		if err := exec.Command("gofmt", "-s", "-w", "-r", "(x) -> x", ofn).Run(); err != nil {
			c.err(errorf("%s: gofmt: %v", ifn, err))
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
			c.defineStruct(c, "\n\n", nil, t)
			delete(c.defineTaggedStructs, k)
		}
	}
	c.w("%s", sep(c.ast.EOF))
	if c.hasMain && c.task.tlsQualifier != "" {
		c.w("\n\nfunc main() { %s%sStart(%smain) }\n", c.task.tlsQualifier, tag(preserve), tag(external))
	}
	var a []string
	for k := range c.externsDefined {
		delete(c.externsDeclared, k)
	}
	for k := range c.externsDeclared {
		if _, ok := c.externsMentioned[k]; ok {
			a = append(a, k)
		}
	}
	sort.Strings(a)
	for _, k := range a {
		switch d := c.externsDeclared[k]; t := d.Type().(type) {
		case *cc.FunctionType:
			c.w("\n\nfunc %s%s%s", tag(meta), k, c.signature(t, false, false))
		default:
			c.w("\n\nvar %s%s %s", tag(meta), k, c.typ(t))
		}
	}
	return nil
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
	var b []string
	w.w("\n\nconst (")
	for _, m := range a {
		r := m.ReplacementList()[0].SrcStr()
		switch x := m.Value().(type) {
		case cc.Int64Value:
			b = append(b, fmt.Sprintf("%s%s%s%s = %v;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), x))
		case cc.UInt64Value:
			b = append(b, fmt.Sprintf("%s%s%s%s = %v;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), x))
		case cc.Float64Value:
			if s := fmt.Sprint(x); s == r {
				b = append(b, fmt.Sprintf("%s%s%s%s = %s;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), s))
			}
		case cc.StringValue:
			b = append(b, fmt.Sprintf("%s%s%s%s = %q;", sep(m.Name), c.posComment(m), tag(macro), m.Name.Src(), x[:len(x)-1]))
		}

		w.w("%s%s%s%s = %q;", sep(m.Name), c.posComment(m), tag(define), m.Name.Src(), r)
	}
	w.w("\n)\n\n")
	if len(b) == 0 {
		return
	}

	w.w("\n\nconst (%s\n)\n\n", strings.Join(b, "\n"))
}

var home = os.Getenv("HOME")

func (c *ctx) posComment(n cc.Node) string {
	if !c.task.positions {
		return ""
	}

	return fmt.Sprintf("\n//  %s:\n", c.pos(n))
}

func (c *ctx) pos(n cc.Node) (r token.Position) {
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
