

j9d *ARGS:
    @echo "Running j9d"
    @echo "ARGS: {{ARGS}}"
    @bin/j9d {{ARGS}}

build:
    @echo "Building justfile"
    @go build -o bin/j9d ./main.go

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

update-bin:
    cp -f bin/j9d ~/.local/bin 
