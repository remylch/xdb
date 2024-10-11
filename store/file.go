package store

import "io"

type File interface {
	Content() io.Reader
}

type XDBIndexFile io.Reader

type XDBFile io.Reader
