package main

import (
	"fmt"
	//"io"
	"strings"
	// "github.com/martinlindhe/notify"
	"github.com/google/go-github/v48/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)
)

type Workflow interface {
	Run(chan FileChanges, int)
}

type FileChanges struct {
	change_type string
	start_line  int
	filename    string
	lines       []string
}

// type LanesGroup struct {
//	Board_id string
//	Lane_ids []string
// }

// type SyncLeankitLaneToOrg struct {
//	Lanes      LanesGroup
//	OrgSection Section
//	Filters    []Filter
// }

// func NewSyncLeankitToOrg(board_id string, lanes []string, section_name string, filters []Filter) SyncLeankitLaneToOrg {
//	return SyncLeankitLaneToOrg{Lanes: LanesGroup{board_id, lanes}, OrgSection: GetOrgSection("gtd.org", section_name), Filters: filters}
// }

// func (wf SyncLeankitLaneToOrg) Run(c chan int, idx int) {
//	cards := getCards(wf.Lanes.Board_id, wf.Lanes.Lane_ids, wf.Filters) // TODO Multiple leankit lane support
//	if len(cards) == 0 {
//		fmt.Println("There are no cards currently in Needs Review lane on Leankit which need reviewed by me!")
//	}
//	for _, card := range cards {
//		SyncCardToSection(card, wf.OrgSection)
//	}
//	c <- idx
// }

// func SyncCardToSection(doc OrgDocument, card Card, section Section) bool {
//	if CheckCardAlreadyInSection(card, section) {
//		return false
//	}
//	AddTODO(doc, section, card)
//	return true
// }

// func (wf SyncLeankitLaneToOrg) PostProcess(n int) {
// }

// func CheckCardAlreadyInSection(card Card, section Section) bool {
//	for _, line_item := range section.Items {
//		if strings.Contains(line_item.FullLine(0), card.Id) {
//			return true
//		}
//	}
//	return false
// }

type PRToOrgBridge struct {
	// Implement the OrgTODO Interface for PRs
	PR *github.PullRequest
}

func (prb PRToOrgBridge) ID() string {
	return fmt.Sprintf("%d", *prb.PR.Number)
}

func (prb PRToOrgBridge) Title() string {
	return *prb.PR.Title
}

func (prb PRToOrgBridge) FullLine(indent_level int) string {
	line := fmt.Sprintf("%s TODO %s\t\t:%s:", strings.Repeat("*", indent_level), prb.Title(), *prb.PR.Head.Repo.Name)
	if *prb.PR.Draft {
		line = line + ":draft:"
	}
	return line
}

func (prb PRToOrgBridge) Summary() string {
	return prb.Title()
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
	details = append(details, fmt.Sprintf("%d\n", prb.PR.GetNumber()))
	details = append(details, fmt.Sprintf("%s\n", prb.PR.GetHTMLURL()))
	details = append(details, fmt.Sprintf("Title: %s\n", prb.PR.GetTitle()))
	details = append(details, fmt.Sprintf("Author: %s\n", prb.PR.GetUser().GetLogin()))
	return details
}

func (prb PRToOrgBridge) String() string {
	return prb.Title()
}

func SyncTODOToSection(doc OrgDocument, pr *github.PullRequest, section Section) FileChanges {
	pr_as_org := PRToOrgBridge{pr}
	if CheckTODOAlreadyInSection(pr_as_org, section) {
		return FileChanges{
			change_type: "No Change"}
	}
	return FileChanges{
		change_type: "Addition",
		start_line:  section.StartLine + 1,
		filename:    doc.Filename,
		lines:       doc.Serializer.Deserialize(pr_as_org, section.IndentLevel),
	}
}

func CheckTODOAlreadyInSection(todo OrgTODO, section Section) bool {
	for _, line_item := range section.Items {
		if strings.Contains(line_item.Summary(), todo.Summary()) {
			return true
		}
		if line_item.Summary() == todo.Summary() {
			return true
		}
		for _, detail := range line_item.Details() {
			if strings.Contains(detail, todo.ID()) {
				return true
			}
		}
	}
	return false
}
