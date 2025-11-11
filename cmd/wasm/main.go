//go:build js && wasm
// +build js,wasm

package main

import (
	"nasheets/internal"
	"syscall/js"
)

func nasheetToJson(this js.Value, args []js.Value) any {
	inputStr := args[0].String() // Convert JS string to Go string
	song, err := internal.ParseSongFromString(inputStr)
	if err != nil {
		return err
	}
	html := internal.RenderSongHtml(
		internal.RenderConfig{
			WithLiveReload: false,
			WholeHtml:      false,
			WithEditor:     true,
		},
		inputStr,
		song,
		"some",
	)
	return string(html)
}

func main() {
	js.Global().Set("go_nasheetToJson", js.FuncOf(nasheetToJson))
	select {}
}
