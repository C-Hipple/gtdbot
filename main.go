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
			owner:   "owner",
			repo:    "repo",
			filters: []PRFilter{
				FilterNotDraft,
				FilterMyTeamRequested,
			},
			org_file: "reviews.org",
			section_title: "Team Reviews",
		},
		// SyncReviewRequestsWorkflow{

		//	owner:   "owner",
		//	repo:    "repo",
		//	filters: []PRFilter{
		//		FilterMyReviewRequested,
		//	},
		//	org_file: "reviews.org",
		//	section_title: "My Review Requests",
		// },
		// SyncReviewRequestsWorkflow{
		//	owner:   "owner",
		//	repo:    "pytest-select-by-coverage",
		//	filters: []PRFilter{
		//		FilterMyReviewRequested,
		//	},
		//	org_file: "reviews.org",
		//	section_title: "Other Repos",
		// },
	}}
	ms.Run()
}
