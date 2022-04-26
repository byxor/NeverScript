package decompiler

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	Byte_EndOfFile = 0x0
	Byte_NewLine = 0x1
	Byte_NewLineWithNumber = 0x2
	Byte_Struct = 0x3
	Byte_EndStruct = 0x4
	Byte_Array = 0x5
	Byte_EndArray = 0x6
	Byte_Equals = 0x7
	Byte_Dot = 0x8
	Byte_Comma = 0x9
	Byte_Minus = 0xA
	Byte_Plus = 0xB
	Byte_Divide = 0xC
	Byte_Multiply = 0xD
	Byte_Parenthesis = 0xE
	Byte_EndParenthesis = 0xF
	Byte_EqualTo = 0x11
	Byte_LessThan = 0x12
	Byte_LessThanEqual = 0x13
	Byte_GreaterThan = 0x14
	Byte_GreaterThanEqual = 0x15
	Byte_Checksum = 0x16
	Byte_Integer = 0x17
	Byte_Float = 0x1A
	Byte_String = 0x1B
	Byte_LocalString = 0x1C
	Byte_Vector = 0x1E
	Byte_Pair = 0x1F
	Byte_While = 0x20
	Byte_EndWhile = 0x21
	Byte_Break = 0x22
	Byte_Script = 0x23
	Byte_EndScript = 0x24
	Byte_If = 0x25
	Byte_Else = 0x26
	Byte_EndIf = 0x28
	Byte_Return = 0x29
	Byte_ChecksumEntry = 0x2B
	Byte_AllArguments = 0x2C
	Byte_Local = 0x2D
	Byte_LongJump = 0x2E
	Byte_Random = 0x2F
	Byte_RandomRange = 0x30
	Byte_Or = 0x32
	Byte_And = 0x33
	Byte_Xor = 0x34
	Byte_Not = 0x39
	Byte_RandomNoRepeat = 0x40
	Byte_Colon = 0x42
)

func Decompile(qb []byte) (string, error) {

	var DecompileExpression func(int, int) (string, int, error)
	var DecompileBodyOfCode func(int, int) (string, int, error)
	var checksumTable map[uint32]string

	GetByte := func(index int) (byte, error) {
		if index >= len(qb) {
			return 0, errors.New(fmt.Sprintf("Index 0x%x out of range", index))
		}
		return qb[index], nil
	}

	GetBytes := func(index, size int) ([]byte, error) {
		if index + size > len(qb) {
			return []byte{}, errors.New(fmt.Sprintf("Index 0x%x out of range", index + size - 1))
		}
		return qb[index:index+size], nil
	}

	GetChecksumTable := func() (map[uint32]string, error) {
		index := len(qb) - 1

		checksumTable := make(map[uint32]string)

		for {
			if index < 0 {
				break
			}

			// find beginning of entry
			b, err := GetByte(index)
			if err != nil {
				return checksumTable, err
			}

			if b == Byte_ChecksumEntry {
				startOfChecksum := index
				// found potential checksum entry
				index++

				checksumBytes, err := GetBytes(index, 4)
				if err != nil {
					return checksumTable, err
				}

				checksum := binary.LittleEndian.Uint32(checksumBytes)
				index += 4
				checksumNameStartIndex := index

				// scan name of checksum
				for {
					nextByte, err := GetByte(index)
					if err != nil {
						return checksumTable, err
					}
					if nextByte == 0 {
						index++
						break
					}
					index++
				}
				checksumName := string(qb[checksumNameStartIndex:index-1])

				// sanity check, may not be a printable checksum
				isPrintable := false
				for i, c := range checksumName {
					if !unicode.IsNumber(c) && !unicode.IsLetter(c) && c != ' ' && c != '_' {
						break
					}
					if i >= len(checksumName) - 1 {
						isPrintable = true
					}
				}
				if isPrintable {
					checksumTable[checksum] = checksumName
				}
				index = startOfChecksum - 1
			}

			index--
		}
		return checksumTable, nil
	}

	Indent := func(indentationLevel int, text string) string {
		return strings.Repeat("    ", indentationLevel) + text
	}

	TrimWhitespace := func(text string) string {
		isEntirelyWhitespace, _ := regexp.MatchString(`^ *$`, text)
		if isEntirelyWhitespace {
			return ""
		}
		return strings.Trim(text, " ")
	}

	FormatFloat := func(floatBytes []byte) string {
		floatAsInteger := binary.LittleEndian.Uint32(floatBytes)
		floatValue := math.Float32frombits(floatAsInteger)
		floatString := strconv.FormatFloat(float64(floatValue), 'f', -1, 32)
		if !strings.Contains(floatString, ".") {
			floatString += ".0"
		}
		return floatString
	}

	DecompilerError := func(message string, b byte, offset int) error {
		return errors.New(fmt.Sprintf("%s - 0x%x byte (offset 0x%x)", message, b, offset))
	}

	DecompileNewLineWithNumber := func(index int) (string, int, error) {
		//lineNumber := int(binary.LittleEndian.Uint32(qb[index : index+4]))
		//result := fmt.Sprintf("\n/* Line number 0x%x */", lineNumber)
		return "\n", 5, nil
	}

	DecompileString := func(index int) (string, int, error) {
		initialIndex := index
		index++
		lengthBytes := qb[index : index+4]
		length := int(binary.LittleEndian.Uint32(lengthBytes))
		index += 4
		stringBytes := qb[index : index+length]
		index += length
		// warning: if string contains new-line bytes, it might produce code that doesn't compile
		return fmt.Sprintf(`"%s"`, string(stringBytes[:length - 1])), index - initialIndex, nil
	}

	DecompileChecksum := func(index int) (string, int, error) {
		index++
		checksumBytes := qb[index : index+4]

		var code string
		checksum := binary.LittleEndian.Uint32(checksumBytes)
		if checksumName, found := checksumTable[checksum]; found {
			code = checksumName
		} else {
			code = fmt.Sprintf("#%02x%02x%02x%02x", checksumBytes[0], checksumBytes[1], checksumBytes[2], checksumBytes[3])
		}

		return code, 5, nil
	}

	DecompilePair := func(index, indentationLevel int) (string, int, error) {
		initialIndex := index
		index++

		xBytes, err := GetBytes(index, 4)
		if err != nil {
			return "", 0, err
		}
		index += 4

		yBytes, err := GetBytes(index, 4)
		if err != nil {
			return "", 0, err
		}
		index += 4

		return fmt.Sprintf("(%s, %s) ", FormatFloat(xBytes), FormatFloat(yBytes)), index - initialIndex, nil
	}

	DecompileBodyOfCode = func(index, indentationLevel int) (string, int, error) {
		var currentLineCode strings.Builder
		var bodyOfCode strings.Builder
		flushCurrentLine := func() {
			bodyOfCode.WriteString(Indent(indentationLevel, currentLineCode.String()))
			currentLineCode.Reset()
		}

		initialIndex := index

		for {
			b, err := GetByte(index)
			if err != nil {
				return "", 0, err
			}

			if b == Byte_NewLine {
				currentLineCode.WriteString("\n")
				flushCurrentLine()
				index++
			} else if b == Byte_NewLineWithNumber {
				newLineCode, bytesRead, err := DecompileNewLineWithNumber(index)
				if err != nil {
					return "", 0, err
				}
				index += bytesRead
				currentLineCode.WriteString(newLineCode)
				flushCurrentLine()
			} else if b == Byte_Return {
				index++
				currentLineCode.WriteString("return ")
			} else if b == Byte_Comma {
				index++
				currentLineCode.WriteString(", ")
			} else if b == Byte_Local || b == Byte_Checksum {
				isLocal := b == Byte_Local
				if isLocal {
					index++
				}

				expressionCode, bytesRead, err := DecompileExpression(index, indentationLevel)
				if err != nil {
					return "", 0, err
				}
				index += bytesRead

				if isLocal {
					// wrap checksum from expression in <>
					// warning: expand this to handle `checksums like this`
					endOfChecksumPart := 0
					for {
						if endOfChecksumPart >= len(expressionCode) {
							break
						}
						c := expressionCode[endOfChecksumPart]
						r := rune(c)
						if c == '#' || unicode.IsLetter(r) || unicode.IsNumber(r) {
							endOfChecksumPart++
						} else {
							break
						}
					}
					expressionCode = "<" + expressionCode[:endOfChecksumPart] + ">" + expressionCode[endOfChecksumPart:]
				}

				nextByte, err := GetByte(index)
				if err != nil {
					currentLineCode.WriteString(expressionCode)
				} else if nextByte != Byte_Equals {
					currentLineCode.WriteString(expressionCode)
				} else {
					index++

					var extraNewLineCode strings.Builder
					for {
						nextByte, err := GetByte(index)
						if err != nil {
							return "", 0, err
						}
						if nextByte == Byte_NewLineWithNumber {
							newLineCode, bytesRead, err := DecompileNewLineWithNumber(index)
							if err != nil {
								return "", 0, err
							}
							index += bytesRead
							extraNewLineCode.WriteString(newLineCode)
						} else if nextByte == Byte_NewLine {
							index++
							extraNewLineCode.WriteString("\n")
						} else {
							break
						}
					}

					secondExpressionCode, bytesRead, err := DecompileExpression(index, indentationLevel)
					if err != nil {
						return "", 0, err
					}
					index += bytesRead

					currentLineCode.WriteString(fmt.Sprintf("%s = %s", expressionCode, secondExpressionCode))
				}
			} else if b == Byte_Break {
				index++
				currentLineCode.WriteString("break")
			} else if b == Byte_If {
				index++

				conditionCode, bytesRead, err := DecompileExpression(index, indentationLevel)
				if err != nil {
					return "", 0, err
				}
				index += bytesRead

				ifBodyCode, bytesRead, err := DecompileBodyOfCode(index, indentationLevel+1)
				if err != nil {
					return "", 0, err
				}
				index += bytesRead

				nextByte, err := GetByte(index)
				if err != nil {
					return "", 0, err
				}

				if nextByte == Byte_Else {
					index++

					elseBodyCode, bytesRead, err := DecompileBodyOfCode(index, indentationLevel+1)
					if err != nil {
						return "", 0, err
					}

					index += bytesRead

					currentLineCode.WriteString(fmt.Sprintf("if %s {%s", conditionCode, ifBodyCode))
					if strings.Contains(currentLineCode.String(), "\n") {
						flushCurrentLine()
					}

					currentLineCode.WriteString(fmt.Sprintf("} else {%s", elseBodyCode))
					if strings.Contains(currentLineCode.String(), "\n") {
						flushCurrentLine()
					}

					currentLineCode.WriteString("}")
				} else {
					currentLineCode.WriteString(fmt.Sprintf("if %s {%s", conditionCode, ifBodyCode))
					if strings.Contains(currentLineCode.String(), "\n") {
						flushCurrentLine()
					}
					currentLineCode.WriteString("}")
				}

				nextByte, err = GetByte(index)
				if err != nil {
					return "", 0, err
				}

				if nextByte != Byte_EndIf {
					return "", 0, DecompilerError("No endif byte", nextByte, index)
				}
				index++
			} else if b == Byte_While {
				index++

				whileBodyCode, bytesRead, err := DecompileBodyOfCode(index, indentationLevel + 1)
				if err != nil {
					return "", 0, err
				}
				index += bytesRead

				currentLineCode.WriteString(fmt.Sprintf("while {%s", whileBodyCode))
				if strings.Contains(currentLineCode.String(), "\n") {
					flushCurrentLine()
				}
				currentLineCode.WriteString("}")

				nextByte, err := GetByte(index)
				if err != nil {
					return "", 0, err
				}

				if nextByte != Byte_EndWhile {
					return "", 0, DecompilerError("No endwhile byte", nextByte, index)
				}
				index++
			} else if expressionCode, bytesRead, err := DecompileExpression(index, indentationLevel); err == nil {
				currentLineCode.WriteString(expressionCode)
				index += bytesRead
			} else if b == Byte_EndScript || b == Byte_EndStruct || b == Byte_EndArray || b == Byte_EndIf || b == Byte_Else || b == Byte_EndWhile || b == Byte_LongJump {
				break
			} else {
				return "", 0, DecompilerError("Byte not recognised in body of code", b, index)
			}
		}

		flushCurrentLine()
		return TrimWhitespace(bodyOfCode.String()), index - initialIndex, nil
	}

	DecompileScript := func(index int) (string, int, error) {
		initialIndex := index
		index++

		scriptNameCode, bytesRead, err := DecompileChecksum(index)
		if err != nil {
			return "", 0, err
		}
		index += bytesRead

		scriptBodyCode, bytesRead, err := DecompileBodyOfCode(index, 1)
		if err != nil {
			return "", 0, err
		}
		index += bytesRead

		nextByte, err := GetByte(index)
		if nextByte != Byte_EndScript {
			return "", 0, DecompilerError("No endscript byte", nextByte, index)
		}
		index++

		return fmt.Sprintf("script %s {%s}", scriptNameCode, scriptBodyCode), index - initialIndex, nil
	}

	DecompileArgument := func(index, indentationLevel int) (string, int, error) {
		initialIndex := index

		firstExpressionCode, bytesRead, err := DecompileExpression(index, indentationLevel)
		if err != nil {
			return "", 0, err
		}
		index += bytesRead

		nextByte, err := GetByte(index)
		if err != nil {
			return firstExpressionCode, index - initialIndex, nil
		}

		if nextByte != Byte_Equals {
			return firstExpressionCode, index - initialIndex, nil
		}
		index++

		var extraNewLineCode strings.Builder
		for {
			nextByte, err := GetByte(index)
			if err != nil {
				return "", 0, err
			}
			if nextByte == Byte_NewLineWithNumber {
				newLineCode, bytesRead, err := DecompileNewLineWithNumber(index)
				if err != nil {
					return "", 0, err
				}
				index += bytesRead
				extraNewLineCode.WriteString(newLineCode)
			} else if nextByte == Byte_NewLine {
				index++
				extraNewLineCode.WriteString("\n")
			} else {
				break
			}
		}

		secondExpressionCode, bytesRead, err := DecompileExpression(index, indentationLevel)
		if err != nil {
			return "", 0, err
		}
		index += bytesRead

		var equalsCode string
		if indentationLevel == 0 {
			equalsCode = " = "
		} else {
			equalsCode = "="
		}

		return fmt.Sprintf("%s%s%s%s", firstExpressionCode, equalsCode, extraNewLineCode.String(), secondExpressionCode), index - initialIndex, nil
	}

	DecompileAtom := func(index, indentationLevel int) (string, int, error) {
		initialIndex := index

		b, err := GetByte(index)
		if err != nil {
			return "", 0, err
		}

		if b == Byte_Local || b == Byte_Checksum {
			isLocal := b == Byte_Local
			if isLocal {
				index++
			}

			firstChecksumCode, bytesRead, err := DecompileChecksum(index)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			if isLocal {
				firstChecksumCode = "<" + firstChecksumCode + ">"
			}

			var argumentCodeArray []string
			for {
				argumentCode, bytesRead, err := DecompileArgument(index, indentationLevel)
				if err != nil {
					break
				}
				argumentCodeArray = append(argumentCodeArray, argumentCode)
				index += bytesRead
			}

			if len(argumentCodeArray) == 0 {
				return firstChecksumCode, index - initialIndex, nil
			} else {
				argumentsCode := strings.Join(argumentCodeArray, " ")
				return fmt.Sprintf("%s %s", firstChecksumCode, argumentsCode), index - initialIndex, nil
			}
		} else if b == Byte_String {
			stringCode, bytesRead, err := DecompileString(index)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead
			return stringCode, index - initialIndex, nil
		} else if b == Byte_LocalString {
			stringCode, bytesRead, err := DecompileString(index)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead
			return stringCode, index - initialIndex, nil
		} else if b == Byte_Integer {
			index++
			integerBytes, err := GetBytes(index, 4)
			if err != nil {
				return "", 0, err
			}
			integer := binary.LittleEndian.Uint32(integerBytes)
			index += 4
			return fmt.Sprintf("%d", integer), index - initialIndex, nil
		} else if b == Byte_Not {
			index++
			expressionCode, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead
			return fmt.Sprintf("! %s", expressionCode), index - initialIndex, nil
		} else if b == Byte_Parenthesis {
			index++

			expressionCode, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			nextByte, err := GetByte(index)
			if err != nil {
				return "", 0, err
			}

			if nextByte != Byte_EndParenthesis {
				return "", 0, DecompilerError("No endparenthesis byte", nextByte, index)
			}
			index++

			return fmt.Sprintf("(%s)", expressionCode), index - initialIndex, nil
		} else if b == Byte_Struct {
			index++

			structBodyCode, bytesRead, err := DecompileBodyOfCode(index, indentationLevel+1)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			nextByte, err := GetByte(index)
			if err != nil {
				return "", 0, err
			}

			if nextByte != Byte_EndStruct {
				return "", 0, DecompilerError("No endstruct byte", nextByte, index)
			}
			index++

			var structClosingBraceCode string
			if strings.HasSuffix(structBodyCode, "\n") {
				structClosingBraceCode = Indent(indentationLevel, "}")
			} else {
				structClosingBraceCode = "}"
			}

			return fmt.Sprintf("{%s%s", structBodyCode, structClosingBraceCode), index - initialIndex, nil
		} else if b == Byte_Array {
			index++

			arrayBodyCode, bytesRead, err := DecompileBodyOfCode(index, indentationLevel+1)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			nextByte, err := GetByte(index)
			if err != nil {
				return "", 0, err
			}

			if nextByte != Byte_EndArray {
				return "", 0, DecompilerError("No endarray byte", nextByte, index)
			}
			index++

			var arrayClosingBracketCode string
			if strings.HasSuffix(arrayBodyCode, "\n") {
				arrayClosingBracketCode = Indent(indentationLevel, "]")
			} else {
				arrayClosingBracketCode = "]"
			}

			return fmt.Sprintf("[%s%s", arrayBodyCode, arrayClosingBracketCode), index - initialIndex, nil
		} else if b == Byte_Float {
			index++

			floatBytes, err := GetBytes(index, 4)
			if err != nil {
				return "", 0, err
			}
			index += 4

			return FormatFloat(floatBytes), index - initialIndex, nil
		} else if b == Byte_Pair {
			pairCode, bytesRead, err := DecompilePair(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return pairCode, index - initialIndex, nil
		} else if b == Byte_Vector {
			index++

			xBytes, err := GetBytes(index, 4)
			if err != nil {
				return "", 0, err
			}
			index += 4

			yBytes, err := GetBytes(index, 4)
			if err != nil {
				return "", 0, err
			}
			index += 4

			zBytes, err := GetBytes(index, 4)
			if err != nil {
				return "", 0, err
			}
			index += 4

			return fmt.Sprintf("(%s, %s, %s) ", FormatFloat(xBytes), FormatFloat(yBytes), FormatFloat(zBytes)), index - initialIndex, nil
		} else if b == Byte_AllArguments {
			index++
			return "<...>", index - initialIndex, nil
		} else if b == Byte_Random || b == Byte_RandomNoRepeat {
			initialIndex := index
			index++

			numberOfBranchesBytes, err := GetBytes(index, 4)
			if err != nil {
				return "", 0, err
			}
			numberOfBranches := int(binary.LittleEndian.Uint32(numberOfBranchesBytes))
			index += 4

			branchOffsets := make([]int, numberOfBranches)
			for i := 0; i < numberOfBranches; i++ {
				branchSizeBytes, err := GetBytes(index, 4)
				if err != nil {
					return "", 0, err
				}
				branchSize := int(binary.LittleEndian.Uint32(branchSizeBytes))
				branchOffsets[i] = branchSize
				index += 4
			}

			branches := make([]string, numberOfBranches)
			lastBranchSize := 0
			for i := 0; i < numberOfBranches; i++ {
				branchIndex := index + branchOffsets[i] - (4 * numberOfBranches) + (4 * (i + 1))
				branchCode, bytesRead, err := DecompileBodyOfCode(branchIndex, indentationLevel)
				if err != nil {
					return "", 0, err
				}
				branches[i] = branchCode
				lastBranchSize = bytesRead
			}

			index += branchOffsets[numberOfBranches - 1]
			index += lastBranchSize

			for i, branch := range branches {
				// wrap each branch in {}
				branches[i] = fmt.Sprintf("{ %s }", branch)
			}

			return fmt.Sprintf("random{ %s }", strings.Join(branches, " ")), index - initialIndex, nil
		} else if b == Byte_RandomRange {
			index++

			pairCode, bytesRead, err := DecompilePair(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}

			index += bytesRead

			return fmt.Sprintf("randomrange%s", pairCode), index - initialIndex, nil
		}

		return "", 0, DecompilerError("Not an atom", b, index)
	}

	DecompileExpression = func(index, indentationLevel int) (string, int, error) {
		initialIndex := index

		atomCode, bytesRead, err := DecompileAtom(index, indentationLevel)
		if err != nil {
			return "", 0, err
		}
		index += bytesRead

		nextByte, err := GetByte(index)
		if err != nil {
			return atomCode, index - initialIndex, nil
		}

		if nextByte == Byte_Plus {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s + %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Minus {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s - %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Multiply {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s * %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Divide {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s / %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_And {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s & %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Or {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s | %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Xor {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s ^ %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Equals || nextByte == Byte_EqualTo {
			index++

			var extraNewLineCode strings.Builder
			for {
				nextByte, err := GetByte(index)
				if err != nil {
					return "", 0, err
				}
				if nextByte == Byte_NewLineWithNumber {
					newLineCode, bytesRead, err := DecompileNewLineWithNumber(index)
					if err != nil {
						return "", 0, err
					}
					index += bytesRead
					extraNewLineCode.WriteString(newLineCode)
				} else if nextByte == Byte_NewLine {
					index++
					extraNewLineCode.WriteString("\n")
				} else {
					break
				}
			}

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s = %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Dot {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s.%s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_Colon {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s:%s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_GreaterThan {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s > %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_GreaterThanEqual {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s >= %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_LessThan {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s < %s", atomCode, nextExpression), index - initialIndex, nil
		} else if nextByte == Byte_LessThanEqual {
			index++

			nextExpression, bytesRead, err := DecompileExpression(index, indentationLevel)
			if err != nil {
				return "", 0, err
			}
			index += bytesRead

			return fmt.Sprintf("%s <= %s", atomCode, nextExpression), index - initialIndex, nil
		}

		return atomCode, index - initialIndex, err
	}

	checksumTable, err := GetChecksumTable()
	if err != nil {
		return "", err
	}

	var output strings.Builder
	index := 0
	for {
		if index >= len(qb) {
			break
		}

		b, err := GetByte(index)
		if err != nil {
			return "", err
		}

		if b == Byte_EndOfFile {
			index++
			break
		} else if b == Byte_NewLine {
			output.WriteString("\n")
			index++
		} else if b == Byte_NewLineWithNumber {
			newLineCode, bytesRead, err := DecompileNewLineWithNumber(index)
			if err != nil {
				return "", err
			}
			index += bytesRead
			output.WriteString(newLineCode)
		} else if b == Byte_Checksum {
			argumentCode, bytesRead, err := DecompileArgument(index, 0)
			if err != nil {
				return "", err
			}
			output.WriteString(argumentCode)
			index += bytesRead
		} else if b == Byte_Script {
			scriptCode, bytesRead, err := DecompileScript(index)
			if err != nil {
				return "", err
			}
			index += bytesRead
			output.WriteString(scriptCode)
		} else if b == Byte_ChecksumEntry {
			index++
			index+=4
			for {
				nextByte, err := GetByte(index)
				if err != nil {
					return "", err
				}
				if nextByte == 0 {
					index++
					break
				}
				index++
			}
		} else {
			break
		}
	}

	if index < len(qb) {
		nextByte, _ := GetByte(index)
		message := fmt.Sprintf("Did not finish decompiling.\n%s\n0x%x/0x%x bytes decompiled.\nnext byte: 0x%x", output.String(), index, len(qb), nextByte)
		return "", errors.New(message)
	}

	return output.String(), nil
}
