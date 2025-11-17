package internal

import (
	"fmt"
	"lesheets/internal/domain"
	"lesheets/internal/timer"
	"os"
)

type Parser struct {
	Lexer              *Lexer
	backtickId         int
	mutilineBacktickId int
	song               *domain.Song
	barsCount          int
}

func (p *Parser) SourceFile() string {
	return p.Lexer.source
}

func NewParser(lex *Lexer) *Parser {
	return &Parser{
		Lexer:      lex,
		backtickId: 0,
	}
}

func ParseSongFromString(s string) (*domain.Song, error) {
	return NewParser(NewLexerFromSource("unknown", s)).ParseSong()
}

func ParseSongFromStringWithFileName(filename string, sourceCode string) (*domain.Song, error) {
	return NewParser(NewLexerFromSource(filename, sourceCode)).ParseSong()
}

func ReadFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

func ParseSongFromFile(file string) (*Parser, *domain.Song, error) {
	defer timer.LogElapsedTime("parsing")()

	data, err := ReadFile(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	parser := NewParser(NewLexerFromSource(file, string(data)))
	res, err := parser.ParseSong()
	return parser, res, err
}

var SURROUNDING_CONTEXT = 15

// Song =
//
//	Frontmatter Body
//	| Body
//	;
func (p *Parser) ParseSong() (*domain.Song, error) {
	song := domain.Song{}
	p.song = &song
	p.Lexer.consumeWhitespacesAndNewLines()

	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	switch tok.Type {
	case domain.TokenFrontMatter:
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
	_, err := p.Lexer.ConsumeNextToken()
	if err != nil {
		return nil, err
	}
	//
	// if tok.Type != domain.TokenFrontMatter {
	// 	return nil, fmt.Errorf("unexpected token. Want TokenFrontmatter, Got: %s", tok.Type)
	// }
	// bytes := []byte(tok.Value)
	// frontmatter := map[string]string{}
	// err = yaml.Unmarshal(bytes, &frontmatter)
	// if err != nil {
	// 	return nil, err
	// }
	// return frontmatter, err
	return map[string]string{}, nil
}

// Body:
// Lines Sections
// |Sections
// ;
func (p *Parser) ParseBody() ([]domain.Section, error) {
	tok, err := p.Lexer.Lookahead()
	sections := []domain.Section{}
	if err != nil {
		return sections, err
	}
	switch tok.Type {
	case domain.TokenHeader, domain.TokenHeaderBreak:
		sections, err := p.ParseSections()
		if err != nil {
			return nil, err
		}
		return sections, nil
	default:
		emptySection := domain.Section{
			Name:  "",
			Lines: []domain.Line{},
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
func (p *Parser) ParseSections() ([]domain.Section, error) {
	sections := []domain.Section{}
	for {
		tok, err := p.Lexer.Lookahead()
		if err != nil {
			return nil, err
		}

		switch tok.Type {
		case domain.TokenEof:
			return sections, nil
		case domain.TokenHeader, domain.TokenHeaderBreak:
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
func (p *Parser) ParseSection() (*domain.Section, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	switch tok.Type {
	case domain.TokenHeader, domain.TokenHeaderBreak:
		section := domain.Section{
			Name:  tok.Value,
			Lines: nil,
			Break: tok.Type == domain.TokenHeaderBreak,
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
func (p *Parser) ParseLines() ([]domain.Line, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	lines := []domain.Line{}

	for {
		switch tok.Type {
		case domain.TokenHeaderBreak, domain.TokenHeader, domain.TokenEof:
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
func (p *Parser) ParseLine() (*domain.Line, error) {
	bars := []domain.Bar{}

	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}

	if tok.Type == domain.TokenBacktickMultiline {
		_, _ = p.Lexer.ConsumeNextToken()
		line := &domain.Line{
			Bars: []domain.Bar{},
			MultilineBacktick: domain.MultilineBacktick{
				Id:            p.mutilineBacktickId,
				Value:         tok.Value,
				DefaultLength: p.song.DefaultLength(),
				SourceFile:    p.SourceFile(),
			},
		}
		p.mutilineBacktickId++
		return line, nil
	}

	var prev *domain.Bar
	for tok.Type != domain.TokenReturn && tok.Type != domain.TokenEof {
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
	return &domain.Line{Bars: bars}, nil
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
func (p *Parser) ParseBar() (*domain.Bar, error) {
	bar := domain.Bar{}
	bar.Id = p.barsCount
	p.barsCount++

	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}

	for tok.Type == domain.TokenBar || tok.Type == domain.TokenBarNote {
		switch tok.Type {
		case domain.TokenBar:
			if tok.Value == "||:" {
				bar.RepeatStart = true
			}
			_, _ = p.Lexer.ConsumeNextToken()
			tok, err = p.Lexer.Lookahead()
			if err != nil {
				return nil, err
			}
		case domain.TokenBarNote:
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
	if tok.Type != domain.TokenChord && tok.Type != domain.TokenAnnotation && tok.Type != domain.TokenBacktick {
		return nil, fmt.Errorf("parsing bar: unexpected token. Want Chord, Annotation or Backtick, got %s %s", tok.Type, p.Lexer.SurroundingString())
	}

	switch tok.Type {
	case domain.TokenBacktick:
		backtick, err := p.ParseBacktick()
		if err != nil {
			return nil, err
		}
		bar.Backtick = *backtick
		tok, err = p.Lexer.Lookahead()
		if err != nil {
			return nil, err
		}
		if tok.Type == domain.TokenBar && tok.Value != "||:" {
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
	case domain.TokenAnnotation, domain.TokenChord:
		chords := []domain.Chord{}
		for tok.Type == domain.TokenChord || tok.Type == domain.TokenAnnotation {
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

		if tok.Type == domain.TokenBar && tok.Value != "||:" {
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

func (p *Parser) ParseBacktick() (*domain.Backtick, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}

	if tok.Type != domain.TokenBacktick {
		return nil, fmt.Errorf("expected backtick, got: %s", tok.Type)
	}

	bt := domain.Backtick{
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
func (p *Parser) ParseChord() (*domain.Chord, error) {
	tok, err := p.Lexer.Lookahead()
	if err != nil {
		return nil, err
	}
	chord := domain.Chord{
		Value:      "",
		Annotation: &domain.Annotation{},
	}
	if tok.Type == domain.TokenAnnotation {
		chord.Annotation.Value = tok.Value
		_, _ = p.Lexer.ConsumeNextToken()
		tok, err = p.Lexer.Lookahead()
		if err != nil {
			return nil, err
		}
	}

	if tok.Type != domain.TokenChord {
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
