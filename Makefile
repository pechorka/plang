GO:=go

.PHONY: repl

repl:
	$(GO) run cmd/repl/main.go
	