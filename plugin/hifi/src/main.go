package main

import (
	"fmt"

	"github.com/DimitriLaPoudre/MusicShackPluginHifi/src/model"
	"github.com/extism/go-pdk"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func main() {}

//go:wasmexport info
func Info() int32 {
	pdk.OutputJSON(map[string]any{
		"name":     "hifi",
		"provider": "tidal",
	})
	return 0
}

//go:wasmexport public_config_schema
func PublicConfigSchema() int32 {
	pdk.OutputJSON(map[string]any{})
	return 0
}

//go:wasmexport update_config
func UpdateConfig() int32 {
	var publicConfig model.PublicConfig
	if err := pdk.InputJSON(&publicConfig); err != nil {
		pdk.SetError(fmt.Errorf("input to JSON: %w", err))
		return 1
	}
	if err := validate.Struct(publicConfig); err != nil {
		pdk.SetError(fmt.Errorf("validate input: %w", err))
		return 1
	}

	pdk.OutputJSON(model.Config{
		Public: publicConfig,
	})
	return 0
}
