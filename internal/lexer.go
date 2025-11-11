package internal

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenFrontMatter       TokenType = "FrontMatter"
	TokenHeader            TokenType = "Header"
	TokenHeaderBreak       TokenType = "HeaderBreak"
	TokenBar               TokenType = "Bar"
	TokenReturn            TokenType = "Return"
	TokenBarNote           TokenType = "BarNote"
	TokenAnnotation        TokenType = "Annotation"
	TokenBacktick          TokenType = "BacktickExpression"
	TokenBacktickMultiline TokenType = "BacktickMultilineExpression"
	TokenUnknown           TokenType = "Unknown"
	TokenEof               TokenType = "EOF"
	TokenChord             TokenType = "Chord"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	source string
	input  string
	pos    int
	line   int
}

func NewLexer(input string) *Lexer {
	return NewLexerFromSource("", input)
}

func NewLexerFromSource(source string, input string) *Lexer {
	return &Lexer{source: source, input: input, pos: 0}
}

func (l *Lexer) nextChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) getPos(pos int, length int) string {
	if pos < 0 {
		return "BeginningOfFile"
	}
	if pos+length >= len(l.input) {
		return "EndOfFile"
	}
	return l.input[pos : pos+length]
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
		if ch == '\n' || ch == '\r' {
			l.line++
		}
		l.advance()
		ch = l.nextChar()
	}
}

func (l *Lexer) consumeWhitespacesAndNewLines() {
	ch := l.nextChar()
	for ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
		if ch == '\n' || ch == '\r' {
			l.line++
		}
		l.advance()
		ch = l.nextChar()
	}
}

func ErrGeneric(context string, want string, got string) error {
	return fmt.Errorf("unexpected string %sWant: <%s> Got: <%s>", context, want, got)
}

func ErrInvalidFrontmatter(context string, want string, got string) error {
	return fmt.Errorf("invalid frontmatter %sWant: <%s> Got: <%s> ", context, want, got)
}

func (l *Lexer) consumeFrontmatter() (*Token, error) {
	start := l.pos
	// opening ---
	for !l.eof() && l.nextChar() == '-' {
		l.advance()
	}
	if l.getPos(start, 3) != "---" {
		return nil, ErrInvalidFrontmatter(l.SurroundingString(), "Opening ---", l.getPos(start, 3))
	}

	// body
	start = l.pos
	for !l.eof() && l.nextChar() != '-' {
		if l.nextChar() == '\n' || l.nextChar() == '\r' {
			l.line++
		}
		l.advance()
	}
	value := l.input[start:l.pos]

	// closing ---
	start = l.pos
	for !l.eof() && l.nextChar() == '-' {
		l.advance()
	}
	if l.input[start:l.pos] != "---" {
		return nil, ErrInvalidFrontmatter(l.SurroundingString(), "Closing ---", l.getPos(start, 3))
	}

	return &Token{
		Type:  TokenFrontMatter,
		Value: value,
	}, nil
}

func (l *Lexer) Lookahead() (*Token, error) {
	prev := l.pos
	line := l.line
	tok, err := l.ConsumeNextToken()
	l.pos = prev
	l.line = line
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func (l *Lexer) ConsumeNextToken() (*Token, error) {
	l.consumeWhitespaces()
	if l.eof() {
		l.pos++
		return &Token{
			Type:  TokenEof,
			Value: "",
		}, nil
	}
	ch := l.nextChar()

	// Frontmatter
	if l.pos+3 < len(l.input) && l.input[l.pos:l.pos+3] == "---" {
		tok, err := l.consumeFrontmatter()
		if err != nil {
			return nil, err
		}

		return tok, nil
	}
	if ch == ':' {
		if l.getPos(l.pos, 3) != ":||" {
			return nil, ErrGeneric(l.SurroundingString(), ":||", l.getPos(l.pos, 3))
		}
		tok := Token{
			Type:  TokenBar,
			Value: ":||",
		}
		l.pos += 3
		return &tok, nil
	}
	// Single bar
	if ch == '|' {
		if l.getPos(l.pos, 3) == "||:" {
			tok := Token{
				Type:  TokenBar,
				Value: "||:",
			}
			l.pos += 3
			return &tok, nil
		} else if l.getPos(l.pos, 2) == "||" {
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

	// Comment
	if l.getPos(l.pos, 2) == "//" {
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '\n' {
			l.advance()
		}
		// consume newline only if the comment started in a new line
		if l.getPos(start-1, 1) == "\n" {
			l.advance()
		}
		return l.ConsumeNextToken()
	}
	// Annotation
	if ch == '!' {
		l.advance()
		start := l.pos
		for l.pos < len(l.input) && !unicode.IsSpace(rune(l.input[l.pos])) && l.input[l.pos] != '!' {
			l.advance()
		}
		if l.nextChar() != '!' {
			return nil, ErrGeneric(l.SurroundingString(), "!", string(l.nextChar()))
		}
		tok := Token{
			Type:  TokenAnnotation,
			Value: l.input[start:l.pos],
		}
		l.advance()
		return &tok, nil
	}

	// BarNote
	if ch == '"' {
		l.advance() // skip opening "
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '"' {
			l.advance()
		}
		tok := Token{
			Type:  TokenBarNote,
			Value: l.input[start:l.pos],
		}

		if l.pos >= len(l.input) || l.input[l.pos] != '"' {
			return nil, ErrGeneric(l.SurroundingString(), "\"", l.getPos(l.pos, 1))
		}
		// consume closing "
		l.advance() // skip closing "
		return &tok, nil
	}

	// Backtick expression
	if ch == '`' {
		// Backtick multiline
		if l.pos+3 < len(l.input) && l.getPos(l.pos, 3) == "```" {
			l.pos += 3
			l.consumeWhitespacesAndNewLines()
			start := l.pos
			for l.pos < len(l.input) && l.input[l.pos] != '`' {
				l.advance()
			}
			if l.pos+3 > len(l.input) || l.getPos(l.pos, 3) != "```" {
				return nil, ErrGeneric(l.SurroundingString(), "Closing ```", l.getPos(l.pos, 3))
			}
			tok := Token{
				Type:  TokenBacktickMultiline,
				Value: l.input[start:l.pos],
			}
			l.pos += 3
			return &tok, nil
		} else {
			// Backtick inline
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
				return nil, ErrGeneric(l.SurroundingString(), "`", l.getPos(l.pos, 1))
			}
			l.advance()
			return &tok, nil
		}
	}

	// Headers
	if ch == '#' && (l.getPos(l.pos-1, 1) == "" || l.getPos(l.pos-1, 1) == "\n") {
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
		if l.nextChar() == '\n' {
			l.line++
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

func (l *Lexer) replaceNewLine(r rune) string {
	if r == '\n' {
		return "\\n"
	} else {
		return string(r)
	}
}

func (l *Lexer) SurroundingString() string {
	context := SURROUNDING_CONTEXT
	start := max(0, l.pos-context)
	end := min(len(l.input), l.pos+context)
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("at pos %d ", l.pos))
	sb.WriteString(fmt.Sprintf("line %d ", l.line))
	sb.WriteString("near:\n")
	pos := start
	for i := start; i < end; i++ {
		str := l.replaceNewLine(rune(l.input[i]))
		sb.WriteString(str)
		if i < l.pos {
			pos += len(str)
		}
	}
	sb.WriteString("\n")
	for i := start; i < pos; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString("^")
	sb.WriteString("\n")
	return sb.String()
}

func (l *Lexer) PrintTokens() {
	for l.pos < len(l.input) {
		tok, err := l.ConsumeNextToken()
		if err != nil {
			return
		}

		fmt.Printf("Token%s: %s\n", tok.Type, tok.Value)
	}
}
