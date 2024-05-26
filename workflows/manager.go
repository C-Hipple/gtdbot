package workflows

import (
	"fmt"
	"gtdbot/org"
	"gtdbot/utils"
	"sync"
	"time"
	"strings"
)

type ManagerService struct {
	Workflows     []Workflow
	workflow_chan chan FileChanges
	sleep_time    time.Duration
	oneoff        bool
}

func ListenChanges(channel chan FileChanges, wg *sync.WaitGroup) {
	for file_change := range channel {
		wg.Add(1)
		if file_change.change_type == "Addition" {
			if strings.Contains(file_change.Lines[0], "draft") {
				fmt.Print("Adding Draft PR: ", file_change.Lines[3])
			} else {
				fmt.Print("Adding PR: ", file_change.Lines[3])
			}
			fmt.Print(file_change.Lines[2])
			utils.InsertLinesToFile(org.GetOrgFile(file_change.filename), file_change.Lines, file_change.start_line)
		}
		wg.Done()
	}
}

func NewManagerService(workflows []Workflow, oneoff bool) ManagerService {

	return ManagerService{
		Workflows:     workflows,
		workflow_chan: make(chan FileChanges),
		sleep_time:    1 * time.Minute,
		oneoff:        oneoff,
	}
}

func (ms ManagerService) RunOnce() {
	var wg sync.WaitGroup
	for _, workflow := range ms.Workflows {
		fmt.Println("Starting Workflow: ", workflow.GetName())
		wg.Add(1)
		go workflow.Run(ms.workflow_chan, &wg)
	}
	wg.Wait()
}

func (ms ManagerService) Run() {
	fmt.Println("Starting Service: ")
	var listener_wg sync.WaitGroup
	listener_wg.Add(1)
	go ListenChanges(ms.workflow_chan, &listener_wg)
	if ms.oneoff {
		fmt.Println("Running Once")
		ms.RunOnce()
	} else {
		for {
			fmt.Println("Cycle")
			ms.RunOnce()
			time.Sleep(ms.sleep_time)
		}
	}
	listener_wg.Done()
	listener_wg.Wait()
	fmt.Println("Exiting Service")
}
