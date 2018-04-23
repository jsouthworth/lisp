package lisp

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
)

var OK = Sym("OK")
var True Expr = Bool(true)
var False Expr = Bool(false)
var Nil Expr = null{}

type null struct{}

func (n null) Eval(e Env) Expr {
	return n
}

func (n null) String() string {
	return "'()"
}

type Primitive func(...Expr) Expr

func (p Primitive) Eval(e Env) Expr {
	return p
}

func (p Primitive) Apply(args ...Expr) Expr {
	return p(args...)
}

type lambda struct {
	fn     func(Env) Expr
	params []Sym
	env    Env
}

func makeLambda(params []Sym, definingEnv Env, fn func(Env) Expr) Expr {
	return lambda{fn: fn, params: params, env: definingEnv}
}

func (p lambda) Eval(e Env) Expr {
	return p.fn(e)
}

func (p lambda) Apply(arguments ...Expr) Expr {
	env := p.env.Extend(p.params, arguments)
	return p.Eval(env)
}

func (p lambda) String() string {
	return fmt.Sprintf("(lambda %v %p)", p.params, p.fn)
}

type Bool bool

func (b Bool) Eval(e Env) Expr {
	return b
}

type String string

func (s String) Eval(e Env) Expr {
	return s
}

func (s String) String() string {
	return fmt.Sprintf("%q", string(s))
}

type Number struct {
	rat *big.Rat
}

func Zero() Number {
	return Number{rat: new(big.Rat)}
}

func Num(rat *big.Rat) Expr {
	return Number{rat: rat}
}

func (n Number) Eval(e Env) Expr {
	return n
}

func (n Number) String() string {
	return n.rat.RatString()
}

type Sym string

func (s Sym) Eval(e Env) Expr {
	return s
}

func (s Sym) String() string {
	return string(s)
}

type quote struct{ Expr }

func Quote(e Expr) quote {
	return quote{e}
}

func Unquote(e Expr) Expr {
	return e.(quote).Expr
}

func (q quote) Eval(e Env) Expr {
	return q
}

func (q quote) String() string {
	return fmt.Sprintf("'%v", q.Expr)
}

type Proc func(Env) Expr

func (p Proc) Eval(e Env) Expr {
	return p(e)
}

func If(predicate, consequent, alternate Expr) Expr {
	return Proc(func(e Env) Expr {
		if IsTrue(predicate.Eval(e)) {
			return consequent.Eval(e)
		}
		return alternate.Eval(e)
	})
}

func Var(name Expr) Expr {
	return Proc(func(e Env) Expr {
		name := name.Eval(e).(Sym)
		return e.Lookup(name)
	})
}

func Set(name Expr, exp Expr) Expr {
	return Proc(func(e Env) Expr {
		name := name.Eval(e).(Sym)
		e.SetValue(name, exp.Eval(e))
		return OK
	})
}

func Define(name Expr, exp Expr) Expr {
	return Proc(func(e Env) Expr {
		name := name.Eval(e).(Sym)
		e.Define(name, exp.Eval(e))
		return OK
	})
}

func Sequence(exprs ...Expr) Expr {
	return Proc(func(e Env) Expr {
		var out Expr
		for _, expr := range exprs {
			out = expr.Eval(e)
		}
		return out
	})
}

func Lambda(parameters Expr, body ...Expr) Expr {
	return Proc(func(defining Env) Expr {
		bdy := Sequence(body...)
		return makeLambda(listToSlice(parameters), defining,
			func(calling Env) Expr {
				return bdy.Eval(calling)
			})
	})
}

func Apply(procedure Expr, arguments ...Expr) Expr {
	return Proc(func(e Env) Expr {
		args := make([]Expr, 0, len(arguments))
		for _, a := range arguments {
			args = append(args, a.Eval(e))
		}
		return procedure.Eval(e).(Applier).Apply(args...)
	})
}

func Load(file Expr) Expr {
	return Proc(func(e Env) Expr {
		name := string(file.(String))
		f, err := os.Open(name)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}
			e.Eval(Analyze(Read(text)))
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		return OK
	})
}

func IsTrue(exp Expr) Bool {
	return exp == True
}

type pair struct {
	head, tail Expr
}

func (c pair) Eval(e Env) Expr {
	return c
}

func (c pair) String() string {
	var buf bytes.Buffer
	p := c
	buf.WriteString("(")
	for {
		if tail, ok := p.tail.(pair); ok {
			fmt.Fprintf(&buf, "%v ", p.head)
			p = tail
		} else if p.tail == Nil {
			fmt.Fprintf(&buf, "%v", p.head)
			break
		} else {
			fmt.Fprintf(&buf, "%v . %v", p.head, p.tail)
			break
		}
	}
	buf.WriteByte(')')
	return buf.String()
}

func Cons(a, b Expr) Expr {
	return pair{head: a, tail: b}
}

func Car(e Expr) Expr {
	cell := e.(pair)
	return cell.head
}

func Cdr(e Expr) Expr {
	cell := e.(pair)
	return cell.tail
}

func List(vals ...Expr) Expr {
	if len(vals) == 0 {
		return Nil
	}
	return Cons(vals[0], List(vals[1:]...))
}
