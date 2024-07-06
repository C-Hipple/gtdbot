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

	return workflows.NewManagerService([]workflows.Workflow{
		workflows.SyncReviewRequestsWorkflow{
			Name:  "CB Team Reviews",
			Owner: "multimediallc",
			Repo:  "chaturbate",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyTeamRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Team Reviews",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:  "CB Personal Reviews",
			Owner: "multimediallc",
			Repo:  "chaturbate",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyReviewRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "My Review Requests",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:  "Core Reviews",
			Owner: "multimediallc",
			Repo:  "chaturbate",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyReviewRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "My Review Requests",
		},
		// workflows.SyncReviewRequestsWorkflow{
		//	Name:         "Select by Coverage Team Reviews",
		//	Owner:        "multimediallc",
		//	Repo:         "pytest-select-by-coverage",
		//	Filters:      []git_tools.PRFilter{},
		//	OrgFileName:  "reviews.org",
		//	SectionTitle: "Other Repos",
		// },
		workflows.SyncReviewRequestsWorkflow{
			Name:  "mm-actions Team Reviews",
			Owner: "multimediallc",
			Repo:  "mm-actions",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyTeamRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:  "mm-cli Team Reviews",
			Owner: "multimediallc",
			Repo:  "mm-cli",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyTeamRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
		workflows.SyncReviewRequestsWorkflow{
			Name:  "mm-cli Team Reviews",
			Owner: "multimediallc",
			Repo:  "mm-cli",
			Filters: []git_tools.PRFilter{
				git_tools.FilterMyTeamRequested,
			},
			OrgFileName:  "reviews.org",
			SectionTitle: "Other Repos",
		},
		workflows.ListMyPRsWorkflow{
			Repos:           []string{"chaturbate", "mm-cli", "cb_billing", "cbjpegstream"},
			Owner:           "multimediallc",
			PRState:         "open",
			OrgFileName:     "reviews.org",
			SectionTitle:    "My Pull Requests",
			ReleasedVersion: release,
		},

		workflows.ListMyPRsWorkflow{
			Repos:           []string{"chaturbate", "mm-cli", "cb_billing", "cbjpegstream"},
			Owner:           "multimediallc",
			PRState:         "closed",
			OrgFileName:     "reviews.org",
			SectionTitle:    "My Closed Pull Requests",
			ReleasedVersion: release,
		},

		// workflows.SyncReviewRequestsWorkflow{
		//	Name: "Coreteam Devkit Reviews"
		//	Owner: "multimediallc",
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
