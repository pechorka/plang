package token

type Type string

const (
	// Identifiers + literals
	IDENT Type = "IDENT" // add, foobar, x, y, ...
	INT   Type = "INT"   // 1343456
	// Operators
	ASSIGN Type = "="
	PLUS   Type = "+"
	// Delimiters
	COMMA     Type = ","
	SEMICOLON Type = ";"
	LPAREN    Type = "("
	RPAREN    Type = ")"
	LBRACE    Type = "{"
	RBRACE    Type = "}"

	// Keywords
	FUNCTION Type = "FUNCTION"
	LET      Type = "LET"
)

type Token struct {
	Type    Type
	Literal string
}
