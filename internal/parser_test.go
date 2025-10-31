package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrontmatter(t *testing.T) {
	p := NewParser(NewLexer(`---
title: something
anotherthing: here
---
A|B|
`))
	song, err := p.ParseSong()
	assert.NoError(t, err)
	fm := song.FrontMatter
	assert.Equal(t, "something", fm["title"])
	assert.Equal(t, "here", fm["anotherthing"])
}

func TestSongSection(t *testing.T) {
	song, err := ParseSongFromString("# a section\nA")
	assert.NoError(t, err)
	assert.False(t, song.Sections[0].IsEmpty())
	assert.Equal(t, 1, len(song.Sections))
	assert.Equal(t, "a section", song.Sections[0].Name)
}

func TestSongWithUnnamedSection(t *testing.T) {
	song, _ := ParseSongFromString("Amaj7 C | B")
	assert.Equal(t, "", song.Sections[0].Name)
	assert.Equal(t, "Amaj7", song.Sections[0].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "C", song.Sections[0].BarsLines[0][0].Chords[1].Value)
	assert.Equal(t, "B", song.Sections[0].BarsLines[0][1].Chords[0].Value)
}

func TestSongWithUnnamedSectionAndNamedSection(t *testing.T) {
	song, _ := ParseSongFromString("A | B \n# verse\nD")
	assert.Equal(t, "", song.Sections[0].Name)
	assert.Equal(t, "A", song.Sections[0].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "B", song.Sections[0].BarsLines[0][1].Chords[0].Value)
	assert.Equal(t, "verse", song.Sections[1].Name)
	assert.Equal(t, "D", song.Sections[1].BarsLines[0][0].Chords[0].Value)
}

func TestSectionWithEmptyBody(t *testing.T) {
	p := NewParser(NewLexer("# a section\n"))
	section, _ := p.ParseSection()
	assert.False(t, section.IsEmpty())
	assert.Equal(t, "a section", section.Name)
}

func TestSongEmptySection(t *testing.T) {
	song, _ := ParseSongFromString("\n#verse\nAmaj7 C | B")
	assert.Equal(t, "verse", song.Sections[0].Name)
	assert.False(t, song.Sections[0].IsEmpty())
	assert.Equal(t, "Amaj7", song.Sections[0].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "C", song.Sections[0].BarsLines[0][0].Chords[1].Value)
	assert.Equal(t, "B", song.Sections[0].BarsLines[0][1].Chords[0].Value)
}

func TestSongNotEmptySectionIfItHasAName(t *testing.T) {
	song, _ := ParseSongFromString("\n#verse\n")
	assert.Equal(t, 1, len(song.Sections))
	assert.Equal(t, "verse", song.Sections[0].Name)
	assert.False(t, song.Sections[0].IsEmpty())
}

func TestSections(t *testing.T) {
	p := NewParser(NewLexer(`
A|B
Bmin
# section1
C
# section2
D|E
`))
	sections, err := p.ParseBody()
	assert.NoError(t, err)

	assert.Equal(t, 3, len(sections))
	assert.Equal(t, "", sections[0].Name)
	assert.Equal(t, "A", sections[0].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "B", sections[0].BarsLines[0][1].Chords[0].Value)
	assert.Equal(t, "Bmin", sections[0].BarsLines[1][0].Chords[0].Value)
	assert.Equal(t, "section1", sections[1].Name)
	assert.Equal(t, "C", sections[1].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "section2", sections[2].Name)
	assert.Equal(t, "D", sections[2].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "E", sections[2].BarsLines[0][1].Chords[0].Value)
}
func TestSectionBreak(t *testing.T) {
	p := NewParser(NewLexer(`
A|B
Bmin
# section1
C
#- section2
D|E
`))
	sections, err := p.ParseBody()
	assert.NoError(t, err)

	assert.Equal(t, "", sections[0].Name)
	assert.False(t, sections[0].Break)
	assert.Equal(t, "section1", sections[1].Name)
	assert.False(t, sections[1].Break)
	assert.Equal(t, "section2", sections[2].Name)
	assert.True(t, sections[2].Break)
}

func TestParseBar(t *testing.T) {
	p := NewParser(NewLexer("!first!Cmaj7 !annotation!D !third!E\n!fourth!F"))
	bar, err := p.ParseBar()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(bar.Chords))
	assert.Equal(t, "Cmaj7", bar.Chords[0].Value)
	assert.Equal(t, "first", bar.Chords[0].Annotation.Value)
	assert.Equal(t, "D", bar.Chords[1].Value)
	assert.Equal(t, "annotation", bar.Chords[1].Annotation.Value)
	assert.Equal(t, "E", bar.Chords[2].Value)
	assert.Equal(t, "third", bar.Chords[2].Annotation.Value)
}

func TestParseBarWithNoChords(t *testing.T) {

	p := NewParser(NewLexer("\"comment\""))
	bar, err := p.ParseBar()
	assert.Nil(t, bar)
	assert.Error(t, err, "expected chord or backtick expression at pos 9")
}

func TestSongRepeatStart(t *testing.T) {
	song, _ := ParseSongFromString("||: A :|| B |")
	assert.Equal(t, "A", song.Sections[0].BarsLines[0][0].Chords[0].Value)
	assert.True(t, song.Sections[0].BarsLines[0][0].RepeatStart)
	assert.True(t, song.Sections[0].BarsLines[0][0].RepeatEnd)

	assert.Equal(t, "B", song.Sections[0].BarsLines[0][1].Chords[0].Value)
	assert.False(t, song.Sections[0].BarsLines[0][1].RepeatStart)
	assert.False(t, song.Sections[0].BarsLines[0][1].RepeatEnd)
}

func TestSongTwoConsecutiveRepeats(t *testing.T) {
	song, err := ParseSongFromString("||: A :|| ||: B |")
	assert.NoError(t, err)
	assert.Equal(t, "A", song.Sections[0].BarsLines[0][0].Chords[0].Value)
	assert.True(t, song.Sections[0].BarsLines[0][0].RepeatStart)
	assert.True(t, song.Sections[0].BarsLines[0][0].RepeatEnd)

	assert.Equal(t, "B", song.Sections[0].BarsLines[0][1].Chords[0].Value)
	assert.True(t, song.Sections[0].BarsLines[0][1].RepeatStart)
	assert.False(t, song.Sections[0].BarsLines[0][1].RepeatEnd)
}

func TestParseBarRepeatEnd(t *testing.T) {
	p := NewParser(NewLexer("C :||"))
	bar, err := p.ParseBar()
	assert.NoError(t, err)
	assert.Equal(t, "C", bar.Chords[0].Value)
	assert.False(t, bar.RepeatStart)
	assert.True(t, bar.RepeatEnd)
}

func TestParseBarRepeatStart(t *testing.T) {

	p := NewParser(NewLexer("||: C"))
	bar, err := p.ParseBar()
	assert.NoError(t, err)
	assert.True(t, bar.RepeatStart)
	assert.False(t, bar.RepeatEnd)
	assert.Equal(t, "C", bar.Chords[0].Value)
}

func TestParseBarsLineEmpty(t *testing.T) {

	p := NewParser(NewLexer("\n\n"))
	lineP, err := p.ParseBarsLine()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(*lineP))
}

func TestParseBarsLine(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{input: "Cmaj7|D\nC"},
		{input: "Cmaj7|D|\nC"},
	}

	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {

			p := NewParser(NewLexer("Cmaj7 | !annotation!D\nC"))
			lineP, err := p.ParseBarsLine()
			assert.NoError(t, err)
			line := *lineP
			assert.Equal(t, 2, len(line))
			assert.Equal(t, "Cmaj7", line[0].Chords[0].Value)
			assert.Equal(t, "D", line[1].Chords[0].Value)
		})
	}
}

func TestParseBarReturn(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{input: "C\nF"},
		{input: "C|\nF"},
		{input: "C\nF|\n"},
		{input: "C\n|F|\n"},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lex := NewLexer(tC.input)

			p := NewParser(lex)
			bar, err := p.ParseBarsLine()
			assert.NoError(t, err)
			assert.Equal(t, 1, len(*bar))
		})
	}
}

func TestParseLines(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{input: "C\nF"},
		{input: "C|\nF"},
		{input: "C\nF|\n"},
		{input: "C\n|F|\n"},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			p := NewParser(NewLexer("C\nF|\n"))
			lines, err := p.ParseLines()
			assert.NoError(t, err)
			assert.Equal(t, "C", lines[0][0].Chords[0].Value)
			assert.Equal(t, "F", lines[1][0].Chords[0].Value)
		})
	}
}

func TestParseBarWithComment(t *testing.T) {
	p := NewParser(NewLexer("\"any comment\"Cmaj7 | \"another comment\"D\nC"))
	barsP, err := p.ParseBarsLine()
	assert.NoError(t, err)
	bars := *barsP
	assert.Equal(t, "any comment", bars[0].Comment)
	assert.Equal(t, 2, len(bars))
	assert.Equal(t, "Cmaj7", bars[0].Chords[0].Value)
	assert.Equal(t, "D", bars[1].Chords[0].Value)
	assert.Equal(t, "another comment", bars[1].Comment)
}

func TestParseBarWithCommentInASeparateLine(t *testing.T) {
	p := NewParser(NewLexer("\"any comment\"\nCmaj7"))
	barsP, err := p.ParseBarsLine()
	assert.NoError(t, err)
	bars := *barsP
	assert.Equal(t, "any comment", bars[0].Comment)
	assert.Equal(t, 1, len(bars))
}

func TestParseChord(t *testing.T) {
	p := NewParser(NewLexer("!annotation!D"))
	chord, err := p.ParseChord()
	assert.NoError(t, err)
	assert.Equal(t, "D", chord.Value)
	assert.Equal(t, "annotation", chord.Annotation.Value)
}

func TestPrettyPrint(t *testing.T) {
	song, _ := ParseSongFromString("Amaj7(#11)|B|")
	assert.Equal(t, "A△⁷(♯¹¹)", song.Sections[0].BarsLines[0][0].Chords[0].PrettyPrint())
}
