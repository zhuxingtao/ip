(define ltNumbers (
        lambda (pivot seq) 
            (if 
                (null? seq) 
                '()
                    (if (< (car seq) pivot)  
                            (cons (car seq)  (ltNumbers pivot  (cdr seq)))   
                           (ltNumbers pivot (cdr seq))
                    )
            )
       )
)

(define equalNumbers (
     lambda (pivot seq) 
            (if 
                (null? seq) 
                '()
                    (if (= (car seq) pivot)  
                            (cons (car seq)  (equalNumbers pivot  (cdr seq)))   
                           (equalNumbers pivot (cdr seq))
                    )
            )
       )
)

(define gtNumbers (
        lambda (pivot seq) 
            (if 
                (null? seq) 
                (list)
                    (if (> (car seq) pivot)  
                            (cons (car seq)  (gtNumbers pivot  (cdr seq)))   
                            (gtNumbers pivot (cdr seq))
                    )
                seq     
            )
       )
)



(define qsort (
    lambda (seq) 
        (if 
            (<= (length seq) 1) 
            seq
            (append (append (qsort (ltNumbers (car seq) seq)) (equalNumbers (car seq) seq)) (qsort (gtNumbers (car seq) seq)))
        )    
  )
)

(define a ( 7 3 4 9 2 2 97 5 6 35 35 43 ))
;(ltNumbers (car a) a)
;(gtNumbers (car a) a)
(qsort a)

