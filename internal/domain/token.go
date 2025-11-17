package domain

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
