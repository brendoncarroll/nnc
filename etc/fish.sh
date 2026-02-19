#!/bin/sh

go run ./cmd/nnc run \
  --mro /usr:/usr \
  --mro /bin:/bin \
  --mro /lib64:/lib64 \
  --mro /root/.config/fish:$HOME/.config/fish \
  --env HOME=/root \
  --mro "/x:$(pwd)" \
  /bin/fish
