local nnc = import "./nnc.libsonnet";

function(spec, caller)
	nnc.merge2(spec, {
		mounts: [
			nnc.mountHostRO("/etc/ssl", "/etc/ssl"),
		],
		net: [
      nnc.netNone("test123")
    ],
    data: [
      nnc.dataLit("/etc/resolv.conf", "nameserver 1.1.1.1"),
    ]
})

