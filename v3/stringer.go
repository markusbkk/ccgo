// Code generated by "stringer -output stringer.go -type=exprMode"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[exprAddrOf-1]
	_ = x[exprBool-2]
	_ = x[exprCall-3]
	_ = x[exprLValue-4]
	_ = x[exprPSelect-5]
	_ = x[exprSelect-6]
	_ = x[exprValue-7]
	_ = x[exprVoid-8]
	_ = x[exprVoidSingle-9]
}

const _exprMode_name = "exprAddrOfexprBoolexprCallexprLValueexprPSelectexprSelectexprValueexprVoidexprVoidSingle"

var _exprMode_index = [...]uint8{0, 10, 18, 26, 36, 47, 57, 66, 74, 88}

func (i exprMode) String() string {
	i -= 1
	if i < 0 || i >= exprMode(len(_exprMode_index)-1) {
		return "exprMode(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _exprMode_name[_exprMode_index[i]:_exprMode_index[i+1]]
}
