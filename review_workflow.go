package main

type SyncReviewRequestsWorkflow struct {
	// Github repo info
	owner   string
	repo    string
	filters []PRFilter

	// org output info
	org_file      string
	section_title string
}

func (w SyncReviewRequestsWorkflow) Run(c chan int, idx int) {
	prs := getPRs(
		GetGithubClient(),
		"open",
		w.owner,
		w.repo,
	)
	prs = ApplyPRFilters(prs, w.filters)
	c <- idx
}
