local nnc = import "./nnc.libsonnet";
local netPreset = import "./net.jsonnet";

function(spec, caller)
  local spec2 = netPreset(spec, caller);
  spec2 + {
    mounts: nnc.mountsMerge([
      spec.mounts,
      [
        nnc.mountHostRO("/bin", "/bin"),
  			nnc.mountHostRO("/lib64", "/lib64"),
  			nnc.mountHostRO("/usr", "/usr"),

        nnc.mountHostRO("/root/.config/goose", nnc.homePath(caller, ".config/goose")),
      ],
    ]),
    env: nnc.envMerge([
      spec.env,
      [
        "HOME=/root",
        "GOOSE_DISABLE_KEYRING=yes",
      ],
      nnc.envSelectKeys(caller, ["OPENROUTER_API_KEY"])
    ]),
}
