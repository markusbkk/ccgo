// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

//TODO CSmith

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/pmezard/go-difflib/difflib"
	"modernc.org/cc/v4"
	"modernc.org/ccorpus2"
	"modernc.org/fileutil"
	"modernc.org/gc/v2"
)

var (
	oDebug      = flag.Bool("debug", false, "")
	oErr1       = flag.Bool("err1", false, "first error line only")
	oKeep       = flag.Bool("keep", false, "keep temp directories")
	oPanic      = flag.Bool("panic", false, "panic on miscompilation")
	oShellTime  = flag.Duration("shelltimeout", 5*time.Minute, "shell() time limit")
	oStackTrace = flag.Bool("trcstack", false, "")
	oTrace      = flag.Bool("trc", false, "Print tested paths.")
	oTraceF     = flag.Bool("trcf", false, "Print test file content")
	oTraceO     = flag.Bool("trco", false, "Print test output")
	oXTags      = flag.String("xtags", "", "passed to go build of TestSQLite")

	cfs    = ccorpus2.FS
	goarch = runtime.GOARCH
	goos   = runtime.GOOS
	re     *regexp.Regexp
	hostCC string
)

func TestMain(m *testing.M) {
	extendedErrors = true
	gc.ExtendedErrors = true
	oRE := flag.String("re", "", "")
	flag.Parse()
	if *oRE != "" {
		re = regexp.MustCompile(*oRE)
	}
	cfg, err := cc.NewConfig(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		panic(err)
	}

	hostCC = cfg.CC
	os.Exit(m.Run())
}

func (p *parallel) close(t *testing.T) {
	p.wg.Wait()
	p.Lock()
	for _, v := range p.errors {
		t.Error(v)
	}
	p.Unlock()
	t.Logf("TOTAL: files %5s, skip %5s, ok %5s, fails %5s: %s", h(p.files), h(p.skips), h(p.oks), h(p.fails), p.resultTag)
}

func h(v interface{}) string {
	switch x := v.(type) {
	case int32:
		return humanize.Comma(int64(x))
	case int64:
		return humanize.Comma(x)
	case uint64:
		if x <= math.MaxInt64 {
			return humanize.Comma(int64(x))
		}
	}
	return fmt.Sprint(v)
}

func cfsWalk(dir string, f func(pth string, fi os.FileInfo) error) error {
	fis, err := cfs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, v := range fis {
		switch {
		case v.IsDir():
			if err = cfsWalk(dir+"/"+v.Name(), f); err != nil {
				return err
			}
		default:
			fi, err := v.Info()
			if err != nil {
				return err
			}

			if err = f(dir+"/"+v.Name(), fi); err != nil {
				return err
			}
		}
	}
	return nil
}

func TestSep(t *testing.T) {
	for i, v := range []struct {
		src         string
		sep         string
		trailingSep string
	}{
		{"int f() {}", "", "\n"},
		{" int f() {}\n", " ", "\n"},
		{"\nint f() {}\n", "\n", "\n"},
		{"/*A*//*B*/int f() {}\n", "/*A*//*B*/", "\n"},
		{"/*A*//*B*/ int f() {}\n", "/*A*//*B*/ ", "\n"},

		{"/*A*//*B*/\nint f() {}\n", "/*A*//*B*/\n", "\n"},
		{"/*A*/ /*B*/int f() {}\n", "/*A*/ /*B*/", "\n"},
		{"/*A*/ /*B*/ int f() {}\n", "/*A*/ /*B*/ ", "\n"},
		{"/*A*/ /*B*/\nint f() {}\n", "/*A*/ /*B*/\n", "\n"},
		{"/*A*/\n/*B*/int f() {}\n", "/*A*/\n/*B*/", "\n"},

		{"/*A*/\n/*B*/ int f() {}\n", "/*A*/\n/*B*/ ", "\n"},
		{"/*A*/\n/*B*/\nint f() {}\n", "/*A*/\n/*B*/\n", "\n"},
		{" /*A*/ /*B*/int f() {}\n", " /*A*/ /*B*/", "\n"},
		{" /*A*/ /*B*/ int f() {}\n", " /*A*/ /*B*/ ", "\n"},
		{" /*A*/ /*B*/\nint f() {}\n", " /*A*/ /*B*/\n", "\n"},

		{" /*A*/\n/*B*/int f() {}\n", " /*A*/\n/*B*/", "\n"},
		{" /*A*/\n/*B*/ int f() {}\n", " /*A*/\n/*B*/ ", "\n"},
		{" /*A*/\n/*B*/\nint f() {}\n", " /*A*/\n/*B*/\n", "\n"},
		{"\n/*A*/ /*B*/int f() {}\n", "\n/*A*/ /*B*/", "\n"},
		{"\n/*A*/ /*B*/ int f() {}\n", "\n/*A*/ /*B*/ ", "\n"},

		{"\n/*A*/ /*B*/\nint f() {}\n", "\n/*A*/ /*B*/\n", "\n"},
		{"\n/*A*/\n/*B*/int f() {}\n", "\n/*A*/\n/*B*/", "\n"},
		{"\n/*A*/\n/*B*/ int f() {}\n", "\n/*A*/\n/*B*/ ", "\n"},
		{"\n/*A*/\n/*B*/\nint f() {}\n", "\n/*A*/\n/*B*/\n", "\n"},
	} {
		ast, err := cc.Parse(
			&cc.Config{},
			[]cc.Source{{Name: "test", Value: v.src + "int __predefined_declarator;"}},
		)
		if err != nil {
			t.Errorf("%v: %v", i, err)
			continue
		}

		t.Logf("%q -> %q", v.src, nodeSource(ast.TranslationUnit))
		var tok cc.Token
		firstToken(ast.TranslationUnit, &tok)
		if g, e := string(tok.Sep()), v.sep; g != e {
			t.Errorf("%v: %q %q", i, g, e)
		}
		if g, e := string(ast.EOF.Sep()), v.trailingSep; g != e {
			t.Errorf("%v: %q %q", i, g, e)
		}
	}
}

func inDir(dir string, f func() error) (err error) {
	var cwd string
	if cwd, err = os.Getwd(); err != nil {
		return err
	}

	defer func() {
		if err2 := os.Chdir(cwd); err2 != nil {
			err = err2
		}
	}()

	if err = os.Chdir(filepath.FromSlash(dir)); err != nil {
		return err
	}

	return f()
}

func absCwd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if wd, err = filepath.Abs(wd); err != nil {
		return "", err
	}

	return wd, nil
}

type echoWriter struct {
	w      bytes.Buffer
	silent bool
}

func (w *echoWriter) Write(b []byte) (int, error) {
	if !w.silent {
		os.Stderr.Write(b)
	}
	return w.w.Write(b)
}

func TestExec(t *testing.T) {
	g := newGolden(t, fmt.Sprintf("testdata/test_exec_%s_%s.golden", runtime.GOOS, runtime.GOARCH))

	defer g.close()

	tmp := t.TempDir()
	if err := inDir(tmp, func() error {
		if out, err := shell(true, "go", "mod", "init", "test"); err != nil {
			return fmt.Errorf("%s\vFAIL: %v", out, err)
		}

		if out, err := shell(true, "go", "get", defaultLibc); err != nil {
			return fmt.Errorf("%s\vFAIL: %v", out, err)
		}

		for _, v := range []struct {
			path string
			exec bool
		}{
			{"CompCert-3.6/test/c", true},
			{"benchmarksgame-team.pages.debian.net", true},
			{"ccgo", true},
			{"gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile", false},
			{"gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute", true},
			{"github.com/AbsInt/CompCert/test/c", true},
			{"github.com/cxgo", true},
			{"github.com/gcc-mirror/gcc/gcc/testsuite", true},
			{"github.com/vnmakarov", true},
			{"tcc-0.9.27/tests/tests2", true},
		} {
			t.Run(v.path, func(t *testing.T) {
				testExec(t, "assets/"+v.path, v.exec, g)
			})
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func testExec(t *testing.T, cfsDir string, exec bool, g *golden) {
	const isolated = "x"
	os.RemoveAll(isolated)
	if err := os.Mkdir(isolated, 0770); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(isolated); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.Chdir(".."); err != nil {
			t.Fatal(err)
		}
	}()

	files, bytes, err := fileutil.CopyDir(cfs, "", cfsDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s: copied %v files, %v bytes", cfsDir, h(files), h(bytes))

	p := newParallel(cfsDir)

	defer func() { p.close(t) }()

	p.err(filepath.Walk(".", func(path string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".c" {
			return nil
		}

		p.file()
		switch {
		case re != nil && !re.MatchString(filepath.Base(path)):
			p.skip()
			return nil
		}

		id := p.id()
		if *oTrace {
			fmt.Fprintln(os.Stderr, filepath.Join(cfsDir, path), id)
		}
		p.exec(func() error { return testExec1(t, p, cfsDir, path, exec, g, id) })
		return nil
	}))
}

func testExec1(t *testing.T, p *parallel, root, path string, exec bool, g *golden, id int) error {
	args, err := getArgs(path)
	fullPath := filepath.Join(root, path)
	var cCompilerFailed, cExecFailed bool
	ofn := fmt.Sprint(id)
	bin := "cbin_" + enforceBinaryExt(ofn)
	if _, err = shell(false, hostCC, "-o", bin, "-w", path, "-lm"); err != nil {
		p.skip()
		return nil
		cCompilerFailed = true
	}

	defer os.Remove(ofn)

	var cOut []byte
	if exec && !cCompilerFailed {
		if cOut, err = shell(false, "./"+bin, args...); err != nil {
			p.skip()
			return nil
			cExecFailed = true
		}
	}

	ofn += ".go"

	defer os.Remove(ofn)

	var out bytes.Buffer
	if err := NewTask(goos, goarch, []string{"ccgo", "-o", ofn, "--prefix-field=F", path}, &out, &out, nil).Main(); err != nil {
		if cCompilerFailed || isTestExecKnownFail(fullPath) {
			p.skip()
			return nil
		}

		trc("`%s`: {}, // COMPILE FAIL", fullPath)
		p.fail()
		return errorf("%s: %s: FAIL: %v", fullPath, out.Bytes(), firstError(err, *oErr1))
	}

	bin = "gobin_" + enforceBinaryExt(ofn)
	if _, err = shell(false, "go", "build", "-o", bin, ofn); err != nil {
		if isTestExecKnownFail(fullPath) {
			p.skip()
			return nil
		}

		trc("`%s`: {}, // BUILD FAIL", fullPath)
		p.fail()
		return firstError(err, *oErr1)
	}

	if !exec {
		p.ok()
		g.w("%s\n", fullPath)
		return nil
	}

	if runtime.GOOS != "windows" {
		bin = "./" + bin
	}
	goOut, err := shell(false, bin, args...)
	if err != nil {
		if cExecFailed || isTestExecKnownFail(fullPath) {
			p.skip()
			return nil
		}

		err := errorf("%s: %s: FAIL: %v", fullPath, goOut, err)
		if *oPanic {
			panic(err)
		}

		trc("`%s`: {}, // EXEC FAIL", fullPath)
		p.fail()
		return firstError(err, *oErr1)
	}

	cOut = bytes.TrimSpace(cOut)
	goOut = bytes.TrimSpace(goOut)
	if bytes.Contains(cOut, []byte("\r\n")) {
		cOut = bytes.ReplaceAll(cOut, []byte("\r\n"), []byte{'\n'})
	}
	if bytes.Contains(goOut, []byte("\r\n")) {
		goOut = bytes.ReplaceAll(goOut, []byte("\r\n"), []byte{'\n'})
	}
	if cCompilerFailed || cExecFailed || bytes.Equal(cOut, goOut) {
		p.ok()
		g.w("%s\n", fullPath)
		return nil
	}

	if isTestExecKnownFail(fullPath) {
		p.skip()
		return nil
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(cOut)),
		B:        difflib.SplitLines(string(goOut)),
		FromFile: "expected",
		ToFile:   "got",
		Context:  0,
	}
	s, _ := difflib.GetUnifiedDiffString(diff)
	err = errorf("%v: output differs\n%v\n--- expexted\n%s\n\n--- got\n%s\n\n--- expected\n%s\n--- got\n%s", path, s, cOut, goOut, hex.Dump(cOut), hex.Dump(goOut))
	if *oPanic {
		panic(err)
	}

	trc("`%s`: {}, // EXEC FAIL", fullPath)
	p.fail()
	return firstError(err, *oErr1)
}

func isTestExecKnownFail(s string) (r bool) {
	_, r = testExecKnownFails[s]
	return r
}

var testExecKnownFails = map[string]struct{}{
	// Won't fix
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20021127-1.c`:                 {}, // EXEC FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20031003-1.c`:                 {}, // EXEC FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20101011-1.c`:                 {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20021127-1.c`: {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20031003-1.c`: {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20101011-1.c`: {}, // EXEC FAIL

	//TODO void func(void) __attribute__((aligned(256)));
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/align-3.c`:                 {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/align-3.c`: {}, // EXEC FAIL

	//TODO
	// struct s1
	// {
	//   int __attribute__ ((aligned (8))) a;
	// };
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr23467.c`:                 {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr23467.c`: {}, // EXEC FAIL

	//TODO flexible array member
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr41395-1.c`:                 {}, // EXEC FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr41395-2.c`:                 {}, // EXEC FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr41463.c`:                   {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr41395-1.c`: {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr41395-2.c`: {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr41463.c`:   {}, // EXEC FAIL

	//TODO missed VLA check
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr41935.c`:                 {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr41935.c`: {}, // EXEC FAIL

	// Other
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pushpop_macro.c`:                 {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pushpop_macro.c`: {}, // EXEC FAIL
	`assets/tcc-0.9.27/tests/tests2/77_push_pop_macro.c`:                                   {}, // EXEC FAIL

	// Miscompiled
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-21.c`:                 {}, // EXEC FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-21.c`: {}, // EXEC FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/do-continue.c`:                       {}, // EXEC FAIL + hangs
	`assets/github.com/vnmakarov/mir/c-tests/lacc/function-pointer.c`:                  {}, // EXEC FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/issue23.c`:                            {}, // EXEC FAIL
	`assets/tcc-0.9.27/tests/tests2/75_array_in_struct_init.c`:                         {}, // EXEC FAIL initializer

	//TODO libc: printf.go:586:parseLengthModifier: TODOTODO
	`assets/github.com/vnmakarov/mir/c-tests/lacc/long-double-load.c`: {}, // EXEC FAIL

	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr69989-2.c`:                  {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr78574.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030909-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040704-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080222-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20111208-1.c`:                 {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-6.c`:                   {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920721-2.c`:                   {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950221-1.c`:                   {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960405-1.c`:                   {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr17078-1.c`:                  {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr34176.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38051.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43269.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr48571-1.c`:                  {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr77766.c`:                    {}, // BUILD FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-5.c`:                   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030909-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040704-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080222-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20111208-1.c`: {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-6.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920721-2.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950221-1.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960405-1.c`:   {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr17078-1.c`:  {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr34176.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38051.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43269.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr48571-1.c`:  {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr77766.c`:    {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93744-1.c`:  {}, // BUILD FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-5.c`:   {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0010-goto1.c`:             {}, // BUILD FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0022-namespaces1.c`:       {}, // BUILD FAIL
	`assets/tcc-0.9.27/tests/tests2/54_goto.c`:                                          {}, // BUILD FAIL
	`assets/tcc-0.9.27/tests/tests2/78_vla_label.c`:                                     {}, // BUILD FAIL

	// longjmp/setjmp
	`assets/github.com/vnmakarov/mir/c-benchmarks/except.c`: {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/setjmp.c`:  {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/setjmp2.c`: {}, // COMPILE FAIL

	`assets/CompCert-3.6/test/c/aes.c`:                                                                  {}, // COMPILE FAIL
	`assets/CompCert-3.6/test/c/fannkuch.c`:                                                             {}, // COMPILE FAIL
	`assets/CompCert-3.6/test/c/fftw.c`:                                                                 {}, // COMPILE FAIL
	`assets/CompCert-3.6/test/c/sha3.c`:                                                                 {}, // COMPILE FAIL
	`assets/CompCert-3.6/test/c/vmach.c`:                                                                {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fannkuchredux-5.c`:                                     {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fannkuchredux.c`:                                       {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-2.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-4.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-5.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-7.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-8.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/fasta-9.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-2.c`:                                        {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-4.c`:                                        {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-8.c`:                                        {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/mandelbrot-9.c`:                                        {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/nbody-4.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/nbody-7.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/nbody-8.c`:                                             {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/reverse-complement-6.c`:                                {}, // COMPILE FAIL
	`assets/benchmarksgame-team.pages.debian.net/spectral-norm.c`:                                       {}, // COMPILE FAIL
	`assets/ccgo/bug/arr.c`:                                                                             {}, // COMPILE FAIL
	`assets/ccgo/bug/bitfield.c`:                                                                        {}, // COMPILE FAIL
	`assets/ccgo/bug/condfn.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/csmith.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/csmith2.c`:                                                                         {}, // COMPILE FAIL
	`assets/ccgo/bug/dereffp.c`:                                                                         {}, // COMPILE FAIL
	`assets/ccgo/bug/enums.c`:                                                                           {}, // COMPILE FAIL
	`assets/ccgo/bug/for.c`:                                                                             {}, // COMPILE FAIL
	`assets/ccgo/bug/for2.c`:                                                                            {}, // COMPILE FAIL
	`assets/ccgo/bug/for3.c`:                                                                            {}, // COMPILE FAIL
	`assets/ccgo/bug/fp.c`:                                                                              {}, // COMPILE FAIL
	`assets/ccgo/bug/incfp.c`:                                                                           {}, // COMPILE FAIL
	`assets/ccgo/bug/incfp2.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/init.c`:                                                                            {}, // COMPILE FAIL
	`assets/ccgo/bug/init2.c`:                                                                           {}, // COMPILE FAIL
	`assets/ccgo/bug/init3.c`:                                                                           {}, // COMPILE FAIL
	`assets/ccgo/bug/init4.c`:                                                                           {}, // COMPILE FAIL
	`assets/ccgo/bug/objv.c`:                                                                            {}, // COMPILE FAIL
	`assets/ccgo/bug/select.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/sizeof2.c`:                                                                         {}, // COMPILE FAIL
	`assets/ccgo/bug/sqlite.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/struct.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/struct2.c`:                                                                         {}, // COMPILE FAIL
	`assets/ccgo/bug/union.c`:                                                                           {}, // COMPILE FAIL
	`assets/ccgo/bug/union2.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/union3.c`:                                                                          {}, // COMPILE FAIL
	`assets/ccgo/bug/union4.c`:                                                                          {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20000105-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20010605-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20011217-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20021108-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/20101217-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/386.c`:                                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/961203-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/981006-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/991229-3.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/BUG12.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/ac.c`:                                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/band.c`:                                       {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/bcopy.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/bt386.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/builtin_constant_p.c`:                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/buns.c`:                                       {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/call.c`:                                       {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/dic.c`:                                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/ex.c`:                                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/foo.c`:                                        {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/loop386.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/memtst.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr21728.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr27863.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr37056.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr38360.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr43417.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr44246.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr53409.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr82052.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr82564.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/pr84136.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/r1.c`:                                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/shft.c`:                                       {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/simd-5.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/sound.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/subcc.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/compile/v.c`:                                          {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000113-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000217-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000519-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000703-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000722-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000801-3.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000815-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000822-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000910-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000914-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20000917-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20001009-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20001101.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20001203-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010122-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010123-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010209-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010518-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010605-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010605-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010904-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20010904-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20011113-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020206-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020206-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020215-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020314-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020320-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020404-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020411-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020412-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020418-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20020611-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20021113-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20021118-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030109-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030222-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030224-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030323-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030330-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030401-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030501-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030714-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030811-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030910-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20030916-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20031201-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20031211-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20031211-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040223-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040302-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040307-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040308-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040331-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040411-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040423-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040520-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040629-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040705-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040705-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040709-3.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20040811-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041124-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041201-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041214-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20041218-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050106-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050121-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050203-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050316-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050316-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050316-3.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050604-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050607-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050613-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20050929-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20051012-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20051113-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20060420-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20061031-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20061220-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070614-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070824-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20070919-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071029-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071120-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071202-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071210-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071211-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071220-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20071220-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080122-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080502-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080519-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20080529-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20081117-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090113-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090113-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090113-3.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20090219-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20100316-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20120111-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20141107-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20180921-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/20181120-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920302-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920415-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920428-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-3.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-4.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-5.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920501-7.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920612-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920625-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920721-4.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920728-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920731-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920908-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920908-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/920929-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921016-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921017-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921202-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921204-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921208-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/921215-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930126-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930406-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930621-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930630-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930718-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/930930-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931002-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-10.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-12.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-14.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-4.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-6.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931004-8.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931031-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/931110-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/941202-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950512-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950628-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/950906-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960301-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960312-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960416-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/960608-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/970217-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/980526-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/980602-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/980929-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990130-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990208-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990222-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990326-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990413-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/990525-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991014-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991112-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/991118-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-2.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-3.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alias-4.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/align-1.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/align-nest.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/alloca-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/anon-1.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf-layout-1.c`:                                {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf-pack-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf-sign-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf-sign-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bf64-1.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-3.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-4.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-5.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-6.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bitfld-7.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/bswap-2.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/built-in-setjmp.c`:                            {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/builtin-bitops-1.c`:                           {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/builtin-prefetch-4.c`:                         {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/builtin-types-compatible-p.c`:                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/call-trap-1.c`:                                {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/comp-goto-1.c`:                                {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/comp-goto-2.c`:                                {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-5.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-6.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/complex-7.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/compndlit-1.c`:                                {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/divconst-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/extzvsi.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ffs-1.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ffs-2.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/fprintf-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/frame-address.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/compare-fp-1.c`:                          {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4f.c`:                             {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4l.c`:                             {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-5.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8f.c`:                             {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8l.c`:                             {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/mzero4.c`:                                {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr38016.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr50310.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr72824-2.c`:                             {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/pr84235.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ieee/unsafe-fp-assoc-1.c`:                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/loop-15.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/lto-tbaa-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/medce-1.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/memchr-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/memcpy-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/memset-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/memset-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/memset-3.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nest-align-1.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nest-stdar-1.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-1.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-2.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-3.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-5.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-6.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/nestfunc-7.c`:                                 {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr17377.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr19449.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr19687.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr19689.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-3.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22061-4.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22098-3.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22141-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr22141-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr23135.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr23324.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr24135.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr28289.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr28982b.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr30185.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr30778.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr31136.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr31169.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr31448-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr31448.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr32244-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr34154.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr34768-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr34768-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr34971.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr35456.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr36038.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr36321.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr37573.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr37780.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr37882.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38151.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38422.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38533.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr38969.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39100.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39228.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr39339.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr40022.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr40404.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr40493.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr40657.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr41239.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr42570.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr42614.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43220.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43385.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43560.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43783.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr43987.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr44164.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr44555.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr44575.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr44852.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr45695.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr46309.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr47148.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr47155.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr47237.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr47925.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr48973-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr48973-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49073.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49123.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49218.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49279.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49390.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49644.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr49768.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51447.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51581-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51581-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51877.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr51933.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52209.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52286.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52979-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr52979-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr53645-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr53645.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr54471.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr55750.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56205.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56837.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56866.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr56982.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57130.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-3.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57344-4.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57568.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57861.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57876.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr57877.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58277-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58277-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58385.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58419.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58431.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58564.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58570.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58726.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58831.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr58984.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr59388.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60003.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60017.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr60960.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr61375.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr61725.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr63302.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr63641.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64006.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64242.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr64756.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65053-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65053-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65170.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65215-3.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65215-4.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65369.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65427.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65648.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr65956.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr66556.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr67037.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr67714.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68185.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68249.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68250.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68321.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68328.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68381.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68506.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr68532.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr69320-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr69320-4.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr69691.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70127.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70460.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70566.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70586.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70602.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr70903.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71083.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71494.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71554.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71626-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71626-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr71700.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr77767.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78170.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78436.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78438.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78477.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78559.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78675.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr78726.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79286.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79354.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79737-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr79737-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr80421.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr80692.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr81423.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr81555.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr81556.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr81588.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82192.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82210.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82387.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82524.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr82954.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr83362.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr83383.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84169.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84478.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84524.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr84748.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85095.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85156.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85169.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85331.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85529-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85582-2.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85582-3.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr85756.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr86492.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr86528.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr87053.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr88714.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr88739.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr88904.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr89195.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr89369.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/pr89434.c`:                                    {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/printf-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/scal-to-vec1.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/scal-to-vec2.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/scal-to-vec3.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-1.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-2.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-4.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-5.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/simd-6.c`:                                     {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/ssad-run.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stdarg-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stdarg-3.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/stkalign.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strcpy-1.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strcpy-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strct-stdarg-1.c`:                             {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strct-varg-1.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/string-opt-18.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-3.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-4.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-6.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/strlen-7.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/struct-ini-1.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/struct-ini-2.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/struct-ini-3.c`:                               {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/usad-run.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/user-printf.c`:                                {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-10.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-14.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-15.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-2.c`:                                   {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-22.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/va-arg-pack-1.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/vla-dealloc-1.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/vrp-7.c`:                                      {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zero-struct-1.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zero-struct-2.c`:                              {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zerolen-1.c`:                                  {}, // COMPILE FAIL
	`assets/gcc-9.1.0/gcc/testsuite/gcc.c-torture/execute/zerolen-2.c`:                                  {}, // COMPILE FAIL
	`assets/github.com/AbsInt/CompCert/test/c/aes.c`:                                                    {}, // COMPILE FAIL
	`assets/github.com/AbsInt/CompCert/test/c/fannkuch.c`:                                               {}, // COMPILE FAIL
	`assets/github.com/AbsInt/CompCert/test/c/fftw.c`:                                                   {}, // COMPILE FAIL
	`assets/github.com/AbsInt/CompCert/test/c/sha3.c`:                                                   {}, // COMPILE FAIL
	`assets/github.com/AbsInt/CompCert/test/c/vmach.c`:                                                  {}, // COMPILE FAIL
	`assets/github.com/cxgo/main no args.c`:                                                             {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000113-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000217-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000519-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000703-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000722-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000801-3.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000815-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000822-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000910-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000914-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20000917-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20001009-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20001101.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20001203-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010122-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010123-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010209-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010518-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010605-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010605-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010904-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20010904-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20011113-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020206-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020206-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020215-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020314-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020320-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020404-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020411-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020412-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020418-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20020611-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20021113-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20021118-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030109-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030222-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030224-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030323-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030330-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030401-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030501-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030714-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030811-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030910-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20030916-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20031201-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20031211-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20031211-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040223-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040302-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040307-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040308-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040331-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040411-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040423-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040520-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040629-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040705-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040705-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040709-3.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20040811-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041124-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041201-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041214-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20041218-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050106-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050121-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050203-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050316-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050316-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050316-3.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050604-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050607-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050613-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20050929-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20051012-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20051113-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20060420-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20061031-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20061220-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070614-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070824-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20070919-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071029-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071120-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071202-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071210-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071211-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071220-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20071220-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080122-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080502-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080519-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20080529-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20081117-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090113-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090113-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090113-3.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20090219-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20100316-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20120111-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20141107-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20180921-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20181120-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20190820-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/20190901-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920302-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920415-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920428-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-3.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-4.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-5.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920501-7.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920612-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920625-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920721-4.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920728-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920731-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920908-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920908-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/920929-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921016-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921017-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921202-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921204-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921208-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/921215-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930126-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930406-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930621-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930630-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930718-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/930930-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931002-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-10.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-12.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-14.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-4.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-6.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931004-8.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931031-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/931110-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/941202-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950512-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950628-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/950906-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960301-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960312-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960416-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/960608-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/970217-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/980526-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/980602-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/980929-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990130-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990208-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990222-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990326-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990413-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/990525-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991014-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991112-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/991118-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-2.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-3.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-4.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alias-access-path-2.c`:        {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/align-1.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/align-nest.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/alloca-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/anon-1.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf-layout-1.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf-pack-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf-sign-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf-sign-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bf64-1.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-3.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-4.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-5.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-6.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-7.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-8.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bitfld-9.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bswap-2.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/bswap-3.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/built-in-setjmp.c`:            {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/builtin-bitops-1.c`:           {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/builtin-prefetch-4.c`:         {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/builtin-types-compatible-p.c`: {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/call-trap-1.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/comp-goto-1.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/comp-goto-2.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-5.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-6.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/complex-7.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/compndlit-1.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/divconst-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/extzvsi.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ffs-1.c`:                      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ffs-2.c`:                      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/fprintf-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/frame-address.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/compare-fp-1.c`:          {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4f.c`:             {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-4l.c`:             {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-5.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8f.c`:             {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/fp-cmp-8l.c`:             {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/mzero4.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr38016.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr50310.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr72824-2.c`:             {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/pr84235.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ieee/unsafe-fp-assoc-1.c`:     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/loop-15.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/lto-tbaa-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/medce-1.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/memchr-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/memcpy-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/memset-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/memset-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/memset-3.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nest-align-1.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nest-stdar-1.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-1.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-2.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-3.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-5.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-6.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/nestfunc-7.c`:                 {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr17377.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr19449.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr19687.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr19689.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-3.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22061-4.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22098-3.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22141-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr22141-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr23135.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr23324.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr24135.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr28289.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr28982b.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr30185.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr30778.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr31136.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr31169.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr31448-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr31448.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr32244-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr34154.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr34768-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr34768-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr34971.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr35456.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr36038.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr36321.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr37573.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr37780.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr37882.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38151.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38422.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38533.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr38969.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39100.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39228.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr39339.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr40022.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr40404.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr40493.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr40657.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr41239.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr42570.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr42614.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43220.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43385.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43560.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr43987.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr44164.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr44555.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr44575.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr44852.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr45695.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr46309.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr47148.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr47155.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr47237.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr47925.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr48973-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr48973-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49073.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49123.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49218.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49279.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49390.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49644.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr49768.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51447.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51581-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51581-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51877.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr51933.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52209.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52286.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52979-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr52979-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr53645-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr53645.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr54471.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr55750.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56205.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56837.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56866.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr56982.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57130.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-3.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57344-4.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57568.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57861.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57876.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr57877.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58277-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58277-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58385.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58419.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58431.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58564.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58570.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58726.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58831.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr58984.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr59388.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60003.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60017.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr60960.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr61375.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr61725.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr63302.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr63641.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64006.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64242.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr64756.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65053-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65053-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65170.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65215-3.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65215-4.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65369.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65427.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65648.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr65956.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr66556.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr67037.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr67714.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68185.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68249.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68250.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68321.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68328.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68381.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68506.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr68532.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr69320-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr69320-4.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr69691.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70127.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70460.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70566.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70586.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70602.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr70903.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71083.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71494.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71554.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71626-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71626-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr71700.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr77767.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78170.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78436.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78438.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78477.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78559.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78675.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr78726.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79286.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79354.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79737-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr79737-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr80421.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr80692.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr81423.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr81555.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr81556.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr81588.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82192.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82210.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82387.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82524.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr82954.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr83362.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr83383.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84169.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84478.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84521.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84524.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr84748.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85095.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85156.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85169.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85331.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85529-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85582-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85582-3.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr85756.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr86492.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr86528.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr86659-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr86659-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr87053.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr88714.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr88739.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr88904.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr89195.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr89369.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr89434.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr90311.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr90949.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91137.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91450-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91450-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91597.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr91635.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr92618.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr92904.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93213.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93249.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93434.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93494.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93908.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr93945.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94130.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94412.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94524-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94524-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94591.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94724.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94734.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr94809.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr96549.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr97325.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr97421-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr97764.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr98366.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/pr98474.c`:                    {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/printf-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/return-addr.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/scal-to-vec1.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/scal-to-vec2.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/scal-to-vec3.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-1.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-2.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-4.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-5.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/simd-6.c`:                     {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/ssad-run.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stdarg-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stdarg-3.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/stkalign.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strcpy-1.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strcpy-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strct-stdarg-1.c`:             {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strct-varg-1.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/string-opt-18.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-3.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-4.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-6.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/strlen-7.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/struct-ini-1.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/struct-ini-2.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/struct-ini-3.c`:               {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/usad-run.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/user-printf.c`:                {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-10.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-14.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-15.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-2.c`:                   {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-22.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/va-arg-pack-1.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/vla-dealloc-1.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/vrp-7.c`:                      {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zero-struct-1.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zero-struct-2.c`:              {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zerolen-1.c`:                  {}, // COMPILE FAIL
	`assets/github.com/gcc-mirror/gcc/gcc/testsuite/gcc.c-torture/execute/zerolen-2.c`:                  {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-benchmarks/funnkuch-reduce.c`:                                    {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-benchmarks/method-call.c`:                                        {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-benchmarks/nbody.c`:                                              {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-benchmarks/strcat.c`:                                             {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0011-switch1.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0015-calls13.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0021-tentativedefs1.c`:                    {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0025-duff.c`:                              {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits06.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits07.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits08.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits09.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits10.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits11.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits12.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits13.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits14.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/andrewchambers_c/0028-inits15.c`:                           {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/anonymous-members.c`:                                  {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-basic.c`:                                     {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-extend.c`:                                    {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/bitfield-pack-next.c`:                                 {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/conditional-void.c`:                                   {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/constant-expression.c`:                                {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/declarator-abstract.c`:                                {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/duffs-device.c`:                                       {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/field-chain-assign.c`:                                 {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/function-pointer-call.c`:                              {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/initialize-call.c`:                                    {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/initialize-string.c`:                                  {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/sizeof.c`:                                             {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/string-addr.c`:                                        {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/struct-padding.c`:                                     {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/union-bitfield.c`:                                     {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/union-zero-init.c`:                                    {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg-complex-1.c`:                                   {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/lacc/vararg.c`:                                             {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/array-size-def.c`:                                      {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/data-than-bss.c`:                                       {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/enum-const-scope.c`:                                    {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/issue117.c`:                                            {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/issue142.c`:                                            {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/matrix-param.c`:                                        {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/test1.c`:                                               {}, // COMPILE FAIL
	`assets/github.com/vnmakarov/mir/c-tests/new/var-size-in-var-initializer.c`:                         {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/73_arm64.c`:                                                         {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/76_dollars_in_identifiers.c`:                                        {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/79_vla_continue.c`:                                                  {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/80_flexarray.c`:                                                     {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/81_types.c`:                                                         {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/85_asm-outside-function.c`:                                          {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/87_dead_code.c`:                                                     {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/88_codeopt.c`:                                                       {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/89_nocode_wanted.c`:                                                 {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/90_struct-init.c`:                                                   {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/92_enum_bitfield.c`:                                                 {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/93_integer_promotion.c`:                                             {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/94_generic.c`:                                                       {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/95_bitfields.c`:                                                     {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/95_bitfields_ms.c`:                                                  {}, // COMPILE FAIL
	`assets/tcc-0.9.27/tests/tests2/98_al_ax_extend.c`:                                                  {}, // COMPILE FAIL
}

func getArgs(src string) (args []string, err error) {
	src = src[:len(src)-len(filepath.Ext(src))] + ".arg"
	b, err := os.ReadFile(src)
	if err != nil {
		return nil, nil
	}

	a := strings.Split(strings.TrimSpace(string(b)), "\n")
	for _, v := range a {
		switch {
		case strings.HasPrefix(v, "\"") || strings.HasPrefix(v, "`"):
			w, err := strconv.Unquote(v)
			if err != nil {
				return nil, fmt.Errorf("%s: %v: %v", src, v, err)
			}

			args = append(args, w)
		default:
			args = append(args, v)
		}
	}
	return args, nil
}

type golden struct {
	a  []string
	f  *os.File
	mu sync.Mutex
	t  *testing.T

	discard bool
}

func newGolden(t *testing.T, fn string) *golden {
	if re != nil {
		return &golden{discard: true}
	}

	f, err := os.Create(filepath.FromSlash(fn))
	if err != nil { // Possibly R/O fs in a VM
		base := filepath.Base(filepath.FromSlash(fn))
		f, err = ioutil.TempFile("", base)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("writing results to %s\n", f.Name())
	}

	return &golden{t: t, f: f}
}

func (g *golden) w(s string, args ...interface{}) {
	if g.discard {
		return
	}

	g.mu.Lock()

	defer g.mu.Unlock()

	if s = strings.TrimRight(s, " \t\n\r"); !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	g.a = append(g.a, fmt.Sprintf(s, args...))
}

func (g *golden) close() {
	if g.discard || g.f == nil {
		return
	}

	defer func() { g.f = nil }()

	sort.Strings(g.a)
	if _, err := g.f.WriteString(strings.Join(g.a, "")); err != nil {
		g.t.Fatal(err)
	}

	if err := g.f.Sync(); err != nil {
		g.t.Fatal(err)
	}

	if err := g.f.Close(); err != nil {
		g.t.Fatal(err)
	}
}

func getCorpusFile(path string) ([]byte, error) {
	f, err := cfs.Open(path)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func TestSQLite(t *testing.T) {
	return //TODO-
	testSQLite(t, "assets/sqlite-amalgamation-3380100")
}

func testSQLite(t *testing.T, dir string) {
	const main = "main.go"
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	defer os.Chdir(wd)

	temp, err := ioutil.TempDir("", "ccgo-test-")
	if err != nil {
		t.Fatal(err)
	}

	switch {
	case *oKeep:
		t.Log(temp)
	default:
		defer os.RemoveAll(temp)
	}

	if err := os.Chdir(temp); err != nil {
		t.Fatal(err)
	}

	ccgoArgs := []string{
		"ccgo",

		"-DHAVE_USLEEP",
		"-DLONGDOUBLE_TYPE=double",
		"-DSQLITE_DEBUG",
		"-DSQLITE_DEFAULT_MEMSTATUS=0",
		"-DSQLITE_ENABLE_DBPAGE_VTAB",
		"-DSQLITE_LIKE_DOESNT_MATCH_BLOBS",
		"-DSQLITE_MEMDEBUG",
		"-DSQLITE_THREADSAFE=0",
		"-o", main,
		path.Join(dir, "shell.c"),
		path.Join(dir, "sqlite3.c"),
	}
	if *oDebug {
		ccgoArgs = append(ccgoArgs, "-DSQLITE_DEBUG_OS_TRACE", "-DSQLITE_FORCE_OS_TRACE", "-DSQLITE_LOCK_TRACE")
	}
	if os.Getenv("GO111MODULE") != "off" {
		if out, err := shell(true, "go", "mod", "init", "example.com/ccgo/v3/lib/sqlite"); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}

		if out, err := shell(true, "go", "get", "modernc.org/libc"); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
	}

	if !func() (r bool) {
		defer func() {
			if err := recover(); err != nil {
				if *oStackTrace {
					fmt.Printf("%s\n", debug.Stack())
				}
				if *oTrace {
					fmt.Println(err)
				}
				t.Errorf("%v", err)
				r = false
			}
			if *oTraceF {
				b, _ := ioutil.ReadFile(main)
				fmt.Printf("\n----\n%s\n----\n", b)
			}
		}()

		if err := NewTask(goos, goarch, ccgoArgs, nil, nil, cfs).Main(); err != nil {
			if *oTrace {
				fmt.Println(err)
			}
			// err = cpp(*oCpp, ccgoArgs, err)
			t.Errorf("%v", err)
			return false
		}

		return true
	}() {
		return
	}

	shell := "./shell"
	if runtime.GOOS == "windows" {
		shell = "shell.exe"
	}
	args := []string{"build"}
	if s := *oXTags; s != "" {
		args = append(args, "-tags", s)
	}
	args = append(args, "-o", shell, main)
	if out, err := exec.Command("go", args...).CombinedOutput(); err != nil {
		s := strings.TrimSpace(string(out))
		if s != "" {
			s += "\n"
		}
		t.Errorf("%s%v", s, err)
		return
	}

	var out []byte
	switch {
	case *oDebug:
		out, err = exec.Command(shell, "tmp", ".log stdout", "create table t(i); insert into t values(42); select 11*i from t;").CombinedOutput()
	default:
		out, err = exec.Command(shell, "tmp", "create table t(i); insert into t values(42); select 11*i from t;").CombinedOutput()
	}
	if err != nil {
		if *oTrace {
			fmt.Printf("%s\n%s\n", out, err)
		}
		t.Errorf("%s\n%v", out, err)
		return
	}

	if g, e := strings.TrimSpace(string(out)), "462"; g != e {
		t.Errorf("got: %s\nexp: %s", g, e)
	}
	if *oTraceO {
		fmt.Printf("%s\n", out)
	}

	if out, err = exec.Command(shell, "tmp", "select 13*i from t;").CombinedOutput(); err != nil {
		if *oTrace {
			fmt.Printf("%s\n%s\n", out, err)
		}
		t.Errorf("%v", err)
		return
	}

	if g, e := strings.TrimSpace(string(out)), "546"; g != e {
		t.Errorf("got: %s\nexp: %s", g, e)
	}
	if *oTraceO {
		fmt.Printf("%s\n", out)
	}
}
