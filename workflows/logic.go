package workflows

import (
	"context"
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/org"
	"gtdbot/utils"
	"regexp"
	"strings"
	"sync"

	"github.com/google/go-github/v48/github"
)

type Workflow interface {
	GetName() string
	Run(chan FileChanges, *sync.WaitGroup) (RunResult, error)
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

func (prb PRToOrgBridge) StartLine() int {
	// This implementation of the interface is for when we pull things from the API and want to compare
	// Therefore the StartLine should't be checked
	panic("Called StartLine for PRToOrgBridge which shouldn't be done.")
	// return -1
}

func (prb PRToOrgBridge) LinesCount() int {
	// This implementation of the interface is for when we pull things from the API and want to compare
	// Therefore the LinesCount should't be checked
	panic("Called LinesCount for PRToOrgBridge which shouldn't be done.")
	// return -1
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

	// TODO: Consider putting these in subsection?
	if prb.PR.MergedAt != nil {
		details = append(details, fmt.Sprintf("Merged at: %s\n", *prb.PR.MergedAt))
		if prb.PR.MergeCommitSHA != nil {
			details = append(details, fmt.Sprintf("Merge Commit SHA: %s\n", *prb.PR.MergeCommitSHA))
		} else {
			details = append(details, "Merged with Empty Merge Commit SHA?")
		}
	} else {
		// fmt.Println("PR is not merged.")
	}
	escaped_body := escapeBody(prb.PR.Body)
	details = append(details, fmt.Sprintf("*** BODY\n %s\n", cleanBody(&escaped_body))) // TODO: Do we need this end newline?
	// comments := getComments(*prb.PR.Head.Repo.Owner.Login, *prb.PR.Head.Repo.Name, *prb.PR.Number)
	// if len(comments) != 0 {
	//	details = append(details, fmt.Sprintf("*** Comments [%v]\n %s\n", len(comments), cleanLines(&comments)))
	// }
	// TODO review comments, see if they're included or not included when we do the above one.
	details = append(details, "END")
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
	return cleanLines(&lines)
}

func cleanEmptyEndingLines(lines *[]string) []string {
	// Removes the empty lines at the end of the details so org collapses prettier
	i := len(*lines) - 1
	for i >= 0 && ((*lines)[i] == "\n" || (*lines)[i] == "") {
		i--
	}
	return (*lines)[:i+1]
}

func cleanLines(lines *[]string) string {
	flat_lines := []string{}
	for _, line := range *lines {
		if strings.Contains(line, "\n") {
			split_lines := strings.Split(line, "\n")
			flat_lines = append(flat_lines, split_lines...)
		} else {
			flat_lines = append(flat_lines, line)
		}
	}

	shorted_lines := cleanEmptyEndingLines(&flat_lines)
	output_lines := make([]string, len(shorted_lines))
	for i, line := range shorted_lines {
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

func getComments(owner string, repo string, number int) []string {
	client := git_tools.GetGithubClient()
	opts := github.PullRequestListCommentsOptions{}
	comments, _, err := client.PullRequests.ListComments(context.Background(), owner, repo, number, &opts)
	if err != nil {
		fmt.Printf("Error getting Comments for PR %v in repo %s: %v", number, repo, err)
		return []string{}
	}
	str_comments := []string{}
	for _, comment := range comments {
		str_comments = append(str_comments, *comment.Body)
	}
	return str_comments
}

// func getReviewComments(owner string, repo string, number int) []string {
//	client := git_tools.GetGithubClient()

//	opts := github.ListOptions{}
//	reviews, _, err := client.PullRequests.ListReviews(context.Background(), owner , repo, number, &opts)
//	// review_comments := [];

//	for _, review := range reviews {
//		review_comments, _, err := client.PullRequests.ListReviewComments(context.Background(), owner, repo, int(*review.ID), &opts)
//	}
//	if err != nil {
//		fmt.Printf("Error getting Comments for PR %v in repo %s: %v", number, repo, err)
//		return []string{}
//	}
//	str_comments := []string{}
//	for _, comment := range comments {
//		str_comments = append(str_comments, *comment.Body)
//	}
//	return str_comments
// }

// func (prb PRToOrgBridge) GetReleased() string {
//	repo_name := *prb.PR.Base.Repo.Name
//	if repo_name == "repo" {
//		released := git_tools.CheckCommitReleased(client, w.ReleasedVersion.SHA, *pr.MergeCommitSHA)
//		fmt.Println("Released: ", released)

//	} else {
//		fmt.Println("Skipping check released due to repo.  PR is for repo: ", repo_name)
//	}
// }

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
