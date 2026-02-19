local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  nnc.merge([spec, {
    data: [
      nnc.copyHostPath("/etc/ssh", "/etc/ssh", 420),
    ],
  }])
