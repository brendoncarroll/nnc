// package nnc provides No Nonsense Containers.
package nnc

import (
	"fmt"
	"os"
	"path/filepath"

	"blobcache.io/blobcache/src/blobcache"
	"lukechampine.com/blake3"
)

type MountSrc struct {
	// TmpFS mounts a tmpfs at the given path
	TmpFS  *struct{} `json:"tmpfs,omitempty"`
	ProcFS *struct{} `json:"procfs,omitempty"`
	SysFS  *struct{} `json:"sysfs,omitempty"`

	// HostRO mounts a host path into the container, as read-only
	HostRO *string `json:"host_ro,omitempty"`
	// HostRW mounts a host path into the container, as read-only
	HostRW *string `json:"host_rw,omitempty"`
}

func (m *MountSrc) Validate() error {
	var set []string
	if m.TmpFS != nil {
		set = append(set, "tmpfs")
	}
	if m.ProcFS != nil {
		set = append(set, "procfs")
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
	Host *string
}

type ContainerSpec struct {
	// Main is the CID of the binary to run as PID 1
	Main blobcache.CID `json:"main"`

	// Args will be passed to the Main
	Args []string `json:"args"`
	// Env will be used as the environment
	Env []string `json:"env"`

	Mounts  []MountSpec   `json:"mounts"`
	Network []NetworkSpec `json:"network"`
}

func (s *ContainerSpec) Validate() error {
	for _, m := range s.Mounts {
		if err := m.Src.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func BinPath(x blobcache.CID) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("nnc-%v", x))
}

func PostBin(x []byte) (blobcache.CID, error) {
	cid := blake3.Sum256(x)
	p := BinPath(cid)
	_, err := os.Stat(p)
	if err != nil && !os.IsNotExist(err) {
		return blobcache.CID{}, err
	} else if os.IsNotExist(err) {
		if err := os.WriteFile(p, x, 0o777); err != nil {
			return blobcache.CID{}, err
		}
	}
	return cid, nil
}

func LoadBin(x blobcache.CID) ([]byte, error) {
	return os.ReadFile(BinPath(x))
}
