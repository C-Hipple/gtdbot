;;;###autoload
(defun run-gtdbot-oneoff ()
  "Runs gtdbot with the oneoff flag to update reviews.org"
  (interactive)
  (async-shell-command "gtdbot --oneoff" "*gtdbot*"))

(define-key evil-normal-state-map (kbd ", r l") 'run-gtdbot-oneoff) ;; l for list?


;;;###autoload
(defun run-gtdbot-service ()
  "Runs gtdbot with the oneoff flag to update reviews.org"
  (interactive)
  (async-shell-command "gtdbot" "*gtdbot*"))

(define-key evil-normal-state-map (kbd ", r S") 'run-gtdbot-service) ;; s I already have bound to review start at url


;; Theese are testing helper functions to make development a little bit easier
;;;###autoload
(defun run-gtdbot-parse-test()
  "Runs gtdbot with the parse flag to check parsing reviews.org"
  (interactive)
  (async-shell-command "gtdbot --parse" "*gtdbot*"))

(define-key evil-normal-state-map (kbd ", r p") 'run-gtdbot-parse-test)


;;;###autoload
(defun reset-reviews-org ()
  (interactive)
  (shell-command "cp ~/gtdbot/reviews_template.org ~/gtd/reviews.org"))
;; (shell-command "cp ~/reviews_template_testdata.org ~/gtd/reviews.org"))

(define-key evil-normal-state-map (kbd ", r b") 'reset-reviews-org)

;; probably don't need this often
(define-key evil-normal-state-map (kbd "SPC g B") (lambda () (interactive) (switch-to-buffer "*gtdbot*" nil t)))
