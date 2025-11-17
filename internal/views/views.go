package views

import (
	"bytes"
	"context"
	"lesheets/internal/domain"
)

type Link struct {
	Name string
	Href string
}

func RenderSong(song *domain.Song, abc string, cfg RenderConfig, buf *bytes.Buffer) error {
	component := base(song, abc, cfg)
	err := component.Render(context.Background(), buf)
	if err != nil {
		return err
	}
	return nil
}

func RenderListOfFiles(files []Link, buf *bytes.Buffer) error {
	component := list(files)
	err := component.Render(context.Background(), buf)
	if err != nil {
		return err
	}
	return nil
}
