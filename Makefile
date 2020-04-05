all: bin/goskel bin/append-template-binary

bin/goskel: pkg/magic/* cmd/goskel/main.go
	go build -o bin/goskel ./cmd/goskel/.

bin/append-template-binary: pkg/magic/* cmd/append-template-binary/main.go
	go build -o bin/append-template-binary ./cmd/append-template-binary/.

append: bin/goskel bin/append-template-binary
	bin/append-template-binary --binary=bin/goskel --template=var/template

clean:
	rm -r bin/goskel bin/append-template-binary

.PHONY: append clean
