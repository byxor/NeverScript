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
				return "", err
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
	case compiler.AstKind_AllArguments:
		return "<...>", nil
	case compiler.AstKind_LogicalNot:
		data := node.Data.(compiler.AstData_UnaryExpression)
		expressionCode, err := DecompileAstNode(data.Node, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return "! " + expressionCode, nil
	case compiler.AstKind_LocalReference:
		expressionCode, err := DecompileAstNode(node.Data.(compiler.AstData_UnaryExpression).Node, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return "<" + expressionCode + ">", nil
	case compiler.AstKind_Script:
		data := node.Data.(compiler.AstData_Script)
		var code strings.Builder
		code.WriteString("script ")
		name, err := DecompileAstNode(data.NameNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		code.WriteString(name)
		for _, defaultParameterNode := range data.DefaultParameterNodes {
			code.WriteString(" ")
			nodeCode, err := DecompileAstNode(defaultParameterNode, indentation, nameTable)
			if err != nil {
				return "", err
			}
			nodeCode = strings.Replace(nodeCode, " = ", "=", 1)
			code.WriteString(nodeCode)
		}
		code.WriteString(" {")
		indentation++
		for i, bodyNode := range data.BodyNodes {
			isFirstNode := i == 0
			decompiledNode, err := DecompileAstNode(bodyNode, indentation, nameTable)
			if err != nil {
				return "", err
			}
			if !isFirstNode && data.BodyNodes[i - 1].Kind == compiler.AstKind_NewLine {
				code.WriteString(strings.Repeat("    ", indentation))
			}
			code.WriteString(decompiledNode)
		}
		indentation--
		code.WriteString("}")
		return code.String(), nil
	case compiler.AstKind_Invocation:
		data := node.Data.(compiler.AstData_Invocation)
		var code strings.Builder
		decompiledName, err := DecompileAstNode(data.ScriptIdentifierNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		code.WriteString(decompiledName)
		if len(data.ParameterNodes) <= 2 { // render params on 1 line
			for i, parameterNode := range data.ParameterNodes {
				decompiledNode, err := DecompileAstNode(parameterNode, indentation, nameTable)
				if err != nil {
					return "", err
				}
				isFirstNode := i == 0
				isLastNode := i == len(data.ParameterNodes)-1
				if isFirstNode {
					if parameterNode.Kind != compiler.AstKind_NewLine {
						code.WriteString(" ")
					}
				}
				if parameterNode.Kind == compiler.AstKind_Assignment {
					decompiledNode = strings.Replace(decompiledNode, " = ", "=", 1)
				}
				code.WriteString(decompiledNode)
				if !isLastNode {
					code.WriteString(" ")
				}
			}
		} else { // render params across multiple lines
			indentation++
			for _, parameterNode := range data.ParameterNodes {
				decompiledNode, err := DecompileAstNode(parameterNode, indentation, nameTable)
				if err != nil {
					return "", err
				}
				code.WriteString(" \\\n")
				code.WriteString(strings.Repeat("    ", indentation))
				code.WriteString(decompiledNode)
			}
			indentation--
		}
		return code.String(), nil
	case compiler.AstKind_UnaryExpression:
		data := node.Data.(compiler.AstData_UnaryExpression)
		nodeCode, err := DecompileAstNode(data.Node, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return "(" + nodeCode + ")", nil
	case compiler.AstKind_Checksum:
		data := node.Data.(compiler.AstData_Checksum)
		if data.ChecksumToken.Data != "" {
			return data.ChecksumToken.Data, nil
		}
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
	case compiler.AstKind_Vector:
		data := node.Data.(compiler.AstData_Vector)
		var left, middle, right float32
		{
			floatData := data.FloatNodeA.Data.(compiler.AstData_Float)
			bits := binary.LittleEndian.Uint32(floatData.FloatBytes)
			left = math.Float32frombits(bits)
		}
		{
			floatData := data.FloatNodeB.Data.(compiler.AstData_Float)
			bits := binary.LittleEndian.Uint32(floatData.FloatBytes)
			middle = math.Float32frombits(bits)
		}
		{
			floatData := data.FloatNodeC.Data.(compiler.AstData_Float)
			bits := binary.LittleEndian.Uint32(floatData.FloatBytes)
			right = math.Float32frombits(bits)
		}
		return fmt.Sprintf("(%s, %s, %s)", RenderFloat(left), RenderFloat(middle), RenderFloat(right)), nil
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
		indentation--
		if len(data.ElementNodes) > 0 {
			lastElement := data.ElementNodes[len(data.ElementNodes)-1]
			if lastElement.Kind == compiler.AstKind_NewLine {
				code.WriteString(strings.Repeat("    ", indentation))
			}
		}
		code.WriteString("}")
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
		indentation--
		if len(data.ElementNodes) > 0 {
			lastElement := data.ElementNodes[len(data.ElementNodes)-1]
			if lastElement.Kind == compiler.AstKind_NewLine {
				code.WriteString(strings.Repeat("    ", indentation))
			}
		}
		code.WriteString("]")
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
	case compiler.AstKind_EqualsExpression:
		data := node.Data.(compiler.AstData_BinaryExpression)
		leftDecompiled, err := DecompileAstNode(data.LeftNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		rightDecompiled, err := DecompileAstNode(data.RightNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return leftDecompiled + " = " + rightDecompiled, nil
	case compiler.AstKind_LessThanExpression:
		data := node.Data.(compiler.AstData_BinaryExpression)
		leftDecompiled, err := DecompileAstNode(data.LeftNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		rightDecompiled, err := DecompileAstNode(data.RightNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return leftDecompiled + " < " + rightDecompiled, nil
	case compiler.AstKind_DotExpression:
		data := node.Data.(compiler.AstData_BinaryExpression)
		leftDecompiled, err := DecompileAstNode(data.LeftNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		rightDecompiled, err := DecompileAstNode(data.RightNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return leftDecompiled + "." + rightDecompiled, nil
	case compiler.AstKind_ColonExpression:
		data := node.Data.(compiler.AstData_BinaryExpression)
		leftDecompiled, err := DecompileAstNode(data.LeftNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		rightDecompiled, err := DecompileAstNode(data.RightNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return leftDecompiled + ":" + rightDecompiled, nil
	case compiler.AstKind_GreaterThanExpression:
		data := node.Data.(compiler.AstData_BinaryExpression)
		leftDecompiled, err := DecompileAstNode(data.LeftNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		rightDecompiled, err := DecompileAstNode(data.RightNode, indentation, nameTable)
		if err != nil {
			return "", err
		}
		return leftDecompiled + " > " + rightDecompiled, nil
	case compiler.AstKind_IfStatement:
		data := node.Data.(compiler.AstData_IfStatement)
		var code strings.Builder
		code.WriteString("if ")
		conditionCode, err := DecompileAstNode(data.Conditions[0], indentation, nameTable)
		if err != nil {
			return "", err
		}
		code.WriteString(conditionCode)
		renderBody := func(body []compiler.AstNode) (string, error) {
			var bodyCode strings.Builder
			bodyCode.WriteString(" {")
			indentation++
			for i, bodyNode := range body {
				isFirstNode := i == 0
				nodeCode, err := DecompileAstNode(bodyNode, indentation, nameTable)
				if err != nil {
					return "", err
				}
				if !isFirstNode && len(body) != 0 && body[len(body)-1].Kind == compiler.AstKind_NewLine {
					bodyCode.WriteString(strings.Repeat("    ", indentation))
				}
				bodyCode.WriteString(nodeCode)
			}
			indentation--
			if len(body) > 0 {
				bodyCode.WriteString(strings.Repeat("    ", indentation))
			}
			bodyCode.WriteString("}")
			return bodyCode.String(), nil
		}
		bodyCode, err := renderBody(data.Bodies[0])
		if err != nil {
			return "", err
		}
		code.WriteString(bodyCode)
		if len(data.Bodies) > 1 { // has 'else'
			code.WriteString(" else")
			elseBodyCode, err := renderBody(data.Bodies[1])
			if err != nil {
				return "", err
			}
			code.WriteString(elseBodyCode)
		}
		return code.String(), nil
	case compiler.AstKind_NameTableEntry:
		return "", nil
	}
	return "", errors.New(WrapLine("Don't know how to produce code for AST node", fmt.Sprintf("%+v", node)))
}

func RenderFloat(f float32) string {
	result := strconv.FormatFloat(float64(f), 'f', -1, 32)
	if !strings.Contains(result, ".") {
		result += ".0"
	}
	return result
}