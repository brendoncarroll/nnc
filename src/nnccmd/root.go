package nnccmd

import (
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
	"run":   runCmd,
	"enter": enterCmd,
})

var enterCmd = star.Command{
	Metadata: star.Metadata{
		Short: "enter a new container with access to the current directory",
	},
	F: func(c star.Context) error {
		return nil
	},
}

var runCmd = star.Command{
	Metadata: star.Metadata{
		Short: "run a container",
	},
	Pos: []star.Positional{mainParam}, Flags: map[string]star.Flag{
		"dr":  droParam,
		"dw":  drwParam,
		"env": envParam,
	},
	F: func(c star.Context) error {
		var cspec nnc.ContainerSpec

		dros := droParam.Load(c)
		drws := drwParam.Load(c)
		cspec.Mounts = append(cspec.Mounts, dros...)
		cspec.Mounts = append(cspec.Mounts, drws...)

		envs := envParam.Load(c)
		cspec.Env = append(cspec.Env, envs...)

		mainBin := mainParam.Load(c)
		mainCID, err := nnc.PostBin(mainBin)
		if err != nil {
			return err
		}
		cspec.Main = mainCID

		if err := cspec.Validate(); err != nil {
			return err
		}

		shimCID, err := nnc.PostBin(shimBin)
		if err != nil {
			return err
		}
		ctx := c.Context
		ec, err := nnc.Run(ctx, shimCID, cspec)
		if err != nil {
			return err
		}
		os.Exit(ec)
		return nil
	}}

var mainParam = star.Required[[]byte]{
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

var envParam = star.Repeated[string]{
	ID:       "env-var",
	Parse:    star.ParseString,
	ShortDoc: "",
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
				HostRO: star.Ptr(parts[1]),
			}
		} else {
			src = nnc.MountSrc{
				HostRW: star.Ptr(parts[1]),
			}
		}
		return nnc.MountSpec{
			Dst: parts[0],
			Src: src,
		}, nil
	}
}
