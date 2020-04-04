package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/rmarianski/goskel/pkg/magic"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Expecting single name")
	}
	if err := run(os.Args[0], os.Args[1]); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func run(binaryPath, name string) error {
	f, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("read binary path: %s", err)
	}
	defer f.Close()
	n := len(magic.Bytes) + 4
	_, err = f.Seek(int64(-n), os.SEEK_END)
	if err != nil {
		return fmt.Errorf("read binary seek: %s", err)
	}
	b := make([]byte, n)
	if nRead, err := io.ReadFull(f, b); nRead != n {
		return fmt.Errorf("read binary trailer %s", err)
	}
	magicBytes := b[4:]
	sizeBytes := b[:4]
	if string(magicBytes) != magic.Bytes {
		return errors.New("magic bytes not found")
	}
	var size uint32
	if err := binary.Read(bytes.NewReader(sizeBytes), binary.BigEndian, &size); err != nil {
		return fmt.Errorf("read size bytes: %s", err)
	}
	_, err = f.Seek(-int64(uint32(n)+size), os.SEEK_END)
	if err != nil {
		return fmt.Errorf("read binary seek: %s", err)
	}
	b = make([]byte, size)
	if nRead, err := io.ReadFull(f, b); uint32(nRead) != size {
		return fmt.Errorf("read binary template: %s", err)
	}
	t, err := template.New("template").Parse(string(b))
	if err != nil {
		return fmt.Errorf("parse template: %s", err)
	}
	if err := os.Mkdir(name, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir %s: %s", name, err)
	}
	if err := os.Mkdir(path.Join(name, "cmd"), os.ModePerm); err != nil {
		return fmt.Errorf("mkdir cmd: %s", err)
	}
	if err := os.Mkdir(path.Join(name, "pkg"), os.ModePerm); err != nil {
		return fmt.Errorf("mkdir pkg: %s", err)
	}
	if err := os.Mkdir(path.Join(name, "cmd", name), os.ModePerm); err != nil {
		return fmt.Errorf("mkdir cmd/%s: %s", name, err)
	}
	mainPath := path.Join(name, "cmd", name, "main.go")
	mainFile, err := os.Create(mainPath)
	if err != nil {
		return fmt.Errorf("create main: %s", err)
	}
	if err := t.Execute(mainFile, struct{}{}); err != nil {
		_ = mainFile.Close()
		return fmt.Errorf("template execute: %s", err)
	}
	if err := mainFile.Close(); err != nil {
		return fmt.Errorf("close main.go: %s", err)
	}
	return nil
}
