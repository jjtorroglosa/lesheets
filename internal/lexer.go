package internal

import (
	"strings"
	"unicode"
)

type TokenType string

const (
	TokenFrontMatterStart TokenType = "FrontMatterStart"
	TokenFrontMatterEnd   TokenType = "FrontMatterEnd"
	TokenYAMLKey          TokenType = "YAMLKey"
	TokenYAMLValue        TokenType = "YAMLValue"
	TokenHeader           TokenType = "Header"
	TokenHeaderBreak      TokenType = "HeaderBreak"
	TokenChord            TokenType = "Chord"
	TokenSymbol           TokenType = "Symbol"
	TokenBar              TokenType = "Bar"
	TokenReturn           TokenType = "Return"
	TokenComment          TokenType = "Comment"
	TokenAnnotation       TokenType = "Annotation"
	TokenBacktick         TokenType = "BacktickExpression"
	TokenUnknown          TokenType = "Unknown"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input         string
	pos           int
	inFrontMatter bool
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

func (l *Lexer) advance() {
	l.pos++
}

// Scan all tokens in the input
func (l *Lexer) Lex() []Token {
	var tokens []Token
	for l.pos < len(l.input) {
		ch := l.nextChar()

		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advance()
			continue
		}

		// Frontmatter
		if strings.HasPrefix(l.input[l.pos:], "---") {
			if !l.inFrontMatter {
				tokens = append(tokens, Token{Type: TokenFrontMatterStart, Value: "---"})
				l.inFrontMatter = true
			} else {
				tokens = append(tokens, Token{Type: TokenFrontMatterEnd, Value: "---"})
				l.inFrontMatter = false
			}
			l.pos += 3
			continue
		}

		// Single bar
		if ch == '|' {
			if strings.HasPrefix(l.input[l.pos:], "||") {
				tokens = append(tokens, Token{Type: TokenBar, Value: "||"})
				l.pos += 2
			} else {
				tokens = append(tokens, Token{Type: TokenBar, Value: "|"})
				l.advance()
			}
			ch := l.nextChar()
			if ch == '\n' {
				tokens = append(tokens, Token{Type: TokenReturn, Value: "RETURN"})
				l.advance()
			}
			continue
		}

		// Annotation
		if ch == '!' {
			l.advance()
			start := l.pos
			for l.pos < len(l.input) && !unicode.IsSpace(rune(l.input[l.pos])) && l.input[l.pos] != '!' {
				l.advance()
			}
			tokens = append(tokens, Token{Type: TokenAnnotation, Value: l.input[start:l.pos]})
			continue
		}

		// Comment
		if ch == '"' {
			start := l.pos
			l.advance() // skip opening `
			for l.pos < len(l.input) && l.input[l.pos] != '"' {
				l.advance()
			}
			if l.pos < len(l.input) && l.input[l.pos] == '"' {
				l.advance() // skip closing `
			}
			tokens = append(tokens, Token{Type: TokenComment, Value: l.input[start:l.pos]})
			continue
		}
		// Backtick expression
		if ch == '`' {
			start := l.pos
			l.advance() // skip opening `
			for l.pos < len(l.input) && l.input[l.pos] != '`' {
				l.advance()
			}
			if l.pos < len(l.input) && l.input[l.pos] == '`' {
				l.advance() // skip closing `
			}
			tokens = append(tokens, Token{Type: TokenBacktick, Value: l.input[start:l.pos]})
			continue
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
			for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
				l.advance()
			}
			end := l.pos
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.advance()
			}
			tokens = append(tokens, Token{Type: tokenType, Value: strings.TrimSpace(l.input[end:l.pos])})
			continue
		}

		// YAML key/value
		if l.inFrontMatter {
			start := l.pos
			end := l.pos
			for end < len(l.input) && l.input[end] != '\n' {
				end++
			}
			line := l.input[start:end]
			parts := strings.SplitN(line, ":", 2)
			tokens = append(tokens, Token{Type: TokenYAMLKey, Value: strings.TrimSpace(parts[0])})
			tokens = append(tokens, Token{Type: TokenYAMLValue, Value: strings.TrimSpace(parts[1])})
			l.pos += len(line)
			continue
		}

		// Chords or symbols like %
		start := l.pos
		for l.pos < len(l.input) && !unicode.IsSpace(rune(l.input[l.pos])) && l.input[l.pos] != '|' {
			l.advance()
		}
		value := l.input[start:l.pos]
		tokens = append(tokens, Token{Type: TokenChord, Value: value})
	}

	return tokens
}
