package workflows

import (
	"errors"
	"fmt"
	"gtdbot/git_tools"
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
		if output.ChangeType == "Replace" {
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
	result := RunResult{}
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		result.Process(&output, c, file_change_wg)
	}
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
		return RunResult{}, errors.New("Section Not Found")
	}
	result := RunResult{}
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		result.Process(&output, c, file_change_wg)
	}
	return result, nil
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

func (w ListMyPRsWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) (RunResult, error) {
	client := git_tools.GetGithubClient()
	prs := git_tools.GetManyRepoPRs(client, w.PRState, w.Owner, w.Repos)

	doc := org.GetOrgDocument(w.OrgFileName, org.MergeInfoOrgSerializer{})
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		fmt.Println("Error getting section: ", err, w.SectionTitle)
		return RunResult{}, errors.New("Section Not Found")
	}
	prs = git_tools.ApplyPRFilters(prs, []git_tools.PRFilter{git_tools.MyPRs})

	result := RunResult{}
	for _, pr := range prs {
		fmt.Printf("Checking My %s PR: %s\n", w.PRState, *pr.Title)
		output := SyncTODOToSection(doc, pr, section)
		result.Process(&output, c, file_change_wg)
		// TODO This is moving to the serializer
		// if pr.MergedAt != nil && output.ChangeType != "No Change" {
		//	repo_name := *pr.Base.Repo.Name
		//	if repo_name == "chaturbate" {
		//		released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
		//		if released {
		//			fmt.Printf("Released PR: %s %t\n", *pr.Title, released)
		//		}
		//		// output.Item.Details() = append(output.Lines, "Released: "+strconv.FormatBool(released))
		//		//output.Lines[0] = strings.Replace(output.Lines[0], "merged", "released", 1)
		//	}
		// }

	}
	return result, nil
}

type ProjectListWorkflow struct {
	Name            string
	Owner           string
	Repo            string
	OrgFileName     string
	SectionTitle    string
	ProjectPRs      []int
	ReleasedVersion git_tools.DeployedVersion
}

func (w ProjectListWorkflow) GetName() string {
	return w.Name
}

func (w ProjectListWorkflow) Run(c chan FileChanges, file_change_wg *sync.WaitGroup) (RunResult, error) {
	client := git_tools.GetGithubClient()
	doc := org.GetOrgDocument(w.OrgFileName, org.MergeInfoOrgSerializer{})
	section, err := doc.GetSection(w.SectionTitle)
	if err != nil {
		return RunResult{}, errors.New("Section Not Found")
	}
	prs := git_tools.GetSpecificPRs(client, w.Owner, w.Repo, w.ProjectPRs)
	result := RunResult{}
	for _, pr := range prs {
		output := SyncTODOToSection(doc, pr, section)
		// TODO This is moving to the serializer
		if pr.MergedAt != nil && output.ChangeType != "No Change" {
			repo_name := *pr.Base.Repo.Name
			if repo_name == "chaturbate" {
				released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
				fmt.Printf("Released PR: %s %t\n", *pr.Title, released)
			}
		}
		result.Process(&output, c, file_change_wg)
	}
	return result, nil
}
