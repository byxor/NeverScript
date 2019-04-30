package NeverScript

type SourceCode struct {
	content string
}

func NewSourceCode(content string) SourceCode {
	return SourceCode{
		content: content,
	}
}
