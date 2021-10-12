package token

type Type string

const (
	EOF     Type = "EOF"
	INVALID Type = "INVALID"
	// Identifiers + literals
	IDENT  Type = "IDENT" // add, foobar, x, y, ...
	INT    Type = "INT"   // 1343456
	STRING Type = "STRING"
	// Operators
	ASSIGN   Type = "="
	PLUS     Type = "+"
	MINUS    Type = "-"
	BANG     Type = "!"
	ASTERISK Type = "*"
	SLASH    Type = "/"
	LT       Type = "<"
	GT       Type = ">"
	EQ       Type = "=="
	NOT_EQ   Type = "!="
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
	TRUE     Type = "TRUE"
	FALSE    Type = "FALSE"
	IF       Type = "IF"
	ELSE     Type = "ELSE"
	RETURN   Type = "RETURN"
)

type Token struct {
	Type    Type
	Literal string
}
