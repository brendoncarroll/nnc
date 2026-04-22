local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    mounts: nnc.mountsMerge([spec.mounts, [
      nnc.mountDev("nvidia0"),
      nnc.mountDev("nvidiactl"),
      nnc.mountDev("nvidia-uvm"),
      nnc.mountDev("nvidia-uvm-tools"),
    ]]),
  }
