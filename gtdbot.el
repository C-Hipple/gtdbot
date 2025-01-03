(defun reset-reviews-org ()
  (interactive)
  (async-shell-command "cp ~/gtdbot/reviews_template.org ~/gtd/reviews.org")
  )

(define-key evil-normal-state-map (kbd ", r b") 'reset-reviews-org)


(defun run-gtdbot-parse-test()
  "Runs gtdbot with the parse flag to check parsing reviews.org"
  (interactive)
  (async-shell-command "gtdbot --parse" "*gtdbot*"))


(define-key evil-normal-state-map (kbd ", r p") 'run-gtdbot-parse-test)
