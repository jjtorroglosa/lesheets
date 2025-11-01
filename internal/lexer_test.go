package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexEmptySection(t *testing.T) {
	lex := NewLexer("Cmaj7 | !annotation!D Caug\nC\nE | F")

	tok, err := lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenChord, tok.Type)
	assert.Equal(t, "Cmaj7", tok.Value)

	tok, err = lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenBar, tok.Type)

	tok, err = lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenAnnotation, tok.Type)

	tok, err = lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenChord, tok.Type)
	assert.Equal(t, "D", tok.Value)

	tok, err = lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenChord, tok.Type)
	assert.Equal(t, "Caug", tok.Value)

	tok, err = lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenReturn, tok.Type)

	tok, err = lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenChord, tok.Type)
	assert.Equal(t, "C", tok.Value)
}

func TestLexBacktick(t *testing.T) {
	lex := NewLexer("`backtick`")
	toks, err := lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenBacktick, toks.Type)
	assert.Equal(t, "backtick", toks.Value)
}

func TestLexBacktickUnclosed(t *testing.T) {
	lex := NewLexer("`backtick")
	_, err := lex.ConsumeNextToken()
	assert.Equal(t, ErrGeneric(lex.SurroundingString(), "`", string(lex.nextChar())), err)
}

func TestLexMultilineBacktick(t *testing.T) {
	lex := NewLexer("```\nmy\nbacktick\n```")
	toks, err := lex.ConsumeNextToken()
	assert.NoError(t, err)
	assert.Equal(t, TokenBacktickMultiline, toks.Type)
	assert.Equal(t, "my\nbacktick\n", toks.Value)
}

func TestLexBacktickMultilineUnclosed(t *testing.T) {
	lex := NewLexer("```backtick``")
	_, err := lex.ConsumeNextToken()
	assert.Equal(t, ErrGeneric(lex.SurroundingString(), "Closing ```", string(lex.nextChar())), err)
}

func TestLexChord(t *testing.T) {
	lex := NewLexer("!annotation!Cmaj7 !second!D")

	tok, _ := lex.ConsumeNextToken()
	assert.Equal(t, TokenAnnotation, tok.Type)
	tok, _ = lex.ConsumeNextToken()
	assert.Equal(t, TokenChord, tok.Type)
	tok, _ = lex.ConsumeNextToken()
	assert.Equal(t, TokenAnnotation, tok.Type)
	tok, _ = lex.ConsumeNextToken()
	assert.Equal(t, TokenChord, tok.Type)
}
