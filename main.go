package main

import (
	"fmt"
	"log"
	"os"
)

func printSong(song *Song) {
	for _, sec := range song.Sections {
		fmt.Println("Section:", sec.Header)
		i := 1
		for _, barline := range sec.BarsLines {
			for _, bar := range barline {
				fmt.Printf("  Bar %d (%s) '%s': ", i+1, bar.Type, bar.Comment)
				for _, t := range bar.Tokens {
					fmt.Printf("%s ", t.Value)
				}
				i++
			}
			fmt.Println()
		}
	}
}

func main() {
	// Read the song file
	data, err := os.ReadFile("sledgehammer.nns")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	lexer := NewLexer(string(data))
	tokens := lexer.Lex()

	// for _, t := range tokens {
	// 	fmt.Printf("%s: %s\n", t.Type, t.Value)
	// }

	song := ParseSong(tokens)

	fmt.Println("Frontmatter:")
	for k, v := range song.FrontMatter {
		fmt.Printf("%s: %s\n", k, v)
	}
	RenderSongHTML(song, "index.html")
	printSong(song)
}
