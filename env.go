package lisp

import (
	"fmt"
	"reflect"
)

type Env struct {
	parent *Env
	m      map[Sym]Expr
}

func EnvNew() Env {
	return Env{
		m: make(map[Sym]Expr),
	}
}

func (e *Env) Extend(vars []Sym, vals []Expr) Env {
	out := EnvNew()
	out.parent = e
	if len(vars) == len(vals) {
		for i, v := range vars {
			out.m[v] = vals[i]
		}
	} else if len(vars) < len(vals) {
		panic(fmt.Errorf("Too many arguments supplied: %v, %v",
			vars, vals))
	} else {
		panic(fmt.Errorf("Too few arguments supplied: %v, %v",
			vars, vals))
	}
	return out
}

func (e *Env) Define(name Sym, val Expr) {
	e.m[name] = val
}

func (e *Env) Lookup(name Sym) Expr {
	out, ok := e.m[name]
	if ok {
		return out
	}
	if e.parent == nil {
		panic("undefined variable: " + string(name))
	}
	return e.parent.Lookup(name)
}

func (e *Env) SetValue(name Sym, val Expr) {
	_, ok := e.m[name]
	if ok {
		e.m[name] = val
		return
	}
	if e.parent == nil {
		panic("undefined variable: " + string(name))
	}
	e.parent.SetValue(name, val)
}

func (e *Env) Eval(exp Expr) Expr {
	return exp.Eval(*e)
}

func InitEnv() Env {
	e := EnvNew()
	e.Define("eval", Primitive(func(args ...Expr) Expr {
		return e.Eval(Analyze(Unquote(args[0])))
	}))
	e.Define("apply", Primitive(func(args ...Expr) Expr {
		return e.Eval(args[0].(Applier).Apply(args[1:]...))
	}))
	e.Define("car", Primitive(func(args ...Expr) Expr {
		return Car(args[0])
	}))
	e.Define("cdr", Primitive(func(args ...Expr) Expr {
		return Cdr(args[0])
	}))
	e.Define("cons", Primitive(func(args ...Expr) Expr {
		return Cons(args[0], args[1])
	}))
	e.Define("list", Primitive(func(args ...Expr) Expr {
		return List(args...)
	}))
	e.Define("+", Primitive(func(args ...Expr) Expr {
		sum := args[0].(Int)
		for _, a := range args[1:] {
			sum = sum + a.(Int)
		}
		return sum
	}))
	e.Define("-", Primitive(func(args ...Expr) Expr {
		diff := args[0].(Int)
		for _, a := range args[1:] {
			diff = diff - a.(Int)
		}
		return diff
	}))
	e.Define("*", Primitive(func(args ...Expr) Expr {
		result := args[0].(Int)
		for _, a := range args[1:] {
			result = result * a.(Int)
		}
		return result
	}))
	e.Define("/", Primitive(func(args ...Expr) Expr {
		result := args[0].(Int)
		for _, a := range args[1:] {
			result = result / a.(Int)
		}
		return result
	}))
	e.Define("pair?", Primitive(func(args ...Expr) Expr {
		_, ok := args[0].(pair)
		return Bool(ok)
	}))
	e.Define("null?", Primitive(func(args ...Expr) Expr {
		return Bool(args[0] == nil)
	}))
	e.Define("equal?", Primitive(func(args ...Expr) Expr {
		return Bool(reflect.DeepEqual(args[0], args[1]))
	}))
	return e
}

var packageEnv = InitEnv()
