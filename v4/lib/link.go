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
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/packages"
	"modernc.org/gc/v2"
)

const (
	objectFile = iota
	objectPkg
)

type object struct {
	externs   nameSet
	id        string // file name or import path
	pkgName   string // for kind == objectPkg
	qualifier string
	static    nameSet

	kind int // {objectFile, objectPkg}

	imported bool
}

func newObject(kind int, id string) *object {
	return &object{
		kind: kind,
		id:   id,
	}
}

// func (o *object) audit(fset *token.FileSet, file *ast.File) (err error) {
// 	var errors errors
// 	ast.Inspect(file, func(n ast.Node) bool {
// 		switch x := n.(type) {
// 		case *ast.UnaryExpr:
// 			if x.Op == token.AND {
// 				switch y := x.X.(type) {
// 				case *ast.Ident:
// 					switch symKind(y.Name) {
// 					case automatic, ccgoAutomatic:
// 						errors.add(errorf("%v: cannot take address of %s", fset.PositionFor(y.Pos(), true), y.Name))
// 					}
// 				case *ast.SelectorExpr:
// 					switch z := y.X.(type) {
// 					case *ast.Ident:
// 						switch symKind(z.Name) {
// 						case automatic, ccgoAutomatic:
// 							errors.add(errorf("%v: cannot take address of %s", fset.PositionFor(z.Pos(), true), z.Name))
// 						}
// 					}
// 				}
// 			}
// 		}
// 		return true
// 	})
// 	return errors.err()
// }

func (o *object) load() (file *gc.SourceFile, err error) {
	if o.kind == objectPkg {
		return nil, errorf("object.load: internal error: wrong kind")
	}

	b, err := os.ReadFile(o.id)
	if err != nil {
		return nil, err
	}

	if file, err = gc.ParseSourceFile(&gc.ParseSourceFileConfig{}, o.id, b); err != nil {
		return nil, err
	}

	return file, nil
}

// link name -> type ID
func (o *object) collectTypes(file *gc.SourceFile) (types map[string]string, err error) {
	var a []string
	in := map[string]gc.Node{}
	for _, decl := range file.TopLevelDecls {
		switch x := decl.(type) {
		case *gc.TypeDecl:
			for _, spec := range x.TypeSpecs {
				ts, ok := spec.(*gc.AliasDecl)
				if !ok {
					continue
				}

				nm := ts.Ident.Src()
				if _, ok := in[nm]; ok {
					return nil, errorf("%v: type %s redeclared", o.id, nm)
				}

				in[nm] = ts.Type
				a = append(a, nm)
			}
		}
	}
	sort.Strings(a)
	types = map[string]string{}
	for _, linkName := range a {
		if _, ok := types[linkName]; !ok {
			if types[linkName], err = typeID(in, types, in[linkName]); err != nil {
				return nil, err
			}
		}
	}
	return types, nil
}

// link name -> const value
func (o *object) collectConsts(file *gc.SourceFile) (consts map[string]string, err error) {
	var a []string
	in := map[string]string{}
	for _, decl := range file.TopLevelDecls {
		switch x := decl.(type) {
		case *gc.ConstDecl:
			for _, spec := range x.ConstSpecs {
				for i, ident := range spec.IdentifierList {
					nm := ident.Ident.Src()
					if _, ok := in[nm]; ok {
						return nil, errorf("%v: const %s redeclared", o.id, nm)
					}

					var b strings.Builder
					b.WriteByte('C') //TODO ?
					b.Write(spec.ExpressionList[i].Expression.Source(true))
					in[nm] = b.String()
					a = append(a, nm)
				}
			}
		}
	}
	sort.Strings(a)
	consts = map[string]string{}
	for _, linkName := range a {
		consts[linkName] = in[linkName]
	}
	return consts, nil
}

func (t *Task) link() (err error) {
	if len(t.inputFiles)+len(t.linkFiles) == 0 {
		return errorf("no input files")
	}

	if !t.keepObjectFiles {
		defer func() {
			for _, v := range t.compiledfFiles {
				os.Remove(v)
			}
		}()
	}

	if len(t.inputFiles) != 0 {
		if err := t.compile(""); err != nil {
			return err
		}
	}

	for i, v := range t.linkFiles {
		if x, ok := t.compiledfFiles[v]; ok {
			t.linkFiles[i] = x
		}
	}

	fset := token.NewFileSet()
	objects := map[string]*object{}
	mode := os.Getenv("GO111MODULE")
	var libc *object
	for _, v := range t.linkFiles {
		var object *object
		switch {
		case strings.HasPrefix(v, "-l="):
			object, err = t.getPkgSymbols(v[len("-l="):], mode)
			if err != nil {
				break
			}

			if object.pkgName == "libc" && libc == nil {
				libc = object
			}
		default:
			object, err = t.getFileSymbols(fset, v)
		}
		if err != nil {
			return err
		}

		if _, ok := objects[v]; !ok {
			objects[v] = object
		}
	}
	fset = nil

	switch {
	case t.o == "":
		return errorf("TODO %v %v %v %v", t.args, t.inputFiles, t.compiledfFiles, t.linkFiles)
	case strings.HasSuffix(t.o, ".go"):
		l, err := newLinker(t, libc)
		if err != nil {
			return err
		}

		r := l.link(t.o, t.linkFiles, objects)
		return r
	default:
		return errorf("TODO %v %v %v %v", t.args, t.inputFiles, t.compiledfFiles, t.linkFiles)
	}
}

func (t *Task) getPkgSymbols(importPath, mode string) (r *object, err error) {
	switch mode {
	case "", "on":
		// ok
	default:
		return nil, errorf("GO111MODULE=%s not supported", mode)
	}

	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedFiles,
			Env:  append(os.Environ(), fmt.Sprintf("GOOS=%s", t.goos), fmt.Sprintf("GOARCH=%s", t.goarch)),
		},
		importPath,
	)
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		return nil, errorf("%s: expected one package, loaded %d", importPath, len(pkgs))
	}

	pkg := pkgs[0]
	if len(pkg.Errors) != 0 {
		var a []string
		for _, v := range pkg.Errors {
			a = append(a, v.Error())
		}
		return nil, errorf("%s", strings.Join(a, "\n"))
	}

	r = newObject(objectPkg, importPath)
	base := fmt.Sprintf("capi_%s_%s.go", t.goos, t.goarch)
	var fn string
	for _, v := range pkg.GoFiles {
		if filepath.Base(v) == base {
			fn = v
			break
		}
	}
	if fn == "" {
		return nil, errorf("%s: file %s not found", importPath, base)
	}

	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, errorf("%s: %v", importPath, err)
	}

	file, err := gc.ParseSourceFile(&gc.ParseSourceFileConfig{}, fn, b)
	if err != nil {
		return nil, errorf("%s: %v", importPath, err)
	}

	var capi gc.Node
out:
	for _, v := range file.TopLevelDecls {
		if x, ok := v.(*gc.VarDecl); ok {
			for _, v := range x.VarSpecs {
				for i, id := range v.IdentifierList {
					if nm := id.Ident.Src(); nm == "CAPI" {
						capi = v.ExpressionList[i].Expression
						break out
					}
				}
			}
		}
	}

	if capi == nil {
		return nil, errorf("%s: CAPI not declared in %s", importPath, fn)
	}

	lit, ok := capi.(*gc.CompositeLit)
	if !ok {
		return nil, errorf("%s: unexpected CAPI node type: %T", importPath, capi)
	}

	if _, ok := lit.LiteralType.(*gc.MapType); !ok {
		return nil, errorf("%s: unexpected CAPI literal type: %T", importPath, lit.LiteralType)
	}

	r.pkgName = file.PackageClause.PackageName.Src()
	for _, v := range lit.LiteralValue.ElementList {
		switch x := v.Key.(type) {
		case gc.Token:
			if x.Ch != gc.STRING_LIT {
				return nil, errorf("%s: invalid CAPI key type: %s", importPath, x)
			}

			var key string
			if key, err = strconv.Unquote(x.Src()); err != nil {
				return nil, errorf("%s: invalid CAPI key value: %s", importPath, x.Src)
			}

			r.externs.add(tag(external) + key)
		default:
			trc("", x)
			panic(todo("%T", x))
		}
	}
	return r, nil
}

func (t *Task) getFileSymbols(fset *token.FileSet, fn string) (r *object, err error) {
	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var pkgName string
	file, err := gc.ParseSourceFile(&gc.ParseSourceFileConfig{
		Accept: func(file *gc.SourceFile) error {
			pkgName = file.PackageClause.PackageName.Src()
			if !strings.HasPrefix(pkgName, objectFilePackageNamePrefix) {
				return errorf("%s: package %s is not a ccgo object file", fn, pkgName)
			}

			version := pkgName[len(objectFilePackageNamePrefix):]
			if !semver.IsValid(version) {
				return errorf("%s: package %s has invalid semantic version", fn, pkgName)
			}

			if semver.Compare(version, objectFileSemver) != 0 {
				return errorf("%s: package %s has incompatible semantic version compared to %s", fn, pkgName, objectFileSemver)
			}

			return nil
		},
	}, fn, b)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(file.PackageClause.Package.Sep(), "\n") {
		if !strings.HasPrefix(line, "//") {
			continue
		}

		x := strings.Index(line, generatedFilePrefix)
		if x < 0 {
			continue
		}

		s := line[x+len(generatedFilePrefix):]
		if len(s) == 0 {
			continue
		}

		if !strings.HasPrefix(s, fmt.Sprintf("%s/%s", t.goos, t.goarch)) {
			return nil, errorf("%s: object file was compiled for different target: %s", fn, line)
		}
	}

	r = newObject(objectFile, fn)
	ex := tag(external)
	si := tag(staticInternal)
	for k, v := range file.TopLevelDecls {
		var a []gc.Token
		switch x := v.(type) {
		case *gc.ConstDecl, *gc.TypeDecl:
			continue
		case *gc.VarDecl:
			for _, v := range x.VarSpecs {
				for _, id := range v.IdentifierList {
					a = append(a, id.Ident)
				}
			}
		case *gc.FunctionDecl:
			a = append(a, x.FunctionName)
		default:
			_ = ex
			_ = si
			_ = k
			panic(todo("%T", x))
		}
		for _, id := range a {
			k := id.Src()
			switch symKind(k) {
			case external:
				if _, ok := r.externs[k]; ok {
					return nil, errorf("invalid object file: multiple defintions of %s", k[len(ex):])
				}

				r.externs.add(k)
			case staticInternal:
				if _, ok := r.static[k]; ok {
					return nil, errorf("invalid object file: multiple defintions of %s", k[len(si):])
				}

				r.static.add(k)
			}
		}
	}
	return r, nil
}

type linker struct {
	errors                errors
	externs               map[string]*object
	fileLinkNames2GoNames dict
	fileLinkNames2IDs     dict
	forceExternalPrefix   nameSet
	fset                  *token.FileSet
	goTags                []string
	goTypeNamesEmited     nameSet
	imports               []*object
	libc                  *object
	out                   io.Writer
	stringLiterals        map[string]int64
	task                  *Task
	textSegment           strings.Builder
	textSegmentName       string
	textSegmentNameP      string
	textSegmentOff        int64
	tld                   nameSpace

	closed bool
}

func newLinker(task *Task, libc *object) (*linker, error) {
	goTags := tags
	for i := range tags {
		switch name(i) {
		case ccgoAutomatic, ccgo:
			goTags[i] = task.prefixCcgoAutomatic
		case define:
			goTags[i] = task.prefixDefine
		case enumConst:
			goTags[i] = task.prefixEnumerator
		case external:
			goTags[i] = task.prefixExternal
		case field:
			goTags[i] = task.prefixField
		case importQualifier:
			goTags[i] = task.prefixImportQualifier
		case macro:
			goTags[i] = task.prefixMacro
		case automatic:
			goTags[i] = task.prefixAutomatic
		case staticInternal:
			goTags[i] = task.prefixStaticInternal
		case staticNone:
			goTags[i] = task.prefixStaticNone
		case preserve:
			goTags[i] = ""
		case taggedEum:
			goTags[i] = task.prefixTaggedEnum
		case taggedStruct:
			goTags[i] = task.prefixTaggedStruct
		case taggedUnion:
			goTags[i] = task.prefixTaggedUnion
		case typename:
			goTags[i] = task.prefixTypename
		//TODO case unpinned:
		//TODO 	goTags[i] = task.prefixUnpinned
		//TODO case externalUnpinned:
		//TODO 	goTags[i] = task.prefixExternalUnpinned
		default:
			return nil, errorf("internal error: %v", name(i))
		}
	}
	return &linker{
		externs:        map[string]*object{},
		fset:           token.NewFileSet(),
		goTags:         goTags[:],
		libc:           libc,
		stringLiterals: map[string]int64{},
		task:           task,
	}, nil
}

func (l *linker) err(err error)                      { l.errors.add(err) }
func (l *linker) rawName(linkName string) (r string) { return linkName[len(tag(symKind(linkName))):] }

func (l *linker) goName(linkName string) (r string) {
	return l.goTags[symKind(linkName)] + l.rawName(linkName)
}

func (l *linker) w(s string, args ...interface{}) {
	if l.closed {
		return
	}

	if _, err := fmt.Fprintf(l.out, s, args...); err != nil {
		l.err(err)
		l.closed = true
	}
}

func (l *linker) link(ofn string, linkFiles []string, objects map[string]*object) (err error) {
	var tld nameSet
	// Build the symbol table.
	for _, linkFile := range linkFiles {
		object := objects[linkFile]
		for nm := range object.externs {
			if _, ok := l.externs[nm]; !ok {
				l.externs[nm] = object
			}
			tld.add(nm)
		}
	}
	l.tld.registerNameSet(l, tld, true)
	l.textSegmentNameP = l.tld.reg.put("ts")
	l.textSegmentName = l.tld.reg.put("ts")

	// Check for unresolved references.
	for _, linkFile := range linkFiles {
		switch object := objects[linkFile]; {
		case object.kind == objectFile:
			file, err := object.load()
			if err != nil {
				return errorf("loading %s: %v", object.id, err)
			}

			for nm, pos := range unresolvedSymbols(file) {
				if !strings.HasPrefix(nm, tag(external)) {
					continue
				}

				lib, ok := l.externs[nm]
				if !ok {
					return errorf("%v: undefined reference to '%s'", pos, l.rawName(nm))
				}

				if lib.kind == objectFile {
					continue
				}

				if l.task.prefixExternal != "X" {
					l.forceExternalPrefix.add(nm)
				}
				if lib.qualifier == "" {
					lib.qualifier = l.tld.registerName(l, tag(importQualifier)+lib.pkgName)
					l.imports = append(l.imports, lib)
					lib.imported = true
				}
			}
		}
	}
	if libc := l.libc; libc != nil && !libc.imported {
		libc.qualifier = l.tld.registerName(l, tag(importQualifier)+libc.pkgName)
		l.imports = append(l.imports, libc)
		libc.imported = true
	}

	f, err := os.Create(ofn)
	if err != nil {
		return errorf("%s", err)
	}

	defer func() {
		if e := f.Close(); e != nil {
			l.err(errorf("%s", e))
		}

		if e := exec.Command("gofmt", "-s", "-w", "-r", "(x) -> x", ofn).Run(); e != nil {
			l.err(errorf("%s: gofmt: %v", ofn, e))
		}
		if *oTraceG {
			b, _ := os.ReadFile(ofn)
			fmt.Fprintf(os.Stderr, "%s\n", b)
		}
		if err != nil {
			l.err(err)
		}
		err = l.errors.err()
	}()

	out := bufio.NewWriter(f)
	l.out = out

	defer func() {
		if err := out.Flush(); err != nil {
			l.err(errorf("%s", err))
		}
	}()

	nm := l.task.packageName
	if nm == "" {
		nm = "main"
	}
	l.prologue(nm)
	l.w("\n\nimport (")
	l.w("\n\t\"reflect\"")
	l.w("\n\t\"unsafe\"")
	if len(l.imports) != 0 {
		l.w("\n")
	}
	for _, v := range l.imports {
		l.w("\n\t")
		if v.pkgName != v.qualifier {
			l.w("%s ", v.qualifier)
		}
		l.w("%q", v.id)
	}
	l.w("\n)")
	l.w(`

var (
	_ reflect.Type
	_ unsafe.Pointer
)

type float128 = struct { __ccgo [2]float64 }`)

	for _, linkFile := range linkFiles {
		object := objects[linkFile]
		if object.kind != objectFile {
			continue
		}

		file, err := object.load()
		if err != nil {
			return errorf("loading %s: %v", object.id, err)
		}

		// types
		fileLinkNames2IDs, err := object.collectTypes(file)
		if err != nil {
			return errorf("loading %s: %v", object.id, err)
		}

		var linkNames []string
		for linkName := range fileLinkNames2IDs {
			linkNames = append(linkNames, linkName)
		}
		sort.Strings(linkNames)
		l.fileLinkNames2GoNames = dict{}
		for _, linkName := range linkNames {
			typeID := fileLinkNames2IDs[linkName]
			associatedTypeID, ok := l.fileLinkNames2IDs[linkName]
			switch {
			case !ok:
				l.fileLinkNames2IDs.put(linkName, typeID)
				goName := l.tld.registerName(l, linkName)
				l.fileLinkNames2GoNames[linkName] = goName
			case ok && associatedTypeID == typeID:
				l.fileLinkNames2GoNames[linkName] = l.tld.dict[linkName]
			default:
				l.err(errorf("TODO obj %s, linkName %s, typeID %s, ok %v, associatedTypeID %s", object.id, linkName, typeID, ok, associatedTypeID))
			}
		}

		// consts
		if fileLinkNames2IDs, err = object.collectConsts(file); err != nil {
			return errorf("loading %s: %v", object.id, err)
		}

		linkNames = linkNames[:0]
		for linkName := range fileLinkNames2IDs {
			linkNames = append(linkNames, linkName)
		}
		sort.Strings(linkNames)
		for _, linkName := range linkNames {
			constID := fileLinkNames2IDs[linkName]
			associatedConstID, ok := l.fileLinkNames2IDs[linkName]
			switch {
			case !ok:
				l.fileLinkNames2IDs.put(linkName, constID)
				goName := l.tld.registerName(l, linkName)
				l.fileLinkNames2GoNames[linkName] = goName
			case ok && associatedConstID == constID:
				l.fileLinkNames2GoNames[linkName] = l.tld.dict[linkName]
			default:
				l.err(errorf("TODO obj %s, linkName %s, constID %s, ok %v, associatedConstID %s", object.id, linkName, constID, ok, associatedConstID))
			}
		}

		// statics
		linkNames = linkNames[:0]
		for linkName := range object.static {
			linkNames = append(linkNames, linkName)
		}
		sort.Strings(linkNames)
		for _, linkName := range linkNames {
			goName := l.tld.registerName(l, linkName)
			l.fileLinkNames2GoNames[linkName] = goName
		}

		for _, n := range file.TopLevelDecls {
			switch x := n.(type) {
			case *gc.ConstDecl:
				l.print(l.newFnInfo(nil), n)
			case *gc.VarDecl:
				l.print(l.newFnInfo(n), n)
			case *gc.TypeDecl:
				if len(x.TypeSpecs) != 1 {
					panic(todo(""))
				}

				spec := x.TypeSpecs[0]
				nm := spec.(*gc.AliasDecl).Ident.Src()
				if _, ok := l.goTypeNamesEmited[nm]; ok {
					break
				}

				l.goTypeNamesEmited.add(nm)
				l.print(l.newFnInfo(nil), n)
			case *gc.FunctionDecl:
				l.funcDecl(x)
			default:
				l.err(errorf("TODO %T", x))
			}
		}
	}
	l.epilogue()
	return l.errors.err()
}

func (l *linker) funcDecl(n *gc.FunctionDecl) {
	l.w("\n\n")
	info := l.newFnInfo(n)
	var static []gc.Node
	w := 0
	for _, stmt := range n.FunctionBody.StatementList {
		if stmt := l.stmtPrune(stmt, info, &static); stmt != nil {
			n.FunctionBody.StatementList[w] = stmt
			w++
		}
	}
	n.FunctionBody.StatementList = n.FunctionBody.StatementList[:w]
	l.print(info, n)
	for _, v := range static {
		l.w("\n\n")
		l.print(info, v)
	}
}

func (l *linker) stmtPrune(n gc.Node, info *fnInfo, static *[]gc.Node) gc.Node {
	switch x := n.(type) {
	case *gc.VarDecl:
		if len(x.VarSpecs) != 1 {
			return n
		}

		vs := x.VarSpecs[0]
		if len(vs.IdentifierList) != 1 {
			return n
		}

		switch nm := vs.IdentifierList[0].Ident.Src(); symKind(nm) {
		case staticInternal, staticNone:
			*static = append(*static, n)
			return nil
		}
	}
	return n
}

func (l *linker) epilogue() {
	if l.textSegment.Len() == 0 {
		return
	}

	l.w("\n\nvar %s = %q\n", l.textSegmentName, l.textSegment.String())
	l.w("\nvar %s = (*reflect.StringHeader)(unsafe.Pointer(&(%s))).Data\n", l.textSegmentNameP, l.textSegmentName)
}

func (l *linker) prologue(nm string) {
	l.w(`// %s%s/%s by '%s %s'%s.

//go:build %[2]s && %[3]s
// +build %[2]s,%[3]s

package %[7]s

`,
		generatedFilePrefix,
		l.task.goos, l.task.goarch,
		filepath.Base(l.task.args[0]),
		strings.Join(l.task.args[1:], " "),
		generatedFileSuffix,
		nm,
	)
}

type fnInfo struct {
	ns        nameSpace
	linkNames nameSet
	linker    *linker
}

func (l *linker) newFnInfo(n gc.Node) (r *fnInfo) {
	r = &fnInfo{linker: l}
	if n != nil {
		walk(n, func(n gc.Node) {
			tok, ok := n.(gc.Token)
			if !ok {
				return
			}

			switch tok.Ch {
			case gc.IDENTIFIER:
				switch nm := tok.Src(); symKind(nm) {
				case staticInternal, field:
					// nop
				default:
					r.linkNames.add(nm)
				}
			case gc.STRING_LIT:
				r.linker.stringLit(tok.Src(), true)
			}
		})
	}
	var linkNames []string
	for k := range r.linkNames {
		linkNames = append(linkNames, k)
	}
	sort.Slice(linkNames, func(i, j int) bool {
		return symKind(linkNames[i]) < symKind(linkNames[j]) || linkNames[i] < linkNames[j]
	})
	r.ns.registerNameSet(l, r.linkNames, false)
	r.linkNames = nil
	return r
}

func (fi *fnInfo) name(linkName string) string {
	switch symKind(linkName) {
	case external:
		if fi.linker.forceExternalPrefix.has(linkName) {
			return linkName
		}

		fallthrough
	case staticInternal, staticNone:
		if goName := fi.linker.tld.dict[linkName]; goName != "" {
			return goName
		}
	case preserve, field:
		return fi.linker.goName(linkName)
	case automatic, ccgoAutomatic, ccgo:
		return fi.ns.dict[linkName]
	case
		typename, taggedEum, taggedStruct, taggedUnion, define, macro, enumConst:

		return fi.linker.fileLinkNames2GoNames[linkName]
	case -1:
		return linkName
	}

	fi.linker.err(errorf("TODO %q %v", linkName, symKind(linkName)))
	return linkName
}

func (l *linker) stringLit(s0 string, reg bool) string {
	s, err := strconv.Unquote(s0)
	if err != nil {
		l.err(errorf("internal error: %v", err))
	}
	off := l.textSegmentOff
	switch x, ok := l.stringLiterals[s]; {
	case ok:
		off = x
	default:
		if !reg {
			return s0
		}

		l.stringLiterals[s] = off
		l.textSegment.WriteString(s)
		l.textSegmentOff += int64(len(s))
	}
	switch {
	case off == 0:
		return l.textSegmentNameP
	default:
		return fmt.Sprintf("(%s%+d)", l.textSegmentNameP, off)
	}
}

func (l *linker) print(fi *fnInfo, n interface{}) {
	if n == nil {
		return
	}

	if x, ok := n.(gc.Token); ok && x.IsValid() {
		l.w("%s", x.Sep())
		switch x.Ch {
		case gc.IDENTIFIER:
			id := x.Src()
			nm := fi.name(id)
			if nm == "" {
				l.w("%s", id)
				return
			}

			if symKind(id) != external {
				l.w("%s", nm)
				return
			}

			obj := fi.linker.externs[id]
			if obj.kind == objectPkg {
				l.w("%s.%s", obj.qualifier, nm)
				return
			}

			l.w("%s", nm)
		case gc.STRING_LIT:
			l.w("%s", l.stringLit(x.Src(), false))
		default:
			l.w("%s", x.Src())
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
		nf := t.NumField()
		for i := 0; i < nf; i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}

			if v == zero || v.IsZero() {
				continue
			}

			l.print(fi, v.Field(i).Interface())
		}
	case reflect.Slice:
		ne := v.Len()
		for i := 0; i < ne; i++ {
			l.print(fi, v.Index(i).Interface())
		}
	default:
		panic(todo("", t.Kind()))
	}
}
