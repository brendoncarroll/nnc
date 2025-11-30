local nnc = import "./nnc.libsonnet";

function(spec, caller)

  spec + {
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
    ])
  }
