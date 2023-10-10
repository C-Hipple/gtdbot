package main

import (
	"fmt"
	//"io"
	"strings"
	"strconv"
	// "github.com/martinlindhe/notify"
	"github.com/google/go-github/v48/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
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

type PRToOrgBridge struct {
	// Implement the OrgTODO Interface for PRs
	PR *github.PullRequest
}

func (prb PRToOrgBridge) Id() string {
	return strconv.FormatInt(*prb.PR.ID, 10)
}

func (prb PRToOrgBridge) Title() string {
	return *prb.PR.Title
}

func (prb PRToOrgBridge) FullLine(indent_level int) string {
	return fmt.Sprintf("%s%s", strings.Repeat("*", indent_level), prb.Title())
}

func (prb PRToOrgBridge) CheckDone() bool {
	return *prb.PR.State == "closed"
}

func (prb PRToOrgBridge) GetStatus() string {
	if prb.CheckDone() {
		return "DONE"
	}
	return "TODO"
}

func (prb PRToOrgBridge) Details() []string {
	details := []string{}
	details = append(details, fmt.Sprintf("%s\n", prb.PR.GetHTMLURL()))
	details = append(details, fmt.Sprintf("Title: %s\n", prb.PR.GetTitle()))
	details = append(details, fmt.Sprintf("Author: %s\n", prb.PR.GetUser().GetLogin()))
	return details
}


func SyncPRToSection(pr *github.PullRequest, section Section) bool {
	if CheckPRAlreadyInSection(pr, section) {
		return false
	}
	AddTODO(GetOrgFile(), section, PRToOrgBridge{pr})
	return true
}


func CheckPRAlreadyInSection(pr *github.PullRequest, section Section) bool {
	for _, line_item := range section.Items {
		// print debugging
		fmt.Println(line_item.FullLine(0))

		if strings.Contains(line_item.FullLine(0), strconv.FormatInt(*pr.ID, 10)) {
			return true
		}
	}
	return false
}
