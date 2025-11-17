package domain

import "lesheets/internal/svg"

type MultilineBacktick struct {
	Value         string `json:"value"`
	Id            int    `json:"id"`
	DefaultLength string `json:"default_length"`
	SourceFile    string `json:"source_file"`
}
type Backtick struct {
	Id            int    `json:"id"`
	Value         string `json:"value"`
	DefaultLength string `json:"default_length"`
}

func (mb *MultilineBacktick) Svg() string {
	html, err := svg.AbcToHtml(mb.SourceFile, mb.DefaultLength, mb.Value)
	if err != nil {
		return "<pre>Error rendering svg</pre>"
	}
	return html
}

func (backtick *Backtick) Svg() string {
	html, err := svg.InlineAbcToHtml("", backtick.DefaultLength, backtick.Value)
	if err != nil {
		return "<pre>Error rendering svg</pre>"
	}
	return html
}
