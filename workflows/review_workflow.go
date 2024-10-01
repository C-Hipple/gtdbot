package workflows

import (
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/org"
	"strconv"
	"strings"
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
	Name            string
	Repos           []string
	Owner           string
	OrgFileName     string
	SectionTitle    string
	PRState         string
	ReleasedVersion git_tools.DeployedVersion
}

func (w ListMyPRsWorkflow) GetName() string {
	return w.Name
}

func (w ListMyPRsWorkflow) Run(c chan FileChanges, wg *sync.WaitGroup) {
	defer wg.Done()
	client := git_tools.GetGithubClient()
	prs := git_tools.GetManyRepoPRs(client, w.PRState, w.Owner, w.Repos)

	doc := org.GetOrgDocument(w.OrgFileName)
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return
	}
	prs = git_tools.ApplyPRFilters(prs, []git_tools.PRFilter{git_tools.MyPRs})
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		if pr.MergedAt != nil {
			// This is a hack, it should be when we make the FileChanges in SyncTODO section
			// but we'd need the released version and repo info for all repos for the workflows.
			repo_name := *pr.Base.Repo.Name
			if repo_name == "chaturbate" {
				released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
				output.Lines = append(output.Lines, "Released: "+strconv.FormatBool(released))
				output.Lines[0] = strings.Replace(output.Lines[0], "merged", "released", 1)

			} else {
				output.Lines = append(output.Lines, "Released: ???")
			}

			//fmt.Printf("PR: %v; Released: %v/n", pr.GetTitle(), released)
		}

		c <- output
	}
}
