package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/plugin"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cache, err := wazero.NewCompilationCacheWithDir("./.wazero-cache")
	if err != nil {
		log.Fatalf("Failed to create compilation cache dir: %v", err)
	}

	configRuntime := wazero.NewRuntimeConfigCompiler().
		WithCompilationCache(cache)
	r := wazero.NewRuntimeWithConfig(ctx, configRuntime)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// _, err = r.NewHostModuleBuilder("env").
	// 	NewFunctionBuilder().
	// 	WithFunc(func(ctx context.Context, m api.Module, ptr, len uint32) uint32 {
	// 		return a + b
	// 	}).
	// 	Export("log").
	// 	Instantiate(ctx)
	// if err != nil {
	// 	log.Fatalf("Failed to create host module: %v", err)
	// }

	res, err := plugin.WASMCall[map[string]any](ctx, r, "hifi", "info", nil)
	if err != nil {
		log.Fatalf("Failed to call plugin: %v", err)
	}

	fmt.Println(res)
}
