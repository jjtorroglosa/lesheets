package views

import (
	"bytes"
	"context"
)

func RenderSong(data BaseData, buf *bytes.Buffer) error {
	component := base(data)
	err := component.Render(context.Background(), buf)
	if err != nil {
		return err
	}
	return nil
}
