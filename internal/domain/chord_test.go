package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChord(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{in: "Cm", out: "C<small>m</small>"},
		{in: "Cm7", out: "C<small>m</small>⁷"},
		{in: "F#min11", out: "F♯<small>m</small>¹¹"},
		{in: "Bbmaj7", out: "B♭△⁷"},
		{in: "Cdim7", out: "C<sup>o</sup>⁷"},
		{in: "Ehalfdim7", out: "E<sup>ø</sup>⁷"},
		{in: "G7b9", out: "G⁷♭⁹"},
		{in: "G7(b9)", out: "G⁷<small>(♭⁹)</small>"},
		{in: "1", out: "1"},
		{in: "b1", out: "♭1"},
		{in: "1sus4", out: "1ˢᵘˢ⁴"},
		{in: "1/2", out: "1<span class=\"over\">/2</span>"},
		{in: "2/3", out: "2<span class=\"over\">/3</span>"},
		{in: "3/4", out: "3<span class=\"over\">/4</span>"},
		{in: "4/5", out: "4<span class=\"over\">/5</span>"},
		{in: "5/6", out: "5<span class=\"over\">/6</span>"},
		{in: "6/7", out: "6<span class=\"over\">/7</span>"},
		{in: "7/1", out: "7<span class=\"over\">/1</span>"},
		{in: "1/2maj7", out: "1<span class=\"over\">/2△7</span>"},
		{in: "1/2min7", out: "1<span class=\"over\">/2<small>m</small>7</span>"},
		{in: "1/b2min7", out: "1<span class=\"over\">/♭2<small>m</small>7</span>"},
	}

	for _, tC := range testCases {
		t.Run(tC.in, func(t *testing.T) {
			assert.Equal(t, tC.out, FormatChord(tC.in))
		})
	}
}
