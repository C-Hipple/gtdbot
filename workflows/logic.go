package workflows

import (
	"fmt"
	"gtdbot/org"
	"strings"
	"sync"

	"github.com/google/go-github/v48/github"
)

type Workflow interface {
	GetName() string
	Run(chan FileChanges, *sync.WaitGroup)
}

type FileChanges struct {
	change_type string
	start_line  int
	filename    string
	Lines       []string
}

type PRToOrgBridge struct {
	// Implement the OrgTODO Interface for PRs
	PR *github.PullRequest
}

func (prb PRToOrgBridge) ID() string {
	return fmt.Sprintf("%d", *prb.PR.Number)
}

func (prb PRToOrgBridge) Title() string {
	return *prb.PR.Title
}

func (prb PRToOrgBridge) FullLine(indent_level int) string {
	line := fmt.Sprintf("%s %s %s\t\t:%s:", strings.Repeat("*", indent_level), prb.GetStatus(), prb.Title(), *prb.PR.Head.Repo.Name)
	//fmt.Println("Here: ", prb.Title(), prb.PR.Merged, prb.PR.MergedAt)
	if *prb.PR.Draft {
		line = line + ":draft:"
	} else if prb.PR.MergedAt != nil {
		line = line + "merged:"
	}
	return line
}

func (prb PRToOrgBridge) Summary() string {
	return prb.Title()
}

func (prb PRToOrgBridge) CheckDone() bool {
	return *prb.PR.State == "closed"
}

func (prb PRToOrgBridge) GetStatus() string {
	if prb.CheckDone() {
		return "DONE"
	}
	return "TODO"
}

func (prb PRToOrgBridge) Details() []string {
	details := []string{}
	details = append(details, fmt.Sprintf("%d\n", prb.PR.GetNumber()))
	details = append(details, fmt.Sprintf("%s\n", prb.PR.GetHTMLURL()))
	details = append(details, fmt.Sprintf("Title: %s\n", prb.PR.GetTitle()))
	details = append(details, fmt.Sprintf("Author: %s\n", prb.PR.GetUser().GetLogin()))
	details = append(details, fmt.Sprintf("Branch: %s\n", *prb.PR.Head.Label))
	return details
}

func (prb PRToOrgBridge) String() string {
	return prb.Title()
}

// func (prb PRToOrgBridge) GetReleased() string {
//	repo_name := *prb.PR.Base.Repo.Name
//	if repo_name == "chaturbate" {
//		released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
//		fmt.Println("Released: ", released)

//	} else {
//		fmt.Println("Skipping check released due to repo.  PR is for repo: ", repo_name)
//	}
// }

// this is the official github package, not our lib, confusing!!
func SyncTODOToSection(doc org.OrgDocument, pr *github.PullRequest, section org.Section) FileChanges {
	pr_as_org := PRToOrgBridge{pr}
	at_line := CheckTODOInSection(pr_as_org, section)
	if at_line != -1 {
		if true {
			// Currently the Replace is broken
			return FileChanges{
				change_type: "No Change"}
		}
		fmt.Println("Replacing / Updating: ", pr_as_org.Title())

		// TODO: TO fix this we need to get teh start_line before we actually write, not when we determine
		// that it can be updated.  We have a race condition :/
		return FileChanges{
			change_type: "Replace",
			start_line: at_line,
			filename: doc.Filename,
			Lines: doc.Serializer.Deserialize(pr_as_org, section.IndentLevel),
		}
	}
	return FileChanges{
		change_type: "Addition",
		start_line:  section.StartLine + 1,
		filename:    doc.Filename,
		Lines:       doc.Serializer.Deserialize(pr_as_org, section.IndentLevel),
	}
}

func CheckTODOInSection(todo org.OrgTODO, section org.Section) int {
	// returns the line number if it's found, otherwise returns -1
	serializer := org.BaseOrgSerializer{}
	at_line := section.StartLine + 1 // account for the section title
	for _, line_item := range section.Items {
		if strings.Contains(line_item.Summary(), todo.Summary()) {
			return at_line
		}
		if line_item.Summary() == todo.Summary() {
			return at_line
		}
		for _, detail := range line_item.Details() {
			if strings.Contains(detail, todo.ID()) {
				return at_line
			}
		}
		at_line += len(serializer.Deserialize(line_item, section.IndentLevel))
	}
	return -1
}
