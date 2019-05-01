package NeverScript

type Checksum struct {
	content uint32
}

func NewChecksum(content uint32) Checksum {
	return Checksum{
		content: content,
	}
}
