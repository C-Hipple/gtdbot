(defun reset-reviews-org ()
  (interactive)
  (async-shell-command "cp ~/gtdbot/reviews_template.org ~/gtd/reviews.org")
  )

(define-key evil-normal-state-map (kbd ", r b") 'reset-reviews-org)
