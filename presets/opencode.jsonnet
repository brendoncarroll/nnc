local nnc = import "./nnc.libsonnet";
local netPreset = import "./net.jsonnet";
local minDev = import "./min-dev.jsonnet";
local etcPasswd = import "./etc-passwd.jsonnet";

function(ctx, spec)
  local applied = nnc.applyAll(ctx, spec, [
    netPreset,
    minDev,
    etcPasswd,
  ]);
  applied + {
    mounts: nnc.mountsMerge([applied.mounts, [
        nnc.mountHostRO("/bin", "/bin"),
        nnc.mountHostRO("/sbin", "/sbin"),
        nnc.mountHostRW("/lib", "/lib"),
        nnc.mountHostRW("/lib64", "/lib64"),
        nnc.mountHostRO("/usr", "/usr"),
        nnc.mountTmpfs("/tmp"),
        nnc.mountProcfs(),

        // users git config
        nnc.mountHostRO("/root/.config/git", nnc.homePath(ctx, ".config/git")),

        nnc.mountHostRW("/root/.local/share/opencode", nnc.homePath(ctx, ".local/share/opencode")),
        nnc.mountHostRW("/root/.local/state/opencode", nnc.homePath(ctx, ".local/state/opencode")),
        nnc.mountHostRW("/root/.config/opencode", nnc.homePath(ctx, ".config/opencode")),

        nnc.mountHostRW("/_", ctx.wd),
    ]]),
    env: nnc.envMerge([
      applied.env,
      [
        "HOME=/root",
        "PATH=/bin:/usr/bin:/sbin:/usr/local/bin",
      ],
      nnc.selectEnvKeys(ctx, [
        "TERM",
      ])
    ]),
    wd: "/_",
  }
