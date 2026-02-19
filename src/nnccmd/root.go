package nnccmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.brendoncarroll.net/nnc/src/nnc"
	"go.brendoncarroll.net/star"
)

func Main() {
	star.Main(rootCmd)
}

var rootCmd = star.NewDir(star.Metadata{
	Short: "No Nonsense Containers",
}, map[string]star.Command{
	"run":        runCmd,
	"enter":      enterCmd,
	"print-spec": printSpecCmd,
})

var enterCmd = star.Command{
	Metadata: star.Metadata{
		Short: "enter a new container with access to the current directory",
	},
	Pos: []star.Positional{mainParam, argsParam},
	Flags: map[string]star.Flag{
		"dr":     droParam,
		"dw":     drwParam,
		"dev":    devParam,
		"env":    envParam,
		"preset": presetsParam,
	},
	F: func(c star.Context) error {
		shellPath := os.Getenv("SHELL")
		if shellPath == "" {
			return fmt.Errorf("SHELL must be set to use enter")
		}
		shellBin, err := os.ReadFile(shellPath)
		if err != nil {
			return err
		}
		mainCID, err := nnc.PostBin(shellBin)
		if err != nil {
			return err
		}
		initSpec := nnc.ContainerSpec{
			Main: mainCID,
		}
		cspec, err := configure(initSpec, c)
		if err != nil {
			return err
		}
		shimCID, err := nnc.PostBin(shimBin)

		if err != nil {
			return err
		}
		ctx := c.Context
		ec, err := nnc.Run(ctx, shimCID, *cspec,
			nnc.RunSetFiles(os.Stdin, os.Stdout, os.Stderr),
		)
		if err != nil {
			return err
		}
		if ec != 0 {
			os.Exit(ec)
		}
		return nil
	},
}

var printSpecCmd = star.Command{
	Metadata: star.Metadata{
		Short: "like run, but prints the container spec instead of running it",
	},

	Pos: []star.Positional{mainParam, argsParam},
	Flags: map[string]star.Flag{
		"dr":     droParam,
		"dw":     drwParam,
		"dev":    devParam,
		"env":    envParam,
		"preset": presetsParam,
	},
	F: func(c star.Context) error {
		cspec, err := configure(nnc.ContainerSpec{}, c)
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(cspec, "", "  ")
		if err != nil {
			return err
		}
		c.Printf("%s\n", data)
		return nil
	},
}

var runCmd = star.Command{
	Metadata: star.Metadata{
		Short: "run a container",
	},
	Pos: []star.Positional{mainParam, argsParam},
	Flags: map[string]star.Flag{
		"dr":     droParam,
		"dw":     drwParam,
		"dev":    devParam,
		"env":    envParam,
		"ldd":    lddParam,
		"preset": presetsParam,
	},
	F: func(c star.Context) error {
		cspec, err := configure(nnc.ContainerSpec{}, c)
		if err != nil {
			return err
		}
		shimCID, err := nnc.PostBin(shimBin)
		if err != nil {
			return err
		}
		ctx := c.Context
		ec, err := nnc.Run(ctx, shimCID, *cspec,
			nnc.RunSetFiles(os.Stdin, os.Stdout, os.Stderr),
		)
		if err != nil {
			return err
		}
		os.Exit(ec)
		return nil
	},
}

func addSysMounts(m []nnc.MountSpec) []nnc.MountSpec {
	// TODO: disable sysfs for now, it doesn't always work
	// depending on if the network namespace is shared.
	// With all namespaces completely fresh, it appears to work.
	// Once we always create a fresh network namespace, then we can add this back.
	// m = append(m, nnc.MountSpec{
	// 	Dst: "sys",
	// 	Src: nnc.MountSrc{
	// 		SysFS: &struct{}{},
	// 	},
	// })
	m = append(m, nnc.MountSpec{
		Dst: "proc",
		Src: nnc.MountSrc{
			ProcFS: &struct{}{},
		},
	})
	m = append(m, nnc.MountSpec{
		Dst: "dev",
		Src: nnc.MountSrc{
			TmpFS: &struct{}{},
		},
	})
	return m
}

// configure configures cspec, using parameters from c
// and returns a copy of cspec with the configuration applied.
func configure(cspec nnc.ContainerSpec, c star.Context) (*nnc.ContainerSpec, error) {
	dros := droParam.Load(c)
	drws := drwParam.Load(c)
	devs := devParam.Load(c)
	cspec.Mounts = addSysMounts(cspec.Mounts)
	cspec.Mounts = append(cspec.Mounts, dros...)
	cspec.Mounts = append(cspec.Mounts, drws...)
	cspec.Mounts = append(cspec.Mounts, devs...)

	envs := envParam.Load(c)
	cspec.Env = append(cspec.Env, envs...)

	if mainBin, ok := mainParam.LoadOpt(c); ok {
		mainCID, err := nnc.PostBin(mainBin)
		if err != nil {
			return nil, err
		}
		cspec.Main = mainCID
		cspec.Args = []string{"main"}
	}

	args := argsParam.Load(c)
	cspec.Args = args

	// apply all presets last
	presets := presetsParam.Load(c)
	cspec2, err := nnc.ApplyPresets(cspec, presets...)
	if err != nil {
		return nil, err
	}
	cspec = *cspec2

	return &cspec, nil
}

var mainParam = star.Optional[[]byte]{
	ID:       "main",
	Parse:    os.ReadFile,
	ShortDoc: "the filepath to the program to run in the container",
}

var droParam = star.Repeated[nnc.MountSpec]{
	ID:       "dir-ro",
	Parse:    parseMountSpec(false),
	ShortDoc: "mount a directory read-only in the container",
}

var drwParam = star.Repeated[nnc.MountSpec]{
	ID:       "dir-rw",
	Parse:    parseMountSpec(true),
	ShortDoc: "mount a directory read-write in the container",
}

var devParam = star.Repeated[nnc.MountSpec]{
	ID:       "dev",
	Parse:    parseDevSpec,
	ShortDoc: "mount a host device into the container (e.g. null, urandom)",
}

var envParam = star.Repeated[string]{
	ID:       "env-var",
	Parse:    star.ParseString,
	ShortDoc: "set and environment variable in a container",
}

var lddParam = star.Repeated[string]{
	ID:       "ldd",
	Parse:    star.ParseString,
	ShortDoc: "include all transitively reachable shared objects",
}

var presetsParam = star.Repeated[nnc.Preset]{
	ID: "preset",
	Parse: func(x string) (nnc.Preset, error) {
		srcs, err := getSources()
		if err != nil {
			return nil, err
		}
		return nnc.NewJsonnetPreset(srcs, x)
	},
	ShortDoc: "specify a preset",
}

var argsParam = star.Repeated[string]{
	ID:       "args",
	Parse:    star.ParseString,
	ShortDoc: "args will be passed on to the container process",
}

func getSources() ([]nnc.Source, error) {
	wdPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	wdRoot, err := os.OpenRoot(wdPath)
	if err != nil {
		return nil, err
	}
	presetDir, err := nnc.OpenPresetDir()
	if err != nil {
		return nil, err
	}
	return []nnc.Source{
		{Prefix: "./", Root: wdRoot},
		{Prefix: "", Root: presetDir},
	}, nil
}

func parseDevSpec(x string) (nnc.MountSpec, error) {
	placeholder := 0 // fd number will be set by spawn()
	return nnc.MountSpec{
		Dst: "dev/" + x,
		Src: nnc.MountSrc{
			HostDev: &placeholder,
		},
	}, nil
}

func parseMountSpec(rw bool) func(x string) (nnc.MountSpec, error) {
	return func(x string) (nnc.MountSpec, error) {
		parts := strings.SplitN(x, ":", 2)
		if len(parts) < 2 {
			return nnc.MountSpec{}, fmt.Errorf("invalid mount spec %q", x)
		}
		var src nnc.MountSrc
		if rw {
			src = nnc.MountSrc{
				HostRW: star.Ptr(parts[1]),
			}
		} else {
			src = nnc.MountSrc{
				HostRO: star.Ptr(parts[1]),
			}
		}
		return nnc.MountSpec{
			Dst: parts[0],
			Src: src,
		}, nil
	}
}
