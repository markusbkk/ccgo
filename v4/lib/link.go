// Copyright 2022 The CCGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ccgo // import "modernc.org/ccgo/v4/lib"

import (
	"bytes"
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
	defs      map[string]gc.Node // extern: node
	externs   nameSet
	id        string // file name or import path
	pkg       *gc.Package
	pkgName   string // for kind == objectPkg
	qualifier string
	static    nameSet

	kind int // {objectFile, objectPkg}

	imported bool
}

func newObject(kind int, id string) *object {
	return &object{
		defs: map[string]gc.Node{},
		kind: kind,
		id:   id,
	}
}

func (o *object) load() (file *gc.SourceFile, err error) {
	if o.kind == objectPkg {
		return nil, errorf("object.load: internal error: wrong kind")
	}

	b, err := os.ReadFile(o.id)
	if err != nil {
		return nil, err
	}

	b = bytes.ReplaceAll(b, []byte("/**/\n"), []byte{'\n'})
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

				in[nm] = ts.TypeNode
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
					b.Write(spec.ExprList[i].Expr.Source(true))
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
						capi = v.ExprList[i].Expr
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

	if _, ok := lit.LiteralType.(*gc.MapTypeNode); !ok {
		return nil, errorf("%s: unexpected CAPI literal type: %T", importPath, lit.LiteralType)
	}

	r.pkgName = file.PackageClause.PackageName.Src()
	for _, v := range lit.LiteralValue.ElementList {
		switch x := v.Key.(type) {
		case *gc.BasicLit:
			if x.Token.Ch != gc.STRING_LIT {
				return nil, errorf("%s: invalid CAPI key type", importPath)
			}

			var key string
			if key, err = strconv.Unquote(x.Token.Src()); err != nil {
				return nil, errorf("%s: invalid CAPI key value: %s", importPath, x.Token.Src())
			}

			r.externs.add(tag(external) + key)
		default:
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
	importsByPath         map[string]*object
	libc                  *object
	out                   io.Writer
	reflectName           string
	stringLiterals        map[string]int64
	synthDecls            map[string][]byte
	task                  *Task
	textSegment           strings.Builder
	textSegmentName       string
	textSegmentNameP      string
	textSegmentOff        int64
	tld                   nameSpace
	unsafeName            string

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
		case meta:
			// nop
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
		importsByPath:  map[string]*object{},
		libc:           libc,
		stringLiterals: map[string]int64{},
		synthDecls:     map[string][]byte{},
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
	l.reflectName = l.tld.reg.put("reflect")
	l.unsafeName = l.tld.reg.put("unsafe")

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

				// trc("extern %q found in %q", nm, lib.id)
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
					l.importsByPath[lib.id] = lib
				}
			}
		}
	}
	if libc := l.libc; libc != nil && !libc.imported {
		libc.qualifier = l.tld.registerName(l, tag(importQualifier)+libc.pkgName)
		l.imports = append(l.imports, libc)
		libc.imported = true
		l.importsByPath[libc.id] = libc
	}

	out := bytes.NewBuffer(nil)
	l.out = out

	nm := l.task.packageName
	if nm == "" {
		nm = "main"
	}
	l.prologue(nm)
	l.w("\n\nimport (")
	switch nm := l.reflectName; nm {
	case "reflect":
		l.w("\n\t\"reflect\"")
	default:
		l.w("\n\t%s \"reflect\"", nm)
	}
	switch nm := l.unsafeName; nm {
	case "unsafe":
		l.w("\n\t\"unsafe\"")
	default:
		l.w("\n\t%s \"unsafe\"", nm)
	}
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
	_ %s.Type
	_ %s.Pointer
)

type float128 = struct { __ccgo [2]float64 }`, l.reflectName, l.unsafeName)

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
				if ln := x.VarSpecs[0].IdentifierList[0].Ident.Src(); l.meta(x, ln) {
					break
				}

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
				fi := l.newFnInfo(nil)
				l.print(fi, n)
				var b buf
				l.print0(&b, fi, n)
				l.synthDecls[nm] = b.bytes()
			case *gc.FunctionDecl:
				if ln := x.FunctionName.Src(); l.meta(x, ln) {
					break
				}

				l.funcDecl(x)
			default:
				l.err(errorf("TODO %T", x))
			}
		}
	}
	l.epilogue()
	if l.task.debugLinkerSave {
		if err := os.WriteFile(ofn, out.Bytes(), 0666); err != nil {
			return errorf("%s", err)
		}
	}

	b, err := l.postProcess(ofn, out.Bytes())
	if err != nil {
		l.err(err)
		return l.errors.err()
	}

	if err := os.WriteFile(ofn, b, 0666); err != nil {
		return errorf("%s", err)
	}

	if e := exec.Command("gofmt", "-s", "-w", "-r", "(x) -> x", ofn).Run(); e != nil {
		l.err(errorf("%s: gofmt: %v", ofn, e))
	}
	if *oTraceG {
		b, _ := os.ReadFile(ofn)
		fmt.Fprintf(os.Stderr, "%s\n", b)
	}
	return l.errors.err()
}

func (l *linker) meta(n gc.Node, linkName string) bool {
	if symKind(linkName) != meta {
		return false
	}

	rawName := l.rawName(linkName)
	if obj := l.externs[tag(external)+rawName]; obj != nil && obj.kind == objectPkg {
		if _, ok := obj.defs[rawName]; !ok {
			obj.defs[rawName] = n
		}
	}
	return true
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

	l.w("\n\nvar %s = (*%s.StringHeader)(%s.Pointer(&(%s))).Data\n", l.textSegmentNameP, l.reflectName, l.unsafeName, l.textSegmentName)
	l.w("\n\nvar %s = %q\n", l.textSegmentName, l.textSegment.String())
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
	case meta:
		return "X" + linkName[len(tag(meta)):]
	case importQualifier:
		switch nm := linkName[len(tag(importQualifier)):]; nm {
		case "libc":
			return fi.linker.libc.qualifier
		case "unsafe":
			return fi.linker.unsafeName
		default:
			fi.linker.err(errorf("TODO %q", nm))
			return linkName
		}
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
	l.print0(l, fi, n)
}

func (l *linker) print0(w writer, fi *fnInfo, n interface{}) {
	if n == nil {
		return
	}

	if x, ok := n.(gc.Token); ok && x.IsValid() {
		w.w("%s", x.Sep())
		switch x.Ch {
		case gc.IDENTIFIER:
			id := x.Src()
			nm := fi.name(id)
			if nm == "" {
				w.w("%s", id)
				return
			}

			if symKind(id) != external {
				w.w("%s", nm)
				return
			}

			obj := fi.linker.externs[id]
			if obj.kind == objectPkg {
				w.w("%s.%s", obj.qualifier, nm)
				return
			}

			w.w("%s", nm)
		case gc.STRING_LIT:
			w.w("%s", l.stringLit(x.Src(), false))
		default:
			w.w("%s", x.Src())
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

			l.print0(w, fi, v.Field(i).Interface())
		}
	case reflect.Slice:
		ne := v.Len()
		for i := 0; i < ne; i++ {
			l.print0(w, fi, v.Index(i).Interface())
		}
	}
}

func (l *linker) postProcess(fn string, b []byte) (r []byte, err error) {
	parserCfg := &gc.ParseSourceFileConfig{}
	sf, err := gc.ParseSourceFile(parserCfg, fn, b)
	if err != nil {
		return nil, errorf("%s", err)
	}

	pkg, err := gc.NewPackage("", []*gc.SourceFile{sf})
	if err != nil {
		return nil, errorf("%s", err)
	}

	if err := pkg.Check(l); err != nil {
		return nil, errorf("%s", err)
	}

	return sf.Source(true), nil
}

// PackageLoader implements gc.Checker.
func (l *linker) PackageLoader(pkg *gc.Package, src *gc.SourceFile, importPath string) (r *gc.Package, err error) {
	switch importPath {
	case "reflect":
		return l.reflectPackage()
	case "unsafe":
		return l.unsafePackage()
	}

	switch obj := l.importsByPath[importPath]; {
	case obj != nil:
		var b buf
		b.w("package %s", obj.pkgName)
		var taken nameSet
		var a []string
		if obj == l.libc {
			l.synthLibc(&b, &taken, pkg)
			b.w("\n")
		}
		for k := range l.synthDecls {
			a = append(a, k)
		}
		sort.Strings(a)
		for _, k := range a {
			if !taken.has(l.rawName(k)) {
				b.Write(l.synthDecls[k])
			}
		}
		a = a[:0]
		for k := range obj.defs {
			a = append(a, k)
		}
		sort.Strings(a)
		fi := l.newFnInfo(nil)
		for _, k := range a {
			l.print0(&b, fi, obj.defs[k])
		}
		// trc("\n%s", b.bytes())
		if r, err = l.syntheticPackage(importPath, importPath, b.bytes()); err != nil {
			return nil, err
		}

		return r, nil
	default:
		return nil, errorf("TODO %s", importPath)
	}
}

func (l *linker) synthLibc(b *buf, taken *nameSet, pkg *gc.Package) {
	b.w(`

type TLS struct{
	Alloc func(int) uintptr
	Free func(int)
}

func Start(func(*TLS, int32, uintptr) int32)

func VaList(p uintptr, args ...interface{}) uintptr

`)
	taken.add("TLS")
	taken.add("Start")
	for _, v := range []string{
		"float32",
		"float64",
		"int16",
		"int32",
		"int64",
		"int8",
		"uint16",
		"uint32",
		"uint64",
		"uint8",
		"uintptr",
	} {
		nm := fmt.Sprintf("Va%s", export(v))
		taken.add(nm)
		b.w("\n\nfunc %s(*uintptr) %s", nm, v)
	}
	for _, v := range []string{
		"int16",
		"int32",
		"int64",
		"int8",
	} {
		nm := fmt.Sprintf("Bool%s", export(v))
		taken.add(nm)
		b.w("\n\nfunc %s(bool) %s", nm, v)
	}
	for _, v := range []string{
		"float32",
		"float64",
		"int16",
		"int32",
		"int64",
		"int8",
		"uint16",
		"uint32",
		"uint64",
		"uint8",
		"uintptr",
	} {
		nm := export(v)
		taken.add(nm)
		b.w("\n\nfunc %s(%s) %[2]s", nm, v)
		for _, w := range []string{
			"float32",
			"float64",
			"int16",
			"int32",
			"int64",
			"int8",
			"uint16",
			"uint32",
			"uint64",
			"uint8",
			"uintptr",
		} {
			nm := fmt.Sprintf("%sFrom%s", export(v), export(w))
			taken.add(nm)
			b.w("\n\nfunc %s(%s) %s", nm, w, v)
		}
	}
}

// SymbolResolver implements gc.Checker.
func (l *linker) SymbolResolver(currentScope, fileScope *gc.Scope, pkg *gc.Package, ident gc.Token) (r gc.Node, err error) {
	// trc("%p %p %q %q %q", currentScope, fileScope, pkg.Name, pkg.ImportPath, ident)
	nm := ident.Src()
	off := ident.Offset()
	for s := currentScope; s != nil; s = s.Parent {
		if s.IsPackage() {
			if r := fileScope.Nodes[nm]; r.Node != nil {
				// trc("defined in file scope")
				return r.Node, nil
			}

			if pkg == l.libc.pkg && nm == "libc" { // rathole
				// trc("defined in libc")
				return pkg, nil
			}
		}

		if r := s.Nodes[nm]; r.Node != nil && r.VisibleFrom <= off {
			// trc("defined in scope %p(%v), parent %p(%v)", s, s.IsPackage(), s.Parent, s.Parent != nil && s.Parent.IsPackage())
			return r.Node, nil
		}
	}

	// trc("undefined: %s", nm)
	return nil, errorf("undefined: %s", nm)
}

// CheckFunctions implements gc.Checker.
func (l *linker) CheckFunctions() bool { return true }

// GOARCG implements gc.Checker.
func (l *linker) GOARCH() string { return l.task.goarch }

var (
	reflectSrc = []byte(`package reflect

type StringHeader struct {
    Data uintptr
    Len  int
}
`)
	unsafeSrc = []byte(`package unsafe

type ArbitraryType int

type Pointer *ArbitraryType

func Alignof(ArbitraryType) uintptr

func Offsetof(ArbitraryType) uintptr

func Sizeof(ArbitraryType) uintptr

func Add(Pointer, int) Pointer
`)
)

func (l *linker) reflectPackage() (r *gc.Package, err error) {
	return l.syntheticPackage("reflect", "<reflect>", reflectSrc)
}

func (l *linker) unsafePackage() (r *gc.Package, err error) {
	return l.syntheticPackage("unsafe", "<unsafe>", unsafeSrc)
}

func (l *linker) syntheticPackage(importPath, fn string, src []byte) (r *gc.Package, err error) {
	cfg := &gc.ParseSourceFileConfig{}
	sf, err := gc.ParseSourceFile(cfg, fn, src)
	if err != nil {
		return nil, err
	}

	r, err = gc.NewPackage(importPath, []*gc.SourceFile{sf})
	if err != nil {
		return nil, err
	}

	if obj := l.importsByPath[importPath]; obj != nil {
		obj.pkg = r
	}

	if err = r.Check(l); err != nil {
		return nil, err
	}

	return r, nil
}
