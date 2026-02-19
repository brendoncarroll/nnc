local nnc = import "./nnc.libsonnet";

function(ctx, spec)
	spec + {
		mounts: nnc.mountsMerge([spec.mounts, [
			nnc.mountHostRO("/bin", "/bin"),
			nnc.mountHostRO("/lib64", "/lib64"),
			nnc.mountHostRO("/usr", "/usr"),
			nnc.mountHostRO("/root/.config/fish", nnc.homePath(ctx, ".config/fish")),

			nnc.mountTmpfs("/dev"),
			nnc.mountProcfs(),

			nnc.mountHostRW("/_", ctx.wd),
		]]),
		env: [
			"HOME=/root",
			"TERM=xterm-256color",
		],
		wd: "/_",
	}
