local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    mounts: nnc.mountsMerge([spec.mounts, [
      nnc.mountHostRW("/dev/fuse", "/dev/fuse"),
    ]]),
  }
