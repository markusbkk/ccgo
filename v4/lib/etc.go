// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"github.com/pbnjay/memory"
	"modernc.org/cc/v4"
	"modernc.org/gc/v2"
)

var (
	extendedErrors bool // true: Errors will include origin info.
	totalRam       = memory.TotalMemory()

	reservedNames = nameSet{
		// Keywords
		"break":       {},
		"case":        {},
		"chan":        {},
		"const":       {},
		"continue":    {},
		"default":     {},
		"defer":       {},
		"else":        {},
		"fallthrough": {},
		"for":         {},
		"func":        {},
		"go":          {},
		"goto":        {},
		"if":          {},
		"import":      {},
		"interface":   {},
		"map":         {},
		"package":     {},
		"range":       {},
		"return":      {},
		"select":      {},
		"struct":      {},
		"switch":      {},
		"type":        {},
		"var":         {},

		// Predeclared identifiers
		"any":        {},
		"append":     {},
		"bool":       {},
		"byte":       {},
		"cap":        {},
		"close":      {},
		"comparable": {},
		"complex":    {},
		"complex128": {},
		"complex64":  {},
		"copy":       {},
		"delete":     {},
		"error":      {},
		"false":      {},
		"float32":    {},
		"float64":    {},
		"imag":       {},
		"int":        {},
		"int16":      {},
		"int32":      {},
		"int64":      {},
		"int8":       {},
		"iota":       {},
		"len":        {},
		"make":       {},
		"new":        {},
		"nil":        {},
		"panic":      {},
		"print":      {},
		"println":    {},
		"real":       {},
		"recover":    {},
		"rune":       {},
		"string":     {},
		"true":       {},
		"uint":       {},
		"uint16":     {},
		"uint32":     {},
		"uint64":     {},
		"uint8":      {},
		"uintptr":    {},

		// Protected identifiers. Names we cannot give up because they have special
		// meaning in Go.
		"_":    {},
		"init": {},
		"main": {},
	}
)

// origin returns caller's short position, skipping skip frames.
func origin(skip int) string {
	pc, fn, fl, _ := runtime.Caller(skip)
	f := runtime.FuncForPC(pc)
	var fns string
	if f != nil {
		fns = f.Name()
		if x := strings.LastIndex(fns, "."); x > 0 {
			fns = fns[x+1:]
		}
		if strings.HasPrefix(fns, "func") {
			num := true
			for _, c := range fns[len("func"):] {
				if c < '0' || c > '9' {
					num = false
					break
				}
			}
			if num {
				return origin(skip + 2)
			}
		}
	}
	return fmt.Sprintf("%s:%d:%s", filepath.Base(fn), fl, fns)
}

// todo prints and return caller's position and an optional message tagged with TODO. Output goes to stderr.
func todo(s string, args ...interface{}) string {
	switch {
	case s == "":
		s = fmt.Sprintf(strings.Repeat("%v ", len(args)), args...)
	default:
		s = fmt.Sprintf(s, args...)
	}
	r := fmt.Sprintf("%s\n\tTODO %s", origin(2), s)
	// fmt.Fprintf(os.Stderr, "%s\n", r)
	// os.Stdout.Sync()
	return r
}

// trc prints and return caller's position and an optional message tagged with TRC. Output goes to stderr.
func trc(s string, args ...interface{}) string {
	switch {
	case s == "":
		s = fmt.Sprintf(strings.Repeat("%v ", len(args)), args...)
	default:
		s = fmt.Sprintf(s, args...)
	}
	r := fmt.Sprintf("%s: TRC %s", origin(2), s)
	fmt.Fprintf(os.Stderr, "%s\n", r)
	os.Stderr.Sync()
	return r
}

type errors []string

// Error implements error.
func (e errors) Error() string { return strings.Join(e, "\n") }

func (e *errors) add(err error) { *e = append(*e, err.Error()) }

func (e errors) err() error {
	w := 0
	for i, v := range e {
		if i != 0 {
			if prev, ok := extractPos(e[i-1]); ok {
				if cur, ok := extractPos(v); ok && prev.Filename == cur.Filename && prev.Line == cur.Line {
					continue
				}
			}
		}
		e[w] = v
		w++
	}
	e = e[:w]
	if len(e) == 0 {
		return nil
	}

	return e
}

// errorf constructs an error value. If ExtendedErrors is true, the error will
// contain its origin.
func errorf(s string, args ...interface{}) error {
	switch {
	case s == "":
		s = fmt.Sprintf(strings.Repeat("%v ", len(args)), args...)
	default:
		s = fmt.Sprintf(s, args...)
	}
	switch {
	case extendedErrors:
		return fmt.Errorf("%s (%v: %v: %v)", s, origin(4), origin(3), origin(2))
	default:
		return fmt.Errorf("%s", s)
	}
}

type parallel struct {
	errors    errors
	limit     chan struct{}
	resultTag string
	sync.Mutex
	wg sync.WaitGroup

	fails int32
	files int32
	ids   int32
	oks   int32
	skips int32
}

func newParallel(resultTag string) *parallel {
	limit := runtime.GOMAXPROCS(0)
	switch runtime.GOARCH {
	case "386", "arm": // 32 bit targets
		limit = 1
	default:
		if totalRam <= 1<<32 {
			limit = 1
		}
	}
	return &parallel{
		limit:     make(chan struct{}, limit),
		resultTag: resultTag,
	}
}

func (p *parallel) eh(msg string, args ...interface{}) { p.err(fmt.Errorf(msg, args...)) }

func (p *parallel) fail()   { atomic.AddInt32(&p.fails, 1) }
func (p *parallel) file()   { atomic.AddInt32(&p.files, 1) }
func (p *parallel) id() int { return int(atomic.AddInt32(&p.ids, 1)) }
func (p *parallel) ok()     { atomic.AddInt32(&p.oks, 1) }
func (p *parallel) skip()   { atomic.AddInt32(&p.skips, 1) }

func (p *parallel) exec(run func() error) {
	p.limit <- struct{}{}
	p.wg.Add(1)

	go func() {
		defer func() {
			p.wg.Done()
			<-p.limit
		}()

		p.err(run())
	}()
}

func (p *parallel) wait() error {
	p.wg.Wait()
	return p.errors.err()
}

func (p *parallel) err(err error) {
	if err == nil {
		return
	}

	p.Lock()
	p.errors.add(err)
	p.Unlock()
}

func extractPos(s string) (p token.Position, ok bool) {
	var prefix string
	if len(s) > 1 && s[1] == ':' { // c:\foo
		prefix = s[:2]
		s = s[2:]
	}
	// "testdata/parser/bug/001.c:1193:6: ..."
	a := strings.SplitN(s, ":", 4)
	// ["testdata/parser/bug/001.c" "1193" "6" "..."]
	if len(a) < 3 {
		return p, false
	}

	line, err := strconv.Atoi(a[1])
	if err != nil {
		return p, false
	}

	col, err := strconv.Atoi(a[2])
	if err != nil {
		return p, false
	}

	return token.Position{Filename: prefix + a[0], Line: line, Column: col}, true
}

func buildDefs(D, U []string) string {
	var a []string
	for _, v := range D {
		v = v[len("-D"):]
		if i := strings.IndexByte(v, '='); i > 0 {
			a = append(a, fmt.Sprintf("#define %s %s", v[:i], v[i+1:]))
			continue
		}

		a = append(a, fmt.Sprintf("#define %s 1", v))
	}
	for _, v := range U {
		v = v[len("-D"):]
		a = append(a, fmt.Sprintf("#undef %s", v))
	}
	return strings.Join(a, "\n")
}

type dict map[string]string

func (d *dict) put(k, v string) {
	if *d == nil {
		*d = map[string]string{}
	}
	(*d)[k] = v
}

type nameSpace struct {
	reg  nameRegister
	dict dict // link name -> Go name
}

func (n *nameSpace) registerNameSet(l *linker, set nameSet, tld bool) {
	var linkNames []string
	for nm := range set {
		linkNames = append(linkNames, nm)
	}
	sort.Slice(linkNames, func(i, j int) bool {
		a, b := linkNames[i], linkNames[j]
		x, y := symKind(a), symKind(b)
		return x < y || x == y && a < b
	})
	for _, linkName := range linkNames {
		switch k := symKind(linkName); k {
		case external:
			n.registerName(l, linkName)
		case typename, taggedEum, taggedStruct, taggedUnion, enumConst:
			if tld {
				panic(todo("", linkName))
			}

			n.registerName(l, linkName)
		case importQualifier:
			if tld {
				panic(todo("", linkName))
			}

			n.registerName(l, linkName)
		case staticInternal, staticNone:
			if tld {
				panic(todo("", linkName))
			}

			n.dict.put(linkName, l.tld.registerName(l, linkName))
			n.registerName(l, linkName)
		case automatic, ccgoAutomatic, ccgo:
			if tld {
				panic(todo("", linkName))
			}

			if _, ok := n.dict[linkName]; !ok {
				n.registerName(l, linkName)
			}
		case field:
			// nop
		case preserve, meta:
			// nop
		default:
			if k >= 0 {
				panic(todo("%q %v", linkName, symKind(linkName)))
			}
		}
	}
}

func (n *nameSpace) registerName(l *linker, linkName string) (goName string) {
	goName = l.goName(linkName)
	goName = n.reg.put(goName)
	n.dict.put(linkName, goName)
	// trc("%p: registered %q -> %q", n, linkName, goName)
	return goName
}

type nameRegister map[string]struct{}

// Colliding names will be adjusted by adding a numeric suffix.
//TODO quadratic when repeatedly adding the same name!
func (n *nameRegister) put(nm string) (r string) {
	// defer func(mn string) { trc("%q -> %q", nm, r) }(nm)
	if *n == nil {
		*n = map[string]struct{}{}
	}
	m := *n
	if !reservedNames.has(nm) && !m.has(nm) {
		m[nm] = struct{}{}
		return nm
	}

	l := 0
	for i := len(nm) - 1; i > 0; i-- {
		if c := nm[i]; c < '0' || c > '9' {
			break
		}

		l++
	}
	num := 0
	if l != 0 {
		if n, err := strconv.Atoi(nm[:len(nm)-l]); err == nil {
			num = n
		}
	}
	for num++; ; num++ {
		s2 := fmt.Sprintf("%s%d", nm, num)
		if _, ok := m[s2]; !ok {
			m[s2] = struct{}{}
			return s2
		}
	}
}

func (n *nameRegister) has(nm string) bool { _, ok := (*n)[nm]; return ok }

type nameSet map[string]struct{}

func (n *nameSet) add(s string) (ok bool) {
	if *n == nil {
		*n = map[string]struct{}{s: {}}
		return true
	}

	m := *n
	if _, ok = m[s]; ok {
		return false
	}

	m[s] = struct{}{}
	return true
}

func (n *nameSet) has(nm string) bool { _, ok := (*n)[nm]; return ok }

func symKind(s string) name {
	for i, v := range tags {
		if strings.HasPrefix(s, v) {
			return name(i)
		}
	}
	return -1
}

func enforceBinaryExt(s string) string {
	ext := filepath.Ext(s)
	s = s[:len(s)-len(ext)]
	switch runtime.GOOS {
	case "windows":
		return s + ".exe"
	}
	return s
}

func roundup(n, to int64) int64 {
	if r := n % to; r != 0 {
		return n + to - r
	}

	return n
}

func bpOff(n int64) string {
	if n != 0 {
		return fmt.Sprintf("%sbp%+d", tag(ccgo), n)
	}

	return fmt.Sprintf("%sbp", tag(ccgo))
}

func fldOff(n int64) string {
	if n != 0 {
		return fmt.Sprintf("%+d", n)
	}

	return ""
}

func export(s string) string {
	r, sz := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(r)) + s[sz:]
}

func unexport(s string) string {
	r, sz := utf8.DecodeRuneInString(s)
	return strings.ToLower(string(r)) + s[sz:]
}

func cpos(n cc.Node) (r token.Position) {
	if n == nil {
		return r
	}

	r = token.Position(n.Position())
	r.Filename = filepath.Join("~/src/modernc.org/ccorpus2", r.Filename)
	return r
}

// Same as cc.NodeSource but keeps the separators.
func nodeSource(s ...cc.Node) string {
	var a []cc.Token
	for _, n := range s {
		nodeTokens(n, &a)
	}
	sort.Slice(a, func(i, j int) bool { return a[i].Seq() < a[j].Seq() })
	var b strings.Builder
	for _, t := range a {
		// fmt.Fprintf(&b, " %d: ", t.Seq())
		b.Write(t.Sep())
		b.Write(t.Src())
	}
	return strings.Trim(b.String(), "\n")
}

func nodeTokens(n cc.Node, a *[]cc.Token) {
	if n == nil {
		return
	}

	t := reflect.TypeOf(n)
	v := reflect.ValueOf(n)
	var zero reflect.Value
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
		if v == zero {
			return
		}
	}

	if t.Kind() != reflect.Struct {
		return
	}

	if x, ok := n.(cc.Token); ok && x.Seq() != 0 {
		*a = append(*a, x)
		return
	}

	nf := t.NumField()
	for i := 0; i < nf; i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		if strings.HasPrefix(f.Name, "Token") {
			if x, ok := v.Field(i).Interface().(cc.Token); ok && x.Seq() != 0 {
				*a = append(*a, x)
			}
			continue
		}

		if v == zero || v.IsZero() {
			continue
		}

		if m, ok := v.Field(i).Interface().(cc.Node); ok {
			nodeTokens(m, a)
		}
	}
}

func sep(n cc.Node) string {
	var t cc.Token
	firstToken(n, &t)
	return strings.ReplaceAll(string(t.Sep()), "\f", "")
}

func firstToken(n cc.Node, r *cc.Token) {
	if n == nil {
		return
	}

	if x, ok := n.(*cc.LabeledStatement); ok && x != nil {
		t := x.Token
		if r.Seq() == 0 || t.Seq() < t.Seq() {
			*r = t
		}
		return
	}

	if x, ok := n.(cc.Token); ok && x.Seq() != 0 {
		if r.Seq() == 0 || x.Seq() < r.Seq() {
			*r = x
		}
		return
	}

	t := reflect.TypeOf(n)
	v := reflect.ValueOf(n)
	var zero reflect.Value
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}
	if v == zero || v.IsZero() || t.Kind() != reflect.Struct {
		return
	}

	nf := t.NumField()
	for i := 0; i < nf; i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		if m, ok := v.Field(i).Interface().(cc.Node); ok {
			firstToken(m, r)
		}
	}
}

func unresolvedSymbols(ast *gc.SourceFile) (r map[string]token.Position) {
	r = map[string]token.Position{}
	def := map[string]struct{}{}
	walk(ast, func(n gc.Node, pre bool, arg interface{}) {
		if !pre {
			return
		}

		switch x := n.(type) {
		case *gc.FunctionDecl:
			def[x.FunctionName.Src()] = struct{}{}
		case *gc.VarSpec:
			for _, v := range x.IdentifierList {
				def[v.Ident.Src()] = struct{}{}
			}
		case gc.Token:
			nm := x.Src()
			if _, ok := r[nm]; !ok {
				r[nm] = x.Position()
			}
		}
	}, nil)
	for k := range def {
		delete(r, k)
	}
	return r
}

func walk(n interface{}, fn func(n gc.Node, pre bool, arg interface{}), arg interface{}) {
	if n == nil {
		return
	}

	if x, ok := n.(gc.Token); ok {
		if x.IsValid() {
			fn(x, true, arg)
		}
		return
	}

	t := reflect.TypeOf(n)
	v := reflect.ValueOf(n)
	var zero reflect.Value
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
		if v == zero {
			return
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		if x, ok := n.(gc.Node); ok {
			fn(x, true, arg)
		}
		nf := t.NumField()
		for i := 0; i < nf; i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}

			if v == zero || v.IsZero() {
				continue
			}

			walk(v.Field(i).Interface(), fn, arg)
		}
		if x, ok := n.(gc.Node); ok {
			fn(x, false, arg)
		}
	case reflect.Slice:
		ne := v.Len()
		for i := 0; i < ne; i++ {
			walk(v.Index(i).Interface(), fn, arg)
		}
	}
}

func argumentExpressionListLen(n *cc.ArgumentExpressionList) (r int) {
	for ; n != nil; n = n.ArgumentExpressionList {
		r++
	}
	return r
}

func unsafe(fn string, arg interface{}) string {
	return fmt.Sprintf("%sunsafe.%s%s(%s)", tag(importQualifier), tag(preserve), fn, arg)
}

func unsafePointer(arg interface{}) string { return unsafe("Pointer", arg) }

func unsafeAddr(arg interface{}) string {
	return fmt.Sprintf("%sunsafe.%sPointer(&(%s))", tag(importQualifier), tag(preserve), arg)
}

func pos(n cc.Node) string {
	if n != nil {
		p := n.Position()
		p.Filename = filepath.Base(p.Filename)
		return p.String()
	}

	return "-"
}

func isOctalString(s string) bool {
	if s == "" || len(s) > 3 {
		return false
	}
	for _, v := range s {
		if v < '0' || v > '7' {
			return false
		}
	}

	return true
}

func firstError(err error, short bool) error {
	if !short || err == nil {
		return nil
	}

	if a := strings.Split(err.Error(), "\n"); len(a) != 0 {
		return fmt.Errorf("%s", a[0])
	}

	return err
}
