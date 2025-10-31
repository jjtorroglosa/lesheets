package internal

import "encoding/json"

type Song struct {
	FrontMatter map[string]string `json:"front_matter"`
	Sections    []Section         `json:"sections"`
}

type Section struct {
	Name      string  `json:"name"`
	BarsLines [][]Bar `json:"bars_lines"`
	Break     bool    `json:"break"`
}

type Annotation struct {
	Value string `json:"value"`
}
type Backtick struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
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
	Comment     string   `json:"comment"`
}

func (section *Section) IsEmpty() bool {
	return len(section.BarsLines) == 0 && section.Name == ""
}

func (bar *Bar) IsEmpty() bool {
	emptyChords := len(bar.Chords) == 0 ||
		(len(bar.Chords) == 1 && bar.Chords[0].Value == "")
	return emptyChords && bar.Backtick.Value == "" && bar.Comment == ""
}

func (chord *Chord) PrettyPrint() string {
	return FormatChord(chord.Value)
}

func (song *Song) Backticks() []Backtick {
	bts := []Backtick{}
	for _, sec := range song.Sections {
		for _, barline := range sec.BarsLines {
			for _, bar := range barline {
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
		for _, barline := range sec.BarsLines {
			for _, bar := range barline {
				Printf("  Bar %d (%s) '%s': ", i+1, bar.Type, bar.Comment)
				for _, t := range bar.Chords {
					Printf("Chord (%s): %s", t.Annotation.Value, t.Value)
				}
				// for _, t := range bar.Tokens {
				// 	Printf("%s ", t.Value)
				// }
				i++
			}
			Printf("\n")
		}
	}
}

func (song *Song) toJson() string {
	j, err := json.MarshalIndent(song, "", "  ")
	if err != nil {
		Fatalf("Error marshalling json: %v", err)
	}
	return string(j)
}
