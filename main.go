package main

func main() {
	ms := ManagerService{Workflows: []Workflow{
		// Several deprecated workflows
		//NewSyncLeankitToOrg(BOARD_CORE, []string{LANE_CHRIS_DOING_NOW}, "Cards", []Filter{MyUserFilter}),
		//NewSyncLeankitToOrg(BOARD_CORE, []string{LANE_NEEDS_REVIEW}, "Code Review", []Filter{NotMeFilter}),
		//PRLinkUpdateService{},
		//SyncTeamAssignedPRsService{"Core Pod"},

		// New workflows
		SyncReviewRequestsWorkflow{
			owner: "multimediallc",
			repo:  "chaturbate",
			filters: []PRFilter{
				FilterNotDraft,
				FilterMyTeamRequested,
			},
			org_file_name:      "reviews.org",
			section_title: "Team Reviews",
		},
		SyncReviewRequestsWorkflow{

			owner:   "multimediallc",
			repo:    "chaturbate",
			filters: []PRFilter{
				FilterMyReviewRequested,
			},
			org_file_name: "reviews.org",
			section_title: "My Review Requests",
		},
		SyncReviewRequestsWorkflow{
			owner:   "multimediallc",
			repo:    "pytest-select-by-coverage",
			filters: []PRFilter{
				FilterMyTeamRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Other Repos",
		},
		SyncReviewRequestsWorkflow{
			owner:   "multimediallc",
			repo:    "pytest-select-by-coverage",
			filters: []PRFilter{
				FilterMyReviewRequested,
			},
			org_file_name: "reviews.org",
			section_title: "Other Repos",
		},
	}}
	ms.Run()
}
