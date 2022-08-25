(display 1)
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



