// Code generated by "stringer -type=Token"; DO NOT EDIT.

package tokens

import "strconv"

const _Token_name = "EndOfFileEndOfLineCommaDotAssignmentLocalReferenceAllLocalReferencesStringLocalStringPairWhileRepeatBreakSubtractionAdditionDivisionMultiplicationNotOrAndEqualityCheckLessThanCheckLessThanOrEqualCheckGreaterThanCheckGreaterThanOrEqualCheckExecuteRandomBlockStartOfExpressionEndOfExpressionStartOfStructEndOfStructStartOfArrayEndOfArrayStartOfSwitchEndOfSwitchSwitchCaseDefaultSwitchCaseStartOfFunctionEndOfFunctionReturnStartOfIfElseElseIfEndOfIfOptimisedIfOptimisedElseIntegerFloatNameShortJumpLongJumpChecksumTableEntryNamespaceAccessInvalid"

var _Token_index = [...]uint16{0, 9, 18, 23, 26, 36, 50, 68, 74, 85, 89, 94, 100, 105, 116, 124, 132, 146, 149, 151, 154, 167, 180, 200, 216, 239, 257, 274, 289, 302, 313, 325, 335, 348, 359, 369, 386, 401, 414, 420, 429, 433, 439, 446, 457, 470, 477, 482, 486, 495, 503, 521, 536, 543}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
