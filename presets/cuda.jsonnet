local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    mounts: nnc.mountsMerge([spec.mounts, [
      nnc.mountSymlinkFD("/dev/nvidia0"),
      nnc.mountSymlinkFD("/dev/nvidiactl"),
      nnc.mountSymlinkFD("/dev/nvidia-uvm"),
      nnc.mountSymlinkFD("/dev/nvidia-uvm-tools"),
    ]]),
  }
