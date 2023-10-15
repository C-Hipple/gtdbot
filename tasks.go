package main

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v48/github"
) // with go modules enabled (GO111MODULE=on or outside GOPATH)

type TaskStatus int

const (
	TODO TaskStatus = iota
	WAITING
	PROGRESS
	IN_REVIEW
	DONE
)

type Task struct {
	card   Card
	pr     *github.PullRequest
	status TaskStatus
}

func NewTask(card Card, pr *github.PullRequest) *Task {
	return &Task{card, pr, WAITING}
}

type PRLinkUpdateService struct{}

func (ts PRLinkUpdateService) Run(c chan int, idx int) {
	tasks := RunMatching()
	SyncMatchedTasks(tasks, GetOrgSection("gtd.org", "Cards"))
	UpdatePRLinksInBody(tasks)
	c <- idx
}

func Match(cards []Card, prs []*github.PullRequest) []*Task {
	tasks := []*Task{}
	for _, card := range cards {
		var matched_pr *github.PullRequest
		for _, pr := range prs {
			if CheckMatches(card, pr) {
				matched_pr = pr
				break
			}
		}
		tasks = append(tasks, NewTask(card, matched_pr))
	}
	return tasks
}

func CheckMatches(card Card, pr *github.PullRequest) bool {
	return card.Title == *pr.Title
}

func RunMatching() []*Task {
	client := GetGithubClient()
	prs := getPRs(client, "open", "multimediallc", "chaturbate")
	prs = FilterPRsByAuthor(prs, "C-Hipple")
	cards := getCards(BOARD_CORE, CORE_ACTIVE_LANES[:], []Filter{MyUserFilter})
	fmt.Printf("Gathered %d PRs and %d Cards", len(prs), len(cards))
	return Match(cards, prs)
}

func SyncMatchedTasks(tasks []*Task, to_section Section) {
	for _, task := range tasks {
		SyncCardToSection(task.card, to_section)
	}
}

func UpdatePRLinksInBody(tasks []*Task) {
	for _, task := range tasks {
		if task.pr == nil {
			continue
		}
		if CheckBodyURLNotYetSet(*task.pr.Body) || CheckBodyNeedsUpdatedToEngCardURL(task) {
			UpdatePRBody(task.pr, ReplaceURLInBody(*task.pr.Body, task.card.Title, task.card.URL()))
			fmt.Println("Updated Body of PR: ", task.pr.GetURL())
		}
	}
}

func CheckBodyNeedsUpdatedToEngCardURL(task *Task) bool {
	if strings.Contains(*task.pr.Body, task.card.Id) {
		children := task.card.GetCardChildren()
		for _, child := range children {
			if child.IsEngineeringCard() {
				return true
			}
		}
	}
	// if there are no engineering children cards or if it's already updated
	return false

}
