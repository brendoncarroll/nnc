//go:build linux

package nnc

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"blobcache.io/blobcache/src/blobcache"
	"github.com/stretchr/testify/require"
	"go.brendoncarroll.net/nnc/src/internal/testutil"
	"lukechampine.com/blake3"
)

func TestRun(t *testing.T) {
	ctx := testutil.Context(t)
	shimBin := setup(t)
	shimBinCID := postExec(t, shimBin)

	scratchDir := t.TempDir()
	testBin := testutil.BuildLinuxAmd64(t, "../internal/testbin")
	testBinCID := postExec(t, testBin)
	require.NoError(t, os.WriteFile(filepath.Join(scratchDir, "testbin"), testBin, 0755))

	pr, pw, err := os.Pipe()
	require.NoError(t, err)
	buf := &bytes.Buffer{}
	go func() {
		_, err := io.Copy(buf, pr)
		require.NoError(t, err)
	}()
	ec, err := Run(ctx, shimBinCID, ContainerSpec{
		Main: testBinCID,
		Mounts: []MountSpec{
			{
				Dst: "/tmp1",
				Src: MountSrc{
					TmpFS: &struct{}{},
				},
			},
			{
				Dst: "/data1",
				// Since scratchDir is on tmpfs, this must be RW
				// Apparently the kernel can't possibly make a tmpfs read-only
				// because $REASONS
				Src: MountSrc{HostRW: &scratchDir},
			},
			{
				Dst: "/proc",
				Src: MountSrc{
					ProcFS: &struct{}{},
				},
			},
			{
				Dst: "/sys",
				Src: MountSrc{
					SysFS: &struct{}{},
				},
			},
		},
	}, RunSetFiles(
		os.Stdin,
		pw,
		os.Stderr,
	))
	require.NoError(t, err)
	require.Equal(t, 0, ec)
	_, err = testutil.ParseProcessSummary(buf.Bytes())
	require.NoError(t, err)
}

func setup(t testing.TB) []byte {
	return testutil.BuildLinuxAmd64(t, "./nnc_shim")
}

func postExec(t testing.TB, data []byte) blobcache.CID {
	cid := blobcache.CID(blake3.Sum256(data))
	_, err := os.Stat(BinPath(cid))
	if err != nil && !os.IsNotExist(err) {
		require.NoError(t, err)
	} else if os.IsNotExist(err) {
		require.NoError(t, os.WriteFile(BinPath(cid), data, 0o755))
	}
	return cid
}
