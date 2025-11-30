// package nnc provides No Nonsense Containers.
package nnc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	return os.StartProcess(shimPath,
		[]string{"", marshalSpec(spec)},
		&os.ProcAttr{
			Sys: &syscall.SysProcAttr{
				Cloneflags: syscall.CLONE_NEWUSER |
					syscall.CLONE_NEWNS |
					syscall.CLONE_NEWNET |
					syscall.CLONE_NEWUTS |
					syscall.CLONE_NEWPID |
					syscall.CLONE_NEWIPC,
				UidMappings: []syscall.SysProcIDMap{
					{ContainerID: 0, HostID: os.Getuid(), Size: 1},
				},
				GidMappings: []syscall.SysProcIDMap{
					{ContainerID: 0, HostID: os.Getgid(), Size: 1},
				},
			},
			Env:   spec.Env,
			Files: rstng.files,
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
		if err := os.WriteFile(p, x, 0o555); err != nil {
			return blobcache.CID{}, err
		}
	}
	return cid, nil
}

func AddCapNetAdmin(cid blobcache.CID) error {
	out, err := exec.Command("sudo", "setcap", "cap_sys_admin+ep cap_net_admin+ep", BinPath(cid)).Output()
	if err != nil {
		return err
	}
	log.Println("output from setcap", string(out))
	return nil
}

func LoadBin(x blobcache.CID) ([]byte, error) {
	return os.ReadFile(BinPath(x))
}
