package main

import (
	//"strings"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v48/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"golang.org/x/oauth2"
	//"strconv"
)

type PullRequest interface {
}

func github_main() {
	// silly little tester function
	ctx := context.Background()
	token := os.Getenv("GTDBOT_GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}
	fmt.Println(len(repos))
	// for _, repo := range repos {
	// 	fmt.Println(*repo.Name)
	// }
	//pr, _, err2 := client.PullRequests.Get(ctx, "C-Hipple", "C-Hipple.github.io", 1)
	pr, _, err2 := client.PullRequests.Get(ctx, "owner", "repo", 9035)
	if err2 != nil {
		fmt.Println("Error on getting PR: ", err2)
	}
	fmt.Println(*pr.User.Login)
	//fmt.Println(*pr.Body)
	//*pr.Body = strings.Replace(*pr.Body, "desc", new string, n int)

	//client.PullRequests.Edit(ctx, "C-Hipple", "C-Hipple.github.io", 1, pr)
	// client.PullRequests.Edit(ctx, "owner", "repo", 9035, pr)
	prs := getPRs(client)
	for _, pr := range prs {
		fmt.Println(*pr.Title)
	}
}

// TODO Implement this so I can do the requested by core pod filter
type PRFilter func([]*github.PullRequest, string) []*github.PullRequest

type PRFilterSet struct {
	filters []PRFilter
}

func getPRs(client *github.Client) []*github.PullRequest {
	options := github.PullRequestListOptions{State: "open", ListOptions: github.ListOptions{PerPage: 300}} // TODO: proper pagination
	prs, _, _ := client.PullRequests.List(context.Background(), "owner", "repo", &options)
	return prs
	//return FilterPRsByAuthor(prs, "C-Hipple") // definitely want to keep for safety.  Noone cares if I gg my own PRs desc
}

func FilterPRsByAuthor(prs []*github.PullRequest, author string) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		if *pr.User.Login == author {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}

func GetGithubClient() *github.Client {
	ctx := context.Background()
	token := os.Getenv("GTDBOT_GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetLeankitCardTitle(pr PullRequest) string {
	return "abc"
}

func CheckBodyURLNotYetSet(body string) bool {
	return strings.Contains(body, "[Card Title]")

}

func ReplaceURLInBody(body string, title string, url string) string {
	lines := strings.Split(body, "\n")
	output_lines := []string{}
	for _, line := range lines {
		if strings.Contains(line, "[Card Title]") {
			output_lines = append(output_lines, fmt.Sprintf("[%s](%s)", title, url))
		} else {
			output_lines = append(output_lines, line)
		}
	}
	return strings.Join(output_lines, "\n")
}

func UpdatePRBody(pr *github.PullRequest, new_body string) bool {
	client := GetGithubClient()
	ctx := context.Background()

	pr.Body = &new_body

	pr, _, err := client.PullRequests.Edit(ctx, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Number, pr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func FilterPRsByAssignedTeam(prs []*github.PullRequest, target_team string) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		for _, team := range pr.RequestedTeams {
			if *team.Name == target_team {
				filtered = append(filtered, pr)
				continue
			}
		}
	}
	return filtered
}

func GetTeamAssignedPrs(team_name string) []*github.PullRequest {
	client := GetGithubClient()
	prs := getPRs(client)
	return FilterPRsByAssignedTeam(prs, team_name)
}

type SyncTeamAssignedPRsWorkflow struct {
	TeamName   string
	OrgSection Section
}

func NewSyncTeamAssignedPRsWorkflow(team_name string, org_section_name string) SyncTeamAssignedPRsWorkflow {
	return SyncTeamAssignedPRsWorkflow{team_name, GetOrgSection(org_section_name)}
}

func (s SyncTeamAssignedPRsWorkflow) Run(c chan int, idx int) {
	prs := GetTeamAssignedPrs(s.TeamName)
	SyncPRsToSection(prs, s.OrgSection)
	CheckFinishedPRsStillInSection(prs, s.OrgSection)
	c <- idx
}

func SyncPRsToSection(prs []*github.PullRequest, section Section)               {}
func CheckFinishedPRsStillInSection(prs []*github.PullRequest, section Section) {}
