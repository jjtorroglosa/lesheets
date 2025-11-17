package views

import (
	"bytes"
	"context"
)

type Link struct {
	Name string
	Href string
}

func RenderSong(data BaseData, buf *bytes.Buffer) error {
	component := base(data)
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
