//go:build tools

package main

import (
	_ "go.bytecodealliance.org/cmd/wit-bindgen-go"
	_ "go.bytecodealliance.org/x/cabi"
	_ "go.wasmcloud.dev/wadge/cmd/wadge-bindgen-go"
)
