

build:
    @echo "Building justfile"
    @go build

test:
    @echo "Running tests"
    @go test ./...