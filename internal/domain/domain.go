package domain

import (
	"encoding/json"
	"errors"
	"lesheets/internal/logger"
	"lesheets/internal/svg"

	"github.com/a-h/templ"
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

func (chord *Chord) PrettyPrintHTML() string {
	return chord.PrettyPrint()
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

func (song *Song) Key() templ.Component {
	key := song.FrontMatter["key"]
	return templ.Raw(FormatChord(key))
}

func (song *Song) PrintSong() {
	logger.Println("Frontmatter:")
	for k, v := range song.FrontMatter {
		logger.Printf("%s: %s\n", k, v)
	}
	i := 1
	for _, sec := range song.Sections {
		logger.Printf("Section: %s\n", sec.Name)
		for _, line := range sec.Lines {
			if line.MultilineBacktick.Value != "" {
				logger.Printf("MultilineBacktick: %s", line.MultilineBacktick.Value)
			} else {
				for _, bar := range line.Bars {
					logger.Printf("  Bar %d (%s) '%s': ", i+1, bar.Type, bar.BarNote)
					for _, t := range bar.Chords {
						logger.Printf("Chord (%s): %s", t.Annotation.Value, t.Value)
					}
					i++
				}
			}
			logger.Printf("\n")
		}
	}
}

func (song *Song) ToJson() (string, error) {
	j, err := json.MarshalIndent(song, "", "  ")
	if err != nil {
		return "", errors.New("error marshalling json: " + err.Error())
	}
	return string(j), nil
}

func (s *Song) DefaultLength() string {
	if s == nil {
		return "1/16"
	}
	defaultLength, ok := s.FrontMatter["L"]
	if !ok || defaultLength == "" {
		return "1/16"
	}
	return defaultLength
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
