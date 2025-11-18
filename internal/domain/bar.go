package domain

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

func (bar *Bar) Number() int {
	return bar.Id + 1
}

func (bar *Bar) IsEmpty() bool {
	emptyChords := len(bar.Chords) == 0 ||
		(len(bar.Chords) == 1 && bar.Chords[0].Value == "")
	return emptyChords && bar.Backtick.Value == "" && bar.BarNote == ""
}
