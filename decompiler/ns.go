package decompiler

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"math"
	"strconv"
	"strings"
)

func DecompileAstNode(node compiler.AstNode, indentation int, nameTable map[uint32]string) (string, error) {
	switch node.Kind {
	case compiler.AstKind_Root:
		data := node.Data.(compiler.AstData_Root)
		text := ""
		for _, bodyNode := range data.BodyNodes {
			decompiledBodyNode, err := DecompileAstNode(bodyNode, indentation, nameTable)
			if err != nil {
				return "", errors.New(fmt.Sprintf("Failed to decompile root body node: %+v", err))
			}
			text += decompiledBodyNode
		}
		return text, nil
	case compiler.AstKind_NewLine:
		return "\n", nil
	case compiler.AstKind_EndOfFile:
		return "", nil
	case compiler.AstKind_Comma:
		return ",", nil
	case compiler.AstKind_Checksum:
		data := node.Data.(compiler.AstData_Checksum)
		checksum := binary.LittleEndian.Uint32(data.ChecksumBytes)
		if resolvedName, ok := nameTable[checksum]; ok {
			return resolvedName, nil
		} else {
			checksumHex := fmt.Sprintf("%#X", data.ChecksumBytes)[2:]
			checksumHex = strings.Repeat("0", 8-len(checksumHex)) + checksumHex
			return "#" + checksumHex, nil
		}
	case compiler.AstKind_Float:
		data := node.Data.(compiler.AstData_Float)
		bits := binary.LittleEndian.Uint32(data.FloatBytes)
		number := math.Float32frombits(bits)
		return RenderFloat(number), nil
	case compiler.AstKind_Integer:
		data := node.Data.(compiler.AstData_Integer)
		number := int32(binary.LittleEndian.Uint32(data.IntegerBytes))
		return fmt.Sprintf("%d", number), nil
	case compiler.AstKind_Pair:
		data := node.Data.(compiler.AstData_Pair)
		var left, right float32
		{
			floatData := data.FloatNodeA.Data.(compiler.AstData_Float)
			bits := binary.LittleEndian.Uint32(floatData.FloatBytes)
			left = math.Float32frombits(bits)
		}
		{
			floatData := data.FloatNodeB.Data.(compiler.AstData_Float)
			bits := binary.LittleEndian.Uint32(floatData.FloatBytes)
			right = math.Float32frombits(bits)
		}
		return fmt.Sprintf("(%s, %s)", RenderFloat(left), RenderFloat(right)), nil
	case compiler.AstKind_String:
		data := node.Data.(compiler.AstData_String)
		return fmt.Sprintf("\"%s\"", data.StringBytes), nil
	case compiler.AstKind_Struct:
		data := node.Data.(compiler.AstData_Struct)
		var code strings.Builder
		code.WriteString("{")
		if len(data.ElementNodes) > 0 {
			code.WriteString(" ")
		}
		indentation++
		for i, element := range data.ElementNodes {
			elementCode, err := DecompileAstNode(element, indentation, nameTable)
			if err != nil {
				return "", err
			}
			isLastElement := i == len(data.ElementNodes)-1
			isFirstElement := i == 0
			if (isFirstElement || data.ElementNodes[i - 1].Kind != compiler.AstKind_NewLine) &&
				element.Kind == compiler.AstKind_Assignment {
				elementCode = strings.Replace(elementCode, " = ", "=", 1)
			}
			code.WriteString(elementCode)
			if element.Kind == compiler.AstKind_NewLine {
				if !isLastElement {
					code.WriteString(strings.Repeat("    ", indentation))
				}
			} else {
				if (!isLastElement && data.ElementNodes[i + 1].Kind != compiler.AstKind_Comma) ||
					(isLastElement && element.Kind != compiler.AstKind_NewLine) {
					code.WriteString(" ")
				}
			}
		}
		code.WriteString("}")
		indentation--
		return code.String(), nil
	case compiler.AstKind_Array:
		data := node.Data.(compiler.AstData_Array)
		var code strings.Builder
		code.WriteString("[")
		if len(data.ElementNodes) > 0 {
			code.WriteString(" ")
		}
		indentation++
		for i, element := range data.ElementNodes {
			elementCode, err := DecompileAstNode(element, indentation, nameTable)
			if err != nil {
				return "", err
			}
			isLastElement := i == len(data.ElementNodes)-1
			code.WriteString(elementCode)
			if element.Kind == compiler.AstKind_NewLine {
				if !isLastElement {
					code.WriteString(strings.Repeat("    ", indentation))
				}
			} else {
				if (!isLastElement && data.ElementNodes[i + 1].Kind != compiler.AstKind_Comma) ||
					(isLastElement && element.Kind != compiler.AstKind_NewLine) {
					code.WriteString(" ")
				}
			}
		}
		code.WriteString("]")
		indentation--
		return code.String(), nil
	case compiler.AstKind_Assignment:
		data := node.Data.(compiler.AstData_Assignment)
		decompiledName, err := DecompileAstNode(data.NameNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		decompiledValue, err := DecompileAstNode(data.ValueNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", decompiledName, decompiledValue), nil
	case compiler.AstKind_NameTableEntry:
		return "", nil
	}
	return "code", errors.New(fmt.Sprintf("Failed to decompile ast node: %+v", node))
}

func RenderFloat(f float32) string {
	result := strconv.FormatFloat(float64(f), 'f', -1, 32)
	if !strings.Contains(result, ".") {
		result += ".0"
	}
	return result
}