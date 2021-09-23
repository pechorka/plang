package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pechorka/plang/lexer"
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
		for l.Next() {
			fmt.Fprintf(w, "%+v\n", l.Token())
		}
		if l.Err() != io.EOF {
			fmt.Fprintf(w, "unxpected error while reading input: %v", l.Err())
		}
	}
}
