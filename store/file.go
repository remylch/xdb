package store

import "io"

type File interface {
	Content() io.Reader
}

type XDBIndexFile struct {
}

type XDBFile struct {
}
