package main

import (
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/workflows"
	"os"
	"path/filepath"

	"github.com/google/go-github/v48/github"
	"github.com/pelletier/go-toml/v2"
)

// Define your classes
type Config struct {
	Repos     []string
	Workflows []workflows.Workflow
}

// This struct implements all possible values a workflow can define, then they're written as-needed.
type RawWorkflow struct {
	WorkflowType        string
	Name                string
	Owner               string
	Repo                string
	Repos               []string
	JiraEpic            string
	Filters             []string
	OrgFileName         string
	SectionTitle        string
	PRState             string
	ReleaseCheckCommand string
}

func LoadConfig() Config {
	// Load TOML config

	var intermediate_config struct {
		Repos      []string
		JiraDomain string
		Workflows  []RawWorkflow
	}
	home_dir, err := os.UserHomeDir()
	the_bytes, err := os.ReadFile(filepath.Join(home_dir, ".config/gtdbot.toml"))
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(the_bytes, &intermediate_config)
	if err != nil {
		panic(err)
	}

	return Config{
		Repos:     intermediate_config.Repos,
		Workflows: MatchWorkflows(intermediate_config.Workflows, &intermediate_config.Repos, intermediate_config.JiraDomain),
	}
}

func MatchWorkflows(workflow_maps []RawWorkflow, repos *[]string, jiraDomain string) []workflows.Workflow {
	workflows := []workflows.Workflow{}
	for _, raw_workflow := range workflow_maps {
		if raw_workflow.WorkflowType == "SyncReviewRequestsWorkflow" {
			workflows = append(workflows, BuildSyncReviewRequestWorkflow(&raw_workflow, repos))
		}
		if raw_workflow.WorkflowType == "SingleRepoSyncReviewRequestsWorkflow" {
			workflows = append(workflows, BuildSingleRepoReviewWorkflow(&raw_workflow, repos))
		}
		if raw_workflow.WorkflowType == "ListMyPRsWorkflow" {
			workflows = append(workflows, BuildListMyPRsWorkflow(&raw_workflow, repos))
		}
		if raw_workflow.WorkflowType == "ProjectListWorkflow" {
			workflows = append(workflows, BuildProjectListWorkflow(&raw_workflow, jiraDomain))
		}
	}
	return workflows
}

func BuildSingleRepoReviewWorkflow(raw *RawWorkflow, repos *[]string) workflows.Workflow {
	wf := workflows.SingleRepoSyncReviewRequestsWorkflow{
		Name:                raw.Name,
		Owner:               raw.Owner,
		Repo:                raw.Repo,
		Filters:             BuildFiltersList(raw.Filters),
		OrgFileName:         raw.OrgFileName,
		SectionTitle:        raw.SectionTitle,
		ReleaseCheckCommand: raw.ReleaseCheckCommand,
	}
	return wf
}

func BuildSyncReviewRequestWorkflow(raw *RawWorkflow, repos *[]string) workflows.Workflow {
	wf := workflows.SyncReviewRequestsWorkflow{
		Name:                raw.Name,
		Owner:               raw.Owner,
		Repos:               *repos,
		Filters:             BuildFiltersList(raw.Filters),
		OrgFileName:         raw.OrgFileName,
		SectionTitle:        raw.SectionTitle,
		ReleaseCheckCommand: raw.ReleaseCheckCommand,
	}
	return wf
}

func BuildListMyPRsWorkflow(raw *RawWorkflow, repos *[]string) workflows.Workflow {
	wf := workflows.ListMyPRsWorkflow{
		Name:                raw.Name,
		Owner:               raw.Owner,
		Repos:               *repos,
		Filters:             BuildFiltersList(raw.Filters),
		PRState:             raw.PRState,
		OrgFileName:         raw.OrgFileName,
		SectionTitle:        raw.SectionTitle,
		ReleaseCheckCommand: raw.ReleaseCheckCommand,
	}
	return wf
}

func BuildProjectListWorkflow(raw *RawWorkflow, jiraDomain string) workflows.Workflow {
	wf := workflows.ProjectListWorkflow{
		Name:                raw.Name,
		Owner:               raw.Owner,
		Repo:                raw.Repo,
		JiraDomain:          jiraDomain,
		JiraEpic:            raw.JiraEpic,
		Filters:             BuildFiltersList(raw.Filters),
		OrgFileName:         raw.OrgFileName,
		SectionTitle:        raw.SectionTitle,
		ReleaseCheckCommand: raw.ReleaseCheckCommand,
	}
	return wf
}

var filter_func_map = map[string]func(prs []*github.PullRequest) []*github.PullRequest{
	"FilterMyReviewRequested": git_tools.FilterMyReviewRequested,
	"FilterNotDraft":          git_tools.FilterNotDraft,
	"FilterIsDraft":           git_tools.FilterIsDraft,
	"FilterMyTeamRequested":   git_tools.FilterMyTeamRequested,
	"FilterNotMyPRs":          git_tools.FilterNotMyPRs,
}

func BuildFiltersList(names []string) []git_tools.PRFilter {
	filters := []git_tools.PRFilter{}
	for _, name := range names {
		filter_func := filter_func_map[name]
		if filter_func == nil {
			fmt.Println("Warning: Unmatched filter function ", name)
			continue
		}
		filters = append(filters, filter_func)
	}
	return filters
}
