package main

import (
	"flag"
	"gtdbot/github"
	"gtdbot/workflows"
)

func get_manager(one_off bool) workflows.ManagerService {
	return workflows.NewManagerService([]workflows.Workflow{
		workflows.SyncReviewRequestsWorkflow{
			Name:  "CB Team Reviews",
			Owner: "owner",
			Repo:  "repo",
			Filters: []github.PRFilter{
				github.FilterMyTeamRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Team Reviews",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:  "CB Personal Reviews",
			Owner: "owner",
			Repo:  "repo",
			Filters: []github.PRFilter{
				github.FilterMyReviewRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "My Review Requests",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:  "Core Reviews",
			Owner: "owner",
			Repo:  "repo",
			Filters: []github.PRFilter{
				github.FilterMyReviewRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "My Review Requests",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:         "Select by Coverage Team Reviews",
			Owner:        "owner",
			Repo:         "pytest-select-by-coverage",
			Filters:      []github.PRFilter{},
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
		// workflows.SyncReviewRequestsWorkflow{
		//	Name: "Coreteam Devkit Reviews"
		//	Owner: "owner",
		//	Repo:  "coreteam-devkit",
		//	Filters: []github.PRFilter{
		//		github.FilterMyReviewRequested,
		//	},
		//	OrgFileName:  "reviews.org",
		//	SectionTitle: "Other Repos",
		// },
	},
		one_off,
	)
}

func main() {
	one_off := flag.Bool("oneoff", false, "Pass oneoff to only run once")
	flag.Parse()
	ms := get_manager(*one_off)
	ms.Run()
}
