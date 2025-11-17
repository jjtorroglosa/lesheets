package internal

import (
	"bytes"
	"fmt"
	"lesheets/internal/domain"
	"lesheets/internal/timer"
	"lesheets/internal/views"
	"os"
	"path/filepath"
	"strings"
)

func dict(values ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		m[key] = values[i+1]
	}
	return m
}

func RenderListHTML(inputFiles []string) error {
	defer timer.LogElapsedTime("RenderList")()
	filename := "output/index.html"
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	files := []views.Link{}
	for _, i := range inputFiles {
		name := strings.TrimSuffix(i, ".nns")
		href := name + ".html"
		href = filepath.Dir(i) + "/" + filepath.Base(href)

		files = append(files, views.Link{
			Name: name,
			Href: href,
		})
	}
	files = append(files, views.Link{
		Name: "editor.html",
		Href: "editor.html",
	})

	defer f.Close()
	var buf bytes.Buffer
	views.RenderListOfFiles(files, &buf)
	// if err := templ.ExecuteTemplate(&buf, "list.html", files); err != nil {
	// 	return err
	// }
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

type RenderConfig struct {
	WithLiveReload bool
	WholeHtml      bool
	WithEditor     bool
}

func RenderSongHtml(cfg RenderConfig, sourceCode string, song *domain.Song, filename string) (string, error) {
	defer timer.LogElapsedTime("RenderHtml")()

	var buf bytes.Buffer

	// if err := templ.ExecuteTemplate(&buf, tmpl, params); err != nil {
	// 	return "", fmt.Errorf("failed to render template: %w", err)
	// }
	err := views.RenderSong(views.BaseData{
		Song:   song,
		Dev:    cfg.WithLiveReload,
		Abc:    sourceCode,
		Whole:  cfg.WholeHtml,
		Editor: cfg.WithEditor,
	}, &buf)
	if err != nil {
		return "", err
	}

	res := buf.String()
	return res, nil
}

func WriteEditorToHtmlFile(dev bool, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create filename %s: %w", filename, err)
	}
	defer f.Close()
	htmlOut, err := RenderSongHtml(RenderConfig{
		WithLiveReload: dev,
		WholeHtml:      true,
		WithEditor:     true,
	}, "", nil, filename)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, []byte(htmlOut), 0644); err != nil {
		return err
	}
	return nil
}

func WriteSongHtmlToFile(dev bool, sourceCode string, song *domain.Song, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %s", filename)
	}
	defer f.Close()
	htmlOut, err := RenderSongHtml(RenderConfig{
		WithLiveReload: dev,
		WholeHtml:      true,
		WithEditor:     false,
	}, sourceCode, song, filename)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, []byte(htmlOut), 0644); err != nil {
		return err
	}
	return nil
}

func RenderError(err error) string {
	buf := bytes.Buffer{}
	// template.HTMLEscape(&buf, []byte(err.Error()))
	return "<pre>" + buf.String() + "</pre>"
}
