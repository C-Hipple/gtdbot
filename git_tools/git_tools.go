package git_tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/google/go-github/v48/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"golang.org/x/oauth2"
)

type PullRequest interface {
}

type PRFilter func([]*github.PullRequest) []*github.PullRequest

func GetPRs(client *github.Client, state string, owner string, repo string) []*github.PullRequest {
	per_page := 100
	options := github.PullRequestListOptions{State: state, ListOptions: github.ListOptions{PerPage: per_page, Page: 1}}
	var prs []*github.PullRequest

	// TODO: Consider if I really want deep lookups.
	// Setting to 0 limits to 1 API call.
	max_additional_calls := 4
	i := 0

	for {
		new_prs, _, err := client.PullRequests.List(context.Background(), owner, repo, &options)
		if err != nil {
			fmt.Println("Error!", err)
			//os.Exit(1)
			break
		}
		prs = append(prs, new_prs...)
		if len(new_prs) != per_page || i >= max_additional_calls {
			break
		}
		options.Page += 1
		i = i + 1
	}
	return prs
}

func GetManyRepoPRs(client *github.Client, state string, owner string, repos []string) []*github.PullRequest {
	var prs []*github.PullRequest
	for _, repo := range repos {
		repo_prs := GetPRs(
			client,
			state,
			owner,
			repo,
		)
		prs = append(prs, repo_prs...)
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

func FilterPRsExcludeAuthor(prs []*github.PullRequest, author string) []*github.PullRequest {
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		if *pr.User.Login != author {
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

func FilterNotMyPRs(prs []*github.PullRequest) []*github.PullRequest {
	return FilterPRsExcludeAuthor(prs, "C-Hipple")
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

func MakeTeamFilters(teams []string) func([]*github.PullRequest) []*github.PullRequest {
	return func(prs []*github.PullRequest) []*github.PullRequest {
		filtered := []*github.PullRequest{}
		for _, pr := range prs {
			for _, team := range pr.RequestedTeams {
				if slices.Contains(teams, *team.Name) {
					filtered = append(filtered, pr)
					break
				}
			}
		}
		return filtered
	}
}

func FilterMyTeamRequested(prs []*github.PullRequest) []*github.PullRequest {
	teams := []string{"growth-pod-review", "purchase-pod-review", "growth-and-purchase-pod", "coreteam-review", "chat-pod-review-backend", "creator-team", "affiliate-program-experts"}
	filtered := []*github.PullRequest{}
	for _, pr := range prs {
		for _, team := range pr.RequestedTeams {
			if slices.Contains(teams, *team.Slug) {
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
		//os.Exit(1)
		// dont' commit me
		token = "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"
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

func CheckReleased(commit_sha string) bool {
	return false
}

type DeployedVersion struct {
	Tag string
	SHA string
}

// TODO cache?
func GetDeployedVersion() (DeployedVersion, error) {
	res, err := http.Get("https://www.chaturbate.com/healthcheck/")
	if err != nil {
		return DeployedVersion{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return DeployedVersion{}, err
	}
	//"v24.5.10.0.post4 (sha: 80c7ed3aab56c4549acc895c292bb3b256b2789c)"

	// THis regex matches 3 things for some reason, so we just index properly
	extractor_regex, _ := regexp.Compile("(v24\\.[\\w+\\.]+)\\s\\(sha:\\s(\\w+)")
	match := extractor_regex.FindStringSubmatch(string(body))
	if len(match) < 3 {
		return DeployedVersion{}, errors.New("Invalid version match")
	}
	return DeployedVersion{
		Tag: match[1],
		SHA: match[2],
	}, nil
}

func CheckCommitReleased(client *github.Client, release_sha string, commit_sha string) bool {
	res, _, err := client.Repositories.CompareCommits(context.Background(), "multimediallc", "chaturbate", release_sha, commit_sha, &github.ListOptions{Page: 1, PerPage: 100})
	if err != nil {
		fmt.Println(err)
	}

	return slices.Contains([]string{"identical", "behind"}, *res.Status)

}

//func GetCommitsForSha(sha string)
