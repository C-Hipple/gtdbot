package main

import "fmt"

type SyncReviewRequestsWorkflow struct {
	// Github repo info
	owner   string
	repo    string
	filters []PRFilter

	// org output info
	org_file_name string
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
	doc := GetOrgDocument(w.org_file_name)

	section, err := doc.GetSection(w.section_title)
	if err != nil {
		fmt.Println("Error getting section: ", err)
		c <- idx
		return
	}
	for _, pr := range prs {
		section.AddTODO(PRToOrgBridge{pr})
	}
	c <- idx

}
