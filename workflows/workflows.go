package workflows

import (
	"errors"
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/jira"
	"gtdbot/org"
	"sync"
)

type RunResult struct {
	Added   int
	Updated int
	Removed int // TODO
}

func (rr *RunResult) Process(output *FileChanges, c chan FileChanges, wg *sync.WaitGroup) {
	if output.ChangeType != "No Change" {
		if output.ChangeType == "Update" {
			rr.Updated += 1
		} else if output.ChangeType == "Addition" {
			rr.Added += 1
		}
		wg.Add(1)
		c <- *output
	}
}

func (rr *RunResult) Report() string {
	return fmt.Sprintf("A: %d; U: %d; R: %d", rr.Added, rr.Updated, rr.Removed)
}

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

func (w SingleRepoSyncReviewRequestsWorkflow) GetOrgFilename() string {
	return w.OrgFileName
}

func (w SingleRepoSyncReviewRequestsWorkflow) GetOrgSectionName() string {
	return w.SectionTitle
}

func (w SingleRepoSyncReviewRequestsWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) (RunResult, error) {
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
		return RunResult{}, errors.New("Section Not Found")
	}

	result := ProcessPRs(prs, c, &doc, &section, file_change_wg, true)
	return result, nil
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

func (w SyncReviewRequestsWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) (RunResult, error) {
	client := git_tools.GetGithubClient()
	prs := git_tools.GetManyRepoPRs(client, "open", w.Owner, w.Repos)
	prs = git_tools.ApplyPRFilters(prs, w.Filters)
	doc := org.GetBaseOrgDocument(w.OrgFileName)
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return RunResult{}, errors.New("Section Not Found")
	}
	result := ProcessPRs(prs, c, &doc, &section, file_change_wg, true)
	return result, nil
}

func (w SyncReviewRequestsWorkflow) GetName() string {
	return w.Name
}

func (w SyncReviewRequestsWorkflow) GetOrgFilename() string {
	return w.OrgFileName
}

func (w SyncReviewRequestsWorkflow) GetOrgSectionName() string {
	return w.SectionTitle
}

type ListMyPRsWorkflow struct {
	Name            string
	Owner           string
	Repos           []string
	Filters         []git_tools.PRFilter
	OrgFileName     string
	SectionTitle    string
	PRState         string
	ReleasedVersion git_tools.DeployedVersion
}

func (w ListMyPRsWorkflow) GetName() string {
	return w.Name
}

func (w ListMyPRsWorkflow) GetOrgFilename() string {
	return w.OrgFileName
}

func (w ListMyPRsWorkflow) GetOrgSectionName() string {
	return w.SectionTitle
}

func (w ListMyPRsWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) (RunResult, error) {
	client := git_tools.GetGithubClient()
	prs := git_tools.GetManyRepoPRs(client, w.PRState, w.Owner, w.Repos)

	prs = git_tools.ApplyPRFilters(prs, w.Filters)
	doc := org.GetOrgDocument(w.OrgFileName, org.MergeInfoOrgSerializer{})
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return RunResult{}, errors.New("Section Not Found")
	}
	prs = git_tools.ApplyPRFilters(prs, []git_tools.PRFilter{git_tools.MyPRs})
	result := ProcessPRs(prs, c, &doc, &section, file_change_wg, false)
	// TODO This is moving to the serializer
	// if pr.MergedAt != nil && output.ChangeType != "No Change" {
	//	repo_name := *pr.Base.Repo.Name
	//	if repo_name == "repo" {
	//		released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
	//		if released {
	//			fmt.Printf("Released PR: %s %t\n", *pr.Title, released)
	//		}
	//		// output.Item.Details() = append(output.Lines, "Released: "+strconv.FormatBool(released))
	//		//output.Lines[0] = strings.Replace(output.Lines[0], "merged", "released", 1)
	//	}
	// }

	return result, nil
}

type ProjectListWorkflow struct {
	Name            string
	Owner           string
	Repo            string
	OrgFileName     string
	Filters         []git_tools.PRFilter
	SectionTitle    string
	JiraDomain      string
	JiraEpic        string
	ReleasedVersion git_tools.DeployedVersion
}

func (w ProjectListWorkflow) GetName() string {
	return w.Name
}

func (w ProjectListWorkflow) GetOrgFilename() string {
	return w.OrgFileName
}

func (w ProjectListWorkflow) GetOrgSectionName() string {
	return w.SectionTitle
}

func (w ProjectListWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) (RunResult, error) {
	client := git_tools.GetGithubClient()
	doc := org.GetOrgDocument(w.OrgFileName, org.MergeInfoOrgSerializer{})
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		return RunResult{}, errors.New("Section Not Found")
	}
	if w.JiraEpic == "" {
		// I used to let just define []int for PR #s in config, could easily bring that back
		return RunResult{}, errors.New("ProjectList requires Jira Epic")
	}
	projectPRs := jira.GetProjectPRKeys(w.JiraDomain, w.JiraEpic, w.Repo)

	prs := git_tools.GetSpecificPRs(client, w.Owner, w.Repo, projectPRs)
	result := ProcessPRs(prs, c, &doc, &section, file_change_wg, false)
	return result, nil
}
