local nnc = import "./nnc.libsonnet";
local etcssl = import "./etc-ssl.jsonnet";

function(ctx, spec)
  local spec2 = nnc.applyAll(ctx, spec, [
    etcssl,
  ]);
  spec2 + {
    mounts: nnc.mountsMerge([spec2.mounts,
      [
        // TODO: use an overlay for this to prevent malicious cache corruption.
        nnc.mountHostRW("/root/.cache/go", nnc.homePath(ctx, ".cache/go")),
        nnc.mountHostRW("/root/.cache/go-build", nnc.homePath(ctx, ".cache/go-build")),
      ]
    ]),
  }
