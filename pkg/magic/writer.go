package magic

import (
	"encoding/binary"
	"io"
)

// NewWriter returns a new io.Writer that will write the payload with the magic
// bytes.
//
// It's expected that the write function will just be called once.
func NewWriter(w io.Writer) io.Writer {
	return &writer{w}
}

type writer struct {
	w io.Writer
}

func (w *writer) Write(b []byte) (n int, err error) {
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(len(b)))
	all := append(append(b, sizeBytes...), Bytes...)
	return w.w.Write(all)
}
