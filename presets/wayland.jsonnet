local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    mounts: nnc.mountsMerge([spec.mounts, [
      nnc.mountHostRW("/dev/dri", "/dev/dri"),
      nnc.mountHostRW("/run/user/0/wayland-1", "/run/user/%d/wayland-1" % ctx.uid),
    ]]),
    env: nnc.envMerge([
      spec.env,
      [
        "XDG_RUNTIME_DIR=/run/user/0",
        "WAYLAND_DISPLAY=wayland-1",
      ],
    ]),
  }
