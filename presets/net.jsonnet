local nnc = import "./nnc.libsonnet";

function(spec, caller)
	spec +
	{
		net: [
      nnc.netNone("test123")
    ],
    data: [
      nnc.dataLit("/etc/resolv.conf", "nameserver 8.8.8.8"),
    ]
  }
