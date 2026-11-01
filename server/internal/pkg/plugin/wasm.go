package plugin

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
)

func WASMCall[T any](ctx context.Context, r wazero.Runtime, pluginName string, functionName string, args any) (T, error) {
	var res T

	pluginWASM := pluginName
	wasmBytes, err := os.ReadFile(fmt.Sprintf("plugin/%s.wasm", pluginWASM))
	if err != nil {
		return res, fmt.Errorf("read WASM file: %w", err)
	}
	compiled, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		return res, fmt.Errorf("compile WASM data: %w", err)
	}

	var stdin bytes.Buffer
	if err := json.MarshalWrite(&stdin, args); err != nil {
		return res, fmt.Errorf("marshal input: %w", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	configModule := wazero.NewModuleConfig().
		WithStdin(&stdin).
		WithStdout(&stdout).
		WithStderr(&stderr).
		// WithArgs("my-wasi-app", "--config", "/app/config.json").
		// WithEnv("API_KEY", "your-secret-key").
		// WithFS(os.DirFS("./app-data")). // Read-only data directory
		// WithFSConfig(
		// 	wazero.NewFSConfig().
		// 		WithDirMount("./cache", "/cache"), // Writable cache directory
		// ).
		WithStartFunctions("_initialize")

	m, err := r.InstantiateModule(ctx, compiled, configModule)
	if err != nil {
		return res, fmt.Errorf("instantiate WASM module: %w", err)
	}
	defer m.Close(ctx)

	f := m.ExportedFunction(functionName)
	if f == nil {
		return res, fmt.Errorf("exported function %v: %w", functionName, errors.New("not found"))
	}

	_, err = f.Call(ctx)
	if err != nil {
		return res, fmt.Errorf("call %v: %w", functionName, err)
	}

	if stderr.Len() != 0 {
		return res, fmt.Errorf("called %v: %w", functionName, errors.New(stderr.String()))
	}

	if stdout.Len() != 0 {
		if err := json.Unmarshal(stdout.Bytes(), &res); err != nil {
			return res, fmt.Errorf("unmarshal output: %w", err)
		}

		return res, nil
	}

	return res, nil
}
