//go:build wasm
// +build wasm

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"github.com/shellyln/go-loose-json-parser/jsonlp"
)

func normalizeJSON(this js.Value, args []js.Value) interface{} {
	src := ""
	indent := 0
	var b []byte

	if 0 < len(args) {
		src = args[0].String()
	}
	if 1 < len(args) {
		indent = args[1].Int()
	}

	parsed, err := jsonlp.ParseJSON(src, jsonlp.Linebreak_Lf, jsonlp.Interop_JSON)
	if err != nil {
		return js.ValueOf(err.Error())
	}

	if 0 < indent {
		b, err = json.MarshalIndent(parsed, "", strings.Repeat(" ", indent))
		if err != nil {
			return js.ValueOf(err.Error())
		}
	} else {
		b, err = json.Marshal(parsed)
		if err != nil {
			return js.ValueOf(err.Error())
		}
	}

	return js.ValueOf(string(b))
}

func normalizeTOML(this js.Value, args []js.Value) interface{} {
	src := ""
	indent := 0
	var b []byte

	if 0 < len(args) {
		src = args[0].String()
	}
	if 1 < len(args) {
		indent = args[1].Int()
	}

	parsed, err := jsonlp.ParseTOML(src, jsonlp.Linebreak_Lf, jsonlp.Interop_JSON)
	if err != nil {
		return js.ValueOf(err.Error())
	}

	if 0 < indent {
		b, err = json.MarshalIndent(parsed, "", strings.Repeat(" ", indent))
		if err != nil {
			return js.ValueOf(err.Error())
		}
	} else {
		b, err = json.Marshal(parsed)
		if err != nil {
			return js.ValueOf(err.Error())
		}
	}

	return js.ValueOf(string(b))
}

func getVersion(this js.Value, args []js.Value) interface{} {
	return js.ValueOf(Version)
}

func main() {
	println("Go WebAssembly Initialized")

	js.Global().Set("normalizeJSON", js.FuncOf(normalizeJSON))
	js.Global().Set("normalizeTOML", js.FuncOf(normalizeTOML))
	js.Global().Set("getVersion", js.FuncOf(getVersion))

	select {}
}
