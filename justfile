

build:
    @echo "Building justfile"
    @go build

test:
    @echo "Running tests"
    @go test ./pkg/cps
    @go test ./pkg/env
    @go test ./pkg/ospaths
    @go test ./pkg/platform
    @go test ./pkg/vaults/sops
    @go test ./pkg/xexec
    @go test ./pkg/xfs
    @go test ./pkg/xrunes
    @go test ./pkg/xstrings