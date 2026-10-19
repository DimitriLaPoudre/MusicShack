package plugin

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func setupParams(ctx context.Context, m api.Module, v any) ([]uint64, error) {
	if v == nil {
		return []uint64{}, nil
	}

	json, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	allocate := m.ExportedFunction("allocate")
	if allocate == nil {
		return nil, fmt.Errorf("exported function %v: %w", "allocate", errors.New("not found"))
	}
	res, err := allocate.Call(ctx, uint64(len(json)))
	if err != nil {
		return nil, err
	}

	ptr := uint32(res[0])
	if !m.Memory().Write(ptr, []byte(json)) {
		return nil, errors.New("memory write: out of range")
	}

	return []uint64{uint64(ptr), uint64(uint32(len(json)))}, nil
}

func collectResult[T any](m api.Module, res []uint64) (T, error) {
	var v T

	ptr := uint32(res[0] >> 32)
	len := uint32(res[0])

	out, ok := m.Memory().Read(ptr, len)
	if !ok {
		return v, errors.New("memory read: out of range")
	}

	err := json.Unmarshal(out, &v)
	if err != nil {
		return v, err
	}

	// check si c'est un json d'erreur et return l'erreur

	return v, nil
}

func WASMCall[T any](ctx context.Context, r wazero.Runtime, pluginName string, functionName string, args any) (T, error) {
	var res T

	pluginWASM := pluginName
	wasmBytes, err := os.ReadFile(fmt.Sprintf("plugin/%s.wasm", pluginWASM))
	if err != nil {
		return res, fmt.Errorf("read WASM file: %w", err)
	}

	configModule := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		// WithArgs("my-wasi-app", "--config", "/app/config.json").
		// WithEnv("API_KEY", "your-secret-key").
		// WithEnv("DEBUG", "true").
		// WithFS(os.DirFS("./app-data")). // Read-only data directory
		// WithFSConfig(
		// 	wazero.NewFSConfig().
		// 		WithDirMount("./cache", "/cache"), // Writable cache directory
		// ).
		WithStartFunctions("_initialize")

	m, err := r.InstantiateWithConfig(ctx, wasmBytes, configModule)
	if err != nil {
		return res, fmt.Errorf("instantiate WASM module: %w", err)
	}
	defer m.Close(ctx)

	f := m.ExportedFunction(functionName)
	if f == nil {
		return res, fmt.Errorf("exported function %v: %w", functionName, errors.New("not found"))
	}

	params, err := setupParams(ctx, m, args)
	if err != nil {
		return res, fmt.Errorf("setup params: %w", err)
	}

	callRes, err := f.Call(ctx, params...)
	if err != nil {
		return res, fmt.Errorf("call %v: %w", functionName, err)
	}

	res, err = collectResult[T](m, callRes)
	if err != nil {
		return res, fmt.Errorf("collect result: %w", err)
	}

	return res, nil
}
