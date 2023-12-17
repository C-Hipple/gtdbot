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


func NewManagerService(workflows []Workflow, oneoff bool) ManagerService {
	var cycle_count int
	if oneoff {
		cycle_count = 1
	} else {
		cycle_count = 9999
	}

	return ManagerService{
		Workflows:     workflows,
		workflow_chan: make(chan FileChanges),
		sleep_time:    1 * time.Minute,
		cycle_count: cycle_count,
	}
}

func (ms ManagerService) Run() {
	fmt.Println("Starting Service")
	for i, workflow := range ms.Workflows {
		go workflow.Run(ms.workflow_chan, i)
	}

	for i := 0; i <= ms.cycle_count; i++ {
		var changes []FileChanges
		for change := range ms.workflow_chan {
			changes = append(changes, change)
			if len(changes) == len(ms.Workflows) {
				break
			}
		}
		ms.ApplyFileChanges(changes)
	}
}

func (ms ManagerService) ApplyFileChanges(changes []FileChanges) {
	// naive solution which opens and closes the file for each one.
	fmt.Println("Applying File Changes")
	for _, file_changes := range changes {
		if file_changes.change_type == "Addition" {
			fmt.Println("Adding PR: ", file_changes.lines[3])
			InsertLinesToFile(GetOrgFile(file_changes.filename), file_changes.lines, file_changes.start_line)
		}
		// TODO: Handle other change types
	}
	fmt.Println("Done Applying File Changes")
}
