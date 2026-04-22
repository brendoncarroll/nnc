local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    mounts: nnc.mountsMerge([spec.mounts, [
      nnc.mountHostRW("/dev/nvidia0", "/dev/nvidia0"),
      nnc.mountHostRW("/dev/nvidiactl", "/dev/nvidiactl"),
      nnc.mountHostRW("/dev/nvidia-uvm", "/dev/nvidia-uvm"),
      nnc.mountHostRW("/dev/nvidia-uvm-tools", "/dev/nvidia-uvm-tools"),
    ]]),
  }
