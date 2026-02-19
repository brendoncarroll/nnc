local nnc = import "./nnc.libsonnet";
local netPreset = import "./net.jsonnet";
local minDev = import "./min-dev.jsonnet";

function(ctx, spec)
  nnc.merge([
    spec,
    netPreset(ctx, spec),
    minDev(ctx, spec),
    {
        mounts: [
            nnc.mountHostRO("/bin", "/bin"),
            nnc.mountHostRO("/sbin", "/sbin"),
            nnc.mountHostRW("/lib", "/lib"),
            nnc.mountHostRW("/lib64", "/lib64"),
            nnc.mountHostRO("/usr", "/usr"),
            nnc.mountTmpfs("/tmp"),
            nnc.mountProcfs(),

            // users git config
            nnc.mountHostRO("/root/.config/git", nnc.homePath(ctx, ".config/git")),

            nnc.mountTmpfs("/root/.amp"),
            nnc.mountHostRW("/root/.local/share/amp", nnc.homePath(ctx, ".local/share/amp")),
            nnc.mountHostRW("/root/.config/amp", nnc.homePath(ctx, ".config/amp")),

            nnc.mountHostRW("/_", ctx.wd),
        ],
        env: nnc.envMerge([
          [
            "HOME=/root",
            "PATH=/root/.amp/bin:/bin:/usr/bin:/sbin:/usr/local/bin",
          ],
          nnc.selectEnvKeys(ctx, [
            "TERM",
          ])
        ]),
        wd: "/_",
    }
  ])
