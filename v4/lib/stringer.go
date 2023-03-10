// Code generated by "stringer -output stringer.go -type mode,name"; DO NOT EDIT.

package ccgo

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[exprBool-1]
	_ = x[exprCall-2]
	_ = x[exprDefault-3]
	_ = x[exprIndex-4]
	_ = x[exprLvalue-5]
	_ = x[exprSelect-6]
	_ = x[exprUintptr-7]
	_ = x[exprVoid-8]
}

const _mode_name = "exprBoolexprCallexprDefaultexprIndexexprLvalueexprSelectexprUintptrexprVoid"

var _mode_index = [...]uint8{0, 8, 16, 27, 36, 46, 56, 67, 75}

func (i mode) String() string {
	i -= 1
	if i < 0 || i >= mode(len(_mode_index)-1) {
		return "mode(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _mode_name[_mode_index[i]:_mode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[external-0]
	_ = x[typename-1]
	_ = x[taggedStruct-2]
	_ = x[taggedUnion-3]
	_ = x[taggedEum-4]
	_ = x[enumConst-5]
	_ = x[importQualifier-6]
	_ = x[macro-7]
	_ = x[define-8]
	_ = x[staticInternal-9]
	_ = x[staticNone-10]
	_ = x[automatic-11]
	_ = x[ccgoAutomatic-12]
	_ = x[ccgo-13]
	_ = x[field-14]
	_ = x[meta-15]
	_ = x[preserve-16]
}

const _name_name = "externaltypenametaggedStructtaggedUniontaggedEumenumConstimportQualifiermacrodefinestaticInternalstaticNoneautomaticccgoAutomaticccgofieldmetapreserve"

var _name_index = [...]uint8{0, 8, 16, 28, 39, 48, 57, 72, 77, 83, 97, 107, 116, 129, 133, 138, 142, 150}

func (i name) String() string {
	if i < 0 || i >= name(len(_name_index)-1) {
		return "name(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _name_name[_name_index[i]:_name_index[i+1]]
}
