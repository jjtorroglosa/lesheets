package internal

import (
	"encoding/json"
	"fmt"
	"html/template"
	"nasheets/internal/svg"
)

type Song struct {
	FrontMatter map[string]string `json:"front_matter"`
	Sections    []Section         `json:"sections"`
	Parser      *Parser           `json:"-"`
}

type Section struct {
	Name  string `json:"name"`
	Lines []Line `json:"lines"`
	Break bool   `json:"break"`
}

type Line struct {
	Bars              []Bar
	MultilineBacktick MultilineBacktick
}

type MultilineBacktick struct {
	Value         string
	Id            int
	DefaultLength string
	SourceFile    string
}

type Annotation struct {
	Value string `json:"value"`
}

type Backtick struct {
	Id            int    `json:"id"`
	Value         string `json:"value"`
	DefaultLength string `json:"default_length"`
}
type Chord struct {
	Value      string      `json:"value"`
	Annotation *Annotation `json:"annotation"`
}

func (p Chord) MarshalJSON() ([]byte, error) {
	type Alias Chord
	return json.Marshal(&struct {
		Alias
		Pretty string `json:"pretty"`
	}{
		Alias:  (Alias)(p),
		Pretty: p.PrettyPrint(),
	})
}

type Bar struct {
	Tokens               []Token  `json:"-"` // chords, symbols, annotations, backticks
	Chords               []Chord  `json:"chords"`
	Backtick             Backtick `json:"backtick"`
	Type                 string   `json:"type"`
	RepeatEnd            bool     `json:"repeat_end"`
	RepeatStart          bool     `json:"repeat_start"`
	DoubleBarEnd         bool     `json:"double_bar_end"`
	BarNote              string   `json:"bar_note"`
	Lyrics               string   `json:"lyrics"`
	Id                   int      `json:"id"`
	PreviousWasRepeatEnd bool     `json:"-"`
}

func (section *Section) IsEmpty() bool {
	return len(section.Lines) == 0 && section.Name == ""
}

func (bar *Bar) Number() int {
	return bar.Id + 1
}

func (bar *Bar) IsEmpty() bool {
	emptyChords := len(bar.Chords) == 0 ||
		(len(bar.Chords) == 1 && bar.Chords[0].Value == "")
	return emptyChords && bar.Backtick.Value == "" && bar.BarNote == ""
}

func (chord *Chord) PrettyPrint() string {
	return FormatChord(chord.Value)
}

func (chord *Chord) PrettyPrintHTML() template.HTML {
	return template.HTML(chord.PrettyPrint())
}

func (song *Song) Backticks() []Backtick {
	bts := []Backtick{}
	for _, sec := range song.Sections {
		for _, line := range sec.Lines {
			for _, bar := range line.Bars {
				if bar.Backtick.Value != "" {
					bts = append(bts, bar.Backtick)
				}
			}
		}
	}
	return bts
}

func (song *Song) Key() template.HTML {
	key := song.FrontMatter["key"]
	return template.HTML(FormatChord(key))
}

func (song *Song) PrintSong() {
	Println("Frontmatter:")
	for k, v := range song.FrontMatter {
		Printf("%s: %s\n", k, v)
	}
	i := 1
	for _, sec := range song.Sections {
		Printf("Section: %s\n", sec.Name)
		for _, line := range sec.Lines {
			if line.MultilineBacktick.Value != "" {
				Printf("MultilineBacktick: %s", line.MultilineBacktick.Value)
			} else {
				for _, bar := range line.Bars {
					Printf("  Bar %d (%s) '%s': ", i+1, bar.Type, bar.BarNote)
					for _, t := range bar.Chords {
						Printf("Chord (%s): %s", t.Annotation.Value, t.Value)
					}
					i++
				}
			}
			Printf("\n")
		}
	}
}

func (song *Song) ToJson() (string, error) {
	j, err := json.MarshalIndent(song, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling json: %w", err)
	}
	return string(j), nil
}

func (mb *MultilineBacktick) Svg() template.HTML {
	html, err := svg.AbcToHtml(mb.SourceFile, mb.DefaultLength, mb.Value)
	if err != nil {
		return template.HTML("<pre>Error rendering svg</pre>")
	}
	return template.HTML(html)
}

func (backtick *Backtick) Svg() template.HTML {
	html, err := svg.InlineAbcToHtml("", backtick.DefaultLength, backtick.Value)
	if err != nil {
		return template.HTML("<pre>Error rendering svg</pre>")
	}
	return html
}
func (a *Annotation) Symbol() template.HTML {
	switch a.Value {
	case "marcato":
		return `<div class="font-bold relative top-[4px] leading-[1.3] text-[1rem]/2 font-music">^</div>`
	case "push":
		return `<span class="text-[10px]/[1.25rem]">❮</span>`
	case "pull", "hold":
		return `<span class="text-[10px]/[1.25rem]">❯</span>`
	case "fermata":
		return `<div class="font-music text-xl leading-none"></div>`
	case "diamond-fermata":
		return `<div class="font-music font-size text-xl leading-none"></div>`
	}
	return ""
}
