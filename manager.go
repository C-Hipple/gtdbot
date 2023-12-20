package main

import (
	"fmt"
	"time"
)

type ManagerService struct {
	Workflows     []Workflow
	workflow_chan chan FileChanges
	sleep_time    time.Duration
	cycle_count   int
}


func ListenChanges(channel chan FileChanges) {
	for file_change := range channel {
		if file_change.change_type == "Addition" {
			fmt.Println("Adding PR: ", file_change.lines[3])
			InsertLinesToFile(GetOrgFile(file_change.filename), file_change.lines, file_change.start_line)
		}
	}
}

func NewManagerService(workflows []Workflow, oneoff bool) ManagerService {
	var cycle_count int
	if oneoff {
		cycle_count = 0
	} else {
		cycle_count = 9999
	}

	return ManagerService{
		Workflows:     workflows,
		workflow_chan: make(chan FileChanges),
		sleep_time:    1 * time.Minute,
		cycle_count:   cycle_count,
	}
}

func (ms ManagerService) Run() {
	fmt.Println("Starting Service")
	go ListenChanges(ms.workflow_chan)

	for i := 0; i <= ms.cycle_count; i++ {

		fmt.Println("Cycle: ", i)

		for i, workflow := range ms.Workflows {
			go workflow.Run(ms.workflow_chan, i)
		}

		var changes []FileChanges
		for change := range ms.workflow_chan {
			changes = append(changes, change)
			if len(changes) == len(ms.Workflows) {
				break
			}
		}
	}
}
