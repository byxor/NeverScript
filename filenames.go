package NeverScript

func QbToNs(filename string) string {
	return withoutExtension(filename) + ".ns"
}

func NsToQb(filename string) string {
	return withoutExtension(filename) + ".qb"
}

func withoutExtension(filename string) string {
	return filename[:len(filename)-3]
}
