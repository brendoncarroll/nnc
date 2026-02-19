local nnc = import "./nnc.libsonnet";

function(ctx, spec)
	spec + {
		data: nnc.mountsMerge([spec.data, [
			nnc.copyHostPath("/etc/resolv.conf", "/etc/resolv.conf"),
			nnc.copyHostPath("/etc/resolvconf.conf", "/etc/resolvconf.conf"),
		]]),
		net: nnc.mountsMerge([spec.net, [
      nnc.netNone("test123")
    ]]),
	}
