local nnc = import "./nnc.libsonnet";

function(ctx, spec)
	nnc.merge([spec, {
		mounts: [
			nnc.mountHostRO("/etc", "/etc"),
		],
		net: [
      nnc.netNone("test123")
    ],
}])

