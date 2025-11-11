package internal

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"nasheets/internal/timer"
	"os"
	"path/filepath"
	"strings"
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

var templ *template.Template

func init() {
	defer timer.LogElapsedTime("InitTmpl")()
	funcs := template.FuncMap{
		"dict": dict,
	}
	templ = template.Must(template.New("").Funcs(funcs).ParseFS(templateFS, "views/*.html"))
}

func RenderListHTML(inputFiles []string) {
	defer timer.LogElapsedTime("RenderList")()
	filename := "output/index.html"
	f, err := os.Create(filename)
	if err != nil {
		Fatalf("Failed to create HTML file: %v", filename)
	}

	type Link struct {
		Name string
		Href string
	}
	files := []Link{}
	for _, i := range inputFiles {
		name := strings.TrimSuffix(i, ".nns")
		href := name + ".html"
		href = filepath.Dir(i) + "/" + filepath.Base(href)

		files = append(files, Link{
			Name: name,
			Href: href,
		})
	}
	files = append(files, Link{
		Name: "editor.html",
		Href: "editor.html",
	})

	defer f.Close()
	var buf bytes.Buffer
	if err := templ.ExecuteTemplate(&buf, "list.html", files); err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		log.Fatalf("Write error: %v", err)
	}
}

type RenderConfig struct {
	WithLiveReload bool
	WholeHtml      bool
	WithEditor     bool
}

func RenderSongHtml(cfg RenderConfig, sourceCode string, song *Song, filename string) string {
	defer timer.LogElapsedTime("RenderHtml")()

	params := map[string]any{
		"Song":   song,
		"Dev":    cfg.WithLiveReload,
		"Abc":    sourceCode,
		"Whole":  cfg.WholeHtml,
		"Editor": cfg.WithEditor,
	}

	var buf bytes.Buffer
	tmpl := "base.html"

	func() {
		if err := templ.ExecuteTemplate(&buf, tmpl, params); err != nil {
			log.Fatalf("Failed to render template: %v", err)
		}
	}()
	res := buf.String()
	return res
}

func WriteEditorToHtmlFile(dev bool, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		Fatalf("Failed to create HTML file: %v", filename)
	}
	defer f.Close()
	htmlOut := []byte(RenderSongHtml(RenderConfig{
		WithLiveReload: dev,
		WholeHtml:      true,
		WithEditor:     true,
	}, "", nil, filename))
	if err := os.WriteFile(filename, htmlOut, 0644); err != nil {
		log.Fatalf("Write error: %v", err)
	}
}

func WriteSongHtmlToFile(dev bool, sourceCode string, song *Song, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		Fatalf("Failed to create HTML file: %v", filename)
	}
	defer f.Close()
	htmlOut := []byte(RenderSongHtml(RenderConfig{
		WithLiveReload: dev,
		WholeHtml:      true,
		WithEditor:     false,
	}, sourceCode, song, filename))
	if err := os.WriteFile(filename, htmlOut, 0644); err != nil {
		log.Fatalf("Write error: %v", err)
	}
}
