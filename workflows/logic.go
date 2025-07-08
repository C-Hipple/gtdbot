package workflows

import (
	"bytes"
	"context"
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/org"
	"gtdbot/utils"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v48/github"
)

type Workflow interface {
	GetName() string
	Run(chan FileChanges, *sync.WaitGroup) (RunResult, error)
	GetOrgSectionName() string
	GetOrgFilename() string
}

type FileChanges struct {
	ChangeType     string
	Filename       string
	Item           org.OrgTODO
	Section        org.Section
	ItemSerializer org.OrgSerializer
}

type SerializedFileChange struct {
	FileChange *FileChanges
	Lines      []string
}

func (fc FileChanges) Log() {
	fmt.Printf("[%s] %-20s - %s\n", fc.ChangeType[:2], fc.Section.Name, fc.Item.Summary())
}

func (fc *FileChanges) Deserialize() SerializedFileChange {
	var lines []string
	if fc.ChangeType != "Delete" {
		lines = fc.ItemSerializer.Deserialize(fc.Item, fc.Section.IndentLevel)
	}
	return SerializedFileChange{
		FileChange: fc,
		Lines:      lines,
	}
}

type PRToOrgBridge struct {
	// Implement the OrgTODO Interface for PRs
	PR *github.PullRequest
}

func (prb PRToOrgBridge) ID() string {
	return fmt.Sprintf("%d", *prb.PR.Number)
}

func (prb PRToOrgBridge) Repo() string {
	return prb.PR.Head.Repo.GetName()

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

func (prb PRToOrgBridge) ItemTitle(indent_level int, release_check_command string) string {

	line := fmt.Sprintf("%s %s %s\t\t:%s:", strings.Repeat("*", indent_level), prb.GetStatus(), prb.Title(), *prb.PR.Head.Repo.Name)
	if *prb.PR.Draft {
		line = line + ":draft:"
	} else if prb.PR.MergedAt != nil {
		if release_check_command != "" {
			status, err := GetReleaseStatus(&release_check_command, prb.PR.Head.Repo.Name, prb.PR.MergeCommitSHA)
			if err != nil {
				line = line + "merged:"
			} else {
				line = line + status + ":"
			}
		} else {
			line = line + "merged:"
		}
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
	if *prb.PR.Draft {
		return "WAITING"
	}
	return "TODO"
}

func (prb PRToOrgBridge) Details() []string {
	details := []string{}
	details = append(details, fmt.Sprintf("%d\n", prb.PR.GetNumber()))
	details = append(details, "Repo: "+*prb.PR.Head.Repo.Name)
	details = append(details, fmt.Sprintf("%s\n", prb.PR.GetHTMLURL()))
	details = append(details, fmt.Sprintf("Title: %s\n", prb.PR.GetTitle()))

	user := prb.PR.GetUser()
	author_string := fmt.Sprintf("Author: %s", user.GetLogin())
	if user.GetName() != "" {
		author_string = author_string + fmt.Sprintf(" (%s)", user.GetName())
	}

	details = append(details, author_string+"\n")

	details = append(details, fmt.Sprintf("Branch: %s\n", *prb.PR.Head.Label))

	reviewers := strings.Join(utils.Map(prb.PR.RequestedReviewers, getReviewerName), ", ")
	if reviewers != "" {
		details = append(details, fmt.Sprintf("Requested Reviewers: %s\n", reviewers))
	} else {
		details = append(details, "Requested Reviewers:\n")
	}
	teams := strings.Join(utils.Map(prb.PR.RequestedTeams, getTeamName), ", ")
	if teams != "" {
		details = append(details, fmt.Sprintf("Requested Teams: %s\n", teams))
	} else {
		details = append(details, "Requested Teams:\n")
	}

	// TODO: Consider putting these in subsection?
	if prb.PR.MergedAt != nil {
		details = append(details, fmt.Sprintf("Merged at: %s\n", *prb.PR.MergedAt))
		if prb.PR.MergeCommitSHA != nil {
			details = append(details, fmt.Sprintf("Merge Commit SHA: %s\n", *prb.PR.MergeCommitSHA))
		} else {
			details = append(details, "Merged with Empty Merge Commit SHA?")
		}
	} else {

		ciStatus, ciHeader := getCIStatus(*prb.PR.Base.Repo.Owner.Login, *prb.PR.Head.Repo.Name, *prb.PR.Head.Label)
		if len(ciStatus) > 0 {
			details = append(details, fmt.Sprintf("*** %s CI Status\n", ciHeader))
			details = append(details, ciStatus...)
		}
	}

	escaped_body := removePRBodySections(prb.PR.Body)
	escaped_body = escapeBody(&escaped_body)
	details = append(details, fmt.Sprintf("*** BODY\n %s\n", escaped_body)) // TODO: Do we need this end newline?
	comments_count, comments := getComments(*prb.PR.Base.Repo.Owner.Login, *prb.PR.Head.Repo.Name, *prb.PR.Number)
	if len(comments) != 0 {
		details = append(details, fmt.Sprintf("*** Comments [%v]\n", comments_count))
		details = append(details, comments...)
	}
	return details
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
	for i >= 0 && strings.TrimSpace((*lines)[i]) == "" {
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

func removePRBodySections(body *string) string {
	// Define the regular expression pattern to match everything between <!-- and -->
	//	re := regexp.MustCompile(`<!--.*?-->`)
	// TODO more empty line cleaning
	re := regexp.MustCompile(`(?s)<!--.*?-->`)

	if body == nil {
		return ""
	}

	// Replace all matches with an empty string
	cleaned := re.ReplaceAllString(*body, "")
	return cleaned
}

func getComments(owner string, repo string, number int) (int, []string) {
	client := git_tools.GetGithubClient()
	opts := github.PullRequestListCommentsOptions{}
	comments, _, err := client.PullRequests.ListComments(context.Background(), owner, repo, number, &opts)
	comments = filterComments(comments)
	trees := buildCommentTrees(comments)
	// debugPrintCommentTree(trees)

	if err != nil {
		fmt.Printf("Error getting Comments for PR %v in repo %s: %v", number, repo, err)
		return 0, []string{}
	}
	str_comments := []string{}
	for _, tree := range trees {
		for i, comment := range tree {
			if i == 0 {
				str_comments = append(str_comments, "**** "+comment.CreatedAt.Format(time.DateTime)+" "+treeAuthors(tree))
				str_comments = append(str_comments, *comment.DiffHunk)
			}
			clean_body := escapeBody(comment.Body)
			str_comments = append(str_comments, fmt.Sprintf("***** (%d) %s %s", i, comment.CreatedAt.Format(time.DateTime), *comment.User.Login))
			str_comments = append(str_comments, clean_body)
		}
	}

	return len(comments), str_comments
}

func ProcessPRs(prs []*github.PullRequest, changes_channel chan FileChanges, doc *org.OrgDocument, section *org.Section, change_wg *sync.WaitGroup, prune_command string) RunResult {
	result := RunResult{}

	// the index for both slices should match
	seen_prs := []*github.PullRequest{}
	pr_strings := []string{}
	changes := []FileChanges{}

	for _, pr := range prs {
		pr_strings = append(pr_strings, fmt.Sprintf("%s-%v", *pr.Head.Repo.Name, pr.GetNumber()))
		seen_prs = append(seen_prs, pr)
		// fmt.Printf("Checking My PR: %s\n", *pr.Title)
		changes = append(changes, SyncTODOToSection(*doc, pr, *section))
	}

	if prune_command == "Delete" || prune_command == "Archive" {
		// prune items that are not seen.  Use the PR string as the comparator
		for _, item := range section.Items {
			check_string := fmt.Sprintf("%s-%s", item.Repo(), item.ID())
			if slices.Contains(pr_strings, check_string) {
				continue
			} else {
				fileChange := FileChanges{
					ChangeType:     prune_command,
					Filename:       doc.Filename,
					Item:           item,
					Section:        *section,
					ItemSerializer: doc.Serializer,
				}
				changes = append(changes, fileChange)
			}
		}
	}

	for _, output := range changes {
		result.Process(&output, changes_channel, change_wg)
	}

	return result
}

func SyncTODOToSection(doc org.OrgDocument, pr *github.PullRequest, section org.Section) FileChanges {
	pr_as_org := PRToOrgBridge{pr}
	at_line, _ := org.CheckTODOInSection(pr_as_org, section)
	changeType := "Addition"
	if at_line != -1 {
		// After a week we stop updating old ones
		if pr_as_org.PR.GetMergedAt().After(time.Now().Add(-7 * 24 * time.Hour)) {
			changeType = "Update"
		} else {
			changeType = "No Change"
		}
	}
	return FileChanges{
		ChangeType:     changeType,
		Filename:       doc.Filename,
		Item:           pr_as_org,
		Section:        section,
		ItemSerializer: doc.Serializer,
	}
}


// Assume git_tools.GetGithubClient() and processWorkflowRuns are defined elsewhere
// and work as intended.

func getCIStatus(owner string, repo string, branch string) ([]string, string) {
	client := git_tools.GetGithubClient()
	branch = strings.Split(branch, ":")[1] // Comes as username:branch_name from github api.
	opts := github.ListWorkflowRunsOptions{Branch: branch}
	runs, _, err := client.Actions.ListRepositoryWorkflowRuns(context.Background(), owner, repo, &opts)

	if err != nil {
		fmt.Printf("Error getting combined status: %v\n", err)
		return []string{}, "TODO"
	}

	var statuses []string
	hasFailure := false
	hasInProgress := false
	allSuccess := true

	processedRuns := processWorkflowRuns(runs.WorkflowRuns)

	for _, run := range processedRuns {
		status := "<nil>" // completed, in_progress
		if run.Status != nil {
			status = "[" + *run.Status + "]"
			if *run.Status == "in_progress" {
				hasInProgress = true
			}
		}
		conclusion := " "
		if run.Conclusion != nil {
			if *run.Conclusion == "success" {
				conclusion = "✅"
				status = "" // We know the status if it was a success
			} else if *run.Conclusion == "failure" {
				conclusion = "❌"
				hasFailure = true
				allSuccess = false
			} else if *run.Conclusion != "success" {
				allSuccess = false // Any conclusion that is not success means not all are success
			}
		} else {
			// If conclusion is nil, it might still be in progress or pending
			allSuccess = false
			if run.Status != nil && *run.Status != "completed" {
				hasInProgress = true
			}
		}

		name := "Unknown Workflow Name"
		if run.Name != nil {
			name = *run.Name
		}

		item := fmt.Sprintf("[%s] %s %s", conclusion, status, name)
		statuses = append(statuses, item)
	}

	overallStatus := ""
	if hasFailure {
		overallStatus = "TODO"
	} else if hasInProgress {
		overallStatus = "WAITING"
	} else if allSuccess {
		overallStatus = "DONE"
	} else {
		// This case implies no failures, no in_progress, and not all successes.
		// This could happen if all runs are completed but some had non-success conclusions other than failure (e.g., cancelled, skipped).
		// Based on the prompt, if there are no failures and no in_progress, and it's not allSuccess, it should be DONE.
		// However, if we strictly interpret "DONE if all of the runs have a success conclusion",
		// then this state would be something else. Given the prompt, we'll default to DONE
		// if no failures and no in_progress, assuming non-success conclusions are also acceptable for DONE.
		// If a more specific state is needed for 'cancelled', 'skipped', etc., the logic here would need adjustment.
		overallStatus = "DONE"
	}

	return statuses, overallStatus
}

func processWorkflowRuns(runs []*github.WorkflowRun) []*github.WorkflowRun {
	latest_per_name := make(map[string]*github.WorkflowRun) // Initialize the map
	for _, run := range runs {
		if run == nil {
			continue
		}
		if run.Name == nil {
			continue
		}
		// fmt.Println("name: ", *run.Name, run)
		lastest_by_name := latest_per_name[*run.Name]
		if lastest_by_name == nil {
			latest_per_name[*run.Name] = run
			continue
		}
		if (*run.CreatedAt).After(lastest_by_name.CreatedAt.Time) {
			latest_per_name[*run.Name] = run
		}
	}

	var output []*github.WorkflowRun
	for _, run := range latest_per_name {
		output = append(output, run)
	}
	return output
}

// If a command was given by the workflow,
func GetReleaseStatus(command *string, repo *string, sha *string) (string, error) {
	cmd := exec.Command(*command, *repo, *sha)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	stdout := outb.String()
	return strings.Replace(stdout, "\n", "", -1), nil

}
func filterComments(comments []*github.PullRequestComment) []*github.PullRequestComment {
	output := []*github.PullRequestComment{}
	for _, comment := range comments {
		if strings.Contains(*comment.User.Login, "advanced") {
			// I don't care about the lint warning stuff
			continue
		}
		output = append(output, comment)
	}
	return output

}

func buildCommentTrees(comments []*github.PullRequestComment) [][]*github.PullRequestComment {
	output := [][]*github.PullRequestComment{}
	for _, comment := range comments {

		replyTo := int64(-1)
		if comment.InReplyTo != nil {
			replyTo = comment.GetInReplyTo()
		}

		if len(output) == 0 || replyTo == -1 {
			output = append(output, []*github.PullRequestComment{comment})
			continue
		}

		for j, commentTree := range output {
			if commentTree[len(commentTree)-1].GetID() == replyTo {
				output[j] = append(commentTree, comment)
				continue
			}
		}
	}
	return output
}

// List all of the authors in a tree for the tree title line.
func treeAuthors(tree []*github.PullRequestComment) string {
	authors := []string{}
	seen := make(map[string]bool)
	for _, comment := range tree {
		login := comment.User.GetLogin()
		if _, ok := seen[login]; !ok {
			authors = append(authors, login)
			seen[login] = true
		}
	}
	return strings.Join(authors, "|")
}

func debugPrintCommentTree(trees [][]*github.PullRequestComment) {
	for i, tree := range trees {
		fmt.Printf("Tree: %d\n", i)
		for j, comment := range tree {
			fmt.Printf("comment: %d - %d  (reply to: %d)\n", j, comment.GetID(), comment.GetInReplyTo())
		}
		fmt.Println("")
	}
}
