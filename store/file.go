package store

type File interface {
	Content() []string
}

type XDBIndexFile struct {
}

type XDBFile struct {
}
