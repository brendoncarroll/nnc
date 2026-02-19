local nnc = import "./nnc.libsonnet";

function (ctx, spec)
  spec + {
    mounts: nnc.mountsMerge([spec.mounts, [
      nnc.mountHostRO("/usr", "/usr"),
      nnc.mountHostRO("/lib64", "/lib64"),
      nnc.mountTmpfs("/dev"),
    ]]),
  }
