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

func TestLexSection(t *testing.T) {
	testCases := []struct {
		input        string
		pos          int
		posNotHeader int
	}{
		{
			input:        "# section",
			pos:          0,
			posNotHeader: -1,
		},
		{
			input:        "\n# section",
			pos:          1,
			posNotHeader: -1,
		},
		{
			input:        " #thisshouldnotbeaheadertoken\n# section",
			pos:          2,
			posNotHeader: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			lex := NewLexer(tC.input)
			toks, err := lex.Lex()
			assert.NoError(t, err)
			assert.Equal(t, "section", toks[tC.pos].Value)
			assert.True(t, tC.posNotHeader == -1 || toks[tC.posNotHeader].Type != TokenHeader)
		})
	}
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

func TestLexIgnoresComments(t *testing.T) {
	lex := NewLexer(`D
// this is a comment
E`)
	lex.PrintTokens()
	lex.pos = 0
	expected := []TokenType{
		TokenChord, TokenReturn, TokenChord, TokenReturn,
	}
	i := 0
	tok, err := lex.ConsumeNextToken()
	for tok.Type != TokenEof {
		assert.NoError(t, err)
		assert.Equalf(t, expected[i], tok.Type, "want %s got %s with val %s; tok: %d", expected[i], tok.Type, tok.Value, i)
		i++
		tok, err = lex.ConsumeNextToken()
	}
}

func TestChordStartingWithSharpShouldNotBeTreatedAsHeaderIfInTheMiddleOfLine(t *testing.T) {
	testCases := []struct {
		line string
	}{
		{line: "|#1"},
		{line: "A #1"},
		{line: " #1"},
	}
	for _, tC := range testCases {
		t.Run(tC.line, func(t *testing.T) {
			lex := NewLexer(tC.line)
			found := false
			tokens, err := lex.Lex()
			assert.NoError(t, err)
			for _, tok := range tokens {
				if tok.Type == TokenChord && tok.Value == "#1" {
					found = true
				}
			}
			assert.True(t, found)
		})
	}
}

func TestLexSongTwoConsecutiveRepeats(t *testing.T) {
	lex := NewLexer("||: A :|| ||: B |")
	lex.PrintTokens()
}

func TestLexCommentDoesNotConsumeTokenReturnIfInline(t *testing.T) {
	lex := NewLexer(`D // this is a comment, it shouldn't consume the TokenReturn
E
`)
	lex.pos = 0
	expected := []TokenType{
		TokenChord, TokenReturn, TokenChord, TokenReturn,
	}
	i := 0
	tok, err := lex.ConsumeNextToken()
	for tok.Type != TokenEof {
		assert.NoError(t, err)
		assert.Equalf(t, expected[i], tok.Type, "want %s got %s with val %s; tok: %d", expected[i], tok.Type, tok.Value, i)
		i++
		tok, err = lex.ConsumeNextToken()
	}
}

func TestLexUnclosedAnnotation(t *testing.T) {
	lex := NewLexer("!unclosedannotation")
	_, err := lex.ConsumeNextToken()
	assert.Error(t, err)
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
