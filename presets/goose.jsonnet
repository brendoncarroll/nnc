local nnc = import "./nnc.libsonnet";
local netPreset = import "./net.jsonnet";

function(spec, caller)
  nnc.merge([
    spec,
    netPreset(spec, caller),
    
    {
        mounts: nnc.mountsMerge([
          spec.mounts,
          [
            nnc.mountHostRO("/bin", "/bin"),
            nnc.mountHostRO("/sbin", "/sbin"),
      			nnc.mountHostRW("/lib", "/lib"),
      			nnc.mountHostRW("/lib64", "/lib64"),
      			nnc.mountHostRO("/usr", "/usr"),
      			nnc.mountTmpfs("/tmp"),

            nnc.mountHostRO("/root/.config/goose", nnc.homePath(caller, ".config/goose")),
      			nnc.mountHostRW("/_", caller.wd),
          ],
        ]),
        env: nnc.envMerge([
          spec.env,
          [
            "HOME=/root",
            "PATH=/bin:/usr/bin:/sbin:/usr/local/bin",
            "GOOSE_DISABLE_KEYRING=yes",
          ],
          nnc.envSelectKeys(caller, [
            "OPENROUTER_API_KEY",
            "TERM",
          ])
        ]),
    		wd: "/_",
    }    
  ])
  
