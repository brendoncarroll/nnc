local nnc = import "./nnc.libsonnet";

function(spec, caller)
	spec +
	{
		mounts: [
			nnc.mountHostRO("/bin", "/bin"),
			nnc.mountHostRO("/lib64", "/lib64"),
			nnc.mountHostRO("/usr", "/usr"),
			nnc.mountHostRW("/root/.config/fish", caller.envKV["HOME"] + "/.config/fish"),
			nnc.mountTmpfs("/dev"),
			nnc.mountHostRW("/_", caller.wd),
		],
		env: [
			"HOME=/root",
			"TERM=xterm-256color",
		],
		wd: "/_",
	}
