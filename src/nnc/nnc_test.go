//go:build linux

package nnc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"blobcache.io/blobcache/src/blobcache"
	"github.com/stretchr/testify/require"
	"go.brendoncarroll.net/nnc/src/internal/testutil"
)

func TestRun(t *testing.T) {
	testBin := testutil.BuildLinuxAmd64(t, "../internal/testbin")
	testBinCID := postExec(t, testBin)
	scratchDir := t.TempDir()
	type testCase struct {
		Name string
		Spec ContainerSpec

		ExitCode int
		Check    func(testing.TB, *testutil.ProcSummary)
	}
	tcs := []testCase{
		{
			Name: "EnvVar1",
			Spec: ContainerSpec{
				Main: testBinCID,
				Env:  []string{"KEY1=VALUE1"},
			},
			Check: func(t testing.TB, ps *testutil.ProcSummary) {
				require.Len(t, ps.Env, 1)
				require.Equal(t, ps.Env[0], "KEY1=VALUE1")
			},
		},
		{
			Spec: ContainerSpec{
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
			}},
	}

	shimBin := setup(t)
	shimBinCID := postExec(t, shimBin)

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			ctx := testutil.Context(t)
			pr, pw, err := os.Pipe()
			require.NoError(t, err)
			buf := &bytes.Buffer{}
			go func() {
				_, err := io.Copy(buf, pr)
				require.NoError(t, err)
			}()
			ec, err := Run(ctx, shimBinCID, tc.Spec, RunSetFiles(os.Stdin, pw, os.Stderr))
			require.NoError(t, err)
			require.Equal(t, ec, tc.ExitCode)

			smry, err := testutil.ParseProcessSummary(buf.Bytes())
			require.NoError(t, err)
			if tc.Check != nil {
				tc.Check(t, smry)
			}
		})
	}
}

func setup(t testing.TB) []byte {
	return testutil.BuildLinuxAmd64(t, "./nnc_shim")
}

func postExec(t testing.TB, data []byte) blobcache.CID {
	cid, err := PostBin(data)
	require.NoError(t, err)
	return cid
}
