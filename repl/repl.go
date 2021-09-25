package repl

import (
	"bufio"
	"fmt"

	"io"

	"github.com/pechorka/plang/lexer"
	"github.com/pechorka/plang/token"
)

const PROMPT = ">> "

func Start(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for {
		fmt.Fprintf(w, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.NewFromString(line)
		for tok := l.Next(); tok.Type != token.EOF; tok = l.Next() {
			fmt.Fprintf(w, "%+v\n", tok)
		}
	}
}
