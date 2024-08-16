package autho

// New creates new auth service
func New() *File {
	return &File{}
}

// File represents auth application service
type File struct{}
