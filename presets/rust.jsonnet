local nnc = import "./nnc.libsonnet";
local etcssl = import "./etc-ssl.jsonnet";

function(ctx, spec)
  local spec2 = nnc.applyAll(ctx, spec, [
    etcssl,
  ]);
  local forwardedEnv = nnc.selectEnvKeys(ctx, [
    "RUSTFLAGS",
    "CARGO_BUILD_TARGET",
    "CC",
    "CXX",
    "AR",
    "PKG_CONFIG_PATH",
  ], true);
  spec2 + {
    mounts: nnc.mountsMerge([spec2.mounts,
      [
        nnc.mountHostRW("/root/.rustup", nnc.homePath(ctx, ".rustup")),
        nnc.mountHostRW("/root/.cargo", nnc.homePath(ctx, ".cargo")),
        nnc.mountHostRW("/root/.cache/sccache", nnc.homePath(ctx, ".cache/sccache")),
      ]
    ]),
    env: nnc.envMerge([
      spec2.env,
      [
        "RUSTUP_HOME=/root/.rustup",
        "CARGO_HOME=/root/.cargo",
        "SCCACHE_DIR=/root/.cache/sccache",
        "RUSTC_WRAPPER=sccache",
      ],
      forwardedEnv,
    ]),
  }
