;the test program writing in scheme


(display 2)
(
    define ;;this is comment;take  (lambda (n seq) (
            if (<= n 0)  (
                quote ()
            ) (
                cons (car seq)  (take (- n 1) (cdr seq))
            )
        )
    )
)
(* pi 3 )
(take 2 (1 2 3))




(define combine (lambda (f)
    (lambda (x y)
      (if (null? x) (quote ())
          (f (list (car x) (car y))
             ((combine f) (cdr x) (cdr y)))))))

(define zip (combine cons))
(quote ())
(zip (list 4 5 6) (list 1 2 3))


(define riff-shuffle (lambda (deck) (begin
    (define take (lambda (n seq) (if (<= n 0) (quote ()) (cons (car seq) (take (- n 1) (cdr seq))))))
    (define drop (lambda (n seq) (if (<= n 0) seq (drop (- n 1) (cdr seq)))))
    (define mid (lambda (seq) (/ (length seq) 2)))
    ((combine append) (take (mid deck) deck) (drop (mid deck) deck)))))


(riff-shuffle (list 1 2 3 4 5 6 7 8))
