package github

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v48/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"golang.org/x/oauth2"
)

type PullRequest interface {
}

func github_main() {
	// silly little tester function
	ctx := context.Background()
	token := os.Getenv("GTDBOT_GITHUB_TOKEN")
	if token == "" {
		panic("Unable to get Github Token!")
	}
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
	//	fmt.Println(*repo.Name)
	// }
	//pr, _, err2 := client.PullRequests.Get(ctx, "C-Hipple", "C-Hipple.github.io", 1)
	pr, _, err2 := client.PullRequests.Get(ctx, "multimediallc", "chaturbate", 9035)
	if err2 != nil {
		fmt.Println("Error on getting PR: ", err2)
	}
	fmt.Println(*pr.User.Login)
	//fmt.Println(*pr.Body)
	//*pr.Body = strings.Replace(*pr.Body, "desc", new string, n int)

	//client.PullRequests.Edit(ctx, "C-Hipple", "C-Hipple.github.io", 1, pr)
	// client.PullRequests.Edit(ctx, "multimediallc", "chaturbate", 9035, pr)
	prs := GetPRs(client, "open", "multimediallc", "chaturbate")
	for _, pr := range prs {
		fmt.Println(*pr.Title)
	}
}

type PRFilter func([]*github.PullRequest) []*github.PullRequest

func GetPRs(client *github.Client, state string, owner string, repo string) []*github.PullRequest {
	per_page := 100
	options := github.PullRequestListOptions{State: state, ListOptions: github.ListOptions{PerPage: per_page, Page: 1}}
	var prs []*github.PullRequest

	for {
		new_prs, _, err := client.PullRequests.List(context.Background(), owner, repo, &options)
		if err != nil {
			fmt.Println("Error!", err)
			//os.Exit(1)
			break
		}
		prs = append(prs, new_prs...)
		if len(new_prs) != per_page {
			break
		}
		options.Page += 1
	}
	return prs
}

func ApplyPRFilters(prs []*github.PullRequest, filters []PRFilter) []*github.PullRequest {
	for _, filter := range filters {
		prs = filter(prs)
	}
	return prs
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

func FilterPRsByState(prs []*github.PullRequest, state string) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		if *pr.State == state {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}

func FilterPRsByLabel(prs []*github.PullRequest, label string) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		for _, pr_label := range pr.Labels {
			if *pr_label.Name == label {
				filtered = append(filtered, pr)
				break
			}
		}
	}
	return filtered
}

func MyPRs(prs []*github.PullRequest) []*github.PullRequest {
	return FilterPRsByAuthor(prs, "C-Hipple")
}

func FilterNotDraft(prs []*github.PullRequest) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		if !*pr.Draft {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}

func FilterMyTeamRequested(prs []*github.PullRequest) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		for _, team := range pr.RequestedTeams {
			if *team.Name == "Core Pod Review Backend" {
				filtered = append(filtered, pr)
				break
			}
		}
	}
	return filtered
}

func FilterMyReviewRequested(prs []*github.PullRequest) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		for _, reviewer := range pr.RequestedReviewers {
			if *reviewer.Login == "C-Hipple" {
				filtered = append(filtered, pr)
				break
			}
		}
	}
	return filtered
}

func GetGithubClient() *github.Client {
	ctx := context.Background()
	token := os.Getenv("GTDBOT_GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Error! No Github Token!")
		os.Exit(1)
	}

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
