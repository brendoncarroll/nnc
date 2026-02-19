package nnc

import (
	"fmt"
	"io/fs"

	"blobcache.io/blobcache/src/blobcache"
)

// MainPath is the path in the container where the initial executable
// is copied.
const MainPath = "/main"

type MountSrc struct {
	// TmpFS mounts a tmpfs at the given path
	TmpFS    *struct{} `json:"tmpfs,omitempty"`
	ProcFS   *struct{} `json:"procfs,omitempty"`
	Devtmpfs *struct{} `json:"devtmpfs,omitempty"`
	SysFS    *struct{} `json:"sysfs,omitempty"`

	// HostRO mounts a host path into the container, as read-only
	HostRO *string `json:"host_ro,omitempty"`
	// HostRW mounts a host path into the container, as read-write
	HostRW *string `json:"host_rw,omitempty"`

	// HostDev passes a pre-opened device fd into the container.
	// The value is the fd number as a string, set by the host process.
	HostDev *int `json:"host_dev,omitempty"`
}

func (m *MountSrc) Validate() error {
	var set []string
	if m.TmpFS != nil {
		set = append(set, "tmpfs")
	}
	if m.ProcFS != nil {
		set = append(set, "procfs")
	}
	if m.Devtmpfs != nil {
		set = append(set, "devtmpfs")
	}
	if m.SysFS != nil {
		set = append(set, "sysfs")
	}
	if m.HostRO != nil {
		set = append(set, "host_ro")
	}
	if m.HostRW != nil {
		set = append(set, "host_rw")
	}
	if m.HostDev != nil {
		set = append(set, "host_dev")
	}
	if len(set) != 1 {
		return fmt.Errorf("exactly one of tmpfs, procfs, or sysfs must be set")
	}
	return nil
}

type MountSpec struct {
	// Dst is the mountpoint, the front-end of the mount
	Dst string `json:"dst"`
	// Src is backend of the mount
	Src MountSrc `json:"src"`
}

func (ms MountSpec) IsSystem() bool {
	switch {
	case ms.Src.ProcFS != nil:
		return true
	case ms.Src.SysFS != nil:
		return true
	default:
		return false
	}
}

type NetworkSpec struct {
	Name string `json:"name"`

	Backend NetBackend `json:"backend"`
}

func (nspec NetworkSpec) Validate() error {
	return nil
}

type NetBackend struct {
	// None is a network interface that is not connected to anything.
	None *struct{}
}

type DataFileSpec struct {
	// Path is where the file should be written in the container.
	Path string      `json:"path"`
	Mode fs.FileMode `json:"mode"`

	Contents DataFileSrc `json:"contents"`
}

type DataFileSrc struct {
	// Literal is the literal contents of the file.
	Literal *string `json:"lit"`
}

type ContainerSpec struct {
	// Main is the CID of the binary to run as PID 1
	Main blobcache.CID `json:"main"`

	// Args will be passed to the Main
	Args []string `json:"args"`
	// Env will be used as the environment
	Env []string `json:"env"`
	// WorkingDir, if not zero, will set the working directory of the main process.
	WorkingDir string `json:"wd"`

	// Mounts is the contents of the container's mount table.
	Mounts []MountSpec `json:"mounts"`
	// Network are the interfaces to create in the container.
	Network []NetworkSpec `json:"net"`
	// Data are files that will be written inside the container.
	Data []DataFileSpec `json:"data"`
}

func (s *ContainerSpec) Validate() error {
	if s.Main.IsZero() {
		return fmt.Errorf("main CID cannot be zero")
	}
	for _, m := range s.Mounts {
		if err := m.Src.Validate(); err != nil {
			return err
		}
	}
	for _, n := range s.Network {
		if err := n.Validate(); err != nil {
			return err
		}
	}
	return nil
}
