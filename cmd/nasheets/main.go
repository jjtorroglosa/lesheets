package main

import (
	"fmt"
	"log"
	"nasheets/internal"
	"os"
)

func main() {
	// Read the song file
	if len(os.Args) <= 1 {
		log.Fatalf("usage: %s <input_file>", os.Args[0])
	}
	file := os.Args[1]
	data, err := os.ReadFile(file)

	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	lexer := internal.NewLexer(string(data))
	tokens := lexer.Lex()

	for _, t := range tokens {
		fmt.Printf("%s: %s\n", t.Type, t.Value)
	}

	song := internal.ParseSong(tokens)

	internal.RenderSongHTML(song, "views/index.html")
	song.PrintSong()
}
