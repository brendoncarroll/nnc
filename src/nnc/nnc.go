// package nnc provides No Nonsense Containers.
package nnc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"blobcache.io/blobcache/src/blobcache"
	"lukechampine.com/blake3"
)

type System struct {
	// nncMain is the executable that runs the main.
	shimCID blobcache.CID
}

func (sys *System) Spawn(spec ContainerSpec) (*os.Process, error) {
	return sys.spawn(spec, runSettings{})
}

func (sys *System) spawn(spec ContainerSpec, rstng runSettings) (*os.Process, error) {
	shimPath := BinPath(sys.shimCID)
	if spec.Env == nil {
		spec.Env = []string{}
	}
	files := append([]*os.File{}, rstng.files...)
	var devFiles []*os.File
	for i := range spec.Mounts {
		m := &spec.Mounts[i]
		if m.Src.HostDev == nil {
			continue
		}
		devPath := "/" + m.Dst
		f, err := os.OpenFile(devPath, os.O_RDWR, 0)
		if err != nil {
			closeAll(devFiles)
			return nil, fmt.Errorf("opening device %s: %w", devPath, err)
		}
		devFiles = append(devFiles, f)
		fd := len(files)
		m.Src.HostDev = &fd
		files = append(files, f)
	}
	proc, err := os.StartProcess(shimPath,
		[]string{"", marshalSpec(spec)},
		&os.ProcAttr{
			Sys: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWUSER |
					syscall.CLONE_NEWPID,
				UidMappings: []syscall.SysProcIDMap{
					{ContainerID: 0, HostID: os.Getuid(), Size: 1},
				},
				GidMappings: []syscall.SysProcIDMap{
					{ContainerID: 0, HostID: os.Getgid(), Size: 1},
				},
			},
			Env:   spec.Env,
			Files: files,
		},
	)
	closeAll(devFiles)
	return proc, err
}

func closeAll(files []*os.File) {
	for _, f := range files {
		f.Close()
	}
}

func marshalSpec(spec ContainerSpec) string {
	b, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func BinPath(x blobcache.CID) string {
	return filepath.Join(os.TempDir(), "nnc", x.String())
}

func PostBin(x []byte) (blobcache.CID, error) {
	cid := blake3.Sum256(x)
	p := BinPath(cid)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return blobcache.CID{}, err
	}
	_, err := os.Stat(p)
	if err != nil && !os.IsNotExist(err) {
		return blobcache.CID{}, err
	} else if os.IsNotExist(err) {
		if err := os.WriteFile(p, x, 0o555); err != nil {
			return blobcache.CID{}, err
		}
	}
	return cid, nil
}

func LoadBin(x blobcache.CID) ([]byte, error) {
	return os.ReadFile(BinPath(x))
}
