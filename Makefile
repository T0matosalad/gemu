.PHONY: all
all: clean
	mkdir dist
	go build -o dist/gemu ./cmd/gemu/main.go

.PHONY: clean
clean:
	rm -rf dist

.PHONY: test
test:
	go test -v ./...