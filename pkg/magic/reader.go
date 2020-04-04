package magic

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// Reader defines how to read a payload from a magic encoded source.
type Reader interface {
	// Read unpacks the magic encoding and returns the underlying bytes.
	Read() ([]byte, error)
}

// NewReader returns a new magic.Reader
func NewReader(readSeeker io.ReadSeeker) Reader {
	return &reader{readSeeker}
}

type reader struct {
	readSeeker io.ReadSeeker
}

var ErrMagicBytesNotFound = errors.New("magic bytes not found")

func (r *reader) Read() ([]byte, error) {
	n := len(Bytes) + 4
	if _, err := r.readSeeker.Seek(int64(-n), os.SEEK_END); err != nil {
		return nil, fmt.Errorf("seek: %s", err)
	}
	b := make([]byte, n)
	if nRead, err := io.ReadFull(r.readSeeker, b); nRead != n {
		return nil, fmt.Errorf("read magic: %s", err)
	}
	magicBytes := b[4:]
	sizeBytes := b[:4]
	if string(magicBytes) != Bytes {
		return nil, ErrMagicBytesNotFound
	}
	var size uint32
	if err := binary.Read(bytes.NewReader(sizeBytes), binary.BigEndian, &size); err != nil {
		return nil, fmt.Errorf("read size bytes: %s", err)
	}
	if _, err := r.readSeeker.Seek(-int64(uint32(n)+size), os.SEEK_END); err != nil {
		return nil, fmt.Errorf("seek: %s", err)
	}
	b = make([]byte, size)
	if nRead, err := io.ReadFull(r.readSeeker, b); uint32(nRead) != size {
		return nil, fmt.Errorf("read bytes: %s", err)
	}
	return b, nil
}
