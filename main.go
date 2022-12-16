package main

func main() {

	ms := ManagerService{Workflows: []Workflow{

		//NewSyncLeankitToOrg(BOARD_CORE, []string{LANE_CHRIS_DOING_NOW}, "Cards", []Filter{MyUserFilter}),
		//NewSyncLeankitToOrg(BOARD_CORE, []string{LANE_NEEDS_REVIEW}, "Code Review", []Filter{NotMeFilter}),
		TaskService{},
	}}
	ms.Start()
}
