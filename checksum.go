package NeverScript

type Checksum struct {
	content uint32
}

func NewChecksum(content uint32) Checksum {
	return Checksum{
		content: content,
	}
}

func NewEmptyChecksum() Checksum {
	return NewChecksum(0)
}

func (this Checksum) ToUint32() uint32 {
	return this.content
}

func (this Checksum) IsEqualTo(other Checksum) bool {
	return this.content == other.content
}
