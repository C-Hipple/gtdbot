package main

import "os"

func getManager() ManagerService {
	return NewManagerService([]Workflow{
		SyncReviewRequestsWorkflow{
			owner: "owner",
			repo:  "repo",
			filters: []PRFilter{
				FilterMyTeamRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Team Reviews",
		},
		SyncReviewRequestsWorkflow{
			owner: "owner",
			repo:  "repo",
			filters: []PRFilter{
				FilterMyReviewRequested,
			},
			org_file_name: "reviews.org",
			section_title: "My Review Requests",
		},
		SyncReviewRequestsWorkflow{
			owner: "owner",
			repo:  "pytest-select-by-coverage",
			filters: []PRFilter{
				FilterMyTeamRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Other Repos",
		},
		SyncReviewRequestsWorkflow{
			owner: "owner",
			repo:  "pytest-select-by-coverage",
			filters: []PRFilter{
				FilterMyReviewRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Other Repos",
		},
	},
	)
}

func main() {
	ms := getManager()

	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "--oneoff" {
			ms.Run(true)
			return
		}
	}
	ms.Run(false)
}
