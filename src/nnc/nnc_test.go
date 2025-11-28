//go:build linux

package nnc

import (
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
				Src: MountSrc{HostRO: &scratchDir},
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
	})
	require.NoError(t, err)
	require.Equal(t, 0, ec)
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
