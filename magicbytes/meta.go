package magicbytes

import "fmt"

// Meta holds the name, magical bytes, and offset of the magical bytes to be searched.
type Meta struct {
	Type   string // name of the file/meta type.
	Bytes  []byte // magical bytes.
	Offset int64  // offset of the magical bytes from the file start position.
}

// OnMatchFunc represents a function to be called when Search function finds a match.
// Returning false must immediately stop Search process.
type OnMatchFunc func(path, metaType string) bool

type ArgumentError struct{ Message string }

func (e *ArgumentError) Error() string {
	return fmt.Sprintf(e.Message)
}

type ContextCancelledError struct{}

func (e *ContextCancelledError) Error() string {
	return fmt.Sprintf("Context cancelled.")
}
