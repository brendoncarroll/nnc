
test: build
    go test -v -count=1 ./...

build:
    CGO_ENABLED=0 go build -o ./src/nnccmd/out/nnc_shim ./src/nnc/nnc_shim
    CGO_ENABLED=0 go build -o ./src/nnccmd/out/testbin ./src/internal/testbin

install:
    CGO_ENABLED=0 go install ./cmd/nnc

fish: build
    ./etc/fish.sh

