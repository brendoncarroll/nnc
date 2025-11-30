local nnc = import "./nnc.libsonnet";

function (spec)
  spec +
  {
    mounts: [
      nnc.mountHostRO("/usr", "/usr"),
      nnc.mountHostRO("/lib64", "/lib64"),
      nnc.mountTmpfs("/dev"),
    ]
  }
