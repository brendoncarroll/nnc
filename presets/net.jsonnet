local nnc = import "./nnc.libsonnet";

function(spec, caller)
	nnc.merge([spec, {
		mounts: [
			nnc.mountHostRO("/etc", "/etc"),
		],
		net: [
      nnc.netNone("test123")
    ],
}])

