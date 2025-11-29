#!/bin/sh

go run ./cmd/nnc run \
  --dr /usr:/usr \
  --dr /bin:/bin \
  --dr /lib64:/lib64 \
  --dr /root/.config/fish:$HOME/.config/fish \
  --env HOME=/root \
  --dr "/x:$(pwd)" \
  /bin/fish
