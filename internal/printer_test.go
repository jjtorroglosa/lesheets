package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintFrontmatter(t *testing.T) {
	input := `---
some: value
some other: diffvalue
---
`
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	assert.NoError(t, err)
	assert.Equal(t, input, output)
}

func TestPrintSection(t *testing.T) {
	input := "\n# section\n\n"
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	assert.Equal(t, input, output)
}

func TestPrintOneLine(t *testing.T) {
	input := "A | Bmaj7\n"
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	assert.Equal(t, input, output)
}

func TestPrintTwoLines(t *testing.T) {
	input := "A | Bmaj7\nD7(b13) | !annotation!F\n"
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	assert.Equal(t, input, output)
}

func TestPrintBacktick(t *testing.T) {
	input := "A | `backtick`\n"
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	assert.Equal(t, input, output)
}

func TestPrintMultilineBacktick(t *testing.T) {
	input := "```\nsomething\n```\n\n"
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	assert.Equal(t, input, output)
}

func TestPrintAll(t *testing.T) {
	bytes, err := os.ReadFile("testdata/all-features.nns")
	input := string(bytes)
	assert.NoError(t, err)
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	assert.NoError(t, err)
	assert.Equal(t, input, output)
}

func TestReflexive(t *testing.T) {
	bytes, err := os.ReadFile("testdata/all-features.nns")
	input := string(bytes)
	assert.NoError(t, err)
	s, err := ParseSongFromString(input)
	assert.NoError(t, err)
	output := PrintLesheet(s)
	s2, err := ParseSongFromString(output)
	assert.NoError(t, err)
	assert.Equal(t, s, s2)
}
