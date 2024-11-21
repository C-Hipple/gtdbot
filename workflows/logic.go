package workflows

import (
	"fmt"
	"github.com/google/go-github/v48/github"
	"gtdbot/org"
	"gtdbot/utils"
	"regexp"
	"strings"
)

type Workflow interface {
	GetName() string
	Run(chan FileChanges)
}

type FileChanges struct {
	ChangeType string
	Filename   string
	Item       org.OrgTODO
	Section    org.Section
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

	user := prb.PR.GetUser()
	author_string := fmt.Sprintf("Author: %s", user.GetLogin())
	if user.GetName() != "" {
		author_string = author_string + fmt.Sprintf(" (%s)", user.GetName())
	}

	details = append(details, author_string+"\n")
	details = append(details, fmt.Sprintf("Branch: %s\n", *prb.PR.Head.Label))
	details = append(details, fmt.Sprintf("Requested Reviewers: %s\n",
		strings.Join(utils.Map(prb.PR.RequestedReviewers, getReviewerName), ", ")))
	details = append(details, fmt.Sprintf("Requested Teams: %s\n",
		strings.Join(utils.Map(prb.PR.RequestedTeams, getTeamName), ", ")))
	escaped_body := escapeBody(prb.PR.Body)
	details = append(details, "*** BODY\n %s\n", cleanBody(&escaped_body))
	return details
}

func (prb PRToOrgBridge) String() string {
	return prb.Title()
}

func getReviewerName(reviewer *github.User) string {
	return *reviewer.Login
}

func getTeamName(reviewer *github.Team) string {
	return *reviewer.Name
}

func escapeBody(body *string) string {
	// Body comes in a single string with newlines and can have things that break orgmode like *
	if body == nil {
		// pretty sure the library uses json:omitempty?
		return ""
	}

	lines := strings.Split(*body, "\n")
	if len(lines) == 0 {
		return ""
	}
	output_lines := make([]string, len(lines))
	for i, line := range lines {
		if strings.HasPrefix(line, "*") {
			output_lines[i] = strings.Replace(line, "*", "-", 1)
		} else {
			output_lines[i] = line
		}
	}
	return strings.Join(output_lines, "\n")
}

func cleanBody(body *string) string {
	// Define the regular expression pattern to match everything between <!-- and -->
	//	re := regexp.MustCompile(`<!--.*?-->`)
	// TODO more empty line cleaning
	re := regexp.MustCompile(`(?s)<!--.*?-->`)

	// Replace all matches with an empty string
	cleaned := re.ReplaceAllString(*body, "")

	return cleaned
}

// func (prb PRToOrgBridge) GetReleased() string {
//	repo_name := *prb.PR.Base.Repo.Name
//	if repo_name == "repo" {
//		released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
//		fmt.Println("Released: ", released)

//	} else {
//		fmt.Println("Skipping check released due to repo.  PR is for repo: ", repo_name)
//	}
// }

// this is the official github package, not our lib, confusing!!
func SyncTODOToSection(doc org.OrgDocument, pr *github.PullRequest, section org.Section) FileChanges {
	pr_as_org := PRToOrgBridge{pr}
	at_line, _ := org.CheckTODOInSection(pr_as_org, section)
	if at_line != -1 {
		// TODO : Determine if actual changes?
		return FileChanges{
			ChangeType: "Replace",
			Filename:   doc.Filename,
			Item:       pr_as_org,
			Section:    section,
		}
	}
	return FileChanges{
		ChangeType: "Addition",
		Filename:   doc.Filename,
		Item:       pr_as_org,
		Section:    section,
	}
}
