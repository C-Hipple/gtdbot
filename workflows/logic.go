package workflows

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"gtdbot/org"
	"strings"
)

type Workflow interface {
	Run(chan FileChanges, int)
}

type FileChanges struct {
	change_type string
	start_line  int
	filename    string
	lines       []string
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
	line := fmt.Sprintf("%s TODO %s\t\t:%s:", strings.Repeat("*", indent_level), prb.Title(), *prb.PR.Head.Repo.Name)
	if *prb.PR.Draft {
		line = line + ":draft:"
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
	return details
}

func (prb PRToOrgBridge) String() string {
	return prb.Title()
}

// this is the official github package, not our lib, confusing!!
func SyncTODOToSection(doc org.OrgDocument, pr *github.PullRequest, section org.Section) FileChanges {
	pr_as_org := PRToOrgBridge{pr}
	if CheckTODOAlreadyInSection(pr_as_org, section) {
		return FileChanges{
			change_type: "No Change"}
	}
	return FileChanges{
		change_type: "Addition",
		start_line:  section.StartLine + 1,
		filename:    doc.Filename,
		lines:       doc.Serializer.Deserialize(pr_as_org, section.IndentLevel),
	}
}

func CheckTODOAlreadyInSection(todo org.OrgTODO, section org.Section) bool {
	for _, line_item := range section.Items {
		if strings.Contains(line_item.Summary(), todo.Summary()) {
			return true
		}
		if line_item.Summary() == todo.Summary() {
			return true
		}
		for _, detail := range line_item.Details() {
			if strings.Contains(detail, todo.ID()) {
				return true
			}
		}
	}
	return false
}
