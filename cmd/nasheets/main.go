package main

import (
	"flag"
	"fmt"
	"log"
	"nasheets/internal"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outputDir := flag.String("d", "output", "Output dir (e.g., )")

	// Parse CLI args
	flag.Parse()

	// Remaining non-flag arguments
	args := flag.Args()

	// Read the song file
	if len(args) != 1 {
		log.Fatalf("usage: %s <input_file>", os.Args[0])
	}
	file := args[0]
	output := ""
	if len(os.Args) == 3 {
		output = os.Args[2]
	} else {
		output = strings.TrimSuffix(file, ".nns") + ".html"
		output = filepath.Base(output)
	}
	data, err := os.ReadFile(file)

	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fmt.Printf("Rendering %s\n", file)
	lexer := internal.NewLexer(string(data))
	tokens := lexer.Lex()

	// for _, t := range tokens {
	// 	fmt.Printf("%s: %s\n", t.Type, t.Value)
	// }

	song := internal.ParseSong(tokens)

	internal.RenderSongHTML(song, *outputDir+"/"+output)
	//song.PrintSong()
}
