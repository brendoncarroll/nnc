package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"go.brendoncarroll.net/nnc/src/internal/testutil"
)

func main() {
	var w io.Writer = os.Stderr
	fmt.Fprintln(w, "Hello, World!")

	fmt.Fprintln(w, "ENV", os.Environ())
	fmt.Fprintln(w, "UID", os.Getuid())
	fmt.Fprintln(w, "GID", os.Getgid())
	fmt.Fprintln(w, "PID", os.Getpid())

	fmt.Fprintln(w, "NET")
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, ifi := range ifs {
		fmt.Fprintln(w, "  ", ifi.Index, ifi.Name, ifi.Flags, ifi.HardwareAddr)
	}

	ls(w, "/")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "WORKDIR:", wd)

	smry := testutil.MakeProcessSummary()
	data, err := json.Marshal(smry)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(data)
}

func ls(w io.Writer, path string) {
	fmt.Fprintln(w, "READDIR", path)
	ents, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, ent := range ents {
		info, err := ent.Info()
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%-16s %-12v %d\n", ent.Name(), info.Mode(), info.Size())
	}
}
