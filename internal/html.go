package internal

import (
	"html/template"
	"log"
	"os"
)

func RenderSongHTML(song *Song, filename string) {
	t := template.Must(template.ParseGlob("views/*.html"))
	f, err := os.Create(filename)
	if err != nil {
		Fatalf("Failed to create HTML file: %v", filename)
	}
	defer f.Close()

	params := map[string]any{
		"Song": song,
		"Dev":  true,
	}
	if err := t.ExecuteTemplate(f, "tmpl.html", params); err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}
}
