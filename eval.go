package lisp

type Expr interface {
	Eval(e Env) Expr
}

type Applier interface {
	Apply(args ...Expr) Expr
}

func Eval(exp Expr) Expr {
	return packageEnv.Eval(exp)
}
