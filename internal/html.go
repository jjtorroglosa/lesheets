package internal

import (
	"html/template"
	"log"
	"os"
)

// Template data structures (reuse your Song/Section/Bar/Token structs)

func RenderSongHTML(song *Song, filename string) {
	t := template.Must(template.ParseFiles("views/tmpl.html"))
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create HTML file: %v", err)
	}
	defer f.Close()

	if err := t.Execute(f, song); err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}
}
