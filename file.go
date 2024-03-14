package merkledag

const (
	FILE = iota
	DIR
)

type Node interface {
	Size() uint64
	Type() int
}

type File interface {
	Node

	Bytes() []byte
	Name() string
}

type Dir interface {
	Node

	It() DirIterator
	Name() string
}

type DirIterator interface {
	Next() bool

	Node() Node
}
