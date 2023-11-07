package main

import (
	"fmt"
	"time"
)

type ManagerService struct {
	Workflows     []Workflow
	workflow_chan chan int
	sleep_time    time.Duration
}

func NewManagerService(workflows []Workflow) ManagerService {
	return ManagerService{
		Workflows:     workflows,
		workflow_chan: make(chan int),
		sleep_time:    1 * time.Minute,
	}
}

func (ms ManagerService) Run(oneoff bool) {
	fmt.Println("Starting Service")
	for i, workflow := range ms.Workflows {
		go workflow.Run(ms.workflow_chan, i)
	}

	if oneoff {
		// just wait untill all workflows are done
		var count int
		for range ms.workflow_chan {
			count++
			if count == len(ms.Workflows) {
				return
			}
		}
	}

	for idx := range ms.workflow_chan {
		go func(idx int) {
			time.Sleep(ms.sleep_time)
			ms.Workflows[idx].Run(ms.workflow_chan, idx)
		}(idx)
	}
}
