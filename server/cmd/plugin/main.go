package main

import (
	"context"
	"fmt"

	"github.com/extism/go-sdk"
)

func main() {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: "./plugin/hifi.wasm",
			},
		},
		Config: map[string]string{
			"vowels": "aeiouyAEIOUY",
		},
	}

	ctx := context.Background()
	config := extism.PluginConfig{
		EnableWasi: true,
	}
	compiledPlugin, err := extism.NewCompiledPlugin(ctx, manifest, config, []extism.HostFunction{})
	if err != nil {
		panic(err)
	}
	defer compiledPlugin.Close(ctx)

	plugin, err := compiledPlugin.Instance(ctx, extism.PluginInstanceConfig{})
	if err != nil {
		fmt.Printf("Failed to initialize plugin: %v\n", err)
		return
	}
	defer plugin.Close(ctx)

	_, out, err := plugin.CallWithContext(ctx, "info", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := string(out)
	fmt.Println(response)
}
