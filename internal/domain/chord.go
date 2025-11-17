package domain

import (
	"encoding/json"
)

type Chord struct {
	Value      string      `json:"value"`
	Annotation *Annotation `json:"annotation"`
}

type Annotation struct {
	Value string `json:"value"`
}

func (chord *Chord) PrettyPrint() string {
	return FormatChord(chord.Value)
}

func (chord *Chord) PrettyPrintHTML() string {
	return chord.PrettyPrint()
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
