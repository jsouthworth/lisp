package lisp

import (
	"fmt"
)

func listToExprSlice(list Expr) []Expr {
	out := []Expr{}
	if list == nil {
		return out
	}
	head, list := Car(list), Cdr(list)
	for {
		out = append(out, head)
		if list == nil {
			break
		}
		head, list = Car(list), Cdr(list)
	}
	return out
}

func listToSlice(list Expr) []Sym {
	out := []Sym{}
	if list == nil {
		return out
	}
	head, list := Car(list), Cdr(list)
	for {
		out = append(out, head.(Sym))
		if list == nil {
			break
		}
		head, list = Car(list), Cdr(list)
	}
	return out
}

func isSelfEvaluating(exp Expr) bool {
	switch exp.(type) {
	case Float, String, Int:
		return true
	default:
		return false
	}
}

func isVariable(exp Expr) bool {
	_, ok := exp.(Sym)
	return ok
}

func isTaggedList(exp Expr, tag Sym) bool {
	_, ok := exp.(pair)
	if ok {
		car := Car(exp)
		str, ok := car.(Sym)
		if ok {
			return str == tag
		}
		return false
	}
	return false
}

func isQuoted(exp Expr) bool {
	return isTaggedList(exp, "quote")
}

func isAssignment(exp Expr) bool {
	return isTaggedList(exp, "set!")
}

func isDefinition(exp Expr) bool {
	return isTaggedList(exp, "define")
}

func isIf(exp Expr) bool {
	return isTaggedList(exp, "if")
}

func isLambda(exp Expr) bool {
	return isTaggedList(exp, "lambda")
}

func isBegin(exp Expr) bool {
	return isTaggedList(exp, "begin")
}

func isApplication(exp Expr) bool {
	_, ok := exp.(pair)
	return ok
}

func analyzeSequence(exp Expr) Expr {
	exprs := listToExprSlice(exp)
	for i, expr := range exprs {
		exprs[i] = Analyze(expr)
	}
	return Sequence(exprs...)
}

func definitionVariable(exp Expr) Expr {
	if _, ok := Car(Cdr(exp)).(Sym); ok {
		return Car(Cdr(exp))
	}
	//Handle special lambda definitions
	return Car(Car(Cdr(exp)))
}

func definitionValue(exp Expr) Expr {
	if _, ok := Car(Cdr(exp)).(Sym); ok {
		return Analyze(Car(Cdr(Cdr(exp))))
	}
	//Handle special lambda definitions
	body := listToExprSlice(Cdr(Cdr(exp)))
	for i, expr := range body {
		body[i] = Analyze(expr)
	}
	return Lambda(Cdr(Car(Cdr(exp))), body...)
}

func Analyze(exp Expr) Expr {
	switch {
	case isSelfEvaluating(exp):
		return exp
	case isVariable(exp):
		return Var(exp)
	case isQuoted(exp):
		return Quote(Car(Cdr(exp)))
	case isAssignment(exp):
		return Set(Car(Cdr(exp)),
			Analyze(Car(Cdr(Cdr(exp)))))
	case isDefinition(exp):
		return Define(definitionVariable(exp),
			definitionValue(exp))
	case isIf(exp):
		consequent := False
		if Cdr(Cdr(Cdr(exp))) != Nil {
			consequent = Car(Cdr(Cdr(Cdr(exp))))
		}
		return If(Analyze(Car(Cdr(exp))),
			Analyze(Car(Cdr(Cdr(exp)))),
			Analyze(consequent))
	case isLambda(exp):
		exprs := listToExprSlice(Cdr(Cdr(exp)))
		for i, expr := range exprs {
			exprs[i] = Analyze(expr)
		}
		return Lambda(Car(Cdr(exp)), exprs...)
	case isBegin(exp):
		return analyzeSequence(Cdr(exp))
	case isApplication(exp):
		operands := listToExprSlice(Cdr(exp))
		for i, operand := range operands {
			operands[i] = Analyze(operand)
		}
		return Apply(Analyze(Car(exp)), operands...)
	default:
		panic(fmt.Errorf("unknown expression type: %v", exp))
	}
}
