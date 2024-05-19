package workflows

import (
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/org"
	"sync"
)

type SyncReviewRequestsWorkflow struct {
	// Github repo info
	Name    string
	Owner   string
	Repo    string
	Filters []git_tools.PRFilter

	// org output info
	OrgFileName  string
	SectionTitle string
}

func (w SyncReviewRequestsWorkflow) Run(c chan FileChanges, wg *sync.WaitGroup) {
	defer wg.Done()
	prs := git_tools.GetPRs(
		git_tools.GetGithubClient(),
		"open",
		w.Owner,
		w.Repo,
	)
	prs = git_tools.ApplyPRFilters(prs, w.Filters)
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

type ListMyPRsWorkflow struct {
	Repos        []string
	Owner        string
	OrgFileName  string
	SectionTitle string
	PRState      string
}

func (w ListMyPRsWorkflow) GetName() string {
	return "List My PRs"
}

func (w ListMyPRsWorkflow) Run(c chan FileChanges, wg *sync.WaitGroup) {
	defer wg.Done()
	client := git_tools.GetGithubClient()
	prs := git_tools.GetManyPrs(client, w.PRState, w.Owner, w.Repos)

	doc := org.GetOrgDocument(w.OrgFileName)
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return
	}
	prs = git_tools.ApplyPRFilters(prs, []git_tools.PRFilter{git_tools.MyPRs})
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		c <- output
	}
}
