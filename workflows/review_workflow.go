package workflows

import (
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/org"
	"sync"
)

type SingleRepoSyncReviewRequestsWorkflow struct {
	// Github repo info
	Name    string
	Owner   string
	Repo    string
	Filters []git_tools.PRFilter

	// org output info
	OrgFileName  string
	SectionTitle string
}

func (w SingleRepoSyncReviewRequestsWorkflow) GetName() string {
	return w.Name
}

func (w SingleRepoSyncReviewRequestsWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) {
	prs := git_tools.GetPRs(
		git_tools.GetGithubClient(),
		"open",
		w.Owner,
		w.Repo,
	)

	prs = git_tools.ApplyPRFilters(prs, w.Filters)
	doc := org.GetBaseOrgDocument(w.OrgFileName)
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return
	}
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		if output.ChangeType != "No Change" {
			file_change_wg.Add(1)
			c <- output
		}
	}
}

type SyncReviewRequestsWorkflow struct {
	// Github repo info
	Name    string
	Owner   string
	Repos   []string
	Filters []git_tools.PRFilter

	// org output info
	OrgFileName  string
	SectionTitle string
}

func (w SyncReviewRequestsWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) {
	client := git_tools.GetGithubClient()
	prs := git_tools.GetManyRepoPRs(client, "open", w.Owner, w.Repos)
	//	git_tools.GetGithubClient(),
	//	"open",
	//	w.Owner,
	//	w.Repo,
	// )
	prs = git_tools.ApplyPRFilters(prs, w.Filters)
	doc := org.GetBaseOrgDocument(w.OrgFileName)
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return
	}
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		if output.ChangeType != "No Change" {
			file_change_wg.Add(1)
			c <- output
		}
	}
}

func (w SyncReviewRequestsWorkflow) GetName() string {
	return w.Name
}

type ListMyPRsWorkflow struct {
	Name            string
	Owner           string
	Repos           []string
	OrgFileName     string
	SectionTitle    string
	PRState         string
	ReleasedVersion git_tools.DeployedVersion
}

func (w ListMyPRsWorkflow) GetName() string {
	return w.Name
}

func (w ListMyPRsWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) {
	client := git_tools.GetGithubClient()
	prs := git_tools.GetManyRepoPRs(client, w.PRState, w.Owner, w.Repos)

	doc := org.GetOrgDocument(w.OrgFileName, org.MergeInfoOrgSerializer{})
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return
	}
	prs = git_tools.ApplyPRFilters(prs, []git_tools.PRFilter{git_tools.MyPRs})
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		// TODO This is moving to the serializer
		if pr.MergedAt != nil && output.ChangeType != "No Change" {
			repo_name := *pr.Base.Repo.Name
			if repo_name == "repo" {
				released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
				fmt.Printf("Released PR: %s %t\n", *pr.Title, released)
				//output.Item.Details() = append(output.Lines, "Released: "+strconv.FormatBool(released))
				//output.Lines[0] = strings.Replace(output.Lines[0], "merged", "released", 1)
			}
		}
		if output.ChangeType != "No Change" {
			file_change_wg.Add(1)
			c <- output
		}

	}
}
