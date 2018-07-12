package nametable

type NameTable struct {
	entries internalTable
}

type internalTable map[string]string
