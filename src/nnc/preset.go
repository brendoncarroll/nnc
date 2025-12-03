package nnc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
)

// Present transforms a ContainerSpec
type Preset interface {
	Apply(x ContainerSpec) (*ContainerSpec, error)
}

func ApplyPresets(init ContainerSpec, presets ...Preset) (*ContainerSpec, error) {
	x := init
	for _, preset := range presets {
		y, err := preset.Apply(x)
		if err != nil {
			return nil, err
		}
		x = *y
	}
	return &x, nil
}

func OpenPresetDir() (*os.Root, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfgDir := filepath.Join(homeDir, ".config/nnc/presets")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return nil, err
	}
	r, err := os.OpenRoot(cfgDir)
	if err != nil {
		return nil, err
	}
	return r, nil
}

const jsonnetExt = ".jsonnet"

type Source struct {
	Prefix string
	Root   *os.Root
}

var _ jsonnet.Importer = &JSImporter{}

type JSImporter struct {
	Sources []Source

	cache map[string]jsonnet.Contents
}

func (jsi *JSImporter) Import(impFrom, impPath string) (contents jsonnet.Contents, foundAt string, _ error) {
	if jsi.cache == nil {
		jsi.cache = make(map[string]jsonnet.Contents)
	}
	var p string
	if impFrom != "" && strings.HasPrefix(impPath, "./") {
		p = filepath.Join(filepath.Dir(impFrom), impPath)
		if strings.HasPrefix(impFrom, "./") {
			p = "./" + p
		}
	} else {
		p = impPath
	}

	for _, src := range jsi.Sources {
		if p2, yes := strings.CutPrefix(p, src.Prefix); yes {
			if c, exists := jsi.cache[p]; exists {
				return c, p, nil
			}
			f, err := src.Root.Open(p2)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return jsonnet.Contents{}, "", err
			}
			if err := func() error {
				defer f.Close()
				data, err := io.ReadAll(f)
				if err != nil {
					return err
				}
				jsi.cache[p] = jsonnet.MakeContents(string(data))
				return nil
			}(); err != nil {
				return jsonnet.Contents{}, "", err
			}
			return jsi.cache[p], p, nil
		}
	}
	return jsonnet.Contents{}, "", fmt.Errorf("could not find preset %q", p)
}

func NewJsonnetPreset(srcs []Source, main string) (*JsonnetPreset, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&JSImporter{
		Sources: srcs,
	})

	if !strings.HasPrefix(main, jsonnetExt) {
		main += jsonnetExt
	}
	_, err := vm.ResolveImport("", main)
	if err != nil {
		return nil, err
	}
	return &JsonnetPreset{
		vm:   vm,
		main: main,
	}, nil
}

// JnsonnetPreset is a preset implemented with Jsonnet.
type JsonnetPreset struct {
	vm   *jsonnet.VM
	main string
}

func (jp JsonnetPreset) Apply(x ContainerSpec) (*ContainerSpec, error) {
	jd, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	jp.vm.TLAReset()
	jp.vm.TLACode("spec", string(jd))
	jp.vm.TLACode("caller", string(jsonMarshal(MkCallerCtx())))
	jout, err := jp.vm.EvaluateFile(jp.main)
	if err != nil {
		return nil, err
	}
	var y ContainerSpec
	if err := json.Unmarshal([]byte(jout), &y); err != nil {
		return nil, err
	}
	return &y, nil
}

// CallerCtx is the resources available to the process
type CallerCtx struct {
	Env   []string          `json:"env"`
	EnvKV map[string]string `json:"envKV"`

	WD  string   `json:"wd"`
	FDs []string `json:"fds"`
}

func MkCallerCtx() CallerCtx {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	env := os.Environ()
	envkv := make(map[string]string)
	for _, pair := range env {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		envkv[parts[0]] = parts[1]
	}
	return CallerCtx{
		Env:   os.Environ(),
		EnvKV: envkv,
		WD:    wd,
	}
}

type FD string

func NewFD(x uintptr) FD {
	ret := fmt.Sprintf("%016x", x)
	return FD(ret)
}

func jsonMarshal(x any) []byte {
	data, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return data
}
