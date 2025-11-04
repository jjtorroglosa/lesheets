package internal

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"nasheets/internal/timer"
	"os"
)

//go:embed views/*.html
var templateFS embed.FS

func dict(values ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		m[key] = values[i+1]
	}
	return m
}

var t *template.Template

func init() {
	defer timer.LogElapsedTime("InitTmpl")()
	t = template.Must(template.New("").Funcs(template.FuncMap{
		"dict": dict,
	}).ParseFS(templateFS, "views/*.html"))
}

func RenderSongHTML(dev bool, song *Song, filename string) {
	defer timer.LogElapsedTime("RenderHtml")()

	f, err := os.Create(filename)
	if err != nil {
		Fatalf("Failed to create HTML file: %v", filename)
	}
	defer f.Close()
	params := map[string]any{
		"Song": song,
		"Dev":  dev,
	}
	var buf bytes.Buffer

	func() {
		if err := t.ExecuteTemplate(&buf, "tmpl.html", params); err != nil {
			log.Fatalf("Failed to render template: %v", err)
		}
	}()
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatalf("Write error: %v", err)
	}
}
