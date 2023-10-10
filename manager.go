package main

import (
	"fmt"
	"time"
)

type ManagerService struct {
	Workflows []Workflow
}

func (ms ManagerService) Run() {
	fmt.Println("Starting Service")
	ch := make(chan int)
	for i, workflow := range ms.Workflows {
		go workflow.Run(ch, i)
	}

	for idx := range ch {
		go func(idx int) {
			time.Sleep(60 * time.Second)
			ms.Workflows[idx].Run(ch, idx)
		}(idx)
	}
}
