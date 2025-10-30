package internal

import (
	"strings"
)

func ParseSong(s string) *Song {
	lexer := NewLexer(s)
	tokens := lexer.Lex()
	return ParseTokens(tokens)
}

func ParseTokens(tokens []Token) *Song {
	song := &Song{FrontMatter: make(map[string]string)}
	var currentSection *Section
	var currentBar *Bar
	var currentLine []*Bar
	currentAnnotation := &Annotation{}
	backtickId := 0
	var currentComment string
	inFrontMatter := false
	var pendingKey string
	currentSection = &Section{
		Name:      "",
		BarsLines: [][]*Bar{},
		Break:     false,
	}
	song.Sections = append(song.Sections, currentSection)

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
		case TokenHeader, TokenHeaderBreak:
			// start a new section
			currentSection = &Section{Name: tok.Value}
			if tok.Type == TokenHeaderBreak {
				currentSection.Break = true
			}
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
				currentBar = &Bar{
					Tokens:      []Token{},
					Type:        "Normal",
					Chords:      []Chord{},
					Backtick:    Backtick{},
					RepeatEnd:   false,
					RepeatStart: false,
					Comment:     "",
				}
			}
			currentBar.Chords = append(currentBar.Chords, Chord{
				Value:      tok.Value,
				Annotation: currentAnnotation,
			})
			currentAnnotation = &Annotation{}
			currentBar.Tokens = append(currentBar.Tokens, tok)
		case TokenComment:
			currentComment = tok.Value
		case TokenRepeatEnd, TokenRepeatStart, TokenBar:
			if tok.Type == TokenRepeatStart && currentBar == nil {
				currentBar = &Bar{
					Tokens:      []Token{},
					Type:        "RepeatStart",
					Chords:      []Chord{},
					Backtick:    Backtick{},
					RepeatEnd:   false,
					RepeatStart: true,
					Comment:     "",
				}
			} else if currentBar != nil && currentSection != nil {
				currentBar.Comment = currentComment
				currentComment = ""
				switch tok.Type {
				case TokenBar, TokenRepeatStart:
					currentBar.Type = "Normal"
					if tok.Value == "||" {
						currentBar.Type = "Double"
					}
				case TokenRepeatEnd:
					currentBar.RepeatEnd = true
					currentBar.Type = "RepeatEnd"
				}
				if !currentBar.IsEmpty() {
					currentLine = append(currentLine, currentBar)
				}
				currentBar = nil
				if tok.Type == TokenRepeatStart {
					// Initialize the next bar
					currentBar = &Bar{
						Tokens:      []Token{},
						Type:        "RepeatStart",
						Chords:      []Chord{},
						Backtick:    Backtick{},
						RepeatEnd:   false,
						RepeatStart: true,
						Comment:     "",
					}
				}
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
		currentLine = append(currentLine, currentBar)
	}

	// flush remaining line
	if len(currentLine) > 0 && currentSection != nil {
		currentSection.BarsLines = append(currentSection.BarsLines, currentLine)
	}

	return song
}
