BINARY = mdmu

.PHONY: build test test-race vet clean run

build:
	go build -o $(BINARY) .

test:
	go test ./...

test-race:
	go test -v -race ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)

run: build
	./$(BINARY) testfile.md
