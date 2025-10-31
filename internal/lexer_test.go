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
