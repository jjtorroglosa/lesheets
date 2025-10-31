package main

import (
	"nasheets/internal"
	"syscall/js"
)

func nasheetToJson(this js.Value, args []js.Value) interface{} {
	s := args[0].String() // Convert JS string to Go string

	println("1 Hi from go: " + s)
	song, err := internal.ParseSongFromString(s)
	if err != nil {
		return err
	}

	//song.PrintSong()
	j := song.ToJson()

	return string(j)
}

func main() {
	js.Global().Set("go_nasheetToJson", js.FuncOf(nasheetToJson))
	select {}
}
