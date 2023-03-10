// Copyright 2017 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo"

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"modernc.org/cc"
	"modernc.org/ccir"
	"modernc.org/internal/buffer"
	"modernc.org/ir"
	"modernc.org/irgo"
	"modernc.org/strutil"
	"modernc.org/xc"
)

func caller(s string, va ...interface{}) {
	if s == "" {
		s = strings.Repeat("%v ", len(va))
	}
	_, fn, fl, _ := runtime.Caller(2)
	fmt.Fprintf(os.Stderr, "# caller: %s:%d: ", path.Base(fn), fl)
	fmt.Fprintf(os.Stderr, s, va...)
	fmt.Fprintln(os.Stderr)
	_, fn, fl, _ = runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "# \tcallee: %s:%d: ", path.Base(fn), fl)
	fmt.Fprintln(os.Stderr)
	os.Stderr.Sync()
}

func dbg(s string, va ...interface{}) {
	if s == "" {
		s = strings.Repeat("%v ", len(va))
	}
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "# dbg %s:%d: ", path.Base(fn), fl)
	fmt.Fprintf(os.Stderr, s, va...)
	fmt.Fprintln(os.Stderr)
	os.Stderr.Sync()
}

func use(...interface{}) {}

func init() {
	use(caller, dbg, TODO) //TODOOK
	flag.BoolVar(&Testing, "testing", false, "")
	flag.BoolVar(&ir.Testing, "irTesting", false, "")
	flag.BoolVar(&irgo.FTrace, "irgoFTrace", false, "")
	flag.BoolVar(&irgo.Testing, "irgoTesting", false, "")
}

// ============================================================================

const (
	crtQ     = "crt"
	prologue = `// Code generated by ccgo DO NOT EDIT.

package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"runtime"
	"unsafe"

	"modernc.org/ccgo/crt"
)

var argv []*int8

func ftrace(s string, args ...interface{}) {
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "# %%s:%%d: %%v\n", path.Base(fn), fl, fmt.Sprintf(s, args...))
	os.Stderr.Sync()
}

func main() {
	os.Args[0] = "./test"
	for _, v := range os.Args {
		argv = append(argv, (*int8)(crt.CString(v)))
	}
	argv = append(argv, nil)
	X_start(%s.NewTLS(), int32(len(os.Args)), &argv[0])
}

%s`
)

var (
	ccTestdata string

	cpp      = flag.Bool("cpp", false, "")
	errLimit = flag.Int("errlimit", 10, "")
	filter   = flag.String("re", "", "")
	ndebug   = flag.Bool("ndebug", false, "")
	noexec   = flag.Bool("noexec", false, "")
	oLog     = flag.Bool("log", false, "")
	trace    = flag.Bool("trc", false, "")
	yydebug  = flag.Int("yydebug", 0, "")
)

func init() {
	ip, err := cc.ImportPath()
	if err != nil {
		panic(err)
	}

	for _, v := range filepath.SplitList(strutil.Gopath()) {
		p := filepath.Join(v, "src", ip, "testdata")
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}

		if fi.IsDir() {
			ccTestdata = p
			break
		}
	}
	if ccTestdata == "" {
		panic("cannot find cc/testdata/")
	}
}

func errStr(err error) string {
	switch x := err.(type) {
	case scanner.ErrorList:
		if len(x) != 1 {
			x.RemoveMultiples()
		}
		var b bytes.Buffer
		for i, v := range x {
			if i != 0 {
				b.WriteByte('\n')
			}
			b.WriteString(v.Error())
			if i == 9 {
				fmt.Fprintf(&b, "\n\t... and %v more errors", len(x)-10)
				break
			}
		}
		return b.String()
	default:
		return err.Error()
	}
}

func parse(src []string, opts ...cc.Opt) (_ *cc.TranslationUnit, err error) {
	defer func() {
		if e := recover(); e != nil && err == nil {
			err = fmt.Errorf("cc.Parse: PANIC: %v\n%s", e, debug.Stack())
		}
	}()

	model, err := ccir.NewModel()
	if err != nil {
		return nil, err
	}

	var ndbg string
	if *ndebug {
		ndbg = "#define NDEBUG 1"
	}
	ast, err := cc.Parse(fmt.Sprintf(`
%s
#define _CCGO 1
#define __arch__ %s
#define __os__ %s
#include <builtin.h>

#define NO_TRAMPOLINES 1
`, ndbg, runtime.GOARCH, runtime.GOOS),
		src,
		model,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("cc.Parse: %v", errStr(err))
	}

	return ast, nil
}

func expect1(wd, match string, hook func(string, string) []string, ccgoOpts []Option, opts ...cc.Opt) (log buffer.Bytes, exitStatus int, err error) {
	var lpos token.Position
	if *cpp {
		opts = append(opts, cc.Cpp(func(toks []xc.Token) {
			if len(toks) != 0 {
				p := toks[0].Position()
				if p.Filename != lpos.Filename {
					fmt.Fprintf(&log, "# %d %q\n", p.Line, p.Filename)
				}
				lpos = p
			}
			for _, v := range toks {
				log.WriteString(cc.TokSrc(v))
			}
			log.WriteByte('\n')
		}))
	}
	if n := *yydebug; n != -1 {
		opts = append(opts, cc.YyDebug(n))
	}
	ast, err := parse([]string{ccir.CRT0Path, match}, opts...)
	if err != nil {
		return log, -1, err
	}

	var out, src buffer.Bytes

	defer func() {
		out.Close()
		src.Close()
	}()

	if err := New([]*cc.TranslationUnit{ast}, &out, ccgoOpts...); err != nil {
		return log, -1, fmt.Errorf("New: %v", err)
	}

	fmt.Fprintf(&src, prologue, crtQ, out.Bytes())
	b, err := format.Source(src.Bytes())
	if err != nil {
		b = src.Bytes()
	}
	fmt.Fprintf(&log, "# ccgo.New\n%s", b)
	if err != nil {
		return log, exitStatus, err
	}

	if *noexec {
		return log, 0, nil
	}

	var stdout, stderr buffer.Bytes

	defer func() {
		stdout.Close()
		stderr.Close()
	}()

	if err := func() (err error) {
		defer func() {
			if e := recover(); e != nil && err == nil {
				err = fmt.Errorf("exec: PANIC: %v", e)
			}
		}()

		vwd, err := ioutil.TempDir("", "ccgo-test-")
		if err != nil {
			return err
		}

		if err := os.Chdir(vwd); err != nil {
			return err
		}

		defer func() {
			os.Chdir(wd)
			os.RemoveAll(vwd)
		}()

		if err := ioutil.WriteFile("main.go", b, 0664); err != nil {
			return err
		}

		args := hook(vwd, match)
		cmd := exec.Command("go", append([]string{"run", "main.go"}, args[1:]...)...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			if b := stdout.Bytes(); b != nil {
				fmt.Fprintf(&log, "stdout:\n%s\n", b)
			}
			if b := stderr.Bytes(); b != nil {
				fmt.Fprintf(&log, "stderr:\n%s\n", b)
			}
			return fmt.Errorf("go run: exit status %v, err %v", exitStatus, err)
		}

		return nil
	}(); err != nil {
		return log, 1, err
	}

	if b := stdout.Bytes(); b != nil {
		fmt.Fprintf(&log, "stdout:\n%s\n", b)
	}
	if b := stderr.Bytes(); b != nil {
		fmt.Fprintf(&log, "stderr:\n%s\n", b)
	}

	expect := match[:len(match)-len(filepath.Ext(match))] + ".expect"
	if _, err := os.Stat(expect); err != nil {
		if !os.IsNotExist(err) {
			return log, 0, err
		}

		return log, 0, nil
	}

	buf, err := ioutil.ReadFile(expect)
	if err != nil {
		return log, 0, err
	}

	if g, e := stdout.Bytes(), buf; !bytes.Equal(g, e) {
		return log, 0, fmt.Errorf("==== %v\n==== got\n%s==== exp\n%s", match, g, e)
	}
	return log, 0, nil
}

func expect(t *testing.T, dir string, skip func(string) bool, hook func(string, string) []string, ccgoOpts []Option, opts ...cc.Opt) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*.c"))
	if err != nil {
		t.Fatal(err)
	}

	seq := 0
	okSeq := 0
	for _, match := range matches {
		if skip(match) {
			continue
		}

		if *trace {
			fmt.Println(match)
		}
		seq++
		doLog := *oLog
		b, err := ioutil.ReadFile(match)
		if err != nil {
			t.Fatal(err)
		}

		co := ccgoOpts
		co = co[:len(co):len(co)]
		if bytes.Contains(b, []byte("#include")) {
			co = append(co, LibcTypes())
		}
		log, exitStatus, err := expect1(wd, match, hook, co, opts...)
		switch {
		case exitStatus <= 0 && err == nil:
			okSeq++
		default:
			if seq-okSeq == 1 {
				t.Logf("%s: FAIL\n%s\n%s", match, errStr(err), log.Bytes())
				doLog = false
			}
		}
		if doLog {
			t.Logf("%s:\n%s", match, log.Bytes())
		}
		log.Close()
	}
	t.Logf("%v/%v ok", okSeq, seq)
	if okSeq != seq {
		t.Errorf("failures: %v", seq-okSeq)
	}
}

func TestTCC(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdata, err := filepath.Rel(wd, ccTestdata)
	if err != nil {
		t.Fatal(err)
	}

	var re *regexp.Regexp
	if s := *filter; s != "" {
		re = regexp.MustCompile(s)
	}

	dir := filepath.Join(testdata, filepath.FromSlash("tcc-0.9.26/tests/tests2/"))
	expect(
		t,
		dir,
		func(match string) bool {
			if re != nil && !re.MatchString(filepath.Base(match)) {
				return true
			}

			return false
		},
		func(wd, match string) []string {
			switch filepath.Base(match) {
			case "31_args.c":
				return []string{"./test", "-", "arg1", "arg2", "arg3", "arg4"}
			case "46_grep.c":
				ioutil.WriteFile(filepath.Join(wd, "test"), []byte("abc\ndef\nghi\n"), 0600)
				return []string{"./grep", "[ea]", "test"}
			default:
				return []string{match}
			}
		},
		nil,
		cc.AllowCompatibleTypedefRedefinitions(),
		cc.EnableEmptyStructs(),
		cc.EnableImplicitFuncDef(),
		cc.ErrLimit(-1),
		cc.KeepComments(),
		cc.SysIncludePaths([]string{ccir.LibcIncludePath}),
	)
}

func TestGCCExec(t *testing.T) {
	blacklist := map[string]struct{}{
		// VLA struct field.
		"20020412-1.c": {},
		"20040308-1.c": {},
		"align-nest.c": {},
		"pr41935.c":    {},

		// Nested function.
		"20010209-1.c":   {},
		"20010605-1.c":   {},
		"20030501-1.c":   {},
		"20040520-1.c":   {},
		"20061220-1.c":   {},
		"20090219-1.c":   {},
		"920612-2.c":     {},
		"921017-1.c":     {},
		"nest-align-1.c": {},
		"nest-stdar-1.c": {},
		"nestfunc-7.c":   {},
		"pr22061-3.c":    {},
		"pr22061-4.c":    {},
		"pr71494.c":      {},

		// __real__, complex integers and and friends.
		"20010605-2.c": {},
		"20020411-1.c": {},
		"20030910-1.c": {},
		"20041124-1.c": {},
		"20041201-1.c": {},
		"20050121-1.c": {},
		"complex-1.c":  {},
		"complex-6.c":  {},
		"pr38151.c":    {},
		"pr38969.c":    {},
		"pr56837.c":    {},

		// Depends on __attribute__((aligned(N)))
		"20010904-1.c": {},
		"20010904-2.c": {},
		"align-3.c":    {},
		"pr23467.c":    {},

		// Depends on __attribute__ ((vector_size (N)))
		"20050316-1.c":   {},
		"20050316-2.c":   {},
		"20050316-3.c":   {},
		"20050604-1.c":   {},
		"20050607-1.c":   {},
		"pr23135.c":      {},
		"pr53645-2.c":    {},
		"pr53645.c":      {},
		"pr60960.c":      {},
		"pr65427.c":      {},
		"pr71626-1.c":    {},
		"pr71626-2.c":    {},
		"scal-to-vec1.c": {},
		"scal-to-vec2.c": {},
		"scal-to-vec3.c": {},
		"simd-1.c":       {},
		"simd-2.c":       {},
		"simd-4.c":       {},
		"simd-5.c":       {},
		"simd-6.c":       {},

		// https://goo.gl/XDxJEL
		"20021127-1.c": {},

		// asm
		"20001009-2.c": {},
		"20020107-1.c": {},
		"20030222-1.c": {},
		"20071211-1.c": {},
		"20071220-1.c": {},
		"20071220-2.c": {},
		"960312-1.c":   {},
		"960830-1.c":   {},
		"990130-1.c":   {},
		"990413-2.c":   {},
		"pr38533.c":    {},
		"pr40022.c":    {},
		"pr40657.c":    {},
		"pr41239.c":    {},
		"pr43385.c":    {},
		"pr43560.c":    {},
		"pr45695.c":    {},
		"pr46309.c":    {},
		"pr49279.c":    {},
		"pr49390.c":    {},
		"pr51877.c":    {},
		"pr51933.c":    {},
		"pr52286.c":    {},
		"pr56205.c":    {},
		"pr56866.c":    {},
		"pr56982.c":    {},
		"pr57344-1.c":  {},
		"pr57344-2.c":  {},
		"pr57344-3.c":  {},
		"pr57344-4.c":  {},
		"pr63641.c":    {},
		"pr65053-1.c":  {},
		"pr65053-2.c":  {},
		"pr65648.c":    {},
		"pr65956.c":    {},
		"pr68328.c":    {},
		"pr69320-2.c":  {},
		"stkalign.c":   {},

		// __label__
		"920415-1.c": {},
		"920721-4.c": {},
		"930406-1.c": {},
		"980526-1.c": {},
		"pr51447.c":  {},

		// attribute alias
		"alias-2.c": {},
		"alias-3.c": {},
		"alias-4.c": {},

		// _Alignas
		"pr68532.c": {},

		// Profiling
		"eeprof-1.c": {},

		// 6.5.16/4: The order of evaluation of the operands is unspecified.
		"pr58943.c": {},
	}

	todolist := map[string]struct{}{
		// long double constant out of range for double.
		"960405-1.c": {},

		// case range
		"pr34154.c": {},

		// VLA. Need to resolve https://gitlab.com/cznic/cc/issues/91 first.
		"20040411-1.c":    {},
		"20040423-1.c":    {},
		"20040811-1.c":    {},
		"20041218-2.c":    {},
		"20070919-1.c":    {},
		"920929-1.c":      {},
		"970217-1.c":      {},
		"pr22061-1.c":     {},
		"pr43220.c":       {},
		"vla-dealloc-1.c": {},

		// Initializer
		"20050613-1.c":        {}, // struct B b = { .a.j = 5 };
		"20050929-1.c":        {}, // struct C e = { &(struct B) { &(struct A) { 1, 2 }, &(struct A) { 3, 4 } }, &(struct A) { 5, 6 } };
		"20071029-1.c":        {}, // t = (T) { { ++i, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 } };
		"921019-1.c":          {}, // void *foo[]={(void *)&("X"[0])};
		"991228-1.c":          {}, // cc.Parse: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/991228-1.c:1:51: invalid designator for type double
		"compndlit-1.c":       {}, // x = (struct S) {b:0, a:0, c:({ struct S o = x; o.a == 1 ? 10 : 20;})};
		"const-addr-expr-1.c": {}, // int *Upgd_minor_ID = (int *) &((Upgrade_items + 1)->uaattrid);
		"pr22098-1.c":         {}, // b = (uintptr_t)(p = &(int []){0, 1, 2}[++a]);
		"pr22098-2.c":         {}, // b = (uintptr_t)(p = &(int []){0, 1, 2}[1]);
		"pr22098-3.c":         {}, // b = (uintptr_t)(p = &(int []){0, f(), 2}[1]);
		"pr70460.c":           {}, // static int b[] = { &&lab1 - &&lab0, &&lab2 - &&lab0 };

		// signal.h
		"20101011-1.c": {},

		// &&label expr
		"comp-goto-1.c": {}, // # [100]: Verify (A): mismatched operand type, got int32, expected uint32; simulator_kernel:0x64: 	lsh             	uint32	; ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/comp-goto-1.c:83:40

		// builtins
		"pr47237.c":       {}, // __builtin_apply, __builtin_apply_args
		"pr64006.c":       {}, // __builtin_mul_overflow
		"pr68381.c":       {}, // __builtin_mul_overflow
		"pr71554.c":       {}, // __builtin_mul_overflow
		"va-arg-pack-1.c": {}, // __builtin_va_arg_pack

		// long double
		"pr39228.c": {},

		// un-flatten (wips wrt cc.0506a942f3efa9b7a0a4b98dbe45bf7e8d06a542)
		"20030714-1.c": {}, // cc.Parse: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20030714-1.c:102:11: assignment from incompatible type ('unsigned' = '<undefined>')
		"anon-1.c":     {}, // cc.Parse: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/anon-1.c:22:7: struct{int; ;} has no member named b

		// irgo -------------------------------------------------------

		// statement expression
		"20000703-1.c": {},
		"20001203-2.c": {},
		"20020206-1.c": {},
		"20020320-1.c": {},
		"20000917-1.c": {},

		// computed goto
		"20040302-1.c": {},
		"20041214-1.c": {},
		"20071210-1.c": {},
		"920302-1.c":   {},
		"920501-3.c":   {},
		"920501-4.c":   {},
		"920501-5.c":   {},

		// __builtin_return_address
		"pr17377.c":    {},
		"20010122-1.c": {},

		// setjmp/longjmp
		"pr60003.c": {},

		// alloca
		"20010209-1.c":      {},
		"20020314-1.c":      {},
		"20020412-1.c":      {},
		"20021113-1.c":      {},
		"20040223-1.c":      {},
		"20040308-1.c":      {},
		"20070824-1.c":      {},
		"921017-1.c":        {},
		"941202-1.c":        {},
		"align-nest.c":      {},
		"alloca-1.c":        {},
		"built-in-setjmp.c": {},
		"frame-address.c":   {},
		"pr22061-4.c":       {},
		"pr36321.c":         {},

		// irgo TODOs
		"20010924-1.c":    {}, // irgo.go:1485: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20010924-1.c:33:3: TODO1247 int8:Int8
		"20020810-1.c":    {}, // irgo.go:1546: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20020810-1.c:17:10: Struct
		"20040307-1.c":    {}, // irgo.go:853: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20040307-1.c:16:11
		"20040331-1.c":    {}, // irgo.go:853: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20040331-1.c:10:10
		"20040629-1.c":    {}, // irgo.go:853: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20040629-1.c:124:1
		"20040705-1.c":    {}, // irgo.go:853: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20040629-1.c:124:1
		"20040705-2.c":    {}, // irgo.go:853: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20040629-1.c:124:1
		"20070614-1.c":    {}, // irgo.go:1601: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/20070614-1.c:3:10: *ir.Complex64Value
		"950628-1.c":      {}, // etc.go:600: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/950628-1.c:28:12: *ir.FieldValue
		"950906-1.c":      {}, // etc.go:790: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/950906-1.c:8:3: *ir.Jz 0x00005
		"980602-2.c":      {}, // irgo.go:853: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/980602-2.c:16:11
		"990208-1.c":      {}, // irgo.go:1116: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/990208-1.c:13:13: *ir.Const
		"bitfld-3.c":      {}, // irgo.go:906: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/bitfld-3.c:51:7
		"complex-2.c":     {}, // etc.go:600: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/complex-2.c:16:7: *ir.ConstC128
		"complex-5.c":     {}, // irgo.go:1601: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/complex-5.c:7:20: *ir.Complex64Value
		"complex-7.c":     {}, // irgo.go:1601: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/complex-7.c:5:25: *ir.Complex64Value
		"pr15296.c":       {}, // irgo.go:1392: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr15296.c:65:12: *struct{}, [<nil> 111]
		"pr23324.c":       {}, // irgo.go:1400: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr23324.c:25:3: int64, *ir.Int32Value(1069379046)
		"pr28865.c":       {}, // irgo.go:1400: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr28865.c:3:9: struct{int32,*int8}, *ir.CompositeValue({1, "123456789012345678901234567890"+0})
		"pr30185.c":       {}, // etc.go:600: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr30185.c:21:18: *ir.FieldValue
		"pr33382.c":       {}, // irgo.go:1499: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr33382.c:6:12: TODO1247 int32:Int32
		"pr38051.c":       {},
		"pr42248.c":       {}, // etc.go:600: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr42248.c:23:15: *ir.ConstC128
		"pr42691.c":       {}, // irgo.go:1400: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr42691.c:36:16: [4]uint16, *ir.CompositeValue({0, 0, 0, 32752})
		"pr44164.c":       {}, // New: runtime error: index out of range
		"pr49644.c":       {}, // irgo.go:1601: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr49644.c:8:34: *ir.Complex128Value
		"pr53084.c":       {}, // irgo.go:1596: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr53084.c:15:21
		"pr55750.c":       {}, // irgo.go:853: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr55750.c:14:3
		"pr58209.c":       {},
		"pr68249.c":       {}, // etc.go:420: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr68249.c:11:11: TODO stack(2): 6:			; ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/pr68249.c:11:11
		"stdarg-2.c":      {}, // irgo.go:1145: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/stdarg-2.c:94:17: *ir.Field *struct{}
		"va-arg-13.c":     {}, // irgo.go:1145: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/va-arg-13.c:24:18: *ir.Field *struct{}
		"wchar_t-1.c":     {}, // irgo.go:1601: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/wchar_t-1.c:3:9: *ir.WideStringValue
		"widechar-2.c":    {}, // irgo.go:1601: ../cc/testdata/gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/widechar-2.c:3:15: *ir.WideStringValue
		"zero-struct-1.c": {}, // New: runtime error: index out of range

		// Incorrect conversion
		"20010123-1.c": {}, // 91:109: expected '==', found '='
		"20010518-2.c": {}, // 90:123: expected '==', found '='
		"20020215-1.c": {}, // 86:126: expected '==', found '='
		"20030224-2.c": {}, // 87:125: expected '==', found '='
		"20031215-1.c": {}, // ./main.go:118: cannot use str(0) (type *int8) as type [3]int8 in field value
		"20041212-1.c": {}, // ./main.go:78: cannot convert Xf (type func(*crt.TLS) unsafe.Pointer) to type unsafe.Pointer
		"20050502-2.c": {}, // ./main.go:82: cannot use (*int8)(unsafe.Pointer(&_x)) (type *int8) as type unsafe.Pointer in argument to crt.Xmemcmp
		"20051012-1.c": {}, // ./main.go:79: too many arguments in call to Xfoo

		"20060929-1.c": {}, // ./main.go:125: *postInc_1018(postInc_20561(&_p, 8), 4) evaluated but not used
		"20071202-1.c": {}, // 108:140: expected '==', found '='
		"20080122-1.c": {}, // ./main.go:82: cannot use str(0) (type *int8) as type [32]uint8 in assignment
		"20120919-1.c": {}, // ./main.go:144: cannot use &Xvd (type *[2]float64) as type *float64 in assignment
		"920501-1.c":   {}, // ./main.go:79: too many arguments in call to Xx
		"921202-1.c":   {}, // ./main.go:94: too many arguments in call to Xmpn_random2
		"921208-2.c":   {}, // ./main.go:95: too many arguments in call to Xg

		"930608-1.c":           {}, // ./main.go:85: cannot convert Xa (type [1]func(*crt.TLS, float64) float64) to type unsafe.Pointer
		"930719-1.c":           {}, // ./main.go:112: invalid indirect of nil
		"941014-1.c":           {}, // ./main.go:82: cannot convert Xf (type func(*crt.TLS, int32, int32) int32) to type unsafe.Pointer
		"960416-1.c":           {}, // ./main.go:106: cannot convert u64(4294967296) (type uint64) to type struct { X [0]struct { X0 uint64; X1 struct { X0 uint64; X1 uint64 } }; U [16]byte }
		"980223.c":             {}, // ./main.go:126: cannot use &Xcons1 (type *[2]struct { X0 *int8; X1 int64 }) as type *int8 in field value
		"991201-1.c":           {}, // ./main.go:110: cannot use &Xa_con (type *struct { X0 uint64; X1 [48]uint8 }) as type struct { X0 unsafe.Pointer } in array or slice literal
		"alias-1.c":            {}, // ./main.go:109: cannot use &Xval (type *int32) as type *float32 in assignment
		"bcp-1.c":              {}, // ./main.go:86: cannot convert Xbad_t0 (type [6]func(*crt.TLS) int32) to type unsafe.Pointer
		"bswap-2.c":            {}, // 97:124: expected '==', found '='
		"builtin-bitops-1.c":   {}, // ./main.go:92: undefined: crt.X__builtin_clz
		"builtin-prefetch-2.c": {}, // ./main.go:98: &Xglob_int_arr evaluated but not used
		"builtin-prefetch-3.c": {}, // ./main.go:101: &Xglob_vol_int_arr evaluated but not used
		"builtin-prefetch-4.c": {}, // ./main.go:227: cannot convert uintptr(unsafe.Pointer(&Xarr)) + 80 (type uintptr) to type *int32
		"builtin-prefetch-5.c": {}, // ./main.go:78: (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(unsafe.Pointer(&Xs))) + uintptr(2))) evaluated but not used
		"builtin-prefetch-6.c": {}, // ./main.go:122: *(**int32)(unsafe.Pointer(uintptr(unsafe.Pointer(unsafe.Pointer(&Xbad_addr))) + 8 * uintptr(_i))) evaluated but not used
		"longlong.c":           {}, // ./main.go:124: cannot use &Xb (type *[32]uint64) as type *uint64 in assignment
		"lto-tbaa-1.c":         {}, // ./main.go:112: cannot use &Xb2 (type *struct { X0 *int32 }) as type **int32 in assignment
		"pr43784.c":            {}, // ./main.go:116: cannot use &_v (type *struct { X [0]struct { X0 struct { X0 struct { X0 [256]uint8 }; X1 int32 }; X1 struct { X0 int32; X1 struct { X0 [256]uint8 } } }; U [260]byte }) as type *struct { X0 [256]uint8 } in assignment

		"pr44555.c":       {}, // Needs a strict-semantic option to pass.
		"pr53160.c":       {}, // ./main.go:86: Xb evaluated but not used
		"pr57130.c":       {}, // 89:111: expected '==', found '='
		"pr57281.c":       {}, // ./main.go:86: Xf evaluated but not used
		"pr57568.c":       {}, // ./main.go:98: cannot convert uintptr(unsafe.Pointer(&Xa)) + 128 (type uintptr) to type *int32
		"pr58277-2.c":     {}, // ./main.go:222: Xd evaluated but not used
		"pr66556.c":       {}, // ./main.go:152: *(*int8)(unsafe.Pointer(uintptr(unsafe.Pointer(unsafe.Pointer(&Xe))) + 1 * uintptr(i32(0)))) evaluated but not used
		"pr67037.c":       {}, // ./main.go:119: too many arguments in call to Xextfunc
		"pr69691.c":       {}, // ./main.go:125: undefined: crt.X__builtin_strchr
		"restrict-1.c":    {}, // 102:130: expected '==', found '='
		"struct-ret-1.c":  {}, // ./main.go:132: cannot use str(64) (type *int8) as type [33]int8 in field value
		"va-arg-4.c":      {}, // ./main.go:120: cannot use str(16) (type *int8) as type [32]int8 in field value
		"zero-struct-2.c": {}, // 84:108: expected '==', found '='

		// Compiles to Go but fails
		"stdarg-3.c": {}, // y = va_arg (ap, int); but passed value is a struct (???)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdata, err := filepath.Rel(wd, ccTestdata)
	if err != nil {
		t.Fatal(err)
	}

	var re *regexp.Regexp
	if s := *filter; s != "" {
		re = regexp.MustCompile(s)
	}

	dir := filepath.Join(testdata, filepath.FromSlash("gcc-6.3.0/gcc/testsuite/gcc.c-torture/execute/"))
	expect(
		t,
		dir,
		func(match string) bool {
			base := filepath.Base(match)
			_, skip := blacklist[base]
			if _, skip2 := todolist[base]; skip2 {
				skip = true
			}
			if re != nil {
				skip = !re.MatchString(base)
			}
			return skip
		},
		func(wd, match string) []string {
			return []string{match}
		},
		nil,
		cc.AllowCompatibleTypedefRedefinitions(),
		cc.EnableAlignOf(),
		cc.EnableAlternateKeywords(),
		cc.EnableAnonymousStructFields(),
		cc.EnableAsm(),
		cc.EnableBuiltinClassifyType(),
		cc.EnableBuiltinConstantP(),
		cc.EnableComputedGotos(),
		cc.EnableDefineOmitCommaBeforeDDD(),
		cc.EnableEmptyDeclarations(),
		cc.EnableEmptyStructs(),
		cc.EnableImaginarySuffix(),
		cc.EnableImplicitFuncDef(),
		cc.EnableImplicitIntType(),
		cc.EnableLegacyDesignators(),
		cc.EnableNonConstStaticInitExpressions(),
		cc.EnableOmitConditionalOperand(),
		cc.EnableOmitFuncArgTypes(),
		cc.EnableOmitFuncRetType(),
		cc.EnableParenthesizedCompoundStatemen(),
		cc.EnableTypeOf(),
		cc.EnableUnsignedEnums(),
		cc.EnableWideBitFieldTypes(),
		cc.ErrLimit(-1),
		cc.SysIncludePaths([]string{ccir.LibcIncludePath}),
	)
}

func build(t *testing.T, predef string, tus [][]string, opts ...cc.Opt) ([]byte, error) {
	ndbg := ""
	if *ndebug {
		ndbg = "#define NDEBUG 1"
	}
	var build []*cc.TranslationUnit
	tus = append(tus, []string{ccir.CRT0Path})
	for _, src := range tus {
		model, err := ccir.NewModel()
		if err != nil {
			t.Fatal(err)
		}

		ast, err := cc.Parse(
			fmt.Sprintf(`
%s
#define _CCGO 1
#define __arch__ %s
#define __os__ %s
#include <builtin.h>
%s
`, ndbg, runtime.GOARCH, runtime.GOOS, predef),
			src,
			model,
			append([]cc.Opt{
				cc.AllowCompatibleTypedefRedefinitions(),
				cc.EnableEmptyStructs(),
				cc.EnableImplicitFuncDef(),
				cc.EnableNonConstStaticInitExpressions(),
				cc.ErrLimit(*errLimit),
				cc.SysIncludePaths([]string{ccir.LibcIncludePath}),
			}, opts...)...,
		)
		if err != nil {
			t.Fatal(errStr(err))
		}

		build = append(build, ast)
	}

	var out, src buffer.Bytes

	defer func() {
		out.Close()
		src.Close()
	}()

	if err := New(build, &out, LibcTypes()); err != nil {
		return nil, fmt.Errorf("New: %v", err)
	}

	fmt.Fprintf(&src, prologue, crtQ, out.Bytes())
	b, err := format.Source(src.Bytes())
	if err != nil {
		return src.Bytes(), err
	}

	out.Close()
	src.Close()
	return b, nil
}

func findRepo(t *testing.T, s string) string {
	s = filepath.FromSlash(s)
	for _, v := range strings.Split(strutil.Gopath(), string(os.PathListSeparator)) {
		p := filepath.Join(v, "src", s)
		fi, err := os.Lstat(p)
		if err != nil {
			continue
		}

		if fi.IsDir() {
			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			if p, err = filepath.Rel(wd, p); err != nil {
				t.Fatal(err)
			}

			return p
		}
	}
	return ""
}

type file struct {
	name string
	data []byte
}

func (f file) String() string { return fmt.Sprintf("%v %v", len(f.data), f.name) }

func run(t *testing.T, src []byte, argv []string, inputFiles []file, errOK bool) (output []byte, resultFiles []file, duration time.Duration) {
	dir, err := ioutil.TempDir("", "ccgo-test-")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	defer os.Chdir(cwd)

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("main.go", src, 0600); err != nil {
		t.Fatal(err)
	}

	for _, v := range inputFiles {
		if err := ioutil.WriteFile(v.name, v.data, 0600); err != nil {
			t.Fatal(err)
		}
	}

	var stdout, stderr buffer.Bytes

	defer func() {
		stdout.Close()
		stderr.Close()
	}()

	cmd := exec.Command("go", append([]string{"run", "main.go"}, argv[1:]...)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	t0 := time.Now()
	err = cmd.Run()
	duration = time.Since(t0)
	if err != nil && !errOK {
		var log bytes.Buffer
		if b := stdout.Bytes(); b != nil {
			fmt.Fprintf(&log, "stdout:\n%s\n", b)
		}
		if b := stderr.Bytes(); b != nil {
			fmt.Fprintf(&log, "stderr:\n%s\n", b)
		}
		t.Fatalf("err %v\n%s", err, log.Bytes())
	}

	glob, err := filepath.Glob("*")
	if err != nil {
		t.Fatal(err)
	}

	for _, m := range glob {
		data, err := ioutil.ReadFile(m)
		if err != nil {
			t.Fatal(err)
		}

		resultFiles = append(resultFiles, file{m, data})
	}

	return bytes.TrimSpace(append(stdout.Bytes(), stderr.Bytes()...)), resultFiles, duration
}

func TestSelfie(t *testing.T) {
	const repo = "github.com/cksystemsteaching/selfie"
	pth := findRepo(t, repo)
	if pth == "" {
		t.Logf("repository not found, skipping: %v", repo)
		return
	}

	src, err := build(t, "", [][]string{{filepath.Join(pth, "selfie.c")}})
	if err != nil {
		t.Fatal(err)
	}

	return //TODO Fails on 32 bit.

	if m, _ := ccir.NewModel(); m.Items[cc.Ptr].Size != 4 {
		return
	}

	args := []string{"./selfie"}
	out, _, d := run(t, src, args, nil, false)
	if g, e := out, []byte("./selfie: usage: selfie { -c { source } | -o binary | -s assembly | -l binary } [ ( -m | -d | -y | -min | -mob ) size ... ]"); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	t.Logf("%s\n%s\n%v", args, out, d)

	args = []string{"./selfie", "-c", "hello.c", "-m", "1"}
	out, _, d = run(t, src, args, []file{{"hello.c", []byte(`
int *foo;

int main() {
	foo = "Hello world!";
	while (*foo!=0) { 
		write(1, foo, 4);
		foo = foo + 1;
	}
	*foo = 10;
	write(1, foo, 1);
}
`)}}, false)
	if g, e := out, []byte(`./selfie: this is selfie's starc compiling hello.c
./selfie: 141 characters read in 12 lines and 0 comments
./selfie: with 102(72.46%) characters in 52 actual symbols
./selfie: 1 global variables, 1 procedures, 1 string literals
./selfie: 2 calls, 3 assignments, 1 while, 0 if, 0 return
./selfie: 660 bytes generated with 159 instructions and 24 bytes of data
./selfie: this is selfie's mipster executing hello.c with 1MB of physical memory
Hello world!
hello.c: exiting with exit code 0 and 0.00MB of mallocated memory
./selfie: this is selfie's mipster terminating hello.c with exit code 0 and 0.01MB of mapped memory
./selfie: profile: total,max(ratio%)@addr(line#),2max(ratio%)@addr(line#),3max(ratio%)@addr(line#)
./selfie: calls: 5,4(80.00%)@0x88(~1),1(20.00%)@0x17C(~5),0(0.00%)
./selfie: loops: 3,3(100.00%)@0x198(~6),0(0.00%),0(0.00%)
./selfie: loads: 32,4(12.50%)@0x88(~1),3(9.38%)@0x1D4(~7),1(3.12%)@0x24(~1)
./selfie: stores: 20,3(15.01%)@0x1D0(~7),1(5.00%)@0x4C(~1),0(0.00%)`); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	t.Logf("%s\n%s\n%v", args, out, d)

	selfie, err := ioutil.ReadFile(filepath.Join(pth, "selfie.c"))
	if err != nil {
		t.Fatal(err)
	}

	args = []string{"./selfie", "-c", "selfie.c"}
	out, _, d = run(t, src, args, []file{{"selfie.c", selfie}}, false)
	if g, e := out, []byte(`./selfie: this is selfie's starc compiling selfie.c
./selfie: 176362 characters read in 7086 lines and 970 comments
./selfie: with 97764(55.55%) characters in 28916 actual symbols
./selfie: 260 global variables, 290 procedures, 450 string literals
./selfie: 1960 calls, 722 assignments, 57 while, 571 if, 241 return
./selfie: 121676 bytes generated with 28783 instructions and 6544 bytes of data`); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	t.Logf("%s\n%s\n%v", args, out, d)
}

func TestSQLite(t *testing.T) {
	const repo = "sqlite.org/sqlite-amalgamation-3190300/"
	pth := findRepo(t, repo)
	if pth == "" {
		t.Logf("repository not found, skipping: %v", repo)
		return
	}

	src, err := build(
		t,
		`
		#define HAVE_MALLOC_H 1
		#define HAVE_MALLOC_USABLE_SIZE 1
		#define SQLITE_DEBUG 1
		#define SQLITE_ENABLE_API_ARMOR 1
		#define SQLITE_WITHOUT_MSIZE 1
		`,
		[][]string{
			{"testdata/sqlite/test.c"},
			{filepath.Join(pth, "sqlite3.c")},
		},
		cc.EnableAnonymousStructFields(),
		cc.EnableWideBitFieldTypes(),
		cc.IncludePaths([]string{pth}),
	)
	if *oLog {
		t.Logf("\n%s", src)
	}
	if err != nil {
		t.Fatal(err)
	}

	args := []string{"./test"}
	out, f, d := run(t, src, args, nil, true)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte(`Usage: ./test DATABASE SQL-STATEMENT
exit status 1`); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	args = []string{"./test", "foo"}
	out, f, d = run(t, src, args, nil, true)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte(`Usage: ./test DATABASE SQL-STATEMENT
exit status 1`); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	args = []string{"./test", "foo", "bar"}
	out, f, d = run(t, src, args, nil, true)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte(`FAIL (1) near "bar": syntax error
SQL error: near "bar": syntax error`); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	args = []string{"./test", "foo", "select * from t"}
	out, f, d = run(t, src, args, nil, false)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte(`FAIL (1) no such table: t
SQL error: no such table: t`); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	args = []string{"./test", "foo", "select name from sqlite_master where type='table'"}
	out, f, d = run(t, src, args, nil, false)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte(""); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	args = []string{"./test", "foo", "create table t(i int)"}
	out, f, d = run(t, src, args, nil, false)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte(""); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	args = []string{"./test", "foo", `
		create table t(i int);
		select name from sqlite_master where type='table';
		`}
	out, f, d = run(t, src, args, nil, false)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte("name = t"); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}

	args = []string{"./test", "foo", `
		create table t(i int);
		select name from sqlite_master where type='table';
		insert into t values(42), (314);
		select * from t order by i asc;
		select * from t order by i desc;
		`}
	out, f, d = run(t, src, args, nil, false)
	t.Logf("%q\n%s\n%v\n%v", args, out, d, f)
	if g, e := out, []byte(`name = t
i = 42
i = 314
i = 314
i = 42`); !bytes.Equal(g, e) {
		t.Fatalf("\ngot\n%s\nexp\n%s", g, e)
	}
}

func TestOther(t *testing.T) {
	var re *regexp.Regexp
	if s := *filter; s != "" {
		re = regexp.MustCompile(s)
	}

	expect(
		t,
		"testdata",
		func(match string) bool {
			if re != nil && !re.MatchString(filepath.Base(match)) {
				return true
			}

			return false
		},
		func(wd, match string) []string {
			return []string{match}
		},
		nil,
		cc.EnableEmptyStructs(),
		cc.EnableImplicitFuncDef(),
		cc.ErrLimit(-1),
		cc.SysIncludePaths([]string{ccir.LibcIncludePath}),
	)
}
