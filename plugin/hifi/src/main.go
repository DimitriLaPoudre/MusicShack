package main

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/DimitriLaPoudre/MusicShackPluginHifi/src/model"
)

func main() {}

//go:wasmexport allocate
func Allocate(size uint32) uint32 {
	data := make([]byte, size)
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(data)))
	return uint32(ptr)
}

func IntoJSON(v any) uint64 {
	data, err := json.Marshal(v)
	if err != nil {
		return ErrorToJSON(err)
	}
	addr := uintptr(unsafe.Pointer(unsafe.SliceData(data)))
	return (uint64(addr) << 32) | uint64(len(data))
}

func FromJSON[T any](ptr uint32, len uint32) (T, error) {
	var t T
	slice := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), len)
	if err := json.Unmarshal(slice, &t); err != nil {
		return t, err
	}
	return t, nil
}

func ErrorToJSON(err error) uint64 {
	data := fmt.Appendf([]byte{}, `{"error": "%v"}`, err.Error())
	addr := uintptr(unsafe.Pointer(unsafe.SliceData(data)))
	return (uint64(addr) << 32) | uint64(len(data))
}

//go:wasmexport info
func Info() uint64 {
	return IntoJSON(map[string]any{
		"name":     "hifi",
		"provider": "tidal",
	})
}

//go:wasmexport public_config_schema
func PublicConfigSchema() uint64 {
	return IntoJSON(map[string]any{})
}

//go:wasmexport update_config
func UpdateConfig(ptr uint32, len uint32) uint64 {
	publicConfig, err := FromJSON[model.PublicConfig](ptr, len)
	if err != nil {
		return ErrorToJSON(err)
	}

	fmt.Println("wasm", publicConfig)

	return IntoJSON(model.Config{
		Public: publicConfig,
	})
}
