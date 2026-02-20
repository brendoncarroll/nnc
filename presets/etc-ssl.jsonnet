local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    mounts: nnc.mountsMerge([spec.mounts,
      [
        nnc.mountHostRO("/etc/ssl", "/etc/ssl"),
      ],
    ]),
  }
