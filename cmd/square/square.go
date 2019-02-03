package main

import (
	"fmt"
	"math/big"

	. "jsouthworth.net/go/lisp"
)

func main() {
	Eval(Define(Sym("square"),
		Lambda(List(Sym("x")),
			Apply(Var(Sym("*")), Var(Sym("x")), Var(Sym("x"))))))
	fmt.Println(Eval(Apply(Var(Sym("square")), Num(big.NewRat(10, 1)))))

}
