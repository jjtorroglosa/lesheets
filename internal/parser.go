package internal

import (
	"fmt"
	"strings"
)

func ParseSong(tokens []Token) *Song {
	song := &Song{FrontMatter: make(map[string]string)}
	var currentSection *Section
	var currentBar *Bar
	var currentLine []*Bar
	currentAnnotation := &Annotation{}
	backtickId := 0
	var currentComment string
	inFrontMatter := false
	var pendingKey string

	for _, tok := range tokens {
		switch tok.Type {
		case TokenFrontMatterStart:
			inFrontMatter = true
		case TokenFrontMatterEnd:
			inFrontMatter = false
		case TokenYAMLKey:
			if inFrontMatter {
				pendingKey = tok.Value
			}
		case TokenYAMLValue:
			if inFrontMatter && pendingKey != "" {
				if strings.ToLower(pendingKey) == "key" {
					song.FrontMatter[pendingKey] = FormatChord(tok.Value)
				} else {
					song.FrontMatter[pendingKey] = tok.Value
				}
				pendingKey = ""
			}
		case TokenHeader:
			// start a new section
			currentSection = &Section{Name: tok.Value}
			song.Sections = append(song.Sections, currentSection)
			currentBar = nil
		case TokenHeaderBreak:
			// start a new section
			currentSection = &Section{Name: tok.Value, Break: true}
			song.Sections = append(song.Sections, currentSection)
			currentBar = nil
		case TokenAnnotation:
			currentAnnotation = &Annotation{Value: tok.Value}
		case TokenBacktick:
			if currentBar == nil {
				currentBar = &Bar{
					Tokens:   []Token{},
					Type:     "Normal",
					Chords:   []Chord{},
					Backtick: Backtick{Id: backtickId, Value: tok.Value},
				}
			}
			currentBar.Backtick = Backtick{Id: backtickId, Value: tok.Value}
			currentBar.Tokens = append(currentBar.Tokens, tok)
			backtickId++
		case TokenChord:
			if currentBar == nil {
				currentBar = &Bar{Tokens: []Token{}, Type: "Normal", Chords: []Chord{}}
			}
			currentBar.Chords = append(currentBar.Chords, Chord{
				Value:      tok.Value,
				Annotation: currentAnnotation,
			})
			currentAnnotation = &Annotation{}
			currentBar.Tokens = append(currentBar.Tokens, tok)
		case TokenComment:
			currentComment = tok.Value
		case TokenRepeatEnd:
			if currentBar != nil && currentSection != nil {
				currentBar.Comment = currentComment
				currentComment = ""
				currentBar.Type = "RepeatEnd"
				currentLine = append(currentLine, currentBar)
				currentBar = nil
			}
		case TokenRepeatStart:
			if currentBar != nil && currentSection != nil {
				currentBar.Comment = currentComment
				currentBar.Type = "Normal"
				currentComment = ""
				currentLine = append(currentLine, currentBar)
				currentBar = &Bar{Tokens: []Token{}, Type: "RepeatStart", Chords: []Chord{}, Backtick: Backtick{}}
			}
		case TokenBar:
			// finish current bar if exists
			if currentBar != nil && currentSection != nil {
				currentBar.Comment = currentComment
				currentComment = ""
				if tok.Value == "||" {
					currentBar.Type = "Double"
				}
				currentLine = append(currentLine, currentBar)
				currentBar = nil
			}
		case TokenReturn:
			// finish the current line of bars
			if len(currentLine) > 0 && currentSection != nil {
				currentSection.BarsLines = append(currentSection.BarsLines, currentLine)
				currentLine = nil
			}
		}
	}

	// flush remaining bar
	if currentBar != nil && !currentBar.IsEmpty() {
		fmt.Printf("is empty %v\n", currentBar.Chords[0].Value)
		if currentSection != nil {
			currentLine = append(currentLine, currentBar)
		}
	}

	// flush remaining line
	if len(currentLine) > 0 && currentSection != nil {
		currentSection.BarsLines = append(currentSection.BarsLines, currentLine)
	}

	return song
}
