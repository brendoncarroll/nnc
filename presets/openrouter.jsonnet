local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    env: nnc.envMerge([
      spec.env,
      nnc.selectEnvKeys(ctx, [
        "OPENROUTER_API_KEY",
      ]),
    ]),
  }
