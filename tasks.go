package main

import "github.com/google/go-github/v48/github" // with go modules enabled (GO111MODULE=on or outside GOPATH)

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

type TaskService struct{}

func (ts TaskService) Run(c chan int, idx int) {
	tasks := RunMatching()
	SyncMatchedTasks(tasks, GetOrgSection("Cards"))
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
	prs := getPRs(client)
	cards := getCards(BOARD_CORE, CORE_ACTIVE_LANES[:], []Filter{MyUserFilter})
	return Match(cards, prs)
}

func SyncMatchedTasks(tasks []*Task, to_section Section) bool {
	var res bool
	for _, task := range tasks {
		res = SyncCardToSection(task.card, to_section)
	}
	return res
}
