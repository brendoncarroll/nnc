local nnc = import "./nnc.libsonnet";

function(ctx, spec)
	nnc.merge([spec, {
		mounts: [
    	nnc.mountTmpfs("/dev"),
			nnc.mountDev("null"),
    	nnc.mountHostRW("/dev/urandom", "/dev/urandom"),
      nnc.mountHostRW("/dev/random", "/dev/random"),
		],
}])
