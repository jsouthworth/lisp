package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jsouthworth/lisp"
)

func tryEval(expr string) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		fmt.Println(r)
		fmt.Println(string(debug.Stack()))
	}()
	if expr == "" {
		return
	}
	out := lisp.Eval(lisp.Analyze(lisp.Read(expr)))
	fmt.Println(out)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("]=> ")
	for scanner.Scan() {
		tryEval(scanner.Text())
		fmt.Print("]=> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Println()
}
