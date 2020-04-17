package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
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
	var binaryReadPath string
	if strings.ContainsRune(binaryPath, os.PathSeparator) {
		binaryReadPath = binaryPath
	} else {
		p, err := findBinaryPath(binaryPath)
		if err != nil {
			return fmt.Errorf("find binary path: %s", err)
		}
		if p == "" {
			return errors.New("binary not found in path")
		}
		binaryReadPath = p
	}
	f, err := os.Open(binaryReadPath)
	if err != nil {
		return fmt.Errorf("read binary path: %s", err)
	}
	defer f.Close()
	magicReader := magic.NewReader(f)
	b, err := magicReader.Read()
	if err != nil {
		return fmt.Errorf("magic: %s", err)
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

func findBinaryPath(binaryName string) (string, error) {
	envPath := os.Getenv("PATH")
	if envPath == "" {
		return "", errors.New("PATH env not set")
	}
	paths := strings.Split(envPath, string(os.PathListSeparator))
	for _, p := range paths {
		binaryPath := path.Join(p, binaryName)
		_, err := os.Stat(binaryPath)
		if err == nil {
			return binaryPath, nil
		}
	}
	return "", nil
}
