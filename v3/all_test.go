// Copyright 2020 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "modernc.org/ccgo/v3"

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"unsafe"

	"github.com/dustin/go-humanize"
	"modernc.org/crt/v3"
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

func TODO(...interface{}) string { //TODOOK
	_, fn, fl, _ := runtime.Caller(1)
	return fmt.Sprintf("# TODO: %s:%d:\n", path.Base(fn), fl) //TODOOK
}

func stack() string { return string(debug.Stack()) }

func use(...interface{}) {}

func init() {
	use(caller, dbg, TODO, stack) //TODOOK
}

// ----------------------------------------------------------------------------

var (
	oDev        = flag.Bool("dev", false, "Enable developer tests/downloads.")
	oDownload   = flag.Bool("download", false, "Download missing testdata. Add -dev to download also 100+ MB of developer resources.")
	oRE         = flag.String("re", "", "")
	oStackTrace = flag.Bool("trcstack", false, "")
	oTrace      = flag.Bool("trc", false, "Print tested paths.")
	oTraceO     = flag.Bool("trco", false, "Print test output")
	oTraceF     = flag.Bool("trcf", false, "Print test file content")

	gccDir    = filepath.FromSlash("testdata/gcc-9.1.0")
	gpsdDir   = filepath.FromSlash("testdata/gpsd-3.20/")
	ntpsecDir = filepath.FromSlash("testdata/ntpsec-master")
	sqliteDir = filepath.FromSlash("testdata/sqlite-amalgamation-3300100")
	tccDir    = filepath.FromSlash("testdata/tcc-0.9.27")

	testWD string

	downloads = []struct {
		dir, url string
		sz       int
		dev      bool
	}{
		{gccDir, "http://mirror.koddos.net/gcc/releases/gcc-9.1.0/gcc-9.1.0.tar.gz", 118000, true},
		{gpsdDir, "http://download-mirror.savannah.gnu.org/releases/gpsd/gpsd-3.20.tar.gz", 3600, false},
		{ntpsecDir, "https://gitlab.com/NTPsec/ntpsec/-/archive/master/ntpsec-master.tar.gz", 2600, false},
		{sqliteDir, "https://www.sqlite.org/2019/sqlite-amalgamation-3300100.zip", 2400, false},
		{tccDir, "http://download.savannah.gnu.org/releases/tinycc/tcc-0.9.27.tar.bz2", 620, false},
	}
)

func TestMain(m *testing.M) {
	defer func() {
		os.Exit(m.Run())
	}()

	flag.BoolVar(&oTraceW, "trcw", false, "Print generator writes")
	flag.BoolVar(&oTraceG, "trcg", false, "Print generator output")
	flag.Parse()
	var err error
	if testWD, err = os.Getwd(); err != nil {
		panic("Cannot determine working dir: " + err.Error())
	}

	if *oDownload {
		download()
	}
}

func download() {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	defer os.RemoveAll(tmp)

	for _, v := range downloads {
		dir := filepath.FromSlash(v.dir)
		root := filepath.Dir(v.dir)
		fi, err := os.Stat(dir)
		switch {
		case err == nil:
			if !fi.IsDir() {
				fmt.Fprintf(os.Stderr, "expected %s to be a directory\n", dir)
			}
			continue
		default:
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "%s", err)
				continue
			}

			if v.dev && !*oDev {
				fmt.Printf("Not downloading (no -dev) %v MB from %s\n", float64(v.sz)/1000, v.url)
				continue
			}

		}

		if err := func() error {
			fmt.Printf("Downloading %v MB from %s\n", float64(v.sz)/1000, v.url)
			resp, err := http.Get(v.url)
			if err != nil {
				return err
			}

			defer resp.Body.Close()

			base := filepath.Base(v.url)
			name := filepath.Join(tmp, base)
			f, err := os.Create(name)
			if err != nil {
				return err
			}

			defer os.Remove(name)

			n, err := io.Copy(f, resp.Body)
			if err != nil {
				return err
			}

			if _, err := f.Seek(0, io.SeekStart); err != nil {
				return err
			}

			switch {
			case strings.HasSuffix(base, ".tar.bz2"):
				b2r := bzip2.NewReader(bufio.NewReader(f))
				tr := tar.NewReader(b2r)
				for {
					hdr, err := tr.Next()
					if err != nil {
						if err != io.EOF {
							return err
						}

						return nil
					}

					switch hdr.Typeflag {
					case tar.TypeDir:
						if err = os.MkdirAll(filepath.Join(root, hdr.Name), 0770); err != nil {
							return err
						}
					case tar.TypeReg, tar.TypeRegA:
						f, err := os.OpenFile(filepath.Join(root, hdr.Name), os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))
						if err != nil {
							return err
						}

						w := bufio.NewWriter(f)
						if _, err = io.Copy(w, tr); err != nil {
							return err
						}

						if err := w.Flush(); err != nil {
							return err
						}

						if err := f.Close(); err != nil {
							return err
						}
					default:
						return fmt.Errorf("unexpected tar header typeflag %#02x", hdr.Typeflag)
					}
				}
			case strings.HasSuffix(base, ".tar.gz"):
				return untar(root, bufio.NewReader(f))
			case strings.HasSuffix(base, ".zip"):
				r, err := zip.NewReader(f, n)
				if err != nil {
					return err
				}

				for _, f := range r.File {
					fi := f.FileInfo()
					if fi.IsDir() {
						if err := os.MkdirAll(filepath.Join(root, f.Name), 0770); err != nil {
							return err
						}

						continue
					}

					if err := func() error {
						rc, err := f.Open()
						if err != nil {
							return err
						}

						defer rc.Close()

						dname := filepath.Join(root, f.Name)
						g, err := os.Create(dname)
						if err != nil {
							return err
						}

						defer g.Close()

						n, err = io.Copy(g, rc)
						return err
					}(); err != nil {
						return err
					}
				}
				return nil
			}
			panic("internal error") //TODOOK
		}(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func untar(root string, r io.Reader) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}

			return nil
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filepath.Join(root, hdr.Name), 0770); err != nil {
				return err
			}
		case tar.TypeSymlink, tar.TypeXGlobalHeader:
			// skip
		case tar.TypeReg, tar.TypeRegA:
			dir := filepath.Dir(filepath.Join(root, hdr.Name))
			if _, err := os.Stat(dir); err != nil {
				if !os.IsNotExist(err) {
					return err
				}

				if err = os.MkdirAll(dir, 0770); err != nil {
					return err
				}
			}

			f, err := os.OpenFile(filepath.Join(root, hdr.Name), os.O_CREATE|os.O_WRONLY, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			w := bufio.NewWriter(f)
			if _, err = io.Copy(w, tr); err != nil {
				return err
			}

			if err := w.Flush(); err != nil {
				return err
			}

			if err := f.Close(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected tar header typeflag %#02x", hdr.Typeflag)
		}
	}
}

type golden struct {
	t *testing.T
	f *os.File
	w *bufio.Writer
}

func newGolden(t *testing.T, fn string) *golden {
	if *oRE != "" {
		return &golden{w: bufio.NewWriter(ioutil.Discard)}
	}

	f, err := os.Create(filepath.FromSlash(fn))
	if err != nil {
		t.Fatal(err)
	}

	w := bufio.NewWriter(f)
	return &golden{t, f, w}
}

func (g *golden) close() {
	if g.f == nil {
		return
	}

	if err := g.w.Flush(); err != nil {
		g.t.Fatal(err)
	}

	if err := g.f.Close(); err != nil {
		g.t.Fatal(err)
	}
}

func h(v interface{}) string {
	switch x := v.(type) {
	case int:
		return humanize.Comma(int64(x))
	case int64:
		return humanize.Comma(x)
	case uint64:
		return humanize.Comma(int64(x))
	case float64:
		return humanize.CommafWithDigits(x, 0)
	default:
		panic(fmt.Errorf("%T", x)) //TODOOK
	}
}

func TestTCC(t *testing.T) {
	root := filepath.Join(testWD, filepath.FromSlash(tccDir))
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("Missing resources in %s. Please run 'go test -download' to fix.", root)
	}

	g := newGolden(t, fmt.Sprintf("testdata/tcc_%s_%s.golden", runtime.GOOS, runtime.GOARCH))

	defer g.close()

	var files, ok int
	const dir = "tests/tests2"
	f, o := testTCCExec(g.w, t, filepath.Join(root, filepath.FromSlash(dir)))
	files += f
	ok += o
	t.Logf("files %s, ok %s", h(files), h(ok))
}

func testTCCExec(w io.Writer, t *testing.T, dir string) (files, ok int) {
	const main = "main.go"
	blacklist := map[string]struct{}{
		"34_array_assignment.c":    {}, // gcc: 16:6: error: assignment to expression with array type
		"60_errors_and_warnings.c": {}, // Not a standalone test. gcc fails.
		"93_integer_promotion.c":   {}, // The expected output does not agree with gcc.
		"95_bitfields.c":           {}, // Included from 95_bitfields_ms.c
		"95_bitfields_ms.c":        {}, // The expected output does not agree with gcc.
		"96_nodata_wanted.c":       {}, // Not a standalone test. gcc fails.
		"99_fastcall.c":            {}, // 386 only

		"40_stdio.c":                {}, //TODO
		"42_function_pointer.c":     {}, //TODO
		"46_grep.c":                 {}, //TODO
		"73_arm64.c":                {}, //TODO struct varargs, not supported by QBE
		"75_array_in_struct_init.c": {}, //TODO flat struct initializer
		"80_flexarray.c":            {}, //TODO Flexible array member
		"85_asm-outside-function.c": {}, //TODO
		"90_struct-init.c":          {}, //TODO cc ... in designator
		"94_generic.c":              {}, //TODO cc _Generic
		"98_al_ax_extend.c":         {}, //TODO
	}
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

	var re *regexp.Regexp
	if s := *oRE; s != "" {
		re = regexp.MustCompile(s)
	}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				err = nil
			}
			return err
		}

		if info.IsDir() {
			return skipDir(path)
		}

		if filepath.Ext(path) != ".c" || info.Mode()&os.ModeType != 0 {
			return nil
		}

		if _, ok := blacklist[filepath.Base(path)]; ok {
			return nil
		}

		files++

		if re != nil && !re.MatchString(path) {
			return nil
		}

		if *oTrace {
			fmt.Fprintln(os.Stderr, files, ok, path)
		}

		if err := os.Remove(main); err != nil && !os.IsNotExist(err) {
			return err
		}

		ccgoArgs := []string{"ccgo", "-o", main}
		var args []string
		switch base := filepath.Base(path); base {
		case "31_args.c":
			args = []string{"arg1", "arg2", "arg3", "arg4", "arg5"}
		case "46_grep.c":
			if err := copyFile(path, filepath.Join(temp, base)); err != nil {
				return err
			}

			args = []string{`[^* ]*[:a:d: ]+\:\*-/: $`, base}
		}
		if !func() (r bool) {
			defer func() {
				if err := recover(); err != nil {
					if *oStackTrace {
						fmt.Printf("%s\n", stack())
					}
					if *oTrace {
						fmt.Println(err)
					}
					t.Errorf("%s: %v", path, err)
					r = false
				}
			}()

			ccgoArgs = append(ccgoArgs, path)
			if err := newTask(ccgoArgs, nil, nil).main(); err != nil {
				if *oTrace {
					fmt.Println(err)
				}
				t.Errorf("%s: %v", path, err)
				return false
			}

			return true
		}() {
			return nil
		}

		out, err := exec.Command("go", append([]string{"run", main}, args...)...).CombinedOutput()
		if err != nil {
			if *oTrace {
				fmt.Println(err)
			}
			b, _ := ioutil.ReadFile(main)
			t.Errorf("\n%s\n%v: %s\n%v", b, path, out, err)
			return nil
		}

		if *oTraceF {
			b, _ := ioutil.ReadFile(main)
			fmt.Printf("\n----\n%s\n----\n", b)
		}
		if *oTraceO {
			fmt.Printf("%s\n", out)
		}
		exp, err := ioutil.ReadFile(noExt(path) + ".expect")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintln(w, filepath.Base(path))
				ok++
				return nil
			}

			return err
		}

		out = trim(out)
		exp = trim(exp)

		switch base := filepath.Base(path); base {
		case "70_floating_point_literals.c": //TODO TCC binary extension
			a := strings.Split(string(exp), "\n")
			exp = []byte(strings.Join(a[:35], "\n"))
		}

		if !bytes.Equal(out, exp) {
			if *oTrace {
				fmt.Println(err)
			}
			t.Errorf("%v: out\n%s\nexp\n%s", path, out, exp)
			return nil
		}

		fmt.Fprintln(w, filepath.Base(path))
		ok++
		return nil
	}); err != nil {
		t.Errorf("%v", err)
	}
	return files, ok
}

func trim(b []byte) (r []byte) {
	b = bytes.TrimLeft(b, "\n")
	b = bytes.TrimRight(b, "\n")
	a := bytes.Split(b, []byte("\n"))
	for i, v := range a {
		a[i] = bytes.TrimSpace(v)
	}
	return bytes.Join(a, []byte("\n"))
}

func noExt(s string) string {
	ext := filepath.Ext(s)
	if ext == "" {
		panic("internal error") //TODOOK
	}
	return s[:len(s)-len(ext)]
}

func copyFile(src, dst string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, b, 0660)
}

func skipDir(path string) error {
	sp := filepath.ToSlash(path)
	if strings.Contains(sp, "/.") {
		return filepath.SkipDir
	}

	return nil
}

func TestCAPI(t *testing.T) {
	var _ crt.Intptr
	task := newTask(nil, nil, nil)
	pkgName, capi, err := task.capi("modernc.org/crt/v3")
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := capi["printf"]; !ok {
		t.Fatal("default crt does not export printf")
	}

	t.Log(pkgName, capi)
}

const text = "abcd\nefgh\x00ijkl"

var (
	text1 = text
	text2 = (*reflect.StringHeader)(unsafe.Pointer(&text1)).Data
)

func TestText(t *testing.T) {
	p := text2
	var b []byte
	for i := 0; i < len(text); i++ {
		b = append(b, *(*byte)(unsafe.Pointer(p)))
		p++
	}
	if g, e := string(b), text; g != e {
		t.Fatalf("%q %q", g, e)
	}
}