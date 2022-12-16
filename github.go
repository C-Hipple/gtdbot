package main

import (
	//"strings"
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v48/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"golang.org/x/oauth2"
	//"strconv"
)

func github_main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
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
	pr, _, err2 := client.PullRequests.Get(ctx, "multimediallc", "chaturbate", 9035)
	if err2 != nil {
		fmt.Println("Error on getting PR: ", err2)
	}
	fmt.Println(*pr.User.Login)
	//fmt.Println(*pr.Body)
	//*pr.Body = strings.Replace(*pr.Body, "desc", new string, n int)

	//client.PullRequests.Edit(ctx, "C-Hipple", "C-Hipple.github.io", 1, pr)
	// client.PullRequests.Edit(ctx, "multimediallc", "chaturbate", 9035, pr)
	prs := getPRs(client)
	prs = FilterPRsByAuthor(prs, "C-Hipple")
	for _, pr := range prs {
		fmt.Println(*pr.Title)
	}
}

func getPRs(client *github.Client) []*github.PullRequest {
	options := github.PullRequestListOptions{State: "open", ListOptions: github.ListOptions{PerPage: 300}} // TODO: proper pagination
	prs, _, _ := client.PullRequests.List(context.Background(), "multimediallc", "chaturbate", &options)
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

func GetGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func GetLeankitCardTitle(pr *github.PullRequest) string {
	return "abc"
}
