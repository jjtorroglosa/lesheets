package main

import (
	"log"
	"nasheets/internal"
	"os"
)

func main() {
	// Read the song file
	data, err := os.ReadFile("sledgehammer.nns")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	lexer := internal.NewLexer(string(data))
	tokens := lexer.Lex()

	// for _, t := range tokens {
	// 	fmt.Printf("%s: %s\n", t.Type, t.Value)
	// }

	song := internal.ParseSong(tokens)

	internal.RenderSongHTML(song, "views/index.html")
	song.PrintSong()
}
