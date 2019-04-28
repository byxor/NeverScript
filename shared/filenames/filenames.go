package filenames

func QbToNs(filename string) string {
	filenameWithoutExtension := filename[:len(filename)-3]
	return filenameWithoutExtension + ".ns"
}
