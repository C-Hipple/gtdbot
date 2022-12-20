package main

import (
	"fmt"
	//"io"
	"strings"
	// "github.com/martinlindhe/notify"
	"github.com/google/go-github/v48/github"
)

type Workflow interface {
	Run(chan int, int)
}

type LanesGroup struct {
	Board_id string
	Lane_ids []string
}

type SyncLeankitLaneToOrg struct {
	Lanes      LanesGroup
	OrgSection Section
	Filters    []Filter
}

func NewSyncLeankitToOrg(board_id string, lanes []string, section_name string, filters []Filter) SyncLeankitLaneToOrg {
	return SyncLeankitLaneToOrg{Lanes: LanesGroup{board_id, lanes}, OrgSection: GetOrgSection(section_name), Filters: filters}
}

func (wf SyncLeankitLaneToOrg) Run(c chan int, idx int) {
	cards := getCards(wf.Lanes.Board_id, wf.Lanes.Lane_ids, wf.Filters) // TODO Multiple leankit lane support
	if len(cards) == 0 {
		fmt.Println("There are no cards currently in Needs Review lane on Leankit which need reviewed by me!")
	}
	for _, card := range cards {
		SyncCardToSection(card, wf.OrgSection)
	}
	c <- idx
}

func SyncCardToSection(card Card, section Section) bool {
	if CheckCardAlreadyInSection(card, section) {
		return false
	}
	AddTODO(GetOrgFile(), section, card)
	return true
}

func (wf SyncLeankitLaneToOrg) PostProcess(n int) {
}

func CheckCardAlreadyInSection(card Card, section Section) bool {
	for _, line_item := range section.Items {
		if strings.Contains(line_item.FullLine(0), card.Id) {
			return true
		}
	}
	return false
}

type UpdatePRLinkWorkflow struct{}

func (uprl UpdatePRLinkWorkflow) Run(c chan int, idx int) {
	client := GetGithubClient()
	prs := getPRs(client)
	my_eng_board_cards := getCards(BOARD_ENGINEERING, []string{LANE_ENG_CORE_NEEDS_REVIEW}, nil)

	// N x N matching
	for _, pr := range prs {
		for _, card := range my_eng_board_cards {
			fmt.Println(pr.Title)
			fmt.Println(card.Title)
		}
	}
}

func UppdatePRLink(pr *github.PullRequest, new_link string) {
	for _, line := range strings.Split(*pr.Body, "\n") {
		if strings.Contains(line, "link_to_card") {
		}
	}

}
