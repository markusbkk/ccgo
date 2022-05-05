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
	ast           *cc.AST
	cfg           *cc.Config
	eh            errHandler
	enumerators   nameSet
	f             *fnCtx
	fields        map[fielder]*nameSpace
	ifn           string
	imports       map[string]string // import path: qualifier
	out           io.Writer
	taggedEnums   nameSet
	taggedStructs nameSet
	taggedUnions  nameSet
	task          *Task
	typenames     nameSet
	void          cc.Type

	nextID int
	pass   int // 0: out of function, 1: func 1st pass, 2: func 2nd pass.

	closed bool
}

func newCtx(task *Task, eh errHandler) *ctx {
	return &ctx{
		cfg:     task.cfg,
		eh:      eh,
		fields:  map[fielder]*nameSpace{},
		imports: map[string]string{},
		task:    task,
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

func (c *ctx) fieldName(t fielder, f *cc.Field) string { return c.fields[t].dict[f.Name()] }

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

	sources := []cc.Source{
		{Name: "<predefined>", Value: c.cfg.Predefined},
		{Name: "<builtin>", Value: cc.Builtin},
	}
	if c.task.defs != "" {
		sources = append(sources, cc.Source{Name: "<command-line>", Value: c.task.defs})
	}
	sources = append(sources, cc.Source{Name: ifn, FS: c.cfg.FS})
	if c.ast, err = cc.Translate(c.cfg, sources); err != nil {
		return err
	}

	c.void = c.ast.Void
	c.ifn = ifn
	c.prologue(c)
	c.defines(c)
	for n := c.ast.TranslationUnit; n != nil; n = n.TranslationUnit {
		c.externalDeclaration(c, n.ExternalDeclaration)
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
			b = append(b, fmt.Sprintf("\n%s%s = %v // %v:", tag(macro), m.Name.Src(), x, c.pos(m.Name)))
		case cc.UInt64Value:
			b = append(b, fmt.Sprintf("\n%s%s = %v // %v:", tag(macro), m.Name.Src(), x, c.pos(m.Name)))
		case cc.Float64Value:
			if s := fmt.Sprint(x); s == r {
				b = append(b, fmt.Sprintf("\n%s%s = %s // %v:", tag(macro), m.Name.Src(), s, c.pos(m.Name)))
			}
		case cc.StringValue:
			b = append(b, fmt.Sprintf("\n%s%s = %q // %v:", tag(macro), m.Name.Src(), x[:len(x)-1], c.pos(m.Name)))
		}

		w.w("\n%s%s = %q // %v:", tag(define), m.Name.Src(), r, c.pos(m.Name))
	}
	w.w("\n)")
	if len(b) == 0 {
		return
	}

	w.w("\n\nconst (\n%s\n)", strings.Join(b, "\n"))
}

var home = os.Getenv("HOME")

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
