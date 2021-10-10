package repl

import (
	"bufio"
	"fmt"

	"io"

	"github.com/pechorka/plang/lexer"
	"github.com/pechorka/plang/parser"
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
		p := parser.New(l)

		program := p.Parse()
		if len(p.Errors()) != 0 {
			printParserErrors(w, p.Errors())
			continue
		}
		io.WriteString(w, program.String())
		io.WriteString(w, "\n")
	}
}

func printParserErrors(w io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(w, "\t"+msg+"\n")
	}
}
