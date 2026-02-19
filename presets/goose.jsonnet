local nnc = import "./nnc.libsonnet";
local netPreset = import "./net.jsonnet";
local waylandPreset = import "./wayland.jsonnet";

function(ctx, spec)
  nnc.applyAll(ctx, spec, [
    netPreset,
    waylandPreset,
    function(ctx, spec) spec + {
        mounts: nnc.mountsMerge([spec.mounts, [
            nnc.mountHostRO("/bin", "/bin"),
            nnc.mountHostRO("/sbin", "/sbin"),
      			nnc.mountHostRW("/lib", "/lib"),
      			nnc.mountHostRW("/lib64", "/lib64"),
      			nnc.mountHostRO("/usr", "/usr"),
      			nnc.mountTmpfs("/tmp"),

            nnc.mountHostRO("/root/.config/goose", nnc.homePath(ctx, ".config/goose")),
      			nnc.mountHostRW("/_", ctx.wd),
        ]]),
        env: nnc.envMerge([
          spec.env,
          [
            "HOME=/root",
            "PATH=/bin:/usr/bin:/sbin:/usr/local/bin",
            "GOOSE_DISABLE_KEYRING=yes",
          ],
          nnc.selectEnvKeys(ctx, [
            "OPENROUTER_API_KEY",
            "TERM",
          ])
        ]),
    		wd: "/_",
    },
  ])
