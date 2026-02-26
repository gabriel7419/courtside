.PHONY: build install run mock clean test

# Build the Courtside binary
build:
	go build -o courtside ./cmd/...

# Install the Courtside binary to $GOPATH/bin
install:
	go install ./cmd/...

# Run Courtside with real NBA API data
run:
	go run ./cmd/...

# Run Courtside with mock data (useful for testing when there are no live games)
mock:
	go run ./cmd/... --mock

# Clean the build output
clean:
	rm -f courtside

# Run tests
test:
	go test ./...
