(define (square x) (* x x))
(define (fib n) (define (helper a b n) (if (equal? n 0) b (helper (+ a b) a (- n 1)))) (helper 0 1 (+ n 1)))
(define (map fn xs) (if (null? xs) '() (cons (fn (car xs)) (map fn (cdr xs)))))
