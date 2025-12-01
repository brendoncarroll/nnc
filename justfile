
test: build
    go test -v -count=1 ./...

build:
    CGO_ENABLED=0 go build -o ./src/nnccmd/out/nnc_shim ./src/nnc/nnc_shim
    CGO_ENABLED=0 go build -o ./src/nnccmd/out/testbin ./src/internal/testbin
    mkdir -p ./build/out
    CGO_ENABLED=0 go build -o ./build/out/nnc ./cmd/nnc

install: build
   cp ./build/out/nnc $HOME/bin/nnc
    
fish: build
    ./etc/fish.sh

