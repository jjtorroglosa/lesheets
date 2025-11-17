//go:build js && wasm
// +build js,wasm

package main

import (
	"lesheets/internal"
	"syscall/js"
)

func lesheetToHtml(this js.Value, args []js.Value) any {
	defer func() {
		if r := recover(); r != nil {
			println("Recovered panic:", r)
		}
	}()
	inputStr := args[0].String() // Convert JS string to Go string
	println("heeeeeeeeeeeeeeeeeeeeeeeeeee")
	song, err := internal.ParseSongFromString(inputStr)
	if err != nil {
		return js.ValueOf(internal.RenderError(err))
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

	return js.ValueOf(html)
}

func main() {
	js.Global().Set("go_lesheetToHtml", js.FuncOf(lesheetToHtml))
	select {}
}
