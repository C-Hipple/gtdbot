;;; gtdbot.el --- Tools for implementing Getting Things Done methodology in emacs -*- lexical-binding: t; -*-

;; Author: Chris Hipple
;; URL: https://github.com/C-Hipple/gtdbot
;; Version: 0.1.0
;; Package-Requires: ((emacs "25.1"))

;; SPDX-License-Identifier: GPL-3.0+

;;; Commentary:

;; This package provides tooling for running gtdbot in (defined in the same repo, built with go) and for reviewing closed / archived tasks
;;
;; Usage:

;; For running gtbot, simply call run-gtdbot-oneoff or run-gtdbot-service depending on what you want.  Service runs the sync every minute (configurable via toml)

;; For using the review functionality, open one of your files called out in org-agenda list and then call (org-agenda).  There you will have a new binding of R to enter your review periods.

;;; Code:
;;;###autoload
(defun delta-wash()
  "interactive of the call magit-delta function if you have the package installed."
  (interactive)
  (if (functionp 'magit-delta-call-delta-and-convert-ansi-escape-sequences)
      (magit-delta-call-delta-and-convert-ansi-escape-sequences)
    (message "You do not have magit-delta installed!")))



(defun gtdbot--callback (x)
  ;; TODO: Allow for setting the reviews file
  (with-current-buffer (find-file-noselect "~/gtd/reviews.org")
    (if (functionp 'magit-delta-call-delta-and-convert-ansi-escape-sequences)
        (progn
          (magit-delta-call-delta-and-convert-ansi-escape-sequences)
          (save-buffer))))
  (message "gtdbot sync complete!"))

;;;###autoload
(defun run-gtdbot-oneoff ()
  "Runs gtdbot with the oneoff flag to update reviews.org"
  (interactive)
  (async-start-process "*gtdbot-async*" "gtdbot" 'gtdbot--callback "--oneoff"
                       :stderr (get-buffer-create "*gtdbot-async::stderr*")))


;;;###autoload
(defun run-gtdbot-service ()

  "Runs gtdbot with the oneoff flag to update reviews.org"
  (interactive)
  (async-shell-command "gtdbot" "*gtdbot*"))

;; I only bind for evil-mode
(if (not (eq evil-normal-state-map nil))
    (progn
      (define-key evil-normal-state-map (kbd ", r S") 'run-gtdbot-service) ;; s I already have bound to review start at url
      (define-key evil-normal-state-map (kbd ", r l") 'run-gtdbot-oneoff) ;; l for list?
      (define-key evil-normal-state-map (kbd ", r d") 'delta-wash)))


;; Theese are testing helper functions to make development a little bit easier
;;;###autoload
(defun run-gtdbot-parse-test()
  "Runs gtdbot with the parse flag to check parsing reviews.org"
  (interactive)
  (async-shell-command "gtdbot --parse" "*gtdbot*"))

(define-key evil-normal-state-map (kbd ", r p") 'run-gtdbot-parse-test)

;; Below this point is the code for doing org agenda reviews in emacs
;; To be honest, I copied this from a blog which seems to have been since taken down and I can't find it.

(setq gtdbot-org-agenda-files '("~/gtd/inbox.org" "~/gtd/gtd.org" "~/gtd/notes.org" "~/gtd/next_actions.org" "~/gtd/reviews.org"))

;; define "R" as the prefix key for reviewing what happened in various time periods
(setq org-agenda-custom-commands
      '(("c" "Work" tags-todo "work" ;; (1) (2) (3) (4)
         ((org-agenda-files gtdbot-org-agenda-files) ) ;; (5)
         (org-agenda-sorting-strategy '(priority-up effort-down))) ;; (5) cont.
        ("~/computer.html"))) ;; (6)

(add-to-list 'org-agenda-custom-commands
             '("R" . "Review" ))

;; Common settings for all reviews
(setq efs/org-agenda-review-settings
      '((org-agenda-files gtdbot-org-agenda-files)
        (org-agenda-show-all-dates t)
        (org-agenda-start-with-log-mode t)
        (org-agenda-start-with-clockreport-mode t)
        (org-agenda-archives-mode t)
        ;; I don't care if an entry was archived
        (org-agenda-hide-tags-regexp
         (concat org-agenda-hide-tags-regexp
                 "\\|ARCHIVE"))
        ))
;; Show the agenda with the log turn on, the clock table show and
;; archived entries shown.  These commands are all the same exept for
;; the time period.
(add-to-list 'org-agenda-custom-commands
             `("Rw" "Week in review"
               agenda ""
               ;; agenda settings
               ,(append
                 efs/org-agenda-review-settings
                 '((org-agenda-span 'week)
                   (org-agenda-start-on-weekday 0)
                   (org-agenda-overriding-header "Week in Review"))
                 )
               ("~//review/week.html")
               ))


(add-to-list 'org-agenda-custom-commands
             `("Rd" "Day in review"
               agenda ""
               ;; agenda settings
               ,(append
                 efs/org-agenda-review-settings
                 '((org-agenda-span 'day)
                   (org-agenda-overriding-header "Day in Review"))
                 )
               ("~//review/day.html")
               ))

(add-to-list 'org-agenda-custom-commands
             `("Rm" "Month in review"
               agenda ""
               ;; agenda settings
               ,(append
                 efs/org-agenda-review-settings
                 '((org-agenda-span 'month)
                   (org-agenda-start-day "01")
                   (org-read-date-prefer-future nil)
                   (org-agenda-overriding-header "Month in Review"))
                 )
               ("~//review/month.html")
               ))

(add-to-list 'org-agenda-custom-commands
             `("Rs" "Sprint in review"
               agenda ""
               ;; agenda settings
               ,(append
                 efs/org-agenda-review-settings
                 '((org-agenda-span 10)
                   (org-agenda-start-day "-6d")
                   (org-agenda-start-on-weekday nil)
                   (org-agenda-overriding-header "Sprint Review"))
                 )
               ("~//review/sprint.html")
               ))
