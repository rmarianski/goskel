package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rmarianski/goskel/pkg/magic"
)

type flags struct {
	binary   string
	template string
}

func main() {
	var flags flags
	flag.StringVar(&flags.binary, "binary", "", "binary path")
	flag.StringVar(&flags.template, "template", "", "template path")
	flag.Parse()
	if err := run(&flags); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func run(flags *flags) (err error) {
	if flags.binary == "" {
		return errors.New("missing --binary")
	}
	if flags.template == "" {
		return errors.New("missing --template")
	}
	templateBytes, err := ioutil.ReadFile(flags.template)
	if err != nil {
		return fmt.Errorf("read template path=%s err=%s", flags.template, err)
	}
	f, err := os.OpenFile(flags.binary, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open binary append: %s", err)
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()
	magicWriter := magic.NewWriter(f)
	if _, err := magicWriter.Write(templateBytes); err != nil {
		return fmt.Errorf("magic: %s", err)
	}
	return err
}
