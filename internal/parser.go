package internal

import (
	"fmt"
	"nasheets/internal/timer"
	"os"

	"gopkg.in/yaml.v3"
)

type Parser struct {
	Lexer              *Lexer
	backtickId         int
	mutilineBacktickId int
	song               *Song
	barsCount          int
}

func (p *Parser) SourceFile() string {
	return p.Lexer.source
}

func (s *Song) DefaultLength() string {
	if s == nil {
		return "1/16"
	}
	defaultLength, ok := s.FrontMatter["L"]
	if !ok || defaultLength == "" {
		return "1/16"
	}
	return defaultLength
}

func NewParser(lex *Lexer) *Parser {
	return &Parser{
		Lexer:      lex,
		backtickId: 0,
	}
}

func ParseSongFromString(s string) (*Song, error) {
	return NewParser(NewLexerFromSource("unknown", s)).ParseSong()
}

func ParseSongFromStringWithFileName(filename string, sourceCode string) (*Song, error) {
	return NewParser(NewLexerFromSource(filename, sourceCode)).ParseSong()
}

func ReadFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(data), nil
}

func ParseSongFromFile(file string) (*Song, error) {
	defer timer.LogElapsedTime("parsing")()

	data, err := ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return NewParser(NewLexerFromSource(file, string(data))).ParseSong()
}

var SURROUNDING_CONTEXT = 15

// Song =
//
//	Frontmatter Body
//	| Body
//	;
func (p *Parser) ParseSong() (*Song, error) {
	song := Song{}
	song.Parser = p
	p.song = &song
	p.Lexer.consumeWhitespacesAndNewLines()

	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	switch tok.Type {
	case TokenFrontMatter:
		fm, err := p.ParseFrontmatter()
		if err != nil {
			return nil, err
		}
		song.FrontMatter = fm
		fallthrough
	default:
		body, err := p.ParseBody()
		if err != nil {
			return nil, err
		}
		song.Sections = body
		return &song, nil
	}
}

// FrontMatter: TokenFrontmatter
func (p *Parser) ParseFrontmatter() (map[string]string, error) {
	tok, err := p.Lexer.ConsumeNextToken()
	if err != nil {
		return nil, err
	}

	if tok.Type != TokenFrontMatter {
		return nil, fmt.Errorf("unexpected token. Want TokenFrontmatter, Got: %s", tok.Type)
	}
	bytes := []byte(tok.Value)
	frontmatter := map[string]string{}
	err = yaml.Unmarshal(bytes, &frontmatter)
	if err != nil {
		return nil, err
	}
	return frontmatter, err
}

// Body:
// Lines Sections
// |Sections
// ;
func (p *Parser) ParseBody() ([]Section, error) {
	tok, err := p.Lexer.Lookahead()
	sections := []Section{}
	if err != nil {
		return sections, err
	}
	switch tok.Type {
	case TokenHeader, TokenHeaderBreak:
		sections, err := p.ParseSections()
		if err != nil {
			return nil, err
		}
		return sections, nil
	default:
		emptySection := Section{
			Name:  "",
			Lines: []Line{},
			Break: false,
		}
		lines, err := p.ParseLines()
		if err != nil {
			return nil, err
		}
		emptySection.Lines = lines

		rest, err := p.ParseSections()
		if err != nil {
			return nil, err
		}
		sections = append(sections, emptySection)
		sections = append(sections, rest...)
		return sections, nil
	}
}

// Sections
// :Section
// |Section Sections
func (p *Parser) ParseSections() ([]Section, error) {
	sections := []Section{}
	for {
		tok, err := p.Lexer.Lookahead()
		if err != nil {
			return nil, err
		}

		switch tok.Type {
		case TokenEof:
			return sections, nil
		case TokenHeader, TokenHeaderBreak:
			section, err := p.ParseSection()
			if err != nil {
				return nil, err
			}
			sections = append(sections, *section)
		default:
			return nil, fmt.Errorf("unexpected token. Expected TokenHeader, TokenHeaderBreak or TokenEof, got %s", tok.Type)
		}
	}
}

// Section
// :Header Lines
func (p *Parser) ParseSection() (*Section, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	switch tok.Type {
	case TokenHeader, TokenHeaderBreak:
		section := Section{
			Name:  tok.Value,
			Lines: nil,
			Break: tok.Type == TokenHeaderBreak,
		}
		_, _ = p.Lexer.ConsumeNextToken()

		lines, err := p.ParseLines()
		if err != nil {
			return nil, err
		}
		section.Lines = lines
		return &section, nil
	default:
		return nil, fmt.Errorf("unexpected token while parsing section: %s at pos %d", tok.Type, p.Lexer.pos)
	}
}

// Lines
// Line
// |Line Lines
// ;
func (p *Parser) ParseLines() ([]Line, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	lines := []Line{}

	for {
		switch tok.Type {
		case TokenHeaderBreak, TokenHeader, TokenEof:
			return lines, nil
		default:
			line, err := p.ParseLine()
			if err != nil {
				return nil, err
			}
			if len(line.Bars) > 0 || line.MultilineBacktick.Value != "" {
				lines = append(lines, *line)
			}
			tok, err = p.Lexer.Lookahead()
			if err != nil {
				return nil, err
			}
		}
	}
}

// Line:
// Bars TokenReturn
// |Bar Bars
func (p *Parser) ParseLine() (*Line, error) {
	bars := []Bar{}

	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}

	if tok.Type == TokenBacktickMultiline {
		_, _ = p.Lexer.ConsumeNextToken()
		line := &Line{
			Bars: []Bar{},
			MultilineBacktick: MultilineBacktick{
				Id:            p.mutilineBacktickId,
				Value:         tok.Value,
				DefaultLength: p.song.DefaultLength(),
				SourceFile:    p.SourceFile(),
			},
		}
		p.mutilineBacktickId++
		return line, nil
	}

	var prev *Bar
	for tok.Type != TokenReturn && tok.Type != TokenEof {
		bar, err := p.ParseBar()
		if err != nil {
			return nil, err
		}
		bar.PreviousWasRepeatEnd = prev != nil && prev.RepeatEnd

		bars = append(bars, *bar)
		tok, err = p.Lexer.Lookahead()
		if err != nil {
			return nil, err
		}
		prev = bar
	}
	_, _ = p.Lexer.ConsumeNextToken()
	return &Line{Bars: bars}, nil
}

// Bar
// :TokenBarNote TokenBar BarBody
// |TokenBar TokenBarNote BarBody
// |BarBody
// ;
//
// BarBody
// :TokenBacktick
// |Chords
// ;
// Chords
// :Chord
// |Chord Chords
func (p *Parser) ParseBar() (*Bar, error) {
	bar := Bar{}
	bar.Id = p.barsCount
	p.barsCount++

	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}

	for tok.Type == TokenBar || tok.Type == TokenBarNote {
		switch tok.Type {
		case TokenBar:
			if tok.Value == "||:" {
				bar.RepeatStart = true
			}
			_, _ = p.Lexer.ConsumeNextToken()
			tok, err = p.Lexer.Lookahead()
			if err != nil {
				return nil, err
			}
		case TokenBarNote:
			bar.BarNote = tok.Value
			// consume it
			_, _ = p.Lexer.ConsumeNextToken()
			p.Lexer.consumeWhitespacesAndNewLines()
			tok, err = p.Lexer.Lookahead()
			if err != nil {
				return nil, err
			}
		}
	}
	// BarBody
	if tok.Type != TokenChord && tok.Type != TokenAnnotation && tok.Type != TokenBacktick {
		return nil, fmt.Errorf("parsing bar: unexpected token. Want Chord, Annotation or Backtick, got %s %s", tok.Type, p.Lexer.SurroundingString())
	}

	switch tok.Type {
	case TokenBacktick:
		backtick, err := p.ParseBacktick()
		if err != nil {
			return nil, err
		}
		bar.Backtick = *backtick
		tok, err = p.Lexer.Lookahead()
		if err != nil {
			return nil, err
		}
		if tok.Type == TokenBar && tok.Value != "||:" {
			switch tok.Value {
			case ":||":
				bar.RepeatEnd = true
			case "||":
				bar.DoubleBarEnd = true
			}
			_, err = p.Lexer.ConsumeNextToken()
			if err != nil {
				return nil, err
			}
		}
		return &bar, nil
	case TokenAnnotation, TokenChord:
		chords := []Chord{}
		for tok.Type == TokenChord || tok.Type == TokenAnnotation {
			chord, err := p.ParseChord()
			if err != nil {
				return nil, err
			}
			chords = append(chords, *chord)
			tok, err = p.Lexer.Lookahead()
			if err != nil {
				return nil, err
			}
		}
		if len(chords) == 0 {
			return nil, fmt.Errorf("found no chords in bar %s", p.Lexer.SurroundingString())
		}
		bar.Chords = chords

		if tok.Type == TokenBar && tok.Value != "||:" {
			switch tok.Value {
			case ":||":
				bar.RepeatEnd = true
			case "||":
				bar.DoubleBarEnd = true
			}
			_, err = p.Lexer.ConsumeNextToken()
			if err != nil {
				return nil, err
			}
		}
		return &bar, nil
	default:
		return nil, fmt.Errorf(
			"expected chord or backtick expression but found %s %s",
			tok.Type,
			p.Lexer.SurroundingString(),
		)
	}
}

func (p *Parser) ParseBacktick() (*Backtick, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}

	if tok.Type != TokenBacktick {
		return nil, fmt.Errorf("expected backtick, got: %s", tok.Type)
	}

	bt := Backtick{
		Id:            p.backtickId,
		Value:         "",
		DefaultLength: p.song.DefaultLength(),
	}
	p.backtickId++
	bt.Value = tok.Value
	_, _ = p.Lexer.ConsumeNextToken()

	return &bt, nil
}

// Chord
// :TokenChord
// |TokenAnnotation TokenChord
func (p *Parser) ParseChord() (*Chord, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	chord := Chord{
		Value:      "",
		Annotation: &Annotation{},
	}
	if tok.Type == TokenAnnotation {
		chord.Annotation.Value = tok.Value
		_, _ = p.Lexer.ConsumeNextToken()
		tok, err = p.Lexer.Lookahead()
		if err != nil {
			return nil, err
		}
	}

	if tok.Type != TokenChord {
		return nil, fmt.Errorf(
			"expected chord but found %s, %s",
			tok.Type,
			p.Lexer.SurroundingString(),
		)
	}
	chord.Value = tok.Value

	if chord.Value == "" {
		return nil, fmt.Errorf("empty chord found at %s", p.Lexer.SurroundingString())
	}

	_, _ = p.Lexer.ConsumeNextToken()
	return &chord, nil
}
