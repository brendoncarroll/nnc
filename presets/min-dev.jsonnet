local nnc = import "./nnc.libsonnet";

function(ctx, spec)
	spec + {
		mounts: nnc.mountsMerge([spec.mounts, [
    	nnc.mountTmpfs("/dev"),
    	nnc.mountHostRW("/dev/null", "/dev/null"),
    	nnc.mountHostRW("/dev/urandom", "/dev/urandom"),
      nnc.mountHostRW("/dev/random", "/dev/random"),
		]]),
	}
