package internal

import "fmt"

type Song struct {
	FrontMatter map[string]string
	Sections    []*Section
}

type Section struct {
	Name      string
	BarsLines [][]*Bar
	Break     bool
}

type Annotation struct {
	Value string
}
type Backtick struct {
	Id    int
	Value string
}
type Chord struct {
	Value      string
	Annotation *Annotation
}

type Bar struct {
	Tokens   []Token // chords, symbols, annotations, backticks
	Chords   []Chord
	Backtick Backtick
	Type     string // "Normal" or "DoubleBar"
	Comment  string // comment
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
	fmt.Println("Frontmatter:")
	for k, v := range song.FrontMatter {
		fmt.Printf("%s: %s\n", k, v)
	}
	i := 1
	for _, sec := range song.Sections {
		fmt.Println("Section:", sec.Name)
		for _, barline := range sec.BarsLines {
			for _, bar := range barline {
				fmt.Printf("  Bar %d (%s) '%s': ", i+1, bar.Type, bar.Comment)
				for _, t := range bar.Chords {
					fmt.Printf("Chord (%s): %s", t.Annotation.Value, t.Value)
				}
				// for _, t := range bar.Tokens {
				// 	fmt.Printf("%s ", t.Value)
				// }
				i++
			}
			fmt.Println()
		}
	}
}
