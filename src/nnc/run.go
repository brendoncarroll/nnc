//go:build linux

package nnc

import (
	"context"
	"os"

	"blobcache.io/blobcache/src/blobcache"
)

type runSettings struct {
	files []*os.File
}

// RunOption configures the Run function
type RunOption = func(stng *runSettings)

func RunSetFiles(fs ...*os.File) RunOption {
	return func(stng *runSettings) {
		stng.files = fs
	}
}

// Run runs a container given a path to the nnc_main binary and a spec for the container.
func Run(ctx context.Context, shimCID blobcache.CID, spec ContainerSpec, opts ...RunOption) (int, error) {
	if err := spec.Validate(); err != nil {
		return 0, err
	}

	var stng runSettings
	for _, opt := range opts {
		opt(&stng)
	}

	sys := &System{
		shimCID: shimCID,
	}
	proc, err := sys.spawn(spec, stng)
	if err != nil {
		return -1, err
	}
	ps, err := proc.Wait()
	if err != nil {
		return -1, err
	}
	return ps.ExitCode(), nil
}
