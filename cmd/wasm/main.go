//go:build js && wasm
// +build js,wasm

package main

import (
	"lesheets/internal"
	"syscall/js"
)

func lesheetToJson(this js.Value, args []js.Value) any {
	inputStr := args[0].String() // Convert JS string to Go string
	song, err := internal.ParseSongFromString(inputStr)
	if err != nil {
		return internal.RenderError(err)
	}
	html, err := internal.RenderSongHtml(
		internal.RenderConfig{
			WithLiveReload: false,
			WholeHtml:      false,
			WithEditor:     true,
		},
		inputStr,
		song,
		"some",
	)

	if err != nil {
		html = internal.RenderError(err)
	}

	return html
}

func main() {
	js.Global().Set("go_lesheetToJson", js.FuncOf(lesheetToJson))
	select {}
}
