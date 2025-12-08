//go:build linux

// Package main is a shim for configuring namespaces before execing
// the application main
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"go.brendoncarroll.net/nnc/src/nnc"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("nnc_shim: 2 args required")
	}
	if err := run(os.Args[1:]); err != nil {
		if err := printDir("/"); err != nil {
			panic(err)
		}
		log.Fatalf("nnc_shim: %s", err)
	}
}

func run(args []string) error {
	spec, err := parseSpec(args[0])
	if err != nil {
		return err
	}
	// read the main while we still can.
	mainBin, err := nnc.LoadBin(spec.Main)
	if err != nil {
		return fmt.Errorf("loading bin: %w", err)
	}

	runtime.LockOSThread()
	if err := syscall.Unshare(syscall.CLONE_NEWNS |
		syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWIPC,
	); err != nil {
		return err
	}
	if len(spec.Network) == 0 {
		// TODO: for now if there is anything in the network spec
		// just use the host's namespace
		if err := syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
			return err
		}
	}
	// Create new tmpfs root
	newRoot, err := os.MkdirTemp("", "newroot")
	if err != nil {
		return err
	}
	if err := prepareMounts(newRoot, spec.Mounts); err != nil {
		return err
	}

	// Set working directory
	if spec.WorkingDir != "" {
		if err := os.Chdir(spec.WorkingDir); err != nil {
			return err
		}
	}

	for _, df := range spec.Data {
		if err := handleDataFile(df); err != nil {
			return err
		}
	}

	// Run the main.
	const mainPath = nnc.MainPath
	if err := os.WriteFile(mainPath, mainBin, 0o555); err != nil {
		return err
	}
	if err := syscall.Exec(mainPath, spec.Args, spec.Env); err != nil {
		return fmt.Errorf("syscall.Exec: %w", err)
	}
	return nil
}

func prepareMounts(newRoot string, mounts []nnc.MountSpec) error {
	// First, ensure we're in a new mount namespace
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("failed to make mount namespace private: %w", err)
	}
	if err := syscall.Mount("tmpfs", newRoot, "tmpfs", 0, ""); err != nil {
		log.Fatalf("mount tmpfs: %v", err)
	}

	// Handle all non-system mounts specified in the container spec
	for _, mount := range mounts {
		if !mount.IsSystem() {
			if err := handleMount("/", newRoot, mount); err != nil {
				return fmt.Errorf("failed to handle mount %s: %w", mount.Dst, err)
			}
		}
	}
	// Make old root a mount point
	putOld := newRoot + "/oldroot"
	if err := os.MkdirAll(putOld, 0755); err != nil {
		log.Fatalf("mkdir oldroot: %v", err)
	}
	// pivot_root: move / to /oldroot and make newRoot the new /
	if err := syscall.PivotRoot(newRoot, putOld); err != nil {
		log.Fatalf("pivot_root: %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		log.Fatalf("chdir /: %v", err)
	}

	for i, mount := range mounts {
		if mount.IsSystem() {
			if err := handleMount("/oldroot", "/", mount); err != nil {
				return fmt.Errorf("failed to handle mount %d %s: %w", i, mount.Dst, err)
			}
		}
	}

	// Unmount old root and remove it
	if err := syscall.Unmount("/oldroot", syscall.MNT_DETACH); err != nil {
		log.Fatalf("unmount oldroot: %v", err)
	}
	if err := os.Remove("/oldroot"); err != nil {
		log.Fatalf("remove oldroot: %v", err)
	}
	return nil
}

func handleMount(oldRoot, newRoot string, mount nnc.MountSpec) error {
	if err := mount.Src.Validate(); err != nil {
		return err
	}
	dst := filepath.Join(newRoot, mount.Dst)

	isFile := false
	switch {
	case mount.Src.HostRO != nil:
		src := filepath.Join(oldRoot, *mount.Src.HostRO)
		srcInfo, err := os.Stat(src)
		if err != nil {
			return err
		}
		isFile = !srcInfo.IsDir()
	case mount.Src.HostRW != nil:
		src := filepath.Join(oldRoot, *mount.Src.HostRW)
		srcInfo, err := os.Stat(src)
		if err != nil {
			return fmt.Errorf("performing stat on %s %w", src, err)
		}
		isFile = !srcInfo.IsDir()
	}

	// Create mount point if it doesn't exist
	if isFile {
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("failed to create mount point: %w", err)
		}
		if err := os.WriteFile(dst, nil, 777); err != nil {
			return err
		}
	} else {
		if err := os.MkdirAll(dst, 0o755); err != nil {
			return fmt.Errorf("failed to create mount point: %w", err)
		}
	}

	switch {
	case mount.Src.TmpFS != nil:
		return syscall.Mount("", dst, "tmpfs", 0, "")
	case mount.Src.ProcFS != nil:
		return syscall.Mount("", dst, "proc", 0, "")
	case mount.Src.SysFS != nil:
		return syscall.Mount("", dst, "sysfs", 0, "")
	case mount.Src.Devtmpfs != nil:
		return syscall.Mount("", dst, "devtmpfs", 0, "mode=755")

	case mount.Src.HostRO != nil:
		// log.Println("mounting", dst, "-ro->", filepath.Join(oldRoot, *mount.Src.HostRO))
		src := filepath.Join(oldRoot, *mount.Src.HostRO)
		if err := syscall.Mount(src, dst, "", syscall.MS_BIND|syscall.MS_RDONLY, ""); err != nil {
			return err
		}
		return syscall.Mount("", dst, "", syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY, "")
	case mount.Src.HostRW != nil:
		// log.Println("mounting", dst, "->", filepath.Join(oldRoot, *mount.Src.HostRW))
		src := filepath.Join(oldRoot, *mount.Src.HostRW)
		return syscall.Mount(src, dst, "", syscall.MS_BIND, "")
	default:
		panic(mount) // Validate should have caught this
	}
}

func handleDataFile(df nnc.DataFileSpec) error {
	switch {
	case df.Contents.Literal != nil:
		if err := os.MkdirAll(filepath.Dir(df.Path), 0o755); err != nil {
			return err
		}
		return os.WriteFile(df.Path, []byte(*df.Contents.Literal), df.Mode)
	default:
		return fmt.Errorf("empty data file source")
	}
}

func parseSpec(x string) (*nnc.ContainerSpec, error) {
	var spec nnc.ContainerSpec
	if err := json.Unmarshal([]byte(x), &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

func printDir(x string) error {
	ents, err := os.ReadDir(x)
	if err != nil {
		return err
	}
	for _, ent := range ents {
		info, err := ent.Info()
		if err != nil {
			return err
		}
		fmt.Printf("%v %v %v\n", ent.Name(), info.Mode(), info.Size())
	}
	return nil
}
