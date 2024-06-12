package workflows

import (
	"fmt"
	"gtdbot/org"
	"strings"
	"sync"
	//"slices"

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
	line := fmt.Sprintf("%s %s %s\t\t%s", strings.Repeat("*", indent_level), prb.GetStatus(), prb.Title(), org.FormatTags(prb.GetTags()))
	//fmt.Println("Here: ", prb.Title(), prb.PR.Merged, prb.PR.MergedAt)
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

func (prb PRToOrgBridge) GetTags() []string {
	// consider setting this in a New method
	tags := []string{*prb.PR.Head.Repo.Name}
	if *prb.PR.Draft {
		tags = append(tags, "draft")
	} else if prb.PR.MergedAt != nil {
		tags = append(tags, "merged")
	}
	return tags
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
	item_in_section := CheckTODOAlreadyInSection(pr_as_org, section)
	if item_in_section != nil {
		return FileChanges{
			change_type: "No Change"}
	}
	return FileChanges{
		change_type: "Addition",
		start_line:  section.StartLine + 1,
		filename:    doc.Filename,
		Lines:       doc.Serializer.Deserialize(pr_as_org, section.IndentLevel),
	}
}

func CheckTODOAlreadyInSection(todo org.OrgTODO, section org.Section) org.OrgTODO {
	for _, line_item := range section.Items {
		if strings.Contains(line_item.Summary(), todo.Summary()) {
			return line_item
		}
		if line_item.Summary() == todo.Summary() {
			return line_item
		}
		for _, detail := range line_item.Details() {
			if strings.Contains(detail, todo.ID()) {
				return line_item
			}
		}
	}
	return nil
}

// func CheckForUpdates(item_in_section org.OrgTODO, pr_as_org PRToOrgBridge) FileChanges {
//	// For now we're only checking tags.
//	// Eventually will check more
//	existing_tags := org.FindOrgTags(item_in_section.Summary())
//	new_tags := pr_as_org.GetTags()
//	for _, tag := range new_tags {
//		if !slices.Contains(existing_tags, tag) {
//			return FileChanges{

//			}
//		}
//	}
// }
