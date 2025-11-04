package internal

import (
	"encoding/json"
	"html/template"
	"nasheets/internal/svg"
)

type Song struct {
	FrontMatter map[string]string `json:"front_matter"`
	Sections    []Section         `json:"sections"`
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
	Tokens      []Token  `json:"-"` // chords, symbols, annotations, backticks
	Chords      []Chord  `json:"chords"`
	Backtick    Backtick `json:"backtick"`
	Type        string   `json:"type"`
	RepeatEnd   bool     `json:"repeat_end"`
	RepeatStart bool     `json:"repeat_start"`
	BarNote     string   `json:"bar_note"`
}

func (section *Section) IsEmpty() bool {
	return len(section.Lines) == 0 && section.Name == ""
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

func (song *Song) ToJson() string {
	j, err := json.MarshalIndent(song, "", "  ")
	if err != nil {
		Fatalf("Error marshalling json: %v", err)
	}
	return string(j)
}

func (mb *MultilineBacktick) Svg() template.HTML {
	return svg.AbcToHtml(mb.DefaultLength, mb.Value)
}

func (backtick *Backtick) Svg() template.HTML {
	return svg.InlineAbcToHtml(backtick.DefaultLength, backtick.Value)
}
