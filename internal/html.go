package internal

import (
	"html/template"
	"log"
	"os"
)

func dict(values ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		m[key] = values[i+1]
	}
	return m
}

func RenderSongHTML(song *Song, filename string) {
	t := template.Must(template.New("").Funcs(template.FuncMap{
		"dict": dict,
	}).ParseGlob("views/*.html"))
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
