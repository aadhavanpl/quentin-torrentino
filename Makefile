.PHONY: run build clean format

BINARY = quentin-torrentino

run:
	go run .

build:
	go build -o $(BINARY)

clean:
	rm -f $(BINARY)

format:
	gofmt -s -w .