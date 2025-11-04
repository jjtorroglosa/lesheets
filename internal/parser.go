package internal

import (
	"fmt"

	yaml "github.com/oasdiff/yaml3"
)

type Parser struct {
	lex                *Lexer
	backtickId         int
	mutilineBacktickId int
	song               *Song
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
		lex:        lex,
		backtickId: 0,
	}
}

func ParseSongFromString(s string) (*Song, error) {
	return NewParser(NewLexer(s)).ParseSong()
}

var SURROUNDING_CONTEXT = 15

// Song =
//
//	Frontmatter Body
//	| Body
//	;
func (p *Parser) ParseSong() (*Song, error) {
	song := Song{}
	p.song = &song
	p.lex.consumeWhitespacesAndNewLines()

	tok, err := p.lex.Lookahead()
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
	tok, err := p.lex.ConsumeNextToken()
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
	tok, err := p.lex.Lookahead()
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
		tok, err := p.lex.Lookahead()
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
	tok, err := p.lex.Lookahead()
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
		_, _ = p.lex.ConsumeNextToken()

		lines, err := p.ParseLines()
		if err != nil {
			return nil, err
		}
		section.Lines = lines
		return &section, nil
	default:
		return nil, fmt.Errorf("unexpected token while parsing section: %s at pos %d", tok.Type, p.lex.pos)
	}
}

// Lines
// Line
// |Line Lines
// ;
func (p *Parser) ParseLines() ([]Line, error) {
	tok, err := p.lex.Lookahead()
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
			tok, err = p.lex.Lookahead()
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

	tok, err := p.lex.Lookahead()
	if err != nil {
		return nil, err
	}

	if tok.Type == TokenBacktickMultiline {
		_, _ = p.lex.ConsumeNextToken()
		line := &Line{
			Bars: []Bar{},
			MultilineBacktick: MultilineBacktick{
				Id:            p.mutilineBacktickId,
				Value:         tok.Value,
				DefaultLength: p.song.DefaultLength(),
			},
		}
		p.mutilineBacktickId++
		return line, nil
	}

	for tok.Type != TokenReturn && tok.Type != TokenEof {
		bar, err := p.ParseBar()
		if err != nil {
			return nil, err
		}
		bars = append(bars, *bar)
		tok, err = p.lex.Lookahead()
		if err != nil {
			return nil, err
		}
		// if tok.Type == TokenBar {
		// 	_, _ = p.lex.ConsumeNextToken()
		// 	tok, err = p.lex.Lookahead()
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// }
	}
	_, _ = p.lex.ConsumeNextToken()
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

	tok, err := p.lex.Lookahead()
	if err != nil {
		return nil, err
	}

	for tok.Type == TokenBar || tok.Type == TokenBarNote {
		switch tok.Type {
		case TokenBar:
			if tok.Value == "||:" {
				bar.RepeatStart = true
			}
			_, _ = p.lex.ConsumeNextToken()
			tok, err = p.lex.Lookahead()
			if err != nil {
				return nil, err
			}
		case TokenBarNote:
			bar.BarNote = tok.Value
			// consume it
			_, _ = p.lex.ConsumeNextToken()
			p.lex.consumeWhitespacesAndNewLines()
			tok, err = p.lex.Lookahead()
			if err != nil {
				return nil, err
			}
		}
	}
	// BarBody
	if tok.Type != TokenChord && tok.Type != TokenAnnotation && tok.Type != TokenBacktick {
		return nil, fmt.Errorf("parsing bar: unexpected token. Want Chord, Annotation or Backtick, got %s %s", tok.Type, p.lex.SurroundingString())
	}

	switch tok.Type {
	case TokenBacktick:
		backtick, err := p.ParseBacktick()
		if err != nil {
			return nil, err
		}
		bar.Backtick = *backtick
		tok, err = p.lex.Lookahead()
		if err != nil {
			return nil, err
		}
		if tok.Type == TokenBar && tok.Value != "||:" {
			if tok.Value == ":||" {
				bar.RepeatEnd = true
			}
			_, _ = p.lex.ConsumeNextToken()
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
			tok, err = p.lex.Lookahead()
			if err != nil {
				return nil, err
			}
		}
		if len(chords) == 0 {
			return nil, fmt.Errorf("found no chords in bar %s", p.lex.SurroundingString())
		}
		bar.Chords = chords

		if tok.Type == TokenBar {
			if tok.Value == ":||" {
				bar.RepeatEnd = true
			}
			// Don't consume repeat start, let the next bar to consume it at the beginning
			if tok.Value != "||:" {
				_, err = p.lex.ConsumeNextToken()
				if err != nil {
					return nil, err
				}
			}
		}
		return &bar, nil
	default:
		return nil, fmt.Errorf(
			"expected chord or backtick expression but found %s %s",
			tok.Type,
			p.lex.SurroundingString(),
		)
	}
}

func (p *Parser) ParseBacktick() (*Backtick, error) {
	tok, err := p.lex.Lookahead()
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
	_, _ = p.lex.ConsumeNextToken()

	return &bt, nil
}

// Chord
// :TokenChord
// |TokenAnnotation TokenChord
func (p *Parser) ParseChord() (*Chord, error) {
	tok, err := p.lex.Lookahead()
	if err != nil {
		return nil, err
	}
	chord := Chord{
		Value:      "",
		Annotation: &Annotation{},
	}
	if tok.Type == TokenAnnotation {
		chord.Annotation.Value = tok.Value
		_, _ = p.lex.ConsumeNextToken()
		tok, err = p.lex.Lookahead()
		if err != nil {
			return nil, err
		}
	}

	if tok.Type != TokenChord {
		return nil, fmt.Errorf(
			"expected chord but found %s, %s",
			tok.Type,
			p.lex.SurroundingString(),
		)
	}
	chord.Value = tok.Value

	if chord.Value == "" {
		return nil, fmt.Errorf("empty chord found at %s", p.lex.SurroundingString())
	}

	_, _ = p.lex.ConsumeNextToken()
	return &chord, nil
}
