package workflows

import (
	"fmt"
	"gtdbot/github"
	"gtdbot/org"
	"sync"
)

type SyncReviewRequestsWorkflow struct {
	// Github repo info
	Name    string
	Owner   string
	Repo    string
	Filters []github.PRFilter

	// org output info
	OrgFileName  string
	SectionTitle string
}

func (w SyncReviewRequestsWorkflow) Run(c chan FileChanges, idx int, wg *sync.WaitGroup) {
	defer wg.Done()
	prs := github.GetPRs(
		github.GetGithubClient(),
		"open",
		w.Owner,
		w.Repo,
	)
	prs = github.ApplyPRFilters(prs, w.Filters)
	doc := org.GetOrgDocument(w.OrgFileName)
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return
	}
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		c <- output
	}
}

func (w SyncReviewRequestsWorkflow) GetName() string {
	return w.Name
}
