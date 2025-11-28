package testutil

import (
	"context"
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
