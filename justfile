
build:
    go build -o ./src/nnccmd/out/nnc_shim ./src/nnc/nnc_shim
    go build -o ./src/nnccmd/out/testbin ./src/internal/testbin

test: build
    go test -v -count=1 ./...

