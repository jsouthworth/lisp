# A silly toy LISP interpreter based on the SICP analyzing interpreter.

This was written to demonstrate how to use closures to build embedded
DSLs in go.

## Run

```
$ go get jsouthworth.net/go/lisp
$ go run jsouthworth.net/go/lisp/cmd/repl
]=> (load "test.scm")
]=> (map square (list 1 2 3 4))
]=> (map fib (10 20 30 40))
```


## License
MIT see [LICENSE](LICENSE)

