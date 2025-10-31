package internal

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenFrontMatter TokenType = "FrontMatter"
	TokenHeader      TokenType = "Header"
	TokenHeaderBreak TokenType = "HeaderBreak"
	TokenChord       TokenType = "Chord"
	TokenBar         TokenType = "Bar"
	TokenReturn      TokenType = "Return"
	TokenComment     TokenType = "Comment"
	TokenAnnotation  TokenType = "Annotation"
	TokenBacktick    TokenType = "BacktickExpression"
	TokenUnknown     TokenType = "Unknown"
	TokenRepeatEnd   TokenType = "RepeatEnd"
	TokenRepeatStart TokenType = "RepeatStart"
	TokenEof         TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) nextChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) eof() bool {
	return l.pos >= len(l.input)
}

func (l *Lexer) advance() {
	l.pos++
}

func (l *Lexer) consumeWhitespaces() {
	ch := l.nextChar()
	for ch == ' ' || ch == '\t' || ch == '\r' {
		l.advance()
		ch = l.nextChar()
	}
}

func (l *Lexer) consumeWhitespacesAndNewLines() {
	ch := l.nextChar()
	for ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
		l.advance()
		ch = l.nextChar()
	}
}

func ErrGeneric(want string, got string) error {
	return fmt.Errorf("unexpected string. Want: %s Got: %s", want, got)
}
func ErrInvalidFrontmatter(want string, got string) error {
	return fmt.Errorf("invalid frontmatter. Want: %s Got: %s", want, got)
}

func (l *Lexer) consumeFrontmatter() (*Token, error) {
	start := l.pos
	// opening ---
	for !l.eof() && l.nextChar() == '-' {
		l.advance()
	}
	if l.input[start:l.pos] != "---" {
		return nil, ErrInvalidFrontmatter("Opening ---", l.input[start:l.pos])
	}

	// body
	start = l.pos
	for !l.eof() && l.nextChar() != '-' {
		l.advance()
	}
	value := l.input[start:l.pos]

	// closing ---
	start = l.pos
	for !l.eof() && l.nextChar() == '-' {
		l.advance()
	}
	if l.input[start:l.pos] != "---" {
		return nil, ErrInvalidFrontmatter("Closing ---", l.input[start:l.pos])
	}

	return &Token{
		Type:  TokenFrontMatter,
		Value: value,
	}, nil
}

func (l *Lexer) Lookahead() (*Token, error) {
	prev := l.pos
	tok, err := l.ConsumeNextToken()
	l.pos = prev
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func (l *Lexer) ConsumeNextToken() (*Token, error) {
	l.consumeWhitespaces()
	if l.eof() {
		return &Token{
			Type:  TokenEof,
			Value: "",
		}, nil
	}
	ch := l.nextChar()

	// Frontmatter
	if ch == '-' {
		tok, err := l.consumeFrontmatter()
		if err != nil {
			return nil, err
		}

		return tok, nil
	}
	if ch == ':' {
		if !strings.HasPrefix(l.input[l.pos:], ":||") {
			return nil, ErrGeneric(":||", l.input[l.pos:3])
		}
		tok := Token{
			Type:  TokenRepeatEnd,
			Value: ":||",
		}
		l.pos += 3
		return &tok, nil
	}
	// Single bar
	if ch == '|' {
		if strings.HasPrefix(l.input[l.pos:], "||:") {
			tok := Token{
				Type:  TokenRepeatStart,
				Value: "||:",
			}
			l.pos += 3
			return &tok, nil
		} else if strings.HasPrefix(l.input[l.pos:], "||") {
			tok := Token{
				Type:  TokenBar,
				Value: "||",
			}
			l.pos += 2
			return &tok, nil
		} else {
			tok := Token{
				Type:  TokenBar,
				Value: "|",
			}
			l.advance()
			return &tok, nil
		}
	}

	// Annotation
	if ch == '!' {
		l.advance()
		start := l.pos
		for l.pos < len(l.input) && !unicode.IsSpace(rune(l.input[l.pos])) && l.input[l.pos] != '!' {
			l.advance()
		}
		tok := Token{
			Type:  TokenAnnotation,
			Value: l.input[start:l.pos],
		}
		l.advance()
		return &tok, nil
	}

	// Comment
	if ch == '"' {
		l.advance() // skip opening "
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '"' {
			l.advance()
		}
		tok := Token{
			Type:  TokenComment,
			Value: l.input[start:l.pos],
		}

		if l.pos >= len(l.input) || l.input[l.pos] != '"' {
			return nil, ErrGeneric("\"", string(l.nextChar()))
		}
		// consume closing "
		l.advance() // skip closing "
		return &tok, nil
	}

	// Backtick expression
	if ch == '`' {
		l.advance() // skip opening `
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '`' {
			l.advance()
		}
		tok := Token{
			Type:  TokenBacktick,
			Value: l.input[start:l.pos],
		}

		if l.pos >= len(l.input) || l.input[l.pos] != '`' {
			return nil, ErrGeneric("`", string(l.nextChar()))
		}
		l.advance()
		return &tok, nil
	}

	// Headers
	if ch == '#' {
		tokenType := TokenHeader
		for l.pos < len(l.input) && l.input[l.pos] == '#' {
			l.advance()
		}
		if l.pos < len(l.input) && l.input[l.pos] == '-' {
			tokenType = TokenHeaderBreak
			l.advance()
		}
		l.consumeWhitespaces()
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '\n' {
			l.advance()
		}
		tok := Token{
			Type:  tokenType,
			Value: strings.TrimSpace(l.input[start:l.pos]),
		}
		l.advance()
		return &tok, nil
	}

	// If after a chord there is a \n, close the bar, and start new line
	if l.nextChar() == '\n' {
		l.consumeWhitespacesAndNewLines()
		// tokens = append(tokens, Token{Type: TokenBar, Value: "|"})
		tok := Token{Type: TokenReturn, Value: "\n"}
		return &tok, nil
	}

	// Chords
	start := l.pos
	for l.pos < len(l.input) && !unicode.IsSpace(rune(l.input[l.pos])) && l.input[l.pos] != '|' && l.input[l.pos] != '\n' {
		l.advance()
	}
	value := l.input[start:l.pos]
	tok := Token{Type: TokenChord, Value: value}
	return &tok, nil
}

// Scan all tokens in the input
func (l *Lexer) Lex() ([]Token, error) {
	var tokens []Token
	for l.pos < len(l.input) {
		tok, err := l.ConsumeNextToken()
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, *tok)
	}

	return tokens, nil
}

func (l *Lexer) SurroundingString(context int) string {
	start := max(0, l.pos-context)
	end := min(len(l.input), l.pos+context)
	return l.input[start:end]
}
func (l *Lexer) PrintTokens() {
	for l.pos < len(l.input) {
		tok, err := l.ConsumeNextToken()
		if err != nil {
			return
		}

		fmt.Println(tok.Type)
	}
}
