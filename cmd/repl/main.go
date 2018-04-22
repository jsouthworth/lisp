package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jsouthworth/lisp"
)

func tryEval(expr string) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		fmt.Println(r)
	}()
	out := lisp.Eval(lisp.Analyze(lisp.Read(expr)))
	fmt.Println(out)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("]=> ")
		text, _ := reader.ReadString('\n')
		if len(text) == 0 {
			fmt.Println()
			return
		}
		expr := text[:len(text)-1]
		if len(expr) == 0 {
			continue
		}
		tryEval(expr)
	}
}
