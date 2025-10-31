package internal

import (
	"fmt"

	yaml "github.com/oasdiff/yaml3"
)

type Parser struct {
	lex        *Lexer
	backtickId int
}

func NewParser(lex *Lexer) *Parser {
	return &Parser{
		lex:        lex,
		backtickId: 0,
	}
}

func ParseSongFromString(s string) (*Song, error) {
	parser := Parser{
		lex:        NewLexer(s),
		backtickId: 0,
	}
	return parser.ParseSong()
}

var SURROUNDING_COUNTEXT = 10

// Song =
//
//	Frontmatter Body
//	| Body
//	;
func (p *Parser) ParseSong() (*Song, error) {
	song := Song{}
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
			Name:      "",
			BarsLines: [][]Bar{},
			Break:     false,
		}
		lines, err := p.ParseLines()
		if err != nil {
			return nil, err
		}
		emptySection.BarsLines = lines

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
			Name:      tok.Value,
			BarsLines: nil,
			Break:     tok.Type == TokenHeaderBreak,
		}
		_, _ = p.lex.ConsumeNextToken()

		lines, err := p.ParseLines()
		if err != nil {
			return nil, err
		}
		section.BarsLines = lines
		return &section, nil
	default:
		return nil, fmt.Errorf("unexpected token while parsing section: %s at pos %d", tok.Type, p.lex.pos)
	}
}

// Lines
// Line
// |Line Lines
// ;
func (p *Parser) ParseLines() ([][]Bar, error) {
	tok, err := p.lex.Lookahead()
	if err != nil {
		return nil, err
	}
	lines := [][]Bar{}

	for {
		switch tok.Type {
		case TokenHeaderBreak, TokenHeader, TokenEof:
			return lines, nil
		default:
			line, err := p.ParseBarsLine()
			if err != nil {
				return nil, err
			}
			if len(*line) > 0 {
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
// Bar TokenReturn
// |Bar Bars
func (p *Parser) ParseBarsLine() (*[]Bar, error) {
	bars := []Bar{}

	tok, err := p.lex.Lookahead()
	if err != nil {
		return nil, err
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
		if tok.Type == TokenBar {
			_, _ = p.lex.ConsumeNextToken()
			tok, err = p.lex.Lookahead()
			if err != nil {
				return nil, err
			}
		}
	}
	_, _ = p.lex.ConsumeNextToken()
	return &bars, nil
}

// Bar
// :TokenComment BarBody
// |BarBody
// ;
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
	switch tok.Type {
	case TokenBar:
		_, _ = p.lex.ConsumeNextToken()
		tok, err = p.lex.Lookahead()
		if err != nil {
			return nil, err
		}
	case TokenRepeatStart:
		bar.RepeatStart = true
		_, _ = p.lex.ConsumeNextToken()
		tok, err = p.lex.Lookahead()
		if err != nil {
			return nil, err
		}
	}
	if tok.Type != TokenChord && tok.Type != TokenAnnotation && tok.Type != TokenComment && tok.Type != TokenBacktick {
		return nil, fmt.Errorf("parsing bar: unexpected token %s", tok.Type)
	}
	if tok.Type == TokenComment {
		bar.Comment = tok.Value
		// consume it
		_, _ = p.lex.ConsumeNextToken()
		p.lex.consumeWhitespacesAndNewLines()
		tok, err = p.lex.Lookahead()
		if err != nil {
			return nil, err
		}
	}

	switch tok.Type {
	case TokenBacktick:
		backtick, err := p.ParseBacktick()
		if err != nil {
			return nil, err
		}
		bar.Backtick = *backtick
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
			return nil, fmt.Errorf("found no chords in bar at pos %d", p.lex.pos)
		}
		if tok.Type == TokenRepeatEnd {
			bar.RepeatEnd = true
			_, err = p.lex.ConsumeNextToken()
			if err != nil {
				return nil, err
			}
		}
		bar.Chords = chords
		return &bar, nil
	default:
		return nil, fmt.Errorf(
			"expected chord or backtick expression but found %s at pos %d, near %s",
			tok.Type,
			p.lex.pos,
			p.lex.SurroundingString(SURROUNDING_COUNTEXT),
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
		Id:    backtickId,
		Value: "",
	}
	backtickId++
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
			"expected chord but found %s, near %s",
			tok.Type,
			p.lex.SurroundingString(SURROUNDING_COUNTEXT),
		)
	}
	chord.Value = tok.Value

	if chord.Value == "" {
		return nil, fmt.Errorf("empty chord found at %s", p.lex.SurroundingString(SURROUNDING_COUNTEXT))
	}

	_, _ = p.lex.ConsumeNextToken()
	return &chord, nil
}

var backtickId = 0
