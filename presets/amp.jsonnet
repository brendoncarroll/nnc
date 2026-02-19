local nnc = import "./nnc.libsonnet";
local netPreset = import "./net.jsonnet";

function(ctx, spec)
  nnc.merge([
    spec,
    netPreset(ctx, spec),
    {
        mounts: [
            nnc.mountHostRO("/bin", "/bin"),
            nnc.mountHostRO("/sbin", "/sbin"),
            nnc.mountHostRW("/lib", "/lib"),
            nnc.mountHostRW("/lib64", "/lib64"),
            nnc.mountHostRO("/usr", "/usr"),
            nnc.mountTmpfs("/tmp"),
            nnc.mountProcfs(),
            nnc.mountTmpfs("/dev"),
            nnc.mountDev("null"),
            nnc.mountDev("urandom"),
            nnc.mountDev("random"),

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
