package newcompiler

import (
    "errors"
    "strings"
    "unicode"
)

type TokenKind int

const (
    TokenKind_Identifier TokenKind = iota
    TokenKind_Equals
    TokenKind_Int
    TokenKind_String
    TokenKind_LeftSquareBracket
    TokenKind_RightSquareBracket
    TokenKind_Comma
    TokenKind_NewLine
    TokenKind_Colon
    TokenKind_AtSymbol
    TokenKind_LeftCurlyBrace
    TokenKind_RightCurlyBrace
    TokenKind_LeftParenthesis
    TokenKind_RightParenthesis
    TokenKind_LeftAngleBracket
    TokenKind_RightAngleBracket
    TokenKind_RawQbKey
    TokenKind_Plus
    TokenKind_Minus
    TokenKind_Asterisk
    TokenKind_ForwardSlash
    TokenKind_BackwardSlash
    TokenKind_EscapedLineBreak
    TokenKind_Exclamation
    TokenKind_Random
    TokenKind_Bytes
    TokenKind_If
    TokenKind_Else
    TokenKind_Loop
    TokenKind_Break
    TokenKind_Return
    TokenKind_Script
    TokenKind_Float
    TokenKind_SingleLineComment
    TokenKind_MultiLineComment
    TokenKind_Dot
    TokenKind_And
    TokenKind_Or
    TokenKind_Switch
    TokenKind_Case
    TokenKind_Default
    TokenKind_Space
    TokenKind_Tab
    TokenKind_CarriageReturn
)

type Token interface {
    Kind() TokenKind
    Data() string
    LineNumber() uint
    CharsConsumed() uint
    LinesConsumed() uint
}

func Lex(sourceCode string) ([]Token, error) {
    var lexer lexer
    lexer.sourceCode = sourceCode
    lexer.preventConsecutiveLineBreaks = true
    lexer.index = 0
    lexer.lineNumber = 1
    lexer.tokens = make([]Token, 0)

    return lexer.lex()
}

// ---------------- internal -------------------

type genericToken struct {
    kind          TokenKind
    data          string
    lineNumber    uint
    charsConsumed uint
    linesConsumed uint
}

func (this genericToken) Kind() TokenKind {
    return this.kind
}

func (this genericToken) Data() string {
    return this.data
}

func (this genericToken) LineNumber() uint {
    return this.lineNumber
}

func (this genericToken) CharsConsumed() uint {
    return this.charsConsumed
}

func (this genericToken) LinesConsumed() uint {
    return this.linesConsumed
}

func newGenericToken(kind TokenKind, data string, lineNumber uint, charsConsumed uint, linesConsumed uint) genericToken {
    return genericToken{
        kind:          kind,
        data:          data,
        lineNumber:    lineNumber,
        charsConsumed: charsConsumed,
        linesConsumed: linesConsumed,
    }
}

type lexer struct {
    sourceCode                   string
    preventConsecutiveLineBreaks bool
    index                        uint
    lineNumber                   uint
    tokens                       []Token
}

func (this *lexer) lex() ([]Token, error) {

    for {
        if this.isOutOfRangeAt(this.index) {
            break
        }

        escapedLineBreak, err := this.tryGetEscapedLineBreak()
        if err != nil {
            return nil, err
        } else if escapedLineBreak != nil {
            _ = this.saveToken(escapedLineBreak)
            continue
        }

        float_, err := this.tryGetFloat()
        if err != nil {
            return nil, err
        } else if float_ != nil {
            _ = this.saveToken(float_)
            continue
        }

        int_, err := this.tryGetInt()
        if err != nil {
            return nil, err
        } else if int_ != nil {
            _ = this.saveToken(int_)
            continue
        }

        string_, err := this.tryGetString()
        if err != nil {
            return nil, err
        } else if string_ != nil {
            _ = this.saveToken(string_)
            continue
        }

        singleLineComment, err := this.tryGetSingleLineComment()
        if err != nil {
            return nil, err
        } else if singleLineComment != nil {
            _ = this.saveToken(singleLineComment)
            continue
        }

        multiLineComment, err := this.tryGetMultiLineComment()
        if err != nil {
            return nil, err
        } else if multiLineComment != nil {
            _ = this.saveToken(multiLineComment)
            continue
        }

        rawQbKey, err := this.tryGetRawQbKey()
        if err != nil {
            return nil, err
        } else if rawQbKey != nil {
            _ = this.saveToken(rawQbKey)
            continue
        }

        tab, err := this.tryGetCharSequence("\t", TokenKind_Tab)
        if err != nil {
            return nil, err
        } else if tab != nil {
            _ = this.saveToken(tab)
            continue
        }

        space, err := this.tryGetCharSequence(" ", TokenKind_Space)
        if err != nil {
            return nil, err
        } else if space != nil {
            _ = this.saveToken(space)
            continue
        }

        carriageReturn, err := this.tryGetCharSequence("\r", TokenKind_CarriageReturn)
        if err != nil {
            return nil, err
        } else if carriageReturn != nil {
            _ = this.saveToken(carriageReturn)
            continue
        }

        newLine, err := this.tryGetCharSequence("\n", TokenKind_NewLine)
        if err != nil {
            return nil, err
        } else if newLine != nil {
            _ = this.saveToken(newLine)
            continue
        }

        equals, err := this.tryGetCharSequence("=", TokenKind_Equals)
        if err != nil {
            return nil, err
        } else if equals != nil {
            _ = this.saveToken(equals)
            continue
        }

        at, err := this.tryGetCharSequence("@", TokenKind_AtSymbol)
        if err != nil {
            return nil, err
        } else if at != nil {
            _ = this.saveToken(at)
            continue
        }

        leftSquareBracket, err := this.tryGetCharSequence("[", TokenKind_LeftSquareBracket)
        if err != nil {
            return nil, err
        } else if leftSquareBracket != nil {
            _ = this.saveToken(leftSquareBracket)
            continue
        }

        rightSquareBracket, err := this.tryGetCharSequence("]", TokenKind_RightSquareBracket)
        if err != nil {
            return nil, err
        } else if rightSquareBracket != nil {
            _ = this.saveToken(rightSquareBracket)
            continue
        }

        leftCurlyBrace, err := this.tryGetCharSequence("{", TokenKind_LeftCurlyBrace)
        if err != nil {
            return nil, err
        } else if leftCurlyBrace != nil {
            _ = this.saveToken(leftCurlyBrace)
            continue
        }

        rightCurlyBrace, err := this.tryGetCharSequence("}", TokenKind_RightCurlyBrace)
        if err != nil {
            return nil, err
        } else if rightCurlyBrace != nil {
            _ = this.saveToken(rightCurlyBrace)
            continue
        }

        leftParenthesis, err := this.tryGetCharSequence("(", TokenKind_LeftParenthesis)
        if err != nil {
            return nil, err
        } else if leftParenthesis != nil {
            _ = this.saveToken(leftParenthesis)
            continue
        }

        rightParenthesis, err := this.tryGetCharSequence(")", TokenKind_RightParenthesis)
        if err != nil {
            return nil, err
        } else if rightParenthesis != nil {
            _ = this.saveToken(rightParenthesis)
            continue
        }

        leftAngleBracket, err := this.tryGetCharSequence("<", TokenKind_LeftAngleBracket)
        if err != nil {
            return nil, err
        } else if leftAngleBracket != nil {
            _ = this.saveToken(leftAngleBracket)
            continue
        }

        rightAngleBracket, err := this.tryGetCharSequence(">", TokenKind_RightAngleBracket)
        if err != nil {
            return nil, err
        } else if rightAngleBracket != nil {
            _ = this.saveToken(rightAngleBracket)
            continue
        }

        plus, err := this.tryGetCharSequence("+", TokenKind_Plus)
        if err != nil {
            return nil, err
        } else if plus != nil {
            _ = this.saveToken(plus)
            continue
        }

        minus, err := this.tryGetCharSequence("-", TokenKind_Minus)
        if err != nil {
            return nil, err
        } else if minus != nil {
            _ = this.saveToken(minus)
            continue
        }

        asterisk, err := this.tryGetCharSequence("*", TokenKind_Asterisk)
        if err != nil {
            return nil, err
        } else if asterisk != nil {
            _ = this.saveToken(asterisk)
            continue
        }

        forwardSlash, err := this.tryGetCharSequence("/", TokenKind_ForwardSlash)
        if err != nil {
            return nil, err
        } else if forwardSlash != nil {
            _ = this.saveToken(forwardSlash)
            continue
        }

        backwardSlash, err := this.tryGetCharSequence("\\", TokenKind_BackwardSlash)
        if err != nil {
            return nil, err
        } else if backwardSlash != nil {
            _ = this.saveToken(backwardSlash)
            continue
        }

        comma, err := this.tryGetCharSequence(",", TokenKind_Comma)
        if err != nil {
            return nil, err
        } else if comma != nil {
            _ = this.saveToken(comma)
            continue
        }

        dot, err := this.tryGetCharSequence(".", TokenKind_Dot)
        if err != nil {
            return nil, err
        } else if dot != nil {
            _ = this.saveToken(dot)
            continue
        }

        exclamation, err := this.tryGetCharSequence("!", TokenKind_Exclamation)
        if err != nil {
            return nil, err
        } else if exclamation != nil {
            _ = this.saveToken(exclamation)
            continue
        }

        colon, err := this.tryGetCharSequence(":", TokenKind_Colon)
        if err != nil {
            return nil, err
        } else if colon != nil {
            _ = this.saveToken(colon)
            continue
        }

        or, err := this.tryGetKeyword("or", TokenKind_Or)
        if err != nil {
            return nil, err
        } else if or != nil {
            _ = this.saveToken(or)
            continue
        }

        if_, err := this.tryGetKeyword("if", TokenKind_If)
        if err != nil {
            return nil, err
        } else if if_ != nil {
            _ = this.saveToken(if_)
            continue
        }

        and, err := this.tryGetKeyword("and", TokenKind_And)
        if err != nil {
            return nil, err
        } else if and != nil {
            _ = this.saveToken(and)
            continue
        }

        else_, err := this.tryGetKeyword("else", TokenKind_Else)
        if err != nil {
            return nil, err
        } else if else_ != nil {
            _ = this.saveToken(else_)
            continue
        }

        loop, err := this.tryGetKeyword("loop", TokenKind_Loop)
        if err != nil {
            return nil, err
        } else if loop != nil {
            _ = this.saveToken(loop)
            continue
        }

        break_, err := this.tryGetKeyword("break", TokenKind_Break)
        if err != nil {
            return nil, err
        } else if break_ != nil {
            _ = this.saveToken(break_)
            continue
        }

        bytes, err := this.tryGetKeyword("bytes", TokenKind_Bytes)
        if err != nil {
            return nil, err
        } else if bytes != nil {
            _ = this.saveToken(bytes)
            continue
        }

        script, err := this.tryGetKeyword("script", TokenKind_Script)
        if err != nil {
            return nil, err
        } else if script != nil {
            _ = this.saveToken(script)
            continue
        }

        random, err := this.tryGetKeyword("random", TokenKind_Random)
        if err != nil {
            return nil, err
        } else if random != nil {
            _ = this.saveToken(random)
            continue
        }

        return_, err := this.tryGetKeyword("return", TokenKind_Return)
        if err != nil {
            return nil, err
        } else if return_ != nil {
            _ = this.saveToken(return_)
            continue
        }

        switch_, err := this.tryGetKeyword("switch", TokenKind_Switch)
        if err != nil {
            return nil, err
        } else if switch_ != nil {
            _ = this.saveToken(switch_)
            continue
        }

        case_, err := this.tryGetKeyword("case", TokenKind_Case)
        if err != nil {
            return nil, err
        } else if case_ != nil {
            _ = this.saveToken(case_)
            continue
        }

        default_, err := this.tryGetKeyword("default", TokenKind_Default)
        if err != nil {
            return nil, err
        } else if default_ != nil {
            _ = this.saveToken(default_)
            continue
        }

        identifier, err := this.tryGetIdentifier()
        if err != nil {
            return nil, err
        } else if identifier != nil {
            _ = this.saveToken(identifier)
            continue
        }

        return this.tokens, errors.New("did not find token")
    }

    return this.tokens, nil
}

func (this *lexer) saveToken(token Token) error {
    kind := token.Kind()
    this.index += token.CharsConsumed()
    this.lineNumber += token.LinesConsumed()
    // exclude certain token types from the output for convenience
    if this.preventConsecutiveLineBreaks {
        if len(this.tokens) > 0 && this.tokens[len(this.tokens)-1].Kind() == TokenKind_NewLine && kind == TokenKind_NewLine {
            return nil
        }
    }
    if kind != TokenKind_SingleLineComment &&
        kind != TokenKind_MultiLineComment &&
        kind != TokenKind_CarriageReturn &&
        kind != TokenKind_Space &&
        kind != TokenKind_Tab &&
        kind != TokenKind_EscapedLineBreak {
        this.tokens = append(this.tokens, token)
    }
    return nil
}

func (this *lexer) isOutOfRangeAt(index uint) bool {
    return index >= uint(len(this.sourceCode))
}

func (this *lexer) isLetterAt(index uint) bool {
    return unicode.IsLetter(rune(this.sourceCode[index]))
}

func (this *lexer) isDigitAt(index uint) bool {
    return unicode.IsDigit(rune(this.sourceCode[index]))
}

func (this *lexer) isHexDigitAt(index uint) bool {
    switch this.sourceCode[index] {
    case '0':
        fallthrough
    case '1':
        fallthrough
    case '2':
        fallthrough
    case '3':
        fallthrough
    case '4':
        fallthrough
    case '5':
        fallthrough
    case '6':
        fallthrough
    case '7':
        fallthrough
    case '8':
        fallthrough
    case '9':
        fallthrough
    case 'a':
        fallthrough
    case 'b':
        fallthrough
    case 'c':
        fallthrough
    case 'd':
        fallthrough
    case 'e':
        fallthrough
    case 'f':
        fallthrough
    case 'A':
        fallthrough
    case 'B':
        fallthrough
    case 'C':
        fallthrough
    case 'D':
        fallthrough
    case 'E':
        fallthrough
    case 'F':
        return true
    default:
        return false
    }
}

func (this *lexer) tryGetEscapedLineBreak() (Token, error) {
    if !this.isOutOfRangeAt(this.index + 2) {
        threeChars := this.sourceCode[this.index : this.index+3]
        if threeChars == "\\\r\n" {
            return newGenericToken(TokenKind_EscapedLineBreak, "\\\r\n", this.lineNumber, 3, 1), nil
        }
    } else if !this.isOutOfRangeAt(this.index + 1) {
        twoChars := this.sourceCode[this.index : this.index+2]
        if twoChars == "\\\n" {
            return newGenericToken(TokenKind_EscapedLineBreak, "\\\n", this.lineNumber, 2, 1), nil
        }
    }
    return nil, nil
}

func (this *lexer) tryGetFloat() (Token, error) {
    startIndex := this.index
    endIndex := startIndex

    state := 0
    // 0 = waiting to scan first digit
    // 1 = first digit found, scanning digits but checking for '.'
    // 2 = '.' found, scanning for first digit after '.'
    // 3 = first digit after '.' found, scanning for more digits
    for {
        switch state {
        case 0:
            if this.isDigitAt(endIndex) {
                state = 1
            } else if this.sourceCode[endIndex] != '-' {
                return nil, nil
            }
        case 1:
            if this.isOutOfRangeAt(endIndex) {
                return nil, nil
            } else if this.sourceCode[endIndex] == '.' {
                state = 2
            } else if !this.isDigitAt(endIndex) {
                return nil, nil
            }
        case 2:
            if this.isOutOfRangeAt(endIndex) {
                return nil, nil
            } else if this.isDigitAt(endIndex) {
                state = 3
            } else {
                return nil, nil
            }
        case 3:
            if this.isOutOfRangeAt(endIndex) || !this.isDigitAt(endIndex) {
                return newGenericToken(TokenKind_Float, this.sourceCode[startIndex:endIndex], this.lineNumber, endIndex-startIndex, 0), nil
            }
        }
        endIndex++
    }
}

func (this *lexer) tryGetInt() (Token, error) {
    startIndex := this.index
    endIndex := startIndex
    size := uint(0)
    if this.sourceCode[endIndex] == '-' {
        endIndex++
        size++
    }
    for {
        if this.isOutOfRangeAt(endIndex) || !this.isDigitAt(endIndex) {
            break
        }
        endIndex++
    }

    if endIndex-startIndex == size {
        return nil, nil
    }

    return newGenericToken(TokenKind_Int, this.sourceCode[startIndex:endIndex], this.lineNumber, endIndex-startIndex, 0), nil
}

func (this *lexer) tryGetString() (Token, error) {
    startIndex := this.index
    endIndex := startIndex

    startLine := this.lineNumber
    endLine := startLine

    state := 0
    // 0 = looking for (")
    // 1 = found ("), scanning string, checking for next (")
    // 2 = found next (")
    for {
        switch state {
        case 0:
            if this.sourceCode[endIndex] != '"' {
                return nil, nil
            } else {
                state = 1
            }
        case 1:
            if this.isOutOfRangeAt(endIndex) {
                return nil, errors.New("EOF while scanning string literal")
            }

            if this.sourceCode[endIndex] == '\\' {
                // TODO lookahead for invalid escape sequence
                endIndex++
            } else if this.sourceCode[endIndex] == '\n' {
                endLine++
            } else if this.sourceCode[endIndex] == '"' {
                endIndex++
                string_ := this.sourceCode[startIndex+1 : endIndex-1]
                string_ = strings.ReplaceAll(string_, "\\\\", "\\")
                string_ = strings.ReplaceAll(string_, "\\n", "\n")
                string_ = strings.ReplaceAll(string_, "\\\"", "\"")
                return newGenericToken(TokenKind_String, string_, this.lineNumber, endIndex-startIndex, endLine-startLine), nil
            }
        }
        endIndex++
    }
}

func (this *lexer) tryGetSingleLineComment() (Token, error) {
    startIndex := this.index
    endIndex := startIndex

    if !this.tryGetCharSequenceAt("//", endIndex) {
        return nil, nil
    }
    endIndex += 2

    for {
        if this.isOutOfRangeAt(endIndex) {
            break
        } else if this.sourceCode[endIndex] == '\n' {
            endIndex++
            break
        }
        endIndex++
    }

    token := newGenericToken(TokenKind_SingleLineComment, this.sourceCode[startIndex:endIndex], this.lineNumber, endIndex-startIndex, 0)
    return token, nil
}

func (this *lexer) tryGetMultiLineComment() (Token, error) {
    startIndex := this.index
    endIndex := startIndex

    startLine := this.lineNumber
    endLine := startLine

    if !this.tryGetCharSequenceAt("/*", endIndex) {
        return nil, nil
    }

    endIndex += 2
    nesting := 1
    for {
        if this.isOutOfRangeAt(endIndex) {
            break
        }
        if this.sourceCode[endIndex] == '\n' {
            endLine++
            endIndex++
        } else if this.tryGetCharSequenceAt("/*", endIndex) {
            nesting++
            endIndex += 2
        } else if this.tryGetCharSequenceAt("*/", endIndex) {
            nesting--
            endIndex += 2
            if nesting <= 0 {
                break
            }
        } else {
            endIndex++
        }
    }

    return newGenericToken(TokenKind_MultiLineComment, this.sourceCode[startIndex:endIndex], this.lineNumber, endIndex-startIndex, endLine-startLine), nil
}

func (this *lexer) tryGetRawQbKey() (Token, error) {
    startIndex := this.index
    endIndex := startIndex

    if this.isOutOfRangeAt(startIndex) {
        return nil, nil
    }

    if this.sourceCode[endIndex] != '#' {
        return nil, nil
    }
    endIndex++

    for i := 0; i < 8; i++ {
        if this.isOutOfRangeAt(endIndex) {
            return nil, nil
        }

        if !this.isHexDigitAt(endIndex) {
            return nil, nil
        }

        endIndex++
    }

    return newGenericToken(TokenKind_RawQbKey, this.sourceCode[startIndex+1:endIndex], this.lineNumber, endIndex-startIndex, 0), nil
}

func (this *lexer) tryGetKeyword(keyword string, tokenKind TokenKind) (Token, error) {
    token, err := this.tryGetCharSequence(keyword, tokenKind)
    if err != nil {
        return nil, err
    } else if token == nil {
        return nil, nil
    }
    indexAfterToken := this.index + token.CharsConsumed()
    if this.isOutOfRangeAt(indexAfterToken) {
        return token, nil
    }
    if this.isLetterAt(indexAfterToken) || this.isDigitAt(indexAfterToken) || this.sourceCode[indexAfterToken] == '_' {
        return nil, nil
    }
    return token, nil
}

func (this *lexer) tryGetCharSequence(sequence string, tokenKind TokenKind) (Token, error) {
    if this.tryGetCharSequenceAt(sequence, this.index) {
        return newGenericToken(tokenKind, sequence, this.lineNumber, uint(len(sequence)), countNewLines(sequence)), nil
    }
    return nil, nil
}

func countNewLines(text string) uint {
    newLines := uint(0)
    for _, c := range text {
        if c == '\n' {
            newLines++
        }
    }
    return newLines
}

func (this *lexer) tryGetCharSequenceAt(sequence string, index uint) bool {
    if this.isOutOfRangeAt(index + uint(len(sequence)-1)) {
        return false
    }
    for i, char := range sequence {
        if this.isOutOfRangeAt(index + uint(i)) {
            return false
        }
        if this.sourceCode[index+uint(i)] != uint8(char) {
            return false
        }
    }
    return true
}

func (this *lexer) tryGetIdentifier() (Token, error) {
    startIndex := this.index
    endIndex := startIndex

    startLine := this.lineNumber
    endLine := startLine

    if this.isOutOfRangeAt(endIndex) {
        return nil, nil
    }

    if this.sourceCode[endIndex] == '`' {
        endIndex++
        for {
            if this.isOutOfRangeAt(endIndex) {
                return nil, errors.New("eof while scanning `identifier")
            } else if this.sourceCode[endIndex] == '`' {
                endIndex++
                return newGenericToken(TokenKind_Identifier, this.sourceCode[startIndex+1:endIndex-1], this.lineNumber, endIndex-startIndex, endLine-startLine), nil
            } else if this.sourceCode[endIndex] == '\n' {
                endLine++
            }
            endIndex++
        }
    } else {
        for {
            if this.isOutOfRangeAt(endIndex) {
                break
            }
            if !(this.isLetterAt(endIndex) || this.isDigitAt(endIndex) || this.sourceCode[endIndex] == '_') {
                break
            }
            endIndex++
        }
        if endIndex-startIndex == 0 {
            return nil, nil
        }
        return newGenericToken(TokenKind_Identifier, this.sourceCode[startIndex:endIndex], this.lineNumber, endIndex-startIndex, 0), nil
    }
}
