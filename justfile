
build:
    CGO_ENABLED=0 go build -o ./src/nnccmd/out/nnc_shim ./src/nnc/nnc_shim
    CGO_ENABLED=0 go build -o ./src/nnccmd/out/testbin ./src/internal/testbin

test: build
    go test -v -count=1 ./...

fish: build
    ./etc/fish.sh
