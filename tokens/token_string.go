// Code generated by "stringer -type=Token"; DO NOT EDIT.

package tokens

import "strconv"

const _Token_name = "EndOfFileEndOfLineAssignmentLocalReferenceSubtractionAdditionDivisionMultiplicationNotEqualityCheckLessThanCheckLessThanOrEqualCheckGreaterThanCheckGreaterThanOrEqualCheckStartOfExpressionEndOfExpressionStartOfStructEndOfStructStartOfArrayEndOfArrayStartOfSwitchEndOfSwitchSwitchCaseDefaultSwitchCaseStartOfFunctionEndOfFunctionReturnBreakStartOfIfElseElseIfEndOfIfOptimisedIfOptimisedElseIntegerFloatNameShortJumpChecksumTableEntryNamespaceAccessInvalid"

var _Token_index = [...]uint16{0, 9, 18, 28, 42, 53, 61, 69, 83, 86, 99, 112, 132, 148, 171, 188, 203, 216, 227, 239, 249, 262, 273, 283, 300, 315, 328, 334, 339, 348, 352, 358, 365, 376, 389, 396, 401, 405, 414, 432, 447, 454}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
