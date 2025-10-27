package internal

import "log"

func ParseSong(tokens []Token) *Song {
	song := &Song{FrontMatter: make(map[string]string)}
	var currentSection *Section
	var currentBar *Bar
	var currentLine []*Bar
	currentAnnotation := &Annotation{}
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
				song.FrontMatter[pendingKey] = tok.Value
				pendingKey = ""
			}
		case TokenHeader:
			// start a new section
			currentSection = &Section{Header: tok.Value}
			song.Sections = append(song.Sections, currentSection)
			currentBar = nil
		case TokenHeaderBreak:
			// start a new section
			currentSection = &Section{Header: tok.Value, Break: true}
			song.Sections = append(song.Sections, currentSection)
			currentBar = nil
		case TokenAnnotation:
			log.Printf("annotation %s\n", tok.Value)
			currentAnnotation = &Annotation{Value: tok.Value}
		case TokenBacktick:
			if currentBar == nil {
				currentBar = &Bar{Tokens: []Token{}, Type: "Normal", Chords: []Chord{}}
			}
			currentBar.Tokens = append(currentBar.Tokens, tok)
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
	if currentBar != nil {
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
