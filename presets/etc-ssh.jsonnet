local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    data: nnc.mountsMerge([spec.data, [
      nnc.copyHostPath("/etc/ssh", "/etc/ssh", 420),
    ]]),
  }
