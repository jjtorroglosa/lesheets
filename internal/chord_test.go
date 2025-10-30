package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChord(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{in: "Cm", out: "Cₘ"},
		{in: "Cm7", out: "Cₘ⁷"},
		{in: "F#min11", out: "F♯ₘ¹¹"},
		{in: "Bbmaj7", out: "B♭△⁷"},
		{in: "Cdim7", out: "C°⁷"},
		{in: "Ehalfdim7", out: "Eø⁷"},
		{in: "G7b9", out: "G⁷♭⁹"},
		{in: "G7(b9)", out: "G⁷(♭⁹)"},
		{in: "1", out: "1"},
		{in: "b1", out: "♭1"},
		{in: "1sus4", out: "1ˢᵘˢ⁴"},
	}

	for _, tC := range testCases {
		t.Run(tC.in, func(t *testing.T) {
			assert.Equal(t, tC.out, FormatChord(tC.in))
		})
	}
}
