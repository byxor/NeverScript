package newcompiler

import (
    "bytes"
    "encoding/binary"
    "encoding/hex"
    "errors"
    "fmt"
    "github.com/byxor/NeverScript/compiler"
    "math"
    "strconv"
)

func ProduceQb(program Node) ([]byte, error) {
    var output output
    output.engineSupportsIf2 = false
    output.writeLineNumbers = false
    output.writeQbKeys = false
    output.nameTable = make(map[string]uint32)
    err := output.writeQb(program)
    return output.qb.Bytes(), err
}

// ------------ internal --------------

type output struct {
    engineSupportsIf2 bool
    writeLineNumbers  bool
    writeQbKeys       bool

    qb        bytes.Buffer
    nameTable map[string]uint32
}

func (this *output) write(bytes ...byte) {
    this.qb.Write(bytes)
}

func (this *output) writeFloat(value float32) {
    var buffer [4]byte
    binary.LittleEndian.PutUint32(buffer[:], math.Float32bits(float32(value)))
    this.write(buffer[:]...)
}

func (this *output) writeLittleEndianUint32(value uint32) {
    var buffer [4]byte
    binary.LittleEndian.PutUint32(buffer[:], value)
    this.write(buffer[:]...)
}

func (this *output) writeQb(node Node) error {
    switch node.Kind() {
    case NodeKind_Program:
        return this.writeProgramQb(node)
    case NodeKind_Script:
        return this.writeScriptQb(node)
    case NodeKind_LocalQbKey:
        return this.writeLocalQbKeyQb(node)
    case NodeKind_QbKey:
        return this.writeQbKeyQb(node)
    case NodeKind_RawQbKey:
        return this.writeRawQbKeyQb(node)
    case NodeKind_String:
        return this.writeStringQb(node)
    case NodeKind_Int:
        return this.writeIntQb(node)
    case NodeKind_Float:
        return this.writeFloatQb(node)
    case NodeKind_LineBreak:
        if this.writeLineNumbers {
            this.write(0x2)
            this.writeLittleEndianUint32(uint32(node.LineNumber()))
        } else {
            this.write(0x1)
        }
        return nil
    case NodeKind_Comma:
        this.write(0x9)
        return nil
    case NodeKind_Loop:
        return this.writeLoopQb(node)
    case NodeKind_IfStatement:
        return this.writeIfStatementQb(node)
    case NodeKind_If:
        return this.writeIfQb(node)
    case NodeKind_Else:
        return this.writeElseQb(node)
    case NodeKind_Return:
        return this.writeReturnQb(node)
    case NodeKind_Random:
        return this.writeRandomQb(node)
    case NodeKind_Pair:
        return this.writePairQb(node)
    case NodeKind_Vector:
        return this.writeVectorQb(node)
    case NodeKind_Array:
        return this.writeArrayQb(node)
    case NodeKind_Struct:
        return this.writeStructQb(node)
    case NodeKind_UnaryMinusOperation:
        return this.writeUnaryMinusOperationQb(node)
    case NodeKind_NotOperation:
        return this.writeNotOperationQb(node)
    case NodeKind_ParenthesisOperation:
        return this.writeParenthesisOperationQb(node)
    case NodeKind_Bytes:
        return this.writeBytesQb(node)
    case NodeKind_AssignmentOperation:
        return this.writeBinaryOperationQb(node, 0x7)
    case NodeKind_PlusOperation:
        return this.writeBinaryOperationQb(node, 0xB)
    case NodeKind_MinusOperation:
        return this.writeBinaryOperationQb(node, 0xA)
    case NodeKind_DivideOperation:
        return this.writeBinaryOperationQb(node, 0xC)
    case NodeKind_MultiplyOperation:
        return this.writeBinaryOperationQb(node, 0xD)
    case NodeKind_LessThanOperation:
        return this.writeBinaryOperationQb(node, 0x12)
    case NodeKind_GreaterThanOperation:
        return this.writeBinaryOperationQb(node, 0x14)
    case NodeKind_AndOperation:
        return this.writeBinaryOperationQb(node, 0x33)
    case NodeKind_OrOperation:
        return this.writeBinaryOperationQb(node, 0x32)
    case NodeKind_ColonOperation:
        return this.writeBinaryOperationQb(node, 0x42)
    case NodeKind_DotOperation:
        return this.writeBinaryOperationQb(node, 0x8)
    case NodeKind_EqualityOperation:
        return this.writeBinaryOperationQb(node, 0x11)
    case NodeKind_LessThanEqualOperation:
        fallthrough
    case NodeKind_GreaterThanEqualOperation:
        fallthrough
    case NodeKind_InequalityOperation:
        fallthrough
    case NodeKind_PlusEqualOperation:
        fallthrough
    case NodeKind_MinusEqualOperation:
        fallthrough
    case NodeKind_MultiplyEqualOperation:
        fallthrough
    case NodeKind_DivideEqualOperation:
        fallthrough
    case NodeKind_ElseIf:
        fallthrough
    default:
        return errors.New(fmt.Sprintf("QB output not implemented for '%s'", node.Kind()))
    }
}

func (this *output) writeProgramQb(node Node) error {
    for _, wrappedNode := range node.(wrappedNodes).nodes {
        err := this.writeQb(wrappedNode)
        if err != nil {
            return err
        }
    }

    if this.writeQbKeys {
        for identifier, qbKey := range this.nameTable {
            this.write(0x2B)
            this.writeLittleEndianUint32(qbKey)
            this.write([]byte(identifier)...)
            this.write(0)
        }
    }

    this.write(0)
    return nil
}

func (this *output) writeBinaryOperationQb(node Node, operators ...byte) error {
    var data manyWrappedNodes
    data, ok := node.(manyWrappedNodes)
    if !ok {
        data = node.(fixedSizeWrappedNode).node
    }

    leftHand := data.nodeLists[0][0]
    rightHand := data.nodeLists[0][1]

    err := this.writeQb(leftHand)
    if err != nil {
        return err
    }

    lineBreaks1 := data.nodeLists[1]
    for _, lineBreak := range lineBreaks1 {
        err := this.writeQb(lineBreak)
        if err != nil {
            return err
        }
    }

    for _, operator := range operators {
        this.write(operator)
    }

    lineBreaks2 := data.nodeLists[2]
    for _, lineBreak := range lineBreaks2 {
        err := this.writeQb(lineBreak)
        if err != nil {
            return err
        }
    }

    if rightHand != nil {
        err = this.writeQb(rightHand)
        if err != nil {
            return err
        }
    }

    return nil
}

func (this *output) writeScriptQb(node Node) error {
    manyWrappedNodes := node.(manyWrappedNodes)
    this.write(0x23)

    qbKeyNode := manyWrappedNodes.nodeLists[0][0]
    err := this.writeQb(qbKeyNode)
    if err != nil {
        return err
    }

    scriptHeaderNodes := manyWrappedNodes.nodeLists[1]
    for _, innerNode := range scriptHeaderNodes {
        err := this.writeQb(innerNode)
        if err != nil {
            return err
        }
    }

    bodyNodes := manyWrappedNodes.nodeLists[2]
    for _, innerNode := range bodyNodes {
        err := this.writeQb(innerNode)
        if err != nil {
            return err
        }
    }

    this.write(0x24)
    return nil
}

func (this *output) writeArrayQb(node Node) error {
    this.write(0x5)
    for _, innerNode := range node.(wrappedNodes).nodes {
        this.writeQb(innerNode)
    }
    this.write(0x6)
    return nil
}

func (this *output) writeStructQb(node Node) error {
    this.write(0x3)
    for _, innerNode := range node.(wrappedNodes).nodes {
        this.writeQb(innerNode)
    }
    this.write(0x4)
    return nil
}

func (this *output) writeStringQb(node Node) error {
    string_ := node.(basicNode).data
    size := uint32(len(string_))
    this.write(0x1B)
    this.writeLittleEndianUint32(size + 1)
    this.write([]byte(string_[:size])...)
    this.write(0x0)
    return nil
}

func (this *output) writeIntQb(node Node) error {
    int_, err := strconv.ParseInt(node.(basicNode).data, 10, 32)
    if err != nil {
        return err
    }
    this.write(0x17)
    this.writeLittleEndianUint32(uint32(int_))
    return nil
}

func (this *output) writeFloatQb(node Node) error {
    floatValue, err := strconv.ParseFloat(node.(basicNode).data, 32)
    if err != nil {
        return err
    }
    this.write(0x1A)
    this.writeFloat(float32(floatValue))
    return nil
}

func (this *output) writeLoopQb(node Node) error {
    loopBodyNodes := node.(manyWrappedNodes).nodeLists[0]
    expressionNode := node.(manyWrappedNodes).nodeLists[1][0]

    this.write(0x20)

    for _, loopBodyNode := range loopBodyNodes {
        err := this.writeQb(loopBodyNode)
        if err != nil {
            return err
        }
    }

    this.write(0x21)

    if expressionNode != nil {
        return this.writeQb(expressionNode)
    }

    return nil
}

func (this *output) writeLocalQbKeyQb(node Node) error {
    this.write(0x2D)
    return this.writeQb(node.(wrappedNode).node)
}

func (this *output) writeBytesQb(node Node) error {
    for _, byte_ := range node.(wrappedNodes).nodes {
        hexString := byte_.(basicNode).data
        hexBytes, err := hex.DecodeString(hexString)
        if err != nil {
            return err
        }
        this.write(hexBytes...)
    }
    return nil
}

func (this *output) writeQbKeyQb(node Node) error {
    this.nameTable[node.(basicNode).data] = compiler.StringToChecksum(node.(basicNode).data)
    basicNode := node.(basicNode)
    this.write(0x16)
    this.write(toQbKey(basicNode.data)...)
    return nil
}

func (this *output) writeRawQbKeyQb(node Node) error {
    rawQbKeyNode := node.(rawQbKeyNode)
    this.write(0x16)
    this.write(rawQbKeyNode.key...)
    return nil
}

func toQbKey(identifier string) []byte {
    key := compiler.StringToChecksum(identifier)
    keyBytes := []byte{0, 0, 0, 0}
    binary.LittleEndian.PutUint32(keyBytes, key)
    return keyBytes
}

func (this *output) writePairQb(node Node) error {
    wrappedNodes := node.(wrappedNodes)

    floatAValue, err := strconv.ParseFloat(wrappedNodes.nodes[0].(basicNode).data, 32)
    if err != nil {
        return err
    }

    floatBValue, err := strconv.ParseFloat(wrappedNodes.nodes[1].(basicNode).data, 32)
    if err != nil {
        return err
    }

    this.write(0x1F)
    this.writeFloat(float32(floatAValue))
    this.writeFloat(float32(floatBValue))
    return nil
}

func (this *output) writeVectorQb(node Node) error {
    wrappedNodes := node.(wrappedNodes)

    floatAValue, err := strconv.ParseFloat(wrappedNodes.nodes[0].(basicNode).data, 32)
    if err != nil {
        return err
    }

    floatBValue, err := strconv.ParseFloat(wrappedNodes.nodes[1].(basicNode).data, 32)
    if err != nil {
        return err
    }

    floatCValue, err := strconv.ParseFloat(wrappedNodes.nodes[2].(basicNode).data, 32)
    if err != nil {
        return err
    }

    this.write(0x1E)
    this.writeFloat(float32(floatAValue))
    this.writeFloat(float32(floatBValue))
    this.writeFloat(float32(floatCValue))
    return nil
}

func (this *output) writeIfStatementQb(node Node) error {
    ifNode := node.(wrappedNodes).nodes[0]
    elseNode := node.(wrappedNodes).nodes[1]

    err := this.writeIfQb(ifNode)
    if err != nil {
        return err
    }

    if elseNode != nil {
        err = this.writeElseQb(elseNode)
        if err != nil {
            return err
        }
    }

    this.write(0x28)
    return nil
}

func (this *output) writeIfQb(node Node) error {
    this.write(0x25)

    conditionNodes := node.(manyWrappedNodes).nodeLists[0]
    bodyNodes := node.(manyWrappedNodes).nodeLists[1]

    for _, conditionNode := range conditionNodes {
        err := this.writeQb(conditionNode)
        if err != nil {
            return err
        }
    }

    for _, bodyNode := range bodyNodes {
        err := this.writeQb(bodyNode)
        if err != nil {
            return err
        }
    }

    return nil
}

func (this *output) writeElseQb(node Node) error {
    this.write(0x26)
    for _, elseBodyNode := range node.(wrappedNodes).nodes {
        if elseBodyNode != nil {
            err := this.writeQb(elseBodyNode)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func (this *output) writeReturnQb(node Node) error {
    this.write(0x29)
    for _, innerNode := range node.(wrappedNodes).nodes {
        err := this.writeQb(innerNode)
        if err != nil {
            return err
        }
    }
    return nil
}

func (this *output) writeRandomQb(node Node) error {
    this.write(0x2F)

    randomEntries := node.(wrappedNodes).nodes
    this.writeLittleEndianUint32(uint32(len(randomEntries)))

    return nil
}

func (this *output) writeUnaryMinusOperationQb(node Node) error {
    this.write(0xA)
    err := this.writeQb(node.(wrappedNode).node)
    return err
}

func (this *output) writeParenthesisOperationQb(node Node) error {
    this.write(0xE)
    err := this.writeQb(node.(wrappedNode).node)
    if err != nil {
        return err
    }
    this.write(0xF)
    return nil
}

func (this *output) writeNotOperationQb(node Node) error {
    this.write(0x39)
    err := this.writeQb(node.(wrappedNode).node)
    return err
}
