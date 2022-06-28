// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

//TODO CSmith

import (
	"bytes"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/pmezard/go-difflib/difflib"
	"modernc.org/cc/v4"
	"modernc.org/ccorpus2"
	"modernc.org/gc/v2"
)

var (
	oDebug      = flag.Bool("debug", false, "")
	oKeep       = flag.Bool("keep", false, "keep temp directories")
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
		if x := strings.Index(v, ":"); x > 0 {
			t.Logf("%q: {}, //TODO", v[:x])
		}
	}
	p.Unlock()
	t.Logf("TOTAL: files %v, skip %v, ok %v, fails %v", h(p.files), h(p.skips), h(p.oks), h(p.fails))
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

func TestCompile(t *testing.T) {
	g := newGolden(t, fmt.Sprintf("testdata/test_compile_%s_%s.golden", runtime.GOOS, runtime.GOARCH))

	defer g.close()

	tmp := t.TempDir()
	//blacklistCompCert := map[string]struct{}{}
	blacklistGCC := map[string]struct{}{
		// Assertions are deprecated, not supported.
		"950919-1.c": {},

		// Alignment
		"20010904-1.c": {},
		"20010904-2.c": {},
		"align-3.c":    {},

		// Alloca
		"20020314-1.c": {},
		"20021113-1.c": {},
		"20040223-1.c": {},
		"20070824-1.c": {},

		// asm
		"20001009-2.c": {},
		"20020107-1.c": {},
		"20030222-1.c": {},
		"20050203-1.c": {},
		"960312-1.c":   {},
		"990130-1.c":   {},
		"990413-2.c":   {},

		// attribute alias
		"alias-3.c": {},
		"alias-4.c": {},

		// __builtin_frame_address
		"frame-address.c": {},

		// __builtin_islessgreater
		"compare-fp-1.c": {},
		"pr50310.c":      {},

		// __builtin_setjmp
		"built-in-setjmp.c": {},

		// __builtin_signbit
		"20080502-1.c": {},

		// __builtin_tanf
		"mzero4.c": {},

		// llabs hardcoded to __builtin_llabs
		"20021127-1.c": {},

		// __builtin_return_address
		"20030323-1.c": {},
		"20030811-1.c": {},

		// DCE
		"20030330-1.c": {},

		// ffs
		"ffs-1.c": {},

		// Flexiable array member
		"20051113-1.c": {},

		// instrumentation
		"eeprof-1.c": {},

		// int128
		"pr54471.c": {},
		"pr65170.c": {},
		"pr84169.c": {},

		// link error
		"fp-cmp-7.c": {},

		// long double
		"960405-1.c": {},

		// old style function declaration
		"920501-1.c": {},
		"921208-2.c": {},

		// OOB/OOL
		"20080222-1.c": {},
		"991014-1.c":   {},

		// sigfpe
		"20101011-1.c": {},

		// tmpnam
		"fprintf-2.c": {},

		//TODO
		"call-trap-1.c":       {},
		"unsafe-fp-assoc-1.c": {},
		"loop-2e.c":           {},

		// uses undefined functions
		"921202-1.c": {},

		// Vector
		"20050316-1.c":   {},
		"20050316-2.c":   {},
		"20050316-3.c":   {},
		"20050604-1.c":   {},
		"20050607-1.c":   {},
		"20060420-1.c":   {},
		"pr23135.c":      {},
		"pr53645-2.c":    {},
		"pr53645.c":      {},
		"pr60960.c":      {},
		"pr65427.c":      {},
		"pr70903.c":      {},
		"pr71626-1.c":    {},
		"pr72824-2.c":    {},
		"pr85169.c":      {},
		"pr85331.c":      {},
		"pr92618.c":      {},
		"pr94412.c":      {},
		"pr94524-1.c":    {},
		"pr94524-2.c":    {},
		"pr94591.c":      {},
		"scal-to-vec1.c": {},
		"scal-to-vec2.c": {},
		"scal-to-vec3.c": {},
		"simd-1.c":       {},
		"simd-2.c":       {},
		"simd-4.c":       {},
		"simd-5.c":       {},
		"simd-6.c":       {},

		// VLA
		"20040411-1.c": {},
		"20040423-1.c": {},
		"20041218-2.c": {},
		"920721-2.c":   {},
		"920929-1.c":   {},
		"970217-1.c":   {},

		// Won't fly
		"20031003-1.c": {},

		// Other
		"20000113-1.c":                 {}, //TODO
		"20000217-1.c":                 {}, //TODO
		"20000519-2.c":                 {}, //TODO
		"20000703-1.c":                 {}, //TODO
		"20000722-1.c":                 {}, //TODO
		"20000815-1.c":                 {}, //TODO
		"20000822-1.c":                 {}, //TODO
		"20000910-1.c":                 {}, //TODO
		"20000914-1.c":                 {}, //TODO
		"20000917-1.c":                 {}, //TODO
		"20001101.c":                   {}, //TODO
		"20001203-2.c":                 {}, //TODO
		"20010122-1.c":                 {}, //TODO
		"20010123-1.c":                 {}, //TODO
		"20010209-1.c":                 {}, //TODO
		"20010518-2.c":                 {}, //TODO
		"20010605-1.c":                 {}, //TODO
		"20010605-2.c":                 {}, //TODO
		"20011113-1.c":                 {}, //TODO
		"20020206-1.c":                 {}, //TODO
		"20020206-2.c":                 {}, //TODO
		"20020215-1.c":                 {}, //TODO
		"20020320-1.c":                 {}, //TODO
		"20020404-1.c":                 {}, //TODO
		"20020411-1.c":                 {}, //TODO
		"20020412-1.c":                 {}, //TODO
		"20020418-1.c":                 {}, //TODO
		"20020611-1.c":                 {}, //TODO
		"20021118-1.c":                 {}, //TODO
		"20030109-1.c":                 {}, //TODO
		"20030224-2.c":                 {}, //TODO
		"20030401-1.c":                 {}, //TODO
		"20030501-1.c":                 {}, //TODO
		"20030714-1.c":                 {}, //TODO
		"20030910-1.c":                 {}, //TODO
		"20030916-1.c":                 {}, //TODO
		"20031201-1.c":                 {}, //TODO
		"20031211-1.c":                 {}, //TODO
		"20031211-2.c":                 {}, //TODO
		"20040302-1.c":                 {}, //TODO
		"20040307-1.c":                 {}, //TODO
		"20040308-1.c":                 {}, //TODO
		"20040331-1.c":                 {}, //TODO
		"20040520-1.c":                 {}, //TODO
		"20040629-1.c":                 {}, //TODO
		"20040705-1.c":                 {}, //TODO
		"20040705-2.c":                 {}, //TODO
		"20040709-1.c":                 {}, //TODO
		"20040709-2.c":                 {}, //TODO
		"20041124-1.c":                 {}, //TODO
		"20041201-1.c":                 {}, //TODO
		"20041214-1.c":                 {}, //TODO
		"20050106-1.c":                 {}, //TODO
		"20050121-1.c":                 {}, //TODO
		"20050613-1.c":                 {}, //TODO
		"20050929-1.c":                 {}, //TODO
		"20051012-1.c":                 {}, //TODO
		"20061031-1.c":                 {}, //TODO
		"20061220-1.c":                 {}, //TODO
		"20070614-1.c":                 {}, //TODO
		"20070919-1.c":                 {}, //TODO
		"20071029-1.c":                 {}, //TODO
		"20071202-1.c":                 {}, //TODO
		"20071210-1.c":                 {}, //TODO
		"20071211-1.c":                 {}, //TODO
		"20071220-1.c":                 {}, //TODO
		"20071220-2.c":                 {}, //TODO
		"20080122-1.c":                 {}, //TODO
		"20080519-1.c":                 {}, //TODO
		"20080529-1.c":                 {}, //TODO
		"20081117-1.c":                 {}, //TODO
		"20090113-1.c":                 {}, //TODO
		"20090113-2.c":                 {}, //TODO
		"20090113-3.c":                 {}, //TODO
		"20090219-1.c":                 {}, //TODO
		"20100316-1.c":                 {}, //TODO
		"20120111-1.c":                 {}, //TODO
		"20141107-1.c":                 {}, //TODO
		"20180131-1.c":                 {}, //TODO
		"20180921-1.c":                 {}, //TODO
		"20181120-1.c":                 {}, //TODO
		"20190820-1.c":                 {}, //TODO
		"20190901-1.c":                 {}, //TODO
		"920302-1.c":                   {}, //TODO
		"920415-1.c":                   {}, //TODO
		"920428-2.c":                   {}, //TODO
		"920501-3.c":                   {}, //TODO
		"920501-4.c":                   {}, //TODO
		"920501-5.c":                   {}, //TODO
		"920501-7.c":                   {}, //TODO
		"920612-2.c":                   {}, //TODO
		"920625-1.c":                   {}, //TODO
		"920721-4.c":                   {}, //TODO
		"920731-1.c":                   {}, //TODO
		"920908-1.c":                   {}, //TODO
		"920908-2.c":                   {}, //TODO
		"921016-1.c":                   {}, //TODO
		"921017-1.c":                   {}, //TODO
		"921019-1.c":                   {}, //TODO
		"921204-1.c":                   {}, //TODO
		"921215-1.c":                   {}, //TODO
		"930126-1.c":                   {}, //TODO
		"930406-1.c":                   {}, //TODO
		"930621-1.c":                   {}, //TODO
		"930630-1.c":                   {}, //TODO
		"930718-1.c":                   {}, //TODO
		"930930-1.c":                   {}, //TODO
		"931002-1.c":                   {}, //TODO
		"931004-10.c":                  {}, //TODO
		"931004-12.c":                  {}, //TODO
		"931004-14.c":                  {}, //TODO
		"931004-2.c":                   {}, //TODO
		"931004-4.c":                   {}, //TODO
		"931004-6.c":                   {}, //TODO
		"931004-8.c":                   {}, //TODO
		"931031-1.c":                   {}, //TODO
		"931110-1.c":                   {}, //TODO
		"941202-1.c":                   {}, //TODO
		"950512-1.c":                   {}, //TODO
		"950628-1.c":                   {}, //TODO
		"950906-1.c":                   {}, //TODO
		"960301-1.c":                   {}, //TODO
		"960326-1.c":                   {}, //TODO
		"960416-1.c":                   {}, //TODO
		"960608-1.c":                   {}, //TODO
		"980526-1.c":                   {}, //TODO
		"980602-2.c":                   {}, //TODO
		"980929-1.c":                   {}, //TODO
		"990208-1.c":                   {}, //TODO
		"990222-1.c":                   {}, //TODO
		"990326-1.c":                   {}, //TODO
		"990525-1.c":                   {}, //TODO
		"991118-1.c":                   {}, //TODO
		"991228-1.c":                   {}, //TODO
		"alias-2.c":                    {}, //TODO
		"alias-access-path-2.c":        {}, //TODO
		"align-nest.c":                 {}, //TODO
		"alloca-1.c":                   {}, //TODO
		"anon-1.c":                     {}, //TODO
		"bcp-1.c":                      {}, //TODO
		"bf-layout-1.c":                {}, //TODO
		"bf-pack-1.c":                  {}, //TODO
		"bf-sign-1.c":                  {}, //TODO
		"bf-sign-2.c":                  {}, //TODO
		"bf64-1.c":                     {}, //TODO
		"bitfld-1.c":                   {}, //TODO
		"bitfld-2.c":                   {}, //TODO
		"bitfld-3.c":                   {}, //TODO
		"bitfld-4.c":                   {}, //TODO
		"bitfld-5.c":                   {}, //TODO
		"bitfld-6.c":                   {}, //TODO
		"bitfld-7.c":                   {}, //TODO
		"bitfld-8.c":                   {}, //TODO
		"bitfld-9.c":                   {}, //TODO
		"bswap-2.c":                    {}, //TODO
		"bswap-3.c":                    {}, //TODO
		"builtin-bitops-1.c":           {}, //TODO
		"builtin-prefetch-4.c":         {}, //TODO
		"builtin-types-compatible-p.c": {}, //TODO
		"comp-goto-1.c":                {}, //TODO
		"comp-goto-2.c":                {}, //TODO
		"complex-1.c":                  {}, //TODO
		"complex-2.c":                  {}, //TODO
		"complex-5.c":                  {}, //TODO
		"complex-6.c":                  {}, //TODO
		"complex-7.c":                  {}, //TODO
		"compndlit-1.c":                {}, //TODO
		"divconst-2.c":                 {}, //TODO
		"extzvsi.c":                    {}, //TODO
		"fp-cmp-4.c":                   {}, //TODO
		"fp-cmp-4l.c":                  {}, //TODO
		"loop-15.c":                    {}, //TODO
		"medce-1.c":                    {}, //TODO
		"memchr-1.c":                   {}, //TODO
		"memcpy-2.c":                   {}, //TODO
		"memset-1.c":                   {}, //TODO
		"memset-2.c":                   {}, //TODO
		"memset-3.c":                   {}, //TODO
		"nest-align-1.c":               {}, //TODO
		"nest-stdar-1.c":               {}, //TODO
		"nestfunc-1.c":                 {}, //TODO
		"nestfunc-2.c":                 {}, //TODO
		"nestfunc-3.c":                 {}, //TODO
		"nestfunc-5.c":                 {}, //TODO
		"nestfunc-6.c":                 {}, //TODO
		"nestfunc-7.c":                 {}, //TODO
		"p18298.c":                     {}, //TODO
		"pr15296.c":                    {}, //TODO
		"pr19449.c":                    {}, //TODO
		"pr19515.c":                    {}, //TODO
		"pr19687.c":                    {}, //TODO
		"pr19689.c":                    {}, //TODO
		"pr22061-3.c":                  {}, //TODO
		"pr22061-4.c":                  {}, //TODO
		"pr22098-1.c":                  {}, //TODO
		"pr22098-2.c":                  {}, //TODO
		"pr22098-3.c":                  {}, //TODO
		"pr22141-1.c":                  {}, //TODO
		"pr22141-2.c":                  {}, //TODO
		"pr23324.c":                    {}, //TODO
		"pr24135.c":                    {}, //TODO
		"pr28865.c":                    {}, //TODO
		"pr28982b.c":                   {}, //TODO
		"pr30185.c":                    {}, //TODO
		"pr30778.c":                    {}, //TODO
		"pr31136.c":                    {}, //TODO
		"pr31169.c":                    {}, //TODO
		"pr31448-2.c":                  {}, //TODO
		"pr31448.c":                    {}, //TODO
		"pr32244-1.c":                  {}, //TODO
		"pr34154.c":                    {}, //TODO
		"pr34971.c":                    {}, //TODO
		"pr36038.c":                    {}, //TODO
		"pr37573.c":                    {}, //TODO
		"pr37882.c":                    {}, //TODO
		"pr38151.c":                    {}, //TODO
		"pr38422.c":                    {}, //TODO
		"pr38969.c":                    {}, //TODO
		"pr39339.c":                    {}, //TODO
		"pr40404.c":                    {}, //TODO
		"pr40493.c":                    {}, //TODO
		"pr41239.c":                    {}, //TODO
		"pr42691.c":                    {}, //TODO
		"pr43987.c":                    {}, //TODO
		"pr44164.c":                    {}, //TODO
		"pr44555.c":                    {}, //TODO
		"pr44575.c":                    {}, //TODO
		"pr46309.c":                    {}, //TODO
		"pr47148.c":                    {}, //TODO
		"pr47155.c":                    {}, //TODO
		"pr48973-1.c":                  {}, //TODO
		"pr48973-2.c":                  {}, //TODO
		"pr49073.c":                    {}, //TODO
		"pr49123.c":                    {}, //TODO
		"pr49218.c":                    {}, //TODO
		"pr49644.c":                    {}, //TODO
		"pr49768.c":                    {}, //TODO
		"pr51447.c":                    {}, //TODO
		"pr52209.c":                    {}, //TODO
		"pr52979-1.c":                  {}, //TODO
		"pr52979-2.c":                  {}, //TODO
		"pr55750.c":                    {}, //TODO
		"pr56837.c":                    {}, //TODO
		"pr57130.c":                    {}, //TODO
		"pr57344-1.c":                  {}, //TODO
		"pr57344-2.c":                  {}, //TODO
		"pr57344-3.c":                  {}, //TODO
		"pr57344-4.c":                  {}, //TODO
		"pr57568.c":                    {}, //TODO
		"pr57861.c":                    {}, //TODO
		"pr57876.c":                    {}, //TODO
		"pr57877.c":                    {}, //TODO
		"pr58277-2.c":                  {}, //TODO
		"pr58385.c":                    {}, //TODO
		"pr58419.c":                    {}, //TODO
		"pr58564.c":                    {}, //TODO
		"pr58570.c":                    {}, //TODO
		"pr58726.c":                    {}, //TODO
		"pr58831.c":                    {}, //TODO
		"pr58984.c":                    {}, //TODO
		"pr59388.c":                    {}, //TODO
		"pr60017.c":                    {}, //TODO
		"pr61375.c":                    {}, //TODO
		"pr63302.c":                    {}, //TODO
		"pr64756.c":                    {}, //TODO
		"pr65215-3.c":                  {}, //TODO
		"pr65215-4.c":                  {}, //TODO
		"pr66556.c":                    {}, //TODO
		"pr67714.c":                    {}, //TODO
		"pr68185.c":                    {}, //TODO
		"pr68249.c":                    {}, //TODO
		"pr68250.c":                    {}, //TODO
		"pr68321.c":                    {}, //TODO
		"pr68506.c":                    {}, //TODO
		"pr68532.c":                    {}, //TODO
		"pr69320-2.c":                  {}, //TODO
		"pr69320-4.c":                  {}, //TODO
		"pr69691.c":                    {}, //TODO
		"pr70127.c":                    {}, //TODO
		"pr70460.c":                    {}, //TODO
		"pr70566.c":                    {}, //TODO
		"pr70586.c":                    {}, //TODO
		"pr70602.c":                    {}, //TODO
		"pr71083.c":                    {}, //TODO
		"pr71494.c":                    {}, //TODO
		"pr71700.c":                    {}, //TODO
		"pr78170.c":                    {}, //TODO
		"pr78436.c":                    {}, //TODO
		"pr78438.c":                    {}, //TODO
		"pr78477.c":                    {}, //TODO
		"pr78559.c":                    {}, //TODO
		"pr78675.c":                    {}, //TODO
		"pr79286.c":                    {}, //TODO
		"pr79737-1.c":                  {}, //TODO
		"pr79737-2.c":                  {}, //TODO
		"pr80692.c":                    {}, //TODO
		"pr81423.c":                    {}, //TODO
		"pr81555.c":                    {}, //TODO
		"pr81556.c":                    {}, //TODO
		"pr82192.c":                    {}, //TODO
		"pr82387.c":                    {}, //TODO
		"pr82524.c":                    {}, //TODO
		"pr83362.c":                    {}, //TODO
		"pr83383.c":                    {}, //TODO
		"pr84235.c":                    {}, //TODO
		"pr84524.c":                    {}, //TODO
		"pr85529-1.c":                  {}, //TODO
		"pr85582-2.c":                  {}, //TODO
		"pr85582-3.c":                  {}, //TODO
		"pr86492.c":                    {}, //TODO
		"pr86659-1.c":                  {}, //TODO
		"pr86659-2.c":                  {}, //TODO
		"pr87053.c":                    {}, //TODO
		"pr88739.c":                    {}, //TODO
		"pr88904.c":                    {}, //TODO
		"pr89195.c":                    {}, //TODO
		"pr89369.c":                    {}, //TODO
		"pr89434.c":                    {}, //TODO
		"pr90949.c":                    {}, //TODO
		"pr91137.c":                    {}, //TODO
		"pr91597.c":                    {}, //TODO
		"pr92904.c":                    {}, //TODO
		"pr93213.c":                    {}, //TODO
		"pr93434.c":                    {}, //TODO
		"pr93908.c":                    {}, //TODO
		"pr93945.c":                    {}, //TODO
		"pr94130.c":                    {}, //TODO
		"pr94724.c":                    {}, //TODO
		"pr94734.c":                    {}, //TODO
		"pr94809.c":                    {}, //TODO
		"pr96549.c":                    {}, //TODO
		"pr97421-1.c":                  {}, //TODO
		"pr97764.c":                    {}, //TODO
		"pr98366.c":                    {}, //TODO
		"pr98474.c":                    {}, //TODO
		"restrict-1.c":                 {}, //TODO
		"return-addr.c":                {}, //TODO
		"ssad-run.c":                   {}, //TODO
		"stdarg-3.c":                   {}, //TODO
		"stkalign.c":                   {}, //TODO
		"strcpy-1.c":                   {}, //TODO
		"strcpy-2.c":                   {}, //TODO
		"strct-stdarg-1.c":             {}, //TODO
		"strct-varg-1.c":               {}, //TODO
		"strlen-2.c":                   {}, //TODO
		"strlen-3.c":                   {}, //TODO
		"strlen-4.c":                   {}, //TODO
		"strlen-5.c":                   {}, //TODO
		"struct-ini-1.c":               {}, //TODO
		"struct-ini-2.c":               {}, //TODO
		"struct-ini-3.c":               {}, //TODO
		"usad-run.c":                   {}, //TODO
		"va-arg-15.c":                  {}, //TODO
		"va-arg-22.c":                  {}, //TODO
		"va-arg-pack-1.c":              {}, //TODO
		"vrp-7.c":                      {}, //TODO
		"widechar-3.c":                 {}, //TODO
		"zero-struct-1.c":              {}, //TODO
		"zero-struct-2.c":              {}, //TODO
	}
	blacklistTCC := map[string]struct{}{
		// asm
		"99_fastcall.c": {},

		"76_dollars_in_identifiers.c": {}, //TODO

		// Init
		"75_array_in_struct_init.c": {}, //TODO

		// Other
		"73_arm64.c":               {}, //TODO
		"80_flexarray.c":           {}, //TODO
		"81_types.c":               {}, //TODO
		"87_dead_code.c":           {}, //TODO
		"88_codeopt.c":             {}, //TODO
		"89_nocode_wanted.c":       {}, //TODO
		"90_struct-init.c":         {}, //TODO
		"92_enum_bitfield.c":       {}, //TODO
		"93_integer_promotion.c":   {}, //TODO
		"94_generic.c":             {}, //TODO
		"95_bitfields.c":           {}, //TODO
		"97_utf8_string_literal.c": {}, //TODO

	}
	switch fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) {
	case "darwin/amd64":
	case "darwin/arm64":
	case "freebsd/386":
	case "freebsd/amd64":
	case "linux/386":
		// _Float128
		blacklistGCC["nest-align-1.c"] = struct{}{}
		blacklistGCC["pr44468.c"] = struct{}{}
		blacklistGCC["strcmp-1.c"] = struct{}{}
		blacklistGCC["strlen-1.c"] = struct{}{}
		blacklistGCC["strncmp-1.c"] = struct{}{}
	case "linux/s390x":
	case "netbsd/amd64":
	case "openbsd/amd64":
	case "windows/386":
	case "windows/amd64":
	case "windows/arm64":
	}
	for _, v := range []struct {
		dir       string
		blacklist map[string]struct{}
	}{
		//TODO {"CompCert-3.6/test/c", blacklistCompCert},
		//TODO {"ccgo", nil},
		//TODO {"gcc-9.1.0/gcc/testsuite/gcc.c-torture", blacklistGCC},
		//TODO {"github.com/AbsInt/CompCert/test/c", blacklistCompCert},
		//TODO {"github.com/cxgo", nil},
		{"github.com/gcc-mirror/gcc/gcc/testsuite", blacklistGCC},
		//TODO {"github.com/vnmakarov", nil},
		//TODO {"sqlite-amalgamation-3380100", nil},
		{"tcc-0.9.27/tests/tests2", blacklistTCC},
		//TODO {"benchmarksgame-team.pages.debian.net", nil},
	} {
		t.Run(v.dir, func(t *testing.T) {
			testCompile(t, tmp, "assets/"+v.dir, v.blacklist, g)
		})
	}
}

func testCompile(t *testing.T, tmp, dir string, blacklist map[string]struct{}, g *golden) {
	p := newParallel()

	defer func() { p.close(t) }()

	p.err(cfsWalk(dir, func(pth string, fi os.FileInfo) error {
		if fi.IsDir() {
			return nil
		}

		if filepath.Ext(pth) != ".c" {
			return nil
		}

		p.file()
		switch {
		case re != nil:
			if !re.MatchString(pth) {
				p.skip()
				return nil
			}
		default:
			if _, ok := blacklist[filepath.Base(pth)]; ok {
				p.skip()
				return nil
			}
		}

		apth := pth
		afi := fi
		p.exec(func() error {
			if *oTrace {
				fmt.Fprintln(os.Stderr, apth)
			}

			func() {
				defer func() {
					if err := recover(); err != nil {
						err = fmt.Errorf("%v: PANIC: %v", filepath.Base(apth), err)
						trc("%v: PANIC: %v\n%s", apth, err, debug.Stack())
						os.Exit(1)
					}
				}()

				ofn := filepath.Join(tmp, fmt.Sprintf("%d.go", p.id()))

				defer os.Remove(ofn)

				var out bytes.Buffer
				task := NewTask(goos, goarch, []string{"ccgo", "-o", ofn, "-c", apth}, &out, &out, cfs)
				ccgoErr := task.Main()
				if ccgoErr == nil {
					p.ok()
					g.w("%s\n", apth)
					return
				}

				checkFailOk(t, p, errorf("%v: %v", filepath.Base(apth), ccgoErr), tmp, apth, ofn, afi, task)
			}()
			return nil
		})
		return nil
	}))
}

func checkFailOk(t *testing.T, p *parallel, ccgoErr error, tmp, src, ofn string, fi os.FileInfo, task *Task) {
	f, err := cfs.Open(src)
	if err != nil {
		p.err(err)
		return
	}

	defer f.Close()

	b := make([]byte, fi.Size())
	if n, _ := f.Read(b); int64(n) != fi.Size() {
		p.err(errorf("%v: short read", src))
		return
	}

	fn := filepath.Join(tmp, filepath.Base(src))
	if err := os.WriteFile(fn, b, 0660); err != nil {
		p.err(errorf("%v: %v", src, err))
		return
	}

	defer os.Remove(fn)

	cfg := task.cfg
	cmd := exec.Command(cfg.CC, "-c", "-o", ofn, "-w", fn)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		t.Logf("%v: skip: %v: %s %v", src, cfg.CC, buf.Bytes(), err)
		p.skip()
		return
	}

	p.fail()
	p.err(ccgoErr)
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

func shell(echo bool, cmd string, args ...string) ([]byte, error) {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}

	wd, err := absCwd()
	if err != nil {
		return nil, err
	}

	if echo {
		fmt.Printf("execute %s %q in %s\n", cmd, args, wd)
	}
	var b echoWriter
	b.silent = !echo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdout = &b
	c.Stderr = &b
	err = c.Run()
	return b.w.Bytes(), err
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

		// blacklistCompCert := map[string]struct{}{}
		blacklistGCC := map[string]struct{}{
			// Assertions are deprecated, not supported.
			"950919-1.c": {},

			// Wrong prototype trips type checking.
			"va-arg-2.c": {},

			// Vector
			"20050316-1.c":   {},
			"20050316-2.c":   {},
			"20050316-3.c":   {},
			"20050604-1.c":   {},
			"20050607-1.c":   {},
			"20060420-1.c":   {},
			"pr23135.c":      {},
			"pr53645-2.c":    {},
			"pr53645.c":      {},
			"pr60960.c":      {},
			"pr65427.c":      {},
			"pr70903.c":      {},
			"pr71626-1.c":    {},
			"pr72824-2.c":    {},
			"pr85169.c":      {},
			"pr85331.c":      {},
			"pr92618.c":      {},
			"pr94412.c":      {},
			"pr94524-1.c":    {},
			"pr94524-2.c":    {},
			"pr94591.c":      {},
			"scal-to-vec1.c": {},
			"scal-to-vec2.c": {},
			"scal-to-vec3.c": {},
			"simd-1.c":       {},
			"simd-2.c":       {},
			"simd-4.c":       {},
			"simd-5.c":       {},
			"simd-6.c":       {},

			// Other
			"20000113-1.c":                 {}, //TODO
			"20000217-1.c":                 {}, //TODO
			"20000519-2.c":                 {}, //TODO
			"20000703-1.c":                 {}, //TODO
			"20000722-1.c":                 {}, //TODO
			"20000815-1.c":                 {}, //TODO
			"20000822-1.c":                 {}, //TODO
			"20000910-1.c":                 {}, //TODO
			"20000914-1.c":                 {}, //TODO
			"20000917-1.c":                 {}, //TODO
			"20001009-2.c":                 {}, //TODO
			"20001101.c":                   {}, //TODO
			"20001203-2.c":                 {}, //TODO
			"20010122-1.c":                 {}, //TODO
			"20010123-1.c":                 {}, //TODO
			"20010209-1.c":                 {}, //TODO
			"20010518-2.c":                 {}, //TODO
			"20010605-1.c":                 {}, //TODO
			"20010605-2.c":                 {}, //TODO
			"20010904-1.c":                 {}, //TODO
			"20010904-2.c":                 {}, //TODO
			"20011113-1.c":                 {}, //TODO
			"20020107-1.c":                 {}, //TODO
			"20020206-1.c":                 {}, //TODO
			"20020206-2.c":                 {}, //TODO
			"20020215-1.c":                 {}, //TODO
			"20020314-1.c":                 {}, //TODO
			"20020320-1.c":                 {}, //TODO
			"20020404-1.c":                 {}, //TODO
			"20020411-1.c":                 {}, //TODO
			"20020412-1.c":                 {}, //TODO
			"20020418-1.c":                 {}, //TODO
			"20020611-1.c":                 {}, //TODO
			"20021113-1.c":                 {}, //TODO
			"20021118-1.c":                 {}, //TODO
			"20021127-1.c":                 {}, //TODO
			"20030109-1.c":                 {}, //TODO
			"20030222-1.c":                 {}, //TODO
			"20030224-2.c":                 {}, //TODO
			"20030323-1.c":                 {}, //TODO
			"20030330-1.c":                 {}, //TODO
			"20030401-1.c":                 {}, //TODO
			"20030408-1.c":                 {}, //TODO
			"20030501-1.c":                 {}, //TODO
			"20030714-1.c":                 {}, //TODO
			"20030811-1.c":                 {}, //TODO
			"20030909-1.c":                 {}, //TODO
			"20030910-1.c":                 {}, //TODO
			"20030916-1.c":                 {}, //TODO
			"20031003-1.c":                 {}, //TODO
			"20031201-1.c":                 {}, //TODO
			"20031211-1.c":                 {}, //TODO
			"20031211-2.c":                 {}, //TODO
			"20040223-1.c":                 {}, //TODO
			"20040302-1.c":                 {}, //TODO
			"20040307-1.c":                 {}, //TODO
			"20040308-1.c":                 {}, //TODO
			"20040331-1.c":                 {}, //TODO
			"20040411-1.c":                 {}, //TODO
			"20040423-1.c":                 {}, //TODO
			"20040520-1.c":                 {}, //TODO
			"20040629-1.c":                 {}, //TODO
			"20040704-1.c":                 {}, //TODO
			"20040705-1.c":                 {}, //TODO
			"20040705-2.c":                 {}, //TODO
			"20040709-1.c":                 {}, //TODO
			"20040709-2.c":                 {}, //TODO
			"20040709-3.c":                 {}, //TODO
			"20040811-1.c":                 {}, //TODO
			"20041124-1.c":                 {}, //TODO
			"20041201-1.c":                 {}, //TODO
			"20041214-1.c":                 {}, //TODO
			"20041218-2.c":                 {}, //TODO
			"20050106-1.c":                 {}, //TODO
			"20050121-1.c":                 {}, //TODO
			"20050203-1.c":                 {}, //TODO
			"20050613-1.c":                 {}, //TODO
			"20050929-1.c":                 {}, //TODO
			"20051012-1.c":                 {}, //TODO
			"20051113-1.c":                 {}, //TODO
			"20061031-1.c":                 {}, //TODO
			"20061220-1.c":                 {}, //TODO
			"20070614-1.c":                 {}, //TODO
			"20070824-1.c":                 {}, //TODO
			"20070919-1.c":                 {}, //TODO
			"20071029-1.c":                 {}, //TODO
			"20071202-1.c":                 {}, //TODO
			"20071210-1.c":                 {}, //TODO
			"20071211-1.c":                 {}, //TODO
			"20071220-1.c":                 {}, //TODO
			"20071220-2.c":                 {}, //TODO
			"20080122-1.c":                 {}, //TODO
			"20080222-1.c":                 {}, //TODO
			"20080502-1.c":                 {}, //TODO
			"20080519-1.c":                 {}, //TODO
			"20080529-1.c":                 {}, //TODO
			"20081117-1.c":                 {}, //TODO
			"20090113-1.c":                 {}, //TODO
			"20090113-2.c":                 {}, //TODO
			"20090113-3.c":                 {}, //TODO
			"20090219-1.c":                 {}, //TODO
			"20100316-1.c":                 {}, //TODO
			"20101011-1.c":                 {}, //TODO
			"20111208-1.c":                 {}, //TODO
			"20120111-1.c":                 {}, //TODO
			"20141107-1.c":                 {}, //TODO
			"20180131-1.c":                 {}, //TODO
			"20180921-1.c":                 {}, //TODO
			"20181120-1.c":                 {}, //TODO
			"20190820-1.c":                 {}, //TODO
			"20190901-1.c":                 {}, //TODO
			"920302-1.c":                   {}, //TODO
			"920415-1.c":                   {}, //TODO
			"920428-2.c":                   {}, //TODO
			"920501-1.c":                   {}, //TODO
			"920501-3.c":                   {}, //TODO
			"920501-4.c":                   {}, //TODO
			"920501-5.c":                   {}, //TODO
			"920501-6.c":                   {}, //TODO
			"920501-7.c":                   {}, //TODO
			"920612-2.c":                   {}, //TODO
			"920625-1.c":                   {}, //TODO
			"920721-2.c":                   {}, //TODO
			"920721-4.c":                   {}, //TODO
			"920728-1.c":                   {}, //TODO
			"920731-1.c":                   {}, //TODO
			"920908-1.c":                   {}, //TODO
			"920908-2.c":                   {}, //TODO
			"920929-1.c":                   {}, //TODO
			"921016-1.c":                   {}, //TODO
			"921017-1.c":                   {}, //TODO
			"921019-1.c":                   {}, //TODO
			"921202-1.c":                   {}, //TODO
			"921204-1.c":                   {}, //TODO
			"921208-2.c":                   {}, //TODO
			"921215-1.c":                   {}, //TODO
			"930126-1.c":                   {}, //TODO
			"930406-1.c":                   {}, //TODO
			"930621-1.c":                   {}, //TODO
			"930630-1.c":                   {}, //TODO
			"930718-1.c":                   {}, //TODO
			"930930-1.c":                   {}, //TODO
			"931002-1.c":                   {}, //TODO
			"931004-10.c":                  {}, //TODO
			"931004-12.c":                  {}, //TODO
			"931004-14.c":                  {}, //TODO
			"931004-2.c":                   {}, //TODO
			"931004-4.c":                   {}, //TODO
			"931004-6.c":                   {}, //TODO
			"931004-8.c":                   {}, //TODO
			"931031-1.c":                   {}, //TODO
			"931110-1.c":                   {}, //TODO
			"941202-1.c":                   {}, //TODO
			"950221-1.c":                   {}, //TODO
			"950512-1.c":                   {}, //TODO
			"950628-1.c":                   {}, //TODO
			"950906-1.c":                   {}, //TODO
			"960301-1.c":                   {}, //TODO
			"960312-1.c":                   {}, //TODO
			"960326-1.c":                   {}, //TODO
			"960405-1.c":                   {}, //TODO
			"960416-1.c":                   {}, //TODO
			"960512-1.c":                   {}, //TODO
			"960608-1.c":                   {}, //TODO
			"970217-1.c":                   {}, //TODO
			"980526-1.c":                   {}, //TODO
			"980602-2.c":                   {}, //TODO
			"980929-1.c":                   {}, //TODO
			"990130-1.c":                   {}, //TODO
			"990208-1.c":                   {}, //TODO
			"990222-1.c":                   {}, //TODO
			"990326-1.c":                   {}, //TODO
			"990413-2.c":                   {}, //TODO
			"990525-1.c":                   {}, //TODO
			"991014-1.c":                   {}, //TODO
			"991112-1.c":                   {}, //TODO
			"991118-1.c":                   {}, //TODO
			"991228-1.c":                   {}, //TODO
			"alias-2.c":                    {}, //TODO
			"alias-3.c":                    {}, //TODO
			"alias-4.c":                    {}, //TODO
			"alias-access-path-2.c":        {}, //TODO
			"align-3.c":                    {}, //TODO
			"align-nest.c":                 {}, //TODO
			"alloca-1.c":                   {}, //TODO
			"anon-1.c":                     {}, //TODO
			"bcp-1.c":                      {}, //TODO
			"bf-layout-1.c":                {}, //TODO
			"bf-pack-1.c":                  {}, //TODO
			"bf-sign-1.c":                  {}, //TODO
			"bf-sign-2.c":                  {}, //TODO
			"bf64-1.c":                     {}, //TODO
			"bitfld-1.c":                   {}, //TODO
			"bitfld-2.c":                   {}, //TODO
			"bitfld-3.c":                   {}, //TODO
			"bitfld-4.c":                   {}, //TODO
			"bitfld-5.c":                   {}, //TODO
			"bitfld-6.c":                   {}, //TODO
			"bitfld-7.c":                   {}, //TODO
			"bitfld-8.c":                   {}, //TODO
			"bitfld-9.c":                   {}, //TODO
			"bswap-2.c":                    {}, //TODO
			"bswap-3.c":                    {}, //TODO
			"built-in-setjmp.c":            {}, //TODO
			"builtin-bitops-1.c":           {}, //TODO
			"builtin-prefetch-4.c":         {}, //TODO
			"builtin-types-compatible-p.c": {}, //TODO
			"call-trap-1.c":                {}, //TODO
			"comp-goto-1.c":                {}, //TODO
			"comp-goto-2.c":                {}, //TODO
			"compare-fp-1.c":               {}, //TODO
			"complex-1.c":                  {}, //TODO
			"complex-2.c":                  {}, //TODO
			"complex-4.c":                  {}, //TODO
			"complex-5.c":                  {}, //TODO
			"complex-6.c":                  {}, //TODO
			"complex-7.c":                  {}, //TODO
			"compndlit-1.c":                {}, //TODO
			"divconst-2.c":                 {}, //TODO
			"eeprof-1.c":                   {}, //TODO
			"extzvsi.c":                    {}, //TODO
			"ffs-1.c":                      {}, //TODO
			"ffs-2.c":                      {}, //TODO
			"fp-cmp-4.c":                   {}, //TODO
			"fp-cmp-4f.c":                  {}, //TODO
			"fp-cmp-4l.c":                  {}, //TODO
			"fp-cmp-5.c":                   {}, //TODO
			"fp-cmp-7.c":                   {}, //TODO
			"fp-cmp-8.c":                   {}, //TODO
			"fp-cmp-8f.c":                  {}, //TODO
			"fp-cmp-8l.c":                  {}, //TODO
			"fprintf-2.c":                  {}, //TODO
			"frame-address.c":              {}, //TODO
			"ipa-sra-2.c":                  {}, //TODO
			"loop-15.c":                    {}, //TODO
			"loop-2e.c":                    {}, //TODO
			"lto-tbaa-1.c":                 {}, //TODO
			"medce-1.c":                    {}, //TODO
			"memchr-1.c":                   {}, //TODO
			"memcpy-2.c":                   {}, //TODO
			"memset-1.c":                   {}, //TODO
			"memset-2.c":                   {}, //TODO
			"memset-3.c":                   {}, //TODO
			"mzero4.c":                     {}, //TODO
			"nest-align-1.c":               {}, //TODO
			"nest-stdar-1.c":               {}, //TODO
			"nestfunc-1.c":                 {}, //TODO
			"nestfunc-2.c":                 {}, //TODO
			"nestfunc-3.c":                 {}, //TODO
			"nestfunc-5.c":                 {}, //TODO
			"nestfunc-6.c":                 {}, //TODO
			"nestfunc-7.c":                 {}, //TODO
			"p18298.c":                     {}, //TODO
			"pr15296.c":                    {}, //TODO
			"pr17078-1.c":                  {}, //TODO
			"pr17377.c":                    {}, //TODO
			"pr19449.c":                    {}, //TODO
			"pr19515.c":                    {}, //TODO
			"pr19687.c":                    {}, //TODO
			"pr19689.c":                    {}, //TODO
			"pr22061-1.c":                  {}, //TODO
			"pr22061-3.c":                  {}, //TODO
			"pr22061-4.c":                  {}, //TODO
			"pr22098-1.c":                  {}, //TODO
			"pr22098-2.c":                  {}, //TODO
			"pr22098-3.c":                  {}, //TODO
			"pr22141-1.c":                  {}, //TODO
			"pr22141-2.c":                  {}, //TODO
			"pr23324.c":                    {}, //TODO
			"pr23467.c":                    {}, //TODO
			"pr24135.c":                    {}, //TODO
			"pr24851.c":                    {}, //TODO
			"pr28289.c":                    {}, //TODO
			"pr28865.c":                    {}, //TODO
			"pr28982b.c":                   {}, //TODO
			"pr30185.c":                    {}, //TODO
			"pr30778.c":                    {}, //TODO
			"pr31136.c":                    {}, //TODO
			"pr31169.c":                    {}, //TODO
			"pr31448-2.c":                  {}, //TODO
			"pr31448.c":                    {}, //TODO
			"pr32244-1.c":                  {}, //TODO
			"pr34154.c":                    {}, //TODO
			"pr34176.c":                    {}, //TODO
			"pr34768-1.c":                  {}, //TODO
			"pr34768-2.c":                  {}, //TODO
			"pr34971.c":                    {}, //TODO
			"pr35456.c":                    {}, //TODO
			"pr36038.c":                    {}, //TODO
			"pr36321.c":                    {}, //TODO
			"pr37573.c":                    {}, //TODO
			"pr37780.c":                    {}, //TODO
			"pr37882.c":                    {}, //TODO
			"pr38016.c":                    {}, //TODO
			"pr38051.c":                    {}, //TODO
			"pr38151.c":                    {}, //TODO
			"pr38422.c":                    {}, //TODO
			"pr38533.c":                    {}, //TODO
			"pr38969.c":                    {}, //TODO
			"pr39100.c":                    {}, //TODO
			"pr39228.c":                    {}, //TODO
			"pr39339.c":                    {}, //TODO
			"pr40022.c":                    {}, //TODO
			"pr40404.c":                    {}, //TODO
			"pr40493.c":                    {}, //TODO
			"pr40657.c":                    {}, //TODO
			"pr41239.c":                    {}, //TODO
			"pr41395-1.c":                  {}, //TODO
			"pr41395-2.c":                  {}, //TODO
			"pr41463.c":                    {}, //TODO
			"pr41935.c":                    {}, //TODO
			"pr42614.c":                    {}, //TODO
			"pr42691.c":                    {}, //TODO
			"pr43220.c":                    {}, //TODO
			"pr43269.c":                    {}, //TODO
			"pr43385.c":                    {}, //TODO
			"pr43560.c":                    {}, //TODO
			"pr43987.c":                    {}, //TODO
			"pr44164.c":                    {}, //TODO
			"pr44555.c":                    {}, //TODO
			"pr44575.c":                    {}, //TODO
			"pr44852.c":                    {}, //TODO
			"pr45695.c":                    {}, //TODO
			"pr46309.c":                    {}, //TODO
			"pr47148.c":                    {}, //TODO
			"pr47155.c":                    {}, //TODO
			"pr47237.c":                    {}, //TODO
			"pr47925.c":                    {}, //TODO
			"pr48571-1.c":                  {}, //TODO
			"pr48973-1.c":                  {}, //TODO
			"pr48973-2.c":                  {}, //TODO
			"pr49073.c":                    {}, //TODO
			"pr49123.c":                    {}, //TODO
			"pr49218.c":                    {}, //TODO
			"pr49279.c":                    {}, //TODO
			"pr49390.c":                    {}, //TODO
			"pr49644.c":                    {}, //TODO
			"pr49768.c":                    {}, //TODO
			"pr50310.c":                    {}, //TODO
			"pr51447.c":                    {}, //TODO
			"pr51581-1.c":                  {}, //TODO
			"pr51581-2.c":                  {}, //TODO
			"pr51877.c":                    {}, //TODO
			"pr51933.c":                    {}, //TODO
			"pr52209.c":                    {}, //TODO
			"pr52286.c":                    {}, //TODO
			"pr52979-1.c":                  {}, //TODO
			"pr52979-2.c":                  {}, //TODO
			"pr54471.c":                    {}, //TODO
			"pr55750.c":                    {}, //TODO
			"pr56205.c":                    {}, //TODO
			"pr56837.c":                    {}, //TODO
			"pr56866.c":                    {}, //TODO
			"pr56962.c":                    {}, //TODO
			"pr56982.c":                    {}, //TODO
			"pr57130.c":                    {}, //TODO
			"pr57344-1.c":                  {}, //TODO
			"pr57344-2.c":                  {}, //TODO
			"pr57344-3.c":                  {}, //TODO
			"pr57344-4.c":                  {}, //TODO
			"pr57568.c":                    {}, //TODO
			"pr57861.c":                    {}, //TODO
			"pr57876.c":                    {}, //TODO
			"pr57877.c":                    {}, //TODO
			"pr58277-1.c":                  {}, //TODO
			"pr58277-2.c":                  {}, //TODO
			"pr58385.c":                    {}, //TODO
			"pr58419.c":                    {}, //TODO
			"pr58431.c":                    {}, //TODO
			"pr58564.c":                    {}, //TODO
			"pr58570.c":                    {}, //TODO
			"pr58726.c":                    {}, //TODO
			"pr58831.c":                    {}, //TODO
			"pr58984.c":                    {}, //TODO
			"pr59388.c":                    {}, //TODO
			"pr60003.c":                    {}, //TODO
			"pr60017.c":                    {}, //TODO
			"pr61375.c":                    {}, //TODO
			"pr61725.c":                    {}, //TODO
			"pr63302.c":                    {}, //TODO
			"pr63641.c":                    {}, //TODO
			"pr63843.c":                    {}, //TODO
			"pr64006.c":                    {}, //TODO
			"pr64242.c":                    {}, //TODO
			"pr64756.c":                    {}, //TODO
			"pr65053-2.c":                  {}, //TODO
			"pr65170.c":                    {}, //TODO
			"pr65215-3.c":                  {}, //TODO
			"pr65215-4.c":                  {}, //TODO
			"pr65648.c":                    {}, //TODO
			"pr65956.c":                    {}, //TODO
			"pr66556.c":                    {}, //TODO
			"pr67037.c":                    {}, //TODO
			"pr67714.c":                    {}, //TODO
			"pr68185.c":                    {}, //TODO
			"pr68249.c":                    {}, //TODO
			"pr68250.c":                    {}, //TODO
			"pr68321.c":                    {}, //TODO
			"pr68328.c":                    {}, //TODO
			"pr68381.c":                    {}, //TODO
			"pr68506.c":                    {}, //TODO
			"pr68532.c":                    {}, //TODO
			"pr69320-2.c":                  {}, //TODO
			"pr69320-4.c":                  {}, //TODO
			"pr69691.c":                    {}, //TODO
			"pr70127.c":                    {}, //TODO
			"pr70460.c":                    {}, //TODO
			"pr70566.c":                    {}, //TODO
			"pr70586.c":                    {}, //TODO
			"pr70602.c":                    {}, //TODO
			"pr71083.c":                    {}, //TODO
			"pr71494.c":                    {}, //TODO
			"pr71554.c":                    {}, //TODO
			"pr71626-2.c":                  {}, //TODO
			"pr71700.c":                    {}, //TODO
			"pr77766.c":                    {}, //TODO
			"pr77767.c":                    {}, //TODO
			"pr78170.c":                    {}, //TODO
			"pr78436.c":                    {}, //TODO
			"pr78438.c":                    {}, //TODO
			"pr78477.c":                    {}, //TODO
			"pr78559.c":                    {}, //TODO
			"pr78675.c":                    {}, //TODO
			"pr78726.c":                    {}, //TODO
			"pr79286.c":                    {}, //TODO
			"pr79354.c":                    {}, //TODO
			"pr79737-1.c":                  {}, //TODO
			"pr79737-2.c":                  {}, //TODO
			"pr80421.c":                    {}, //TODO
			"pr80692.c":                    {}, //TODO
			"pr81423.c":                    {}, //TODO
			"pr81555.c":                    {}, //TODO
			"pr81556.c":                    {}, //TODO
			"pr81588.c":                    {}, //TODO
			"pr82192.c":                    {}, //TODO
			"pr82210.c":                    {}, //TODO
			"pr82387.c":                    {}, //TODO
			"pr82524.c":                    {}, //TODO
			"pr82954.c":                    {}, //TODO
			"pr83362.c":                    {}, //TODO
			"pr83383.c":                    {}, //TODO
			"pr84169.c":                    {}, //TODO
			"pr84235.c":                    {}, //TODO
			"pr84478.c":                    {}, //TODO
			"pr84521.c":                    {}, //TODO
			"pr84524.c":                    {}, //TODO
			"pr84748.c":                    {}, //TODO
			"pr85095.c":                    {}, //TODO
			"pr85156.c":                    {}, //TODO
			"pr85529-1.c":                  {}, //TODO
			"pr85582-2.c":                  {}, //TODO
			"pr85582-3.c":                  {}, //TODO
			"pr86492.c":                    {}, //TODO
			"pr86528.c":                    {}, //TODO
			"pr86659-1.c":                  {}, //TODO
			"pr86659-2.c":                  {}, //TODO
			"pr87053.c":                    {}, //TODO
			"pr88714.c":                    {}, //TODO
			"pr88739.c":                    {}, //TODO
			"pr88904.c":                    {}, //TODO
			"pr89195.c":                    {}, //TODO
			"pr89369.c":                    {}, //TODO
			"pr89434.c":                    {}, //TODO
			"pr90311.c":                    {}, //TODO
			"pr90949.c":                    {}, //TODO
			"pr91137.c":                    {}, //TODO
			"pr91450-1.c":                  {}, //TODO
			"pr91450-2.c":                  {}, //TODO
			"pr91597.c":                    {}, //TODO
			"pr91635.c":                    {}, //TODO
			"pr92904.c":                    {}, //TODO
			"pr93213.c":                    {}, //TODO
			"pr93249.c":                    {}, //TODO
			"pr93434.c":                    {}, //TODO
			"pr93494.c":                    {}, //TODO
			"pr93744-1.c":                  {}, //TODO
			"pr93908.c":                    {}, //TODO
			"pr93945.c":                    {}, //TODO
			"pr94130.c":                    {}, //TODO
			"pr94724.c":                    {}, //TODO
			"pr94734.c":                    {}, //TODO
			"pr94809.c":                    {}, //TODO
			"pr96549.c":                    {}, //TODO
			"pr97325.c":                    {}, //TODO
			"pr97421-1.c":                  {}, //TODO
			"pr97764.c":                    {}, //TODO
			"pr98366.c":                    {}, //TODO
			"pr98474.c":                    {}, //TODO
			"printf-2.c":                   {}, //TODO
			"ptr-arith-1.c":                {}, //TODO
			"pushpop_macro.c":              {}, //TODO
			"restrict-1.c":                 {}, //TODO
			"return-addr.c":                {}, //TODO
			"ssad-run.c":                   {}, //TODO
			"stdarg-1.c":                   {}, //TODO
			"stdarg-3.c":                   {}, //TODO
			"stkalign.c":                   {}, //TODO
			"strcpy-1.c":                   {}, //TODO
			"strcpy-2.c":                   {}, //TODO
			"strct-stdarg-1.c":             {}, //TODO
			"strct-varg-1.c":               {}, //TODO
			"string-opt-18.c":              {}, //TODO
			"strlen-2.c":                   {}, //TODO
			"strlen-3.c":                   {}, //TODO
			"strlen-4.c":                   {}, //TODO
			"strlen-5.c":                   {}, //TODO
			"strlen-6.c":                   {}, //TODO
			"strlen-7.c":                   {}, //TODO
			"struct-ini-1.c":               {}, //TODO
			"struct-ini-2.c":               {}, //TODO
			"struct-ini-3.c":               {}, //TODO
			"struct-ini-4.c":               {}, //TODO
			"unsafe-fp-assoc-1.c":          {}, //TODO
			"usad-run.c":                   {}, //TODO
			"user-printf.c":                {}, //TODO
			"va-arg-10.c":                  {}, //TODO
			"va-arg-14.c":                  {}, //TODO
			"va-arg-15.c":                  {}, //TODO
			"va-arg-21.c":                  {}, //TODO
			"va-arg-22.c":                  {}, //TODO
			"va-arg-pack-1.c":              {}, //TODO
			"vla-dealloc-1.c":              {}, //TODO
			"vrp-7.c":                      {}, //TODO
			"widechar-3.c":                 {}, //TODO
			"zero-struct-1.c":              {}, //TODO
			"zero-struct-2.c":              {}, //TODO
			"zerolen-1.c":                  {}, //TODO
		}
		blacklistTCC := map[string]struct{}{
			// asm
			"99_fastcall.c": {},

			// Init
			"75_array_in_struct_init.c": {}, //TODO

			"76_dollars_in_identifiers.c": {}, //TODO
			"77_push_pop_macro.c":         {}, //TODO

			// Other
			"34_array_assignment.c":     {}, //TODO
			"46_grep.c":                 {}, //TODO
			"54_goto.c":                 {}, //TODO
			"60_errors_and_warnings.c":  {}, //TODO
			"73_arm64.c":                {}, //TODO
			"78_vla_label.c":            {}, //TODO
			"79_vla_continue.c":         {}, //TODO
			"80_flexarray.c":            {}, //TODO
			"81_types.c":                {}, //TODO
			"85_asm-outside-function.c": {}, //TODO
			"87_dead_code.c":            {}, //TODO
			"88_codeopt.c":              {}, //TODO
			"89_nocode_wanted.c":        {}, //TODO
			"90_struct-init.c":          {}, //TODO
			"92_enum_bitfield.c":        {}, //TODO
			"93_integer_promotion.c":    {}, //TODO
			"94_generic.c":              {}, //TODO
			"95_bitfields.c":            {}, //TODO
			"95_bitfields_ms.c":         {}, //TODO
			"96_nodata_wanted.c":        {}, //TODO
			"97_utf8_string_literal.c":  {}, //TODO
			"98_al_ax_extend.c":         {}, //TODO

		}
		switch fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) {
		case "darwin/amd64":
		case "darwin/arm64":
		case "freebsd/386":
			// asm
			blacklistGCC["960830-1.c"] = struct{}{}
			blacklistGCC["rbug.c"] = struct{}{}

			// Other
			blacklistGCC["copysign2.c"] = struct{}{} //TODO
		case "freebsd/amd64":
			// Other
			blacklistGCC["copysign2.c"] = struct{}{} //TODO
		case "linux/386":
			// asm
			blacklistGCC["960830-1.c"] = struct{}{}

			// _Float128
			blacklistGCC["nest-align-1.c"] = struct{}{}
			blacklistGCC["pr44468.c"] = struct{}{}
			blacklistGCC["strcmp-1.c"] = struct{}{}
			blacklistGCC["strlen-1.c"] = struct{}{}
			blacklistGCC["strncmp-1.c"] = struct{}{}

			// Needs -D_FILE_OFFSET_BITS=64.
			blacklistGCC["loop-2f.c"] = struct{}{} //TODO
			blacklistGCC["loop-2g.c"] = struct{}{} //TODO
		case "linux/arm":
		case "linux/s390x":
			// asm
			blacklistGCC["pr58574.c"] = struct{}{}
		case "netbsd/amd64":
			// Other
			blacklistGCC["copysign2.c"] = struct{}{} //TODO
			blacklistGCC["fp-cmp-7.c"] = struct{}{}
		case "openbsd/amd64":
			// Other
			blacklistGCC["copysign2.c"] = struct{}{} //TODO
		case "windows/386":
		case "windows/amd64":
		case "windows/arm64":
		}
		for _, v := range []struct {
			dir       string
			blacklist map[string]struct{}
		}{
			//TODO {"CompCert-3.6/test/c", blacklistCompCert},
			//TODO {"ccgo", nil},
			//TODO {"gcc-9.1.0/gcc/testsuite/gcc.c-torture", blacklistGCC},
			//TODO {"github.com/AbsInt/CompCert/test/c", blacklistCompCert},
			//TODO {"github.com/cxgo", nil},
			{"github.com/gcc-mirror/gcc/gcc/testsuite", blacklistGCC},
			//TODO {"github.com/vnmakarov", nil},
			//TODO {"sqlite-amalgamation-3380100", nil},
			{"tcc-0.9.27/tests/tests2", blacklistTCC},
			//TODO {"benchmarksgame-team.pages.debian.net", nil},
		} {
			t.Run(v.dir, func(t *testing.T) {
				testExec(t, "assets/"+v.dir, v.blacklist, g)
			})
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func testExec(t *testing.T, dir string, blacklist map[string]struct{}, g *golden) {
	p := newParallel()

	defer func() { p.close(t) }()

	p.err(cfsWalk(dir, func(pth string, fi os.FileInfo) error {
		if fi.IsDir() {
			return nil
		}

		if filepath.Ext(pth) != ".c" {
			return nil
		}

		p.file()
		switch {
		case re != nil:
			if !re.MatchString(pth) {
				p.skip()
				return nil
			}
		default:
			if _, ok := blacklist[filepath.Base(pth)]; ok {
				p.skip()
				return nil
			}
		}

		apth := pth
		p.exec(func() error {
			id := p.id()
			if *oTrace {
				fmt.Fprintln(os.Stderr, apth, id)
			}

			func() {
				defer func() {
					if err := recover(); err != nil {
						err = fmt.Errorf("%v: PANIC: %v", filepath.Base(apth), err)
						trc("%v: PANIC: %v\n%s", apth, err, debug.Stack())
						os.Exit(1)
					}
				}()

				b, err := getCorpusFile(apth)
				if err != nil {
					p.err(errorf("", err))
					p.fail()
					return
				}

				cfn := fmt.Sprintf("%d.c", id)
				if err := os.WriteFile(cfn, b, 0660); err != nil {
					p.err(errorf("", err))
					p.fail()
					return
				}

				var cfailed bool
				ofn := fmt.Sprintf("%d", id)
				if _, err := shell(false, hostCC, "-o", enforceBinaryExt(ofn), "-w", cfn); err != nil {
					cfailed = true
				}

				defer os.Remove(ofn)

				var cOut []byte
				if !cfailed {
					cOut, err = shell(false, "./"+enforceBinaryExt(ofn))
					if err != nil {
						cfailed = true
					}
				}

				ofn += ".go"

				defer os.Remove(ofn)

				var out bytes.Buffer
				if err := NewTask(goos, goarch, []string{"ccgo", "-o", ofn, "--prefix-field=F", apth}, &out, &out, cfs).Main(); err != nil {
					p.err(errorf("%s: %s: FAIL: %v", filepath.Base(apth), out.Bytes(), err))
					if cfailed {
						p.skip()
						return
					}

					p.fail()
					return
				}

				goOut, err := exec.Command("go", "run", ofn).CombinedOutput()
				if err != nil {
					p.err(errorf("%s: %s: FAIL: %v", filepath.Base(apth), goOut, err))
					if cfailed {
						p.skip()
						return
					}

					p.fail()
					return
				}

				if cfailed || bytes.Equal(cOut, goOut) {
					p.ok()
					g.w("%s\n", apth)
					return
				}

				cOut = bytes.TrimSpace(cOut)
				goOut = bytes.TrimSpace(goOut)
				if bytes.Equal(cOut, goOut) {
					p.ok()
					g.w("%s\n", apth)
					return
				}

				if bytes.Contains(cOut, []byte("\r\n")) {
					cOut = bytes.ReplaceAll(cOut, []byte("\r\n"), []byte{'\n'})
				}
				if bytes.Contains(goOut, []byte("\r\n")) {
					goOut = bytes.ReplaceAll(goOut, []byte("\r\n"), []byte{'\n'})
				}
				if bytes.Equal(cOut, goOut) {
					p.ok()
					g.w("%s\n", apth)
					return
				}

				diff := difflib.UnifiedDiff{
					A:        difflib.SplitLines(string(cOut)),
					B:        difflib.SplitLines(string(goOut)),
					FromFile: "expected",
					ToFile:   "got",
					Context:  0,
				}
				s, _ := difflib.GetUnifiedDiffString(diff)
				t.Errorf("%v:\n%v\n--- expexted\n%s\n\n--- got\n%s\n\n--- expected\n%s\n--- got\n%s", filepath.Base(apth), s, cOut, goOut, hex.Dump(cOut), hex.Dump(goOut))
				p.fail()
			}()
			return nil
		})
		return nil
	}))
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
