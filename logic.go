package main

import (
	"fmt"
	//"io"
	"strings"
	// "github.com/martinlindhe/notify"
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
