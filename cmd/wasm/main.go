//go:build js && wasm
// +build js,wasm

package main

import (
	"nasheets/internal"
	"syscall/js"
)

func nasheetToJson(this js.Value, args []js.Value) any {
	s := args[0].String() // Convert JS string to Go string
	song, err := internal.ParseSongFromString(s)
	if err != nil {
		return err
	}
	j := internal.RenderSongHtml(false, false, s, song, "some")
	return string(j)
}

func main() {
	js.Global().Set("go_nasheetToJson", js.FuncOf(nasheetToJson))
	select {}
}
