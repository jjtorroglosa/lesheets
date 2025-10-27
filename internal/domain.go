package internal

import "fmt"

type Song struct {
	FrontMatter map[string]string
	Sections    []*Section
}

type Section struct {
	Header    string
	BarsLines [][]*Bar
	Break     bool
}

type Annotation struct {
	Value string
}
type Chord struct {
	Value      string
	Annotation *Annotation
}

type Bar struct {
	Tokens  []Token // chords, symbols, annotations, backticks
	Chords  []Chord // chords, symbols, annotations, backticks
	Type    string  // "Normal" or "DoubleBar"
	Comment string  // comment
}

func (song *Song) PrintSong() {
	fmt.Println("Frontmatter:")
	for k, v := range song.FrontMatter {
		fmt.Printf("%s: %s\n", k, v)
	}
	for _, sec := range song.Sections {
		fmt.Println("Section:", sec.Header)
		i := 1
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
