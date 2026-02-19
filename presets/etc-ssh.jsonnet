local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    data: nnc.mountsMerge([spec.data, [
      nnc.copyHostPath("/etc/ssh/ssh_config", "/etc/ssh/ssh_config", 420),
      nnc.copyHostPath("/etc/ssh/ssh_config.d", "/etc/ssh/ssh_config.d", 420),
      nnc.copyHostPath("/root/.ssh/known_hosts", nnc.homePath(ctx, ".ssh/known_hosts"), 420)
    ]]),
  }
