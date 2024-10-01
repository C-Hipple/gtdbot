package main

import (
	"flag"
	"gtdbot/git_tools"
	"gtdbot/workflows"
)

func get_manager(one_off bool) workflows.ManagerService {
	release, err := git_tools.GetDeployedVersion()
	if err != nil {
		panic(err)
	}
	repos := []string{"repo", "mm-cli", "mock-affiliate", "devsync", "mock-identity", "cb_billing", "billingv2"}

	return workflows.NewManagerService([]workflows.Workflow{
		workflows.SyncReviewRequestsWorkflow{
			Name:  "CB Team Reviews",
			Owner: "owner",
			Repo:  "repo",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyTeamRequested,
				git_tools.FilterNotDraft,
				git_tools.FilterNotMyPRs,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Team Reviews",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:  "CB Personal Reviews",
			Owner: "owner",
			Repo:  "repo",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyReviewRequested,
				git_tools.FilterNotDraft,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "My Review Requests",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:         "Falcon Nest Helper",
			Owner:        "owner",
			Repo:         "falcon-nest",
			Filters:      nil,
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
		// workflows.SyncReviewRequestsWorkflow{
		//	Name:  "Core Reviews",
		//	Owner: "owner",
		//	Repo:  "repo",
		//	Filters: []git_tools.PRFilter{
		//		git_tools.FilterMyReviewRequested,
		//	},
		//	OrgFileName:  "reviews.org",
		//	SectionTitle: "My Review Requests",
		// },
		// workflows.SyncReviewRequestsWorkflow{
		//	Name:         "Select by Coverage Team Reviews",
		//	Owner:        "owner",
		//	Repo:         "pytest-select-by-coverage",
		//	Filters:      []git_tools.PRFilter{},
		//	OrgFileName:  "reviews.org",
		//	SectionTitle: "Other Repos",
		// },
		// workflows.SyncReviewRequestsWorkflow{
		//	Name:  "mm-actions Team Reviews",
		//	Owner: "owner",
		//	Repo:  "mm-actions",
		//	Filters: []git_tools.PRFilter{
		//		git_tools.FilterMyTeamRequested,
		//	},
		//	OrgFileName:  "reviews.org",
		//	SectionTitle: "Other Repos",
		// },
		// workflows.SyncReviewRequestsWorkflow{
		//	Name:  "mm-cli Team Reviews",
		//	Owner: "owner",
		//	Repo:  "mm-cli",
		//	Filters: []git_tools.PRFilter{
		//		git_tools.FilterMyTeamRequested,
		//	},
		//	OrgFileName:  "reviews.org",
		//	SectionTitle: "Other Repos",
		// },
		// workflows.SyncReviewRequestsWorkflow{
		//	Name:  "mm-cli Team Reviews",
		//	Owner: "owner",
		//	Repo:  "mm-cli",
		//	Filters: []git_tools.PRFilter{
		//		git_tools.FilterMyTeamRequested,
		//	},
		//	OrgFileName:  "reviews.org",
		//	SectionTitle: "Other Repos",
		// },

		workflows.ListMyPRsWorkflow{
			Name:            "List my Open PRs",
			Repos:           repos,
			Owner:           "owner",
			PRState:         "open",
			OrgFileName:     "reviews.org",
			SectionTitle:    "My Pull Requests",
			ReleasedVersion: release,
		},

		workflows.ListMyPRsWorkflow{
			Name:            "List my Closed PRs",
			Repos:           repos,
			Owner:           "owner",
			PRState:         "closed",
			OrgFileName:     "reviews.org",
			SectionTitle:    "My Closed Pull Requests",
			ReleasedVersion: release,
		},

		// workflows.SyncReviewRequestsWorkflow{
		//	Name: "Coreteam Devkit Reviews"
		//	Owner: "owner",
		//	Repo:  "coreteam-devkit",
		//	Filters: []git_tools.PRFilter{
		//		git_tools.FilterMyReviewRequested,
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
