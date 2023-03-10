// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
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
	"modernc.org/mathutil"
)

const (
	csmithBitfields   = "--bitfields"    // --bitfields | --no-bitfields: enable | disable full-bitfields structs (enabled by default). // Was disabled by default in older versions,
	csmithNoBitfields = "--no-bitfields" // --bitfields | --no-bitfields: enable | disable full-bitfields structs (enabled by default).
)

var (
	oBlackBox     = flag.String("blackbox", "", "Record CSmith file to this file")
	oCSmith       = flag.Duration("csmith", 2*time.Hour, "")
	oCSmithClimit = flag.Duration("csmithc", 1*time.Minute, "")
	oDebug        = flag.Bool("debug", false, "")
	oErr1         = flag.Bool("err1", false, "first error line only")
	oKeep         = flag.Bool("keep", false, "keep temp directories")
	oPanic        = flag.Bool("panic", false, "panic on miscompilation")
	oShellTime    = flag.Duration("shelltimeout", 3600*time.Second, "shell() time limit")
	oStackTrace   = flag.Bool("trcstack", false, "")
	oTrace        = flag.Bool("trc", false, "print tested paths.")
	oTraceC       = flag.Bool("trcc", false, "trace TestExec transiple errors")
	oTraceF       = flag.Bool("trcf", false, "print test file content")
	oTraceO       = flag.Bool("trco", false, "print test output")
	oXTags        = flag.String("xtags", "", "passed to go build of TestSQLite")

	cfs    fs.FS
	goarch = runtime.GOARCH
	goos   = runtime.GOOS
	hostCC string
	re     *regexp.Regexp
	testWD string

	csmithDefaultArgs = strings.Join([]string{
		"--max-nested-struct-level", "10", // --max-nested-struct-level <num>: limit maximum nested level of structs to <num>(default 0). Only works in the exhaustive mode.
		"--no-const-pointers",    // --const-pointers | --no-const-pointers: enable | disable const pointers (enabled by default).
		"--no-consts",            // --consts | --no-consts: enable | disable const qualifier (enabled by default).
		"--no-packed-struct",     // --packed-struct | --no-packed-struct: enable | disable packed structs by adding #pragma pack(1) before struct definition (disabled by default).
		"--no-volatile-pointers", // --volatile-pointers | --no-volatile-pointers: enable | disable volatile pointers (enabled by default).
		"--no-volatiles",         // --volatiles | --no-volatiles: enable | disable volatiles (enabled by default).
		"--paranoid",             // --paranoid | --no-paranoid: enable | disable pointer-related assertions (disabled by default).
	}, " ")
)

type diskFS string

func newDiskFS(root string) diskFS { return diskFS(root) }

func (f diskFS) Open(name string) (fs.File, error) { return os.Open(filepath.Join(string(f), name)) }

type overlayFS struct {
	fs   fs.FS
	over fs.FS
}

func newOverlayFS(fs, over fs.FS) *overlayFS { return &overlayFS{fs, over} }

func (f *overlayFS) Open(name string) (fs.File, error) {
	fi, err := fs.Stat(f.over, name)
	if err == nil && !fi.IsDir() {
		if f, err := f.over.Open(name); err == nil {
			return f, nil
		}
	}

	return f.fs.Open(name)
}

func TestMain(m *testing.M) {
	isTesting = true
	testWD, err := filepath.Abs("testdata")
	if err != nil {
		panic(todo("", err))
	}

	overlay := filepath.Join(testWD, "overlay")
	cfs = newOverlayFS(ccorpus2.FS, newDiskFS(overlay))
	extendedErrors = true
	gc.ExtendedErrors = true
	oRE := flag.String("re", "", "")
	flag.BoolVar(&trcTODOs, "trctodo", false, "")
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
	if err := p.wait(); err != nil {
		a := strings.Split(err.Error(), "\n")
		for _, v := range a {
			t.Error(v)
		}
	}
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
	fis, err := fs.ReadDir(cfs, dir)
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

func TestGoAlign(t *testing.T) {
	for _, osarch := range []string{
		"darwin/amd64",
		"darwin/arm64",
		"freebsd/386",
		"freebsd/amd64",
		"freebsd/arm",
		"linux/386",
		"linux/amd64",
		"linux/arm",
		"linux/arm64",
		"linux/ppc64le",
		"linux/riscv64",
		"linux/s390x",
		"netbsd/amd64",
		"netbsd/arm",
		"openbsd/amd64",
		"openbsd/arm64",
		"windows/386",
		"windows/amd64",
		"windows/arm64",
	} {
		a := strings.Split(osarch, "/")
		os := a[0]
		arch := a[1]
		cabi, err := cc.NewABI(os, arch)
		if err != nil {
			t.Errorf("%s: %v", osarch, err)
			continue
		}

		var ks []cc.Kind
		for ck := range cabi.Types {
			ks = append(ks, ck)
		}
		sort.Slice(ks, func(i, j int) bool { return ks[i] < ks[j] })

		abi, err := gc.NewABI(os, arch)
		if err != nil {
			t.Errorf("%s: %v", osarch, err)
			continue
		}

		for _, ck := range ks {
			gk := gcKind(ck, cabi)
			if gk < 0 {
				continue
			}

			ct := cabi.Types[ck]
			gt := abi.Types[gk]
			if g, e := gt.Size, ct.Size; g != e {
				t.Errorf("%s: Go %v size %d, C %v size %d", osarch, gk, g, ck, e)
			}
			if g, e := gt.Align, ct.Align; g != e {
				t.Logf("%s: warning: Go %v align %d, C %v align %d", osarch, gk, g, ck, e)
			}
			if g, e := gt.FieldAlign, ct.FieldAlign; g != e {
				t.Logf("%s: warning: Go %v field align %d, C %v field align %d", osarch, gk, g, ck, e)
			}
		}
	}
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
			{"github.com/cxgo", false},
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
		args, err := getArgs(path)
		if err != nil {
			return err
		}

		if *oTrace {
			fmt.Fprintf(os.Stderr, "%5d %s %v\n", id, filepath.Join(cfsDir, path), args)
		}
		p.exec(func() error { return testExec1(t, p, cfsDir, path, exec, g, id, args) })
		return nil
	}))
}

func testExec1(t *testing.T, p *parallel, root, path string, execute bool, g *golden, id int, args []string) (err error) {
	// p.Lock()
	// trc("ADD %q", path)
	// if p.paths == nil {
	// 	p.paths = map[string]struct{}{}
	// }
	// p.paths[path] = struct{}{}
	// p.Unlock()
	defer func() {
		// p.Lock()
		// delete(p.paths, path)
		// trc("DONE %q %v", path, p.paths)
		// p.Unlock()
		if err != nil {
			p.fail()
		}
	}()

	fullPath := filepath.ToSlash(filepath.Join(root, path))
	var cCompilerFailed, cExecFailed bool
	ofn := fmt.Sprint(id)
	bin := "cbin_" + enforceBinaryExt(ofn)
	switch {
	case !execute:
		if _, err = shell(false, hostCC, "-c", "-w", path, "-lm"); err != nil {
			cCompilerFailed = true
		}
	default:
		if _, err = shell(false, hostCC, "-o", bin, "-w", path, "-lm"); err != nil {
			cCompilerFailed = true
		}
	}

	defer os.Remove(ofn)

	cbinRC := -1
	var cOut []byte
	if execute && !cCompilerFailed {
		switch cOut, err = shell(false, "./"+bin, args...); {
		case err != nil:
			cbinRC = exitCode(err)
			cExecFailed = true
		default:
			cbinRC = 0
		}
	}

	ofn += ".go"

	defer os.Remove(ofn)

	var out bytes.Buffer
	switch {
	case !execute:
		err = NewTask(
			goos,
			goarch,
			[]string{
				"ccgo",
				"-c",
				"-verify-types",
				"--prefix-field=F",
				"-ignore-asm-errors",
				"-ignore-vector-functions",
				"-ignore-unsupported-alignment",
				path,
			},
			&out, &out, nil).Main()
	default:
		err = NewTask(
			goos,
			goarch,
			[]string{
				"ccgo",
				"-o", ofn,
				"-verify-types",
				"--prefix-field=F",
				"-ignore-asm-errors",
				"-ignore-vector-functions",
				"-ignore-unsupported-alignment",
				path,
			},
			&out, &out, nil).Main()
	}
	if err != nil {
		if *oTraceC {
			trc("ccgo %v %v", fullPath, err)
		}
		if cCompilerFailed || isTestExecKnownFail(fullPath) {
			p.skip()
			return nil
		}

		trc("`%s`: {}, // COMPILE FAIL: %v", fullPath, firstError(err, true))
		p.err(err)
		return errorf("%s: %s: FAIL: %v", fullPath, out.Bytes(), firstError(err, *oErr1))
	}

	if !execute {
		p.ok()
		g.w("%s\n", fullPath)
		return nil
	}

	bin = "gobin_" + enforceBinaryExt(ofn)
	if _, err = shell(false, "go", "build", "-o", bin, ofn); err != nil {
		// trc("gc %v %v", path, err)
		if cCompilerFailed || isTestExecKnownFail(fullPath) {
			p.skip()
			return nil
		}

		trc("`%s`: {}, // BUILD FAIL: %v", fullPath, firstError(err, true))
		p.err(err)
		return firstError(err, *oErr1)
	}

	goOut, err := shell(false, "./"+bin, args...)
	if err != nil {
		// trc("gobin %v %v", path, err)
		gobinRC := exitCode(err)
		//trc("", cbinRC, gobinRC)
		switch {
		case gobinRC == cbinRC && gobinRC > 0:
			// makarov et al
			cExecFailed = false
		default:
			if cExecFailed || isTestExecKnownFail(fullPath) {
				p.skip()
				return nil
			}

			err := errorf("%s: %s: FAIL: %v", fullPath, goOut, err)
			if *oPanic {
				panic(err)
			}

			trc("`%s`: {}, // EXEC FAIL: %v", fullPath, firstError(err, true))
			p.err(err)
			return firstError(err, *oErr1)
		}
	}

	cOut = bytes.TrimSpace(cOut)
	goOut = bytes.TrimSpace(goOut)
	if *oTraceO {
		fmt.Printf("C out\n==== (A)\n%s\n==== (Z)\n", cOut)
		fmt.Printf("Go out\n==== (A)\n%s\n==== (Z)\n", goOut)
	}
	if bytes.Contains(cOut, []byte("\r\n")) {
		cOut = bytes.ReplaceAll(cOut, []byte("\r"), nil)
	}
	if bytes.Contains(goOut, []byte("\r\n")) {
		goOut = bytes.ReplaceAll(goOut, []byte("\r"), nil)
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
	p.err(err)
	return firstError(err, *oErr1)
}

func exitCode(err error) int {
	switch x := err.(type) {
	case *exec.ExitError:
		return x.ProcessState.ExitCode()
	default:
		trc("%T %s", x, x)
		return -1
	}
}

func isTestExecKnownFail(s string) (r bool) {
	_, r = testExecKnownFails[s]
	return r
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

func TestCSmith(t *testing.T) {
	abi, err := cc.NewABI(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Fatal(err)
	}

	if testing.Short() {
		t.Skip("skipped: -short")
	}

	csmith, err := exec.LookPath("csmith")
	if err != nil {
		switch runtime.GOOS {
		case "windows": //TODO
			t.Skip(err)
			return
		default:
			t.Fatal(err)
		}
	}

	bigEndian := abi.ByteOrder == binary.BigEndian
	binaryName := filepath.FromSlash("./a.out")
	goBinaryName := filepath.FromSlash("./main")
	if runtime.GOOS == "windows" {
		binaryName = filepath.FromSlash("./a.exe")
		goBinaryName = filepath.FromSlash("./main.exe")
	}
	mainName := filepath.FromSlash("main.go")
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	defer os.Chdir(wd)

	temp, err := ioutil.TempDir("", "ccgo-test-")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(temp)

	if err := os.Chdir(temp); err != nil {
		t.Fatal(err)
	}

	if os.Getenv("GO111MODULE") != "off" {
		if out, err := shell(true, "go", "mod", "init", "example.com/ccgo/v4/lib/csmith"); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}

		if out, err := shell(true, "go", "get", "modernc.org/libc"); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
	}

	//TODO report the problem at http://www.flux.utah.edu/mailman/listinfo/csmith-bugs
	bigEndianBlacklist := []string{
		"-s 2949258094",
		"-s 3329111231",
		"-s 4101947480",
	}

	fixedBugs := []string{
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 1110506964",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 1338573550",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 1416441494",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 15739796933983044010",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 169375684",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 1833258637",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 1885311141",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 2205128324",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 2273393378",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 241244373",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 2517344771",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 2648215054",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 2876930815",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 2877850218",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 2949258094",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3043990076",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3100949894",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3126091077",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3130410542",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3329111231",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3363122597",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3365074920",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3578720023",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3645367888",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3919255949",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3980073540",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 4058772172",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 4101947480",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 4130344133",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 4146870674",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 424465590",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 517639208",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 56498550",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 890611563",
		"--bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 963985971",
		"--bitfields --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid --max-nested-struct-level 10 -s 1236173074",
		"--bitfields --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid --max-nested-struct-level 10 -s 1906742816",
		"--bitfields --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid --max-nested-struct-level 10 -s 3629008936",
		"--bitfields --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid --max-nested-struct-level 10 -s 612971101",
		"--no-bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 1302111308",
		"--no-bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3285852464",
		"--no-bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3609090094",
		"--no-bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 3720922579",
		"--no-bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 4263172072",
		"--no-bitfields --max-nested-struct-level 10 --no-const-pointers --no-consts --no-packed-struct --no-volatile-pointers --no-volatiles --paranoid -s 572192313",
	}
	var ch <-chan time.Time
	t0 := time.Now()
	var files, ok int
	var size int64
out:
	for i := 0; ; i++ {
		extra := ""
		var args string
		switch {
		case i < len(fixedBugs):
			s := fixedBugs[i]
			if re != nil && !re.MatchString(s) {
				continue
			}

			if bigEndian {
				for _, v := range bigEndianBlacklist {
					if strings.Contains(s, v) {
						continue out
					}
				}
			}

			args += fixedBugs[i]
			a := strings.Split(s, " ")
			extra = strings.Join(a[len(a)-2:], " ")
			t.Log(args)
		default:
			if ch == nil {
				ch = time.After(*oCSmith)
			}
			select {
			case <-ch:
				break out
			default:
			}

			args += csmithDefaultArgs
			switch {
			case bigEndian:
				args += " " + csmithNoBitfields
			default:
				args += " " + csmithBitfields
			}
		}
		csOut, err := exec.Command(csmith, strings.Split(args, " ")...).Output()
		if err != nil {
			t.Fatalf("%v\n%s", err, csOut)
		}

		if fn := *oBlackBox; fn != "" {
			if err := ioutil.WriteFile(fn, csOut, 0660); err != nil {
				t.Fatal(err)
			}
		}

		if err := ioutil.WriteFile("main.c", csOut, 0660); err != nil {
			t.Fatal(err)
		}

		// Shorten the log output, the CSmith header captures all info.  Full C code in
		// blackbox, if turned on.
		lines := strings.Split(strings.TrimSpace(string(csOut)), "\n")
		lines = lines[:mathutil.Min(len(lines), 8)]
		csOut = []byte(strings.Join(lines, "\n"))

		csp := fmt.Sprintf("-I%s", filepath.FromSlash("/usr/include/csmith"))
		if s := os.Getenv("CSMITH_PATH"); s != "" {
			csp = fmt.Sprintf("-I%s", s)
		}

		ccOut, err := exec.Command(hostCC, "-o", binaryName, "main.c", csp).CombinedOutput()
		if err != nil {
			t.Logf("%s\n%s\ncc: %v", extra, ccOut, err)
			continue
		}

		ctime0 := time.Now()
		binOutA, err := func() ([]byte, error) {
			ctx, cancel := context.WithTimeout(context.Background(), *oCSmithClimit)
			defer cancel()

			return exec.CommandContext(ctx, binaryName).CombinedOutput()
		}()
		if err != nil {
			continue
		}

		ctime := time.Since(ctime0)
		if *oTrace {
			fmt.Fprintf(os.Stderr, "[%s %s]:  C binary real %s\n", time.Now().Format("15:04:05"), durationStr(time.Since(t0)), ctime)
		}
		if ctime > *oCSmithClimit {
			continue
		}

		size += int64(len(csOut))

		if err := os.Remove(binaryName); err != nil {
			t.Fatal(err)
		}

		files++
		var stdout, stderr bytes.Buffer
		j := NewTask(
			goos,
			goarch,
			[]string{
				"ccgo",

				"-o", mainName,
				"-extended-errors",
				"-verify-types",
				"--prefix-field=F",
				"-ignore-asm-errors",
				"-ignore-vector-functions",
				"-ignore-unsupported-alignment",
				"main.c",
				csp,
			},
			&stdout,
			&stderr,
			nil)

		func() {

			defer func() {
				if err := recover(); err != nil {
					t.Errorf("%s\n%s\nccgo: %s\n%s\n%s", extra, csOut, stdout.Bytes(), stderr.Bytes(), debug.Stack())
					t.Fatal(err)
				}
			}()

			if err := j.Main(); err != nil || stdout.Len() != 0 {
				t.Errorf("%s\n%s\nccgo: %s\n%s", extra, csOut, stdout.Bytes(), stderr.Bytes())
				t.Fatal(err)
			}
		}()

		out, err := shell(false, "go", "build", "-o", goBinaryName, mainName)
		if err != nil {
			t.Errorf("%s\n%s\n%s\nccgo: %v", extra, csOut, out, err)
			break
		}

		goLimit := 10 * ctime
		if goLimit < 20*time.Minute {
			goLimit = 20 * time.Minute
		}
		goTime0 := time.Now()
		binOutB, err := func() ([]byte, error) {
			ctx, cancel := context.WithTimeout(context.Background(), goLimit)
			defer cancel()

			return exec.CommandContext(ctx, goBinaryName).CombinedOutput()
		}()
		if g, e := binOutB, binOutA; !bytes.Equal(g, e) {
			t.Errorf("%s\n%s\nccgo: %v\ngot: %s\nexp: %s", extra, csOut, err, g, e)
			break
		}

		ok++
		if *oTrace {
			fmt.Fprintf(os.Stderr, "[%s %s]: Go binary real %s\tfiles %v, ok %v, \n", time.Now().Format("15:04:05"), durationStr(time.Since(t0)), time.Since(goTime0), files, ok)
		}

		if err := os.Remove(mainName); err != nil {
			t.Fatal(err)
		}
	}
	d := time.Since(t0)
	t.Logf("files %v, bytes %v, ok %v in %v", h(files), h(size), h(ok), d)
}

func durationStr(d time.Duration) string {
	secs := d / time.Second
	mins := secs / 60
	hours := mins / 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, mins%60, secs%60)
}

func TestSQLite(t *testing.T) {
	testSQLite(t, "assets/sqlite-amalgamation")
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
		"--prefix-field=F",
		"-positions",
		"-full-paths",
		"-verify-types",
		"-ignore-asm-errors",
		"-ignore-unsupported-alignment",
		"-ignore-vector-functions",
		"-o", main,
		filepath.Join(dir, "shell.c"),
		filepath.Join(dir, "sqlite3.c"),
		filepath.Join(dir, "patch.c"),
	}
	if *oKeep {
		ccgoArgs = append(ccgoArgs, "-keep-object-files", "-extended-errors", "-debug-linker-save")
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

		task := NewTask(goos, goarch, ccgoArgs, nil, nil, cfs)
		if err := task.Main(); err != nil {
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
