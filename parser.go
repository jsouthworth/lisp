package lisp

import (
	"strconv"
	"unicode/utf8"
)

func readRun(input string) (string, string) {
	var tlen int
	var result string = input
Loop:
	for len(result) != 0 {
		r, sz := utf8.DecodeRuneInString(result)
		switch r {
		case ' ', '(', ')':
			break Loop
		default:
			tlen += sz
			result = result[sz:]
		}
	}
	return input[0:tlen], input[tlen:]
}

func readString(input string) (string, string) {
	var tlen int
	var result string = input
	_, sz := utf8.DecodeRuneInString(result)
	tlen += sz
	result = result[sz:]
Loop:
	for len(result) != 0 {
		r, sz := utf8.DecodeRuneInString(result)
		switch r {
		case '"':
			tlen += sz
			break Loop
		default:
			tlen += sz
			result = result[sz:]
		}
	}
	return input[0:tlen], input[tlen:]
}

func readToken(input string) (string, string, bool) {
	var token string
	r, sz := utf8.DecodeRuneInString(input)
	switch r {
	case ' ': //skip all spaces
		input = input[sz:]
		return "", input, false
	case '(', ')', '\'': //bracket
		token = input[0:sz]
		input = input[sz:]
		return token, input, true
	case '"':
		token, input = readString(input)
		return token, input, true
	default:
		token, input = readRun(input)
		return token, input, true
	}
}

func tokenize(input string) []string {
	var token string
	var ok bool
	tokens := []string{}
	for len(input) != 0 {
		token, input, ok = readToken(input)
		if ok {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func parseTokens(tokens []string) Expr {
	if len(tokens) == 0 {
		return Nil
	}
	elem, tokens := parseToken(tokens)
	if len(tokens) == 0 {
		return elem
	}
	return Cons(elem, parseTokens(tokens))
}

func parseList(end Expr, tokens []string) (Expr, []string) {
	if len(tokens) == 0 {
		return end, tokens
	}
	if tokens[0] == ")" {
		return end, tokens[1:]
	}
	out, tokens := parseToken(tokens)
	li, tokens := parseList(end, tokens)
	return Cons(out, li), tokens
}

func parseQuote(tokens []string) (Expr, []string) {
	token, rest := parseToken(tokens)
	if token == Nil {
		return Nil, rest
	}
	token = List(Sym("quote"), token)
	return token, rest
}

func parseToken(tokens []string) (Expr, []string) {
	if len(tokens) == 0 {
		panic("unexpected EOF")
	}
	token, tokens := tokens[0], tokens[1:]
	switch token {
	case "(":
		return parseList(Nil, tokens)
	case "'":
		return parseQuote(tokens)
	case ")":
		panic("unexpected )")
	default:
		return parseAtom(token), tokens
	}
}

func parseAtom(str string) Expr {
	i, err := strconv.ParseInt(str, 10, 32)
	if err == nil {
		return Int(i)
	}
	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return Float(f)
	}
	if len(str) > 0 && str[0] == '"' {
		if str[len(str)-1] != '"' {
			panic("malformed string " + str)
		}
		return String(str[1 : len(str)-1])
	}
	return Sym(str)
}

func Read(program string) Expr {
	return parseTokens(tokenize(program))
}
