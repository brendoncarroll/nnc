//go:build linux

package nnc

import (
	"context"
	"encoding/json"
	"os"
	"syscall"

	"blobcache.io/blobcache/src/blobcache"
)

// Run runs a container given a path to the nnc_main binary and a spec for the container.
func Run(ctx context.Context, shimCID blobcache.CID, spec ContainerSpec) (int, error) {
	sys := &System{
		shimCID: shimCID,
	}
	proc, err := sys.Start(spec)
	if err != nil {
		return -1, err
	}
	ps, err := proc.Wait()
	if err != nil {
		return -1, err
	}
	return ps.ExitCode(), nil
}

type System struct {
	// nncMain is the executable that runs the main.
	shimCID blobcache.CID
}

func (sys *System) Start(spec ContainerSpec) (*os.Process, error) {
	shimPath := BinPath(sys.shimCID)

	return os.StartProcess(shimPath,
		[]string{"", marshalSpec(spec)},
		&os.ProcAttr{
			Sys: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWNS |
					syscall.CLONE_NEWUTS |
					syscall.CLONE_NEWPID |
					syscall.CLONE_NEWUSER |
					syscall.CLONE_NEWIPC |
					syscall.CLONE_NEWNET,
				UidMappings: []syscall.SysProcIDMap{
					{ContainerID: 0, HostID: os.Getuid(), Size: 1},
				},
				GidMappings: []syscall.SysProcIDMap{
					{ContainerID: 0, HostID: os.Getgid(), Size: 1},
				},
			},
			Env: []string{},
			Files: []*os.File{
				os.Stdin,
				os.Stdout,
				os.Stderr,
			},
		},
	)
}

func marshalSpec(spec ContainerSpec) string {
	b, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}
	return string(b)
}
