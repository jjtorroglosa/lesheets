//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"nasheets/internal"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outputDir := flag.String("d", "output", "Output dir")
	outputFilename := flag.String("o", "", "Output filename")
	printSong := flag.Bool("p", false, "Print song")
	printTokens := flag.Bool("t", false, "Print tokens")
	input := flag.String("i", "", "Input filedir")

	// Parse CLI args
	flag.Parse()

	// Remaining non-flag arguments
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		internal.Fatalf("invalid args")
	}
	cmd := args[0]

	// Read the song file
	if *input == "" {
		internal.Fatalf("usage: %s -i <input_file>", os.Args[0])
	}
	inputFile := *input
	data, err := os.ReadFile(inputFile)

	if *outputFilename == "" {
		*outputFilename = strings.TrimSuffix(inputFile, ".nns") + ".html"
		*outputFilename = filepath.Base(*outputFilename)
	}

	if err != nil {
		internal.Fatalf("Failed to read file: %v", err)
	}

	song, err := internal.ParseSongFromString(string(data))
	if err != nil {
		internal.Fatalf("error parsing song: %v", err)
	}

	switch cmd {
	case "json":
		j, err := json.Marshal(song)
		if err != nil {
			log.Fatalf("Error marshalling json: %v", err)
		}
		fmt.Println(string(j))
	case "html":
		fmt.Printf("Rendering %s to %s\n", inputFile, *outputDir+"/"+*outputFilename)

		if *printTokens {
			lexer := internal.NewLexer(string(data))
			_, err := lexer.Lex()
			if err != nil {
				log.Fatalf("lexer error: %v", err)
			}

			// for _, t := range tokens {
			// 	fmt.Printf("%s: %s\n", t.Type, t.Value)
			// }
		}

		if *printSong {
			song.PrintSong()
		}
		internal.RenderSongHTML(song, *outputDir+"/"+*outputFilename)
	}
}
