package main

import "os"

func getManager(oneoff bool) ManagerService {
	return NewManagerService([]Workflow{
		SyncReviewRequestsWorkflow{
			owner: "multimediallc",
			repo:  "chaturbate",
			filters: []PRFilter{
				FilterMyTeamRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Team Reviews",
		},
		SyncReviewRequestsWorkflow{
			owner: "multimediallc",
			repo:  "chaturbate",
			filters: []PRFilter{
				FilterMyReviewRequested,
			},
			org_file_name: "reviews.org",
			section_title: "My Review Requests",
		},
		SyncReviewRequestsWorkflow{
			owner: "multimediallc",
			repo:  "pytest-select-by-coverage",
			filters: []PRFilter{
				FilterMyTeamRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Other Repos",
		},
		SyncReviewRequestsWorkflow{
			owner: "multimediallc",
			repo:  "pytest-select-by-coverage",
			filters: []PRFilter{
				FilterMyReviewRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Other Repos",
		},
		SyncReviewRequestsWorkflow{
			owner: "multimediallc",
			repo:  "coreteam-devkit",
			filters: []PRFilter{
				FilterMyReviewRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Other Repos",
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
