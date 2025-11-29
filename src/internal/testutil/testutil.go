package testutil

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Context(t testing.TB) context.Context {
	return context.Background()
}

func BuildLinuxAmd64(t testing.TB, srcDir string) []byte {
	outPath := filepath.Join(t.TempDir(), "main-bin")
	defer os.Remove(outPath)
	cmd := exec.Command("go", "build",
		"-o", outPath,
		srcDir)
	cmd.Env = []string{
		"GOOS=linux",
		"GOARCH=amd64",
		"CGO_ENABLED=0",
	}
	for _, key := range []string{
		"GOPATH",
		"GOCACHE",
		"GOROOT",
		"HOME",
	} {
		if val := os.Getenv(key); val != "" {
			cmd.Env = append(cmd.Env, key+"="+val)
		}
	}
	cmdOut, err := cmd.CombinedOutput()
	if len(cmdOut) != 0 {
		t.Log("cmd out: ", string(cmdOut))
	}
	require.NoError(t, err)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)
	return data
}

func MakeProcessSummary() ProcSummary {
	files, err := MakeFileSummaries(os.DirFS("/"), ".")
	if err != nil {
		panic(err)
	}
	// TODO: load files
	return ProcSummary{
		PID:  os.Getpid(),
		Args: os.Args,
		Env:  os.Environ(),
		UID:  os.Getuid(),
		GID:  os.Getgid(),

		Files: files,
	}
}

func ParseProcessSummary(x []byte) (*ProcSummary, error) {
	var ret ProcSummary
	err := json.Unmarshal(x, &ret)
	return &ret, err
}

type ProcSummary struct {
	PID  int
	Env  []string
	Args []string

	UID int
	GID int

	Files []FileSummary
}

func MakeFileSummaries(fsys fs.FS, p string) (ret []FileSummary, _ error) {
	err := fs.WalkDir(fsys, p, func(dirp string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}
		finfo, err := d.Info()
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		ret = append(ret, FileSummary{
			Path: filepath.Join(dirp, d.Name()),
			Size: finfo.Size(),
			Mode: finfo.Mode(),
		})
		return nil
	})
	return ret, err
}

type FileSummary struct {
	Path string
	Size int64
	Mode fs.FileMode
}
