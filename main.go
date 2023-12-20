package main

import (
	"gtdbot/github"
	"gtdbot/workflows"
	"os"
)

func getManager(oneoff bool) workflows.ManagerService {
	return workflows.NewManagerService([]workflows.Workflow{
		workflows.SyncReviewRequestsWorkflow{
			Owner: "owner",
			Repo:  "repo",
			Filters: []github.PRFilter{
				github.FilterMyTeamRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Team Reviews",
		},
		workflows.SyncReviewRequestsWorkflow{
			Owner: "owner",
			Repo:  "repo",
			Filters: []github.PRFilter{
				github.FilterMyReviewRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "My Review Requests",
		},
		workflows.SyncReviewRequestsWorkflow{
			Owner: "owner",
			Repo:  "pytest-select-by-coverage",
			Filters: []github.PRFilter{
				github.FilterMyTeamRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
		workflows.SyncReviewRequestsWorkflow{
			Owner: "owner",
			Repo:  "pytest-select-by-coverage",
			Filters: []github.PRFilter{
				github.FilterMyReviewRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
		workflows.SyncReviewRequestsWorkflow{
			Owner: "owner",
			Repo:  "coreteam-devkit",
			Filters: []github.PRFilter{
				github.FilterMyReviewRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
	},
		oneoff,
	)
}

func main() {
	oneoff := false
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "--oneoff" {
			oneoff = true
		}
	}
	ms := getManager(oneoff)
	ms.Run()
}
