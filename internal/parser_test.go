package internal

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSection(t *testing.T) {
	song := ParseSong("# a section")
	assert.True(t, song.Sections[0].IsEmpty())
	assert.Equal(t, "a section", song.Sections[1].Name)
}

func TestEmptySection(t *testing.T) {
	song := ParseSong("\n#verse\nAmaj7 C | B")
	assert.Equal(t, "", song.Sections[0].Name)
	assert.True(t, song.Sections[0].IsEmpty())
	assert.Equal(t, "Amaj7", song.Sections[1].BarsLines[0][0].Chords[0].Value)
}

func TestNotEmptySectionIfItHasAName(t *testing.T) {
	song := ParseSong("\n#verse\n")
	assert.Equal(t, "", song.Sections[0].Name)
	assert.True(t, song.Sections[0].IsEmpty())
	assert.Equal(t, "verse", song.Sections[1].Name)
	assert.False(t, song.Sections[1].IsEmpty())
}

func TestRepeatStart(t *testing.T) {
	song := ParseSong("||: A :|| B |")
	assert.True(t, song.Sections[0].BarsLines[0][0].RepeatStart)
	assert.True(t, song.Sections[0].BarsLines[0][0].RepeatEnd)
	assert.Equal(t, "A", song.Sections[0].BarsLines[0][0].Chords[0].Value)

	assert.False(t, song.Sections[0].BarsLines[0][1].RepeatStart)
	assert.False(t, song.Sections[0].BarsLines[0][1].RepeatEnd)
}

func TestUnnamedSection(t *testing.T) {
	song := ParseSong("Amaj7 C | B")
	assert.Equal(t, "", song.Sections[0].Name)
	assert.Equal(t, "Amaj7", song.Sections[0].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "C", song.Sections[0].BarsLines[0][0].Chords[1].Value)
	assert.Equal(t, "B", song.Sections[0].BarsLines[0][1].Chords[0].Value)
}

func TestUnnamedSection2(t *testing.T) {
	song := ParseSong("A | B |\n# verse\nD")
	log.Println(song.toJson())
	assert.Equal(t, "", song.Sections[0].Name)
	assert.Equal(t, "A", song.Sections[0].BarsLines[0][0].Chords[0].Value)
	assert.Equal(t, "B", song.Sections[0].BarsLines[0][1].Chords[0].Value)
	assert.Equal(t, "verse", song.Sections[1].Name)
	assert.Equal(t, "D", song.Sections[1].BarsLines[0][0].Chords[0].Value)
}

func TestFrontmatter(t *testing.T) {
	song := ParseSong(`
---
title: something
anotherthing: here
---
A|B|
`)
	assert.Equal(t, "something", song.FrontMatter["title"])
	assert.Equal(t, "here", song.FrontMatter["anotherthing"])
}

func TestPrettyPrint(t *testing.T) {
	song := ParseSong("Amaj7(#11)|B|")
	assert.Equal(t, "A△⁷(♯¹¹)", song.Sections[0].BarsLines[0][0].Chords[0].PrettyPrint())
}

func TestFinishBarAndLineWhenNewLine(t *testing.T) {
	song := ParseSong("A\nB")
	assert.Equal(t, "B", song.Sections[0].BarsLines[1][0].Chords[0].Value)
}

func TestComment(t *testing.T) {
	song := ParseSong("\"some comment\" A|B|")
	assert.Equal(t, "some comment", song.Sections[0].BarsLines[0][0].Comment)
}

func TestCommentSecondBar(t *testing.T) {
	song := ParseSong("A|\"second bar\" B|")
	assert.Equal(t, "second bar", song.Sections[0].BarsLines[0][1].Comment)
}
