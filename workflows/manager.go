package workflows

import (
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/org"
	"strings"
	"sync"
	"time"
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
		if file_change.ChangeType != "No Change" {
			doc := org.GetBaseOrgDocument(file_change.Filename)
			change_lines := doc.Serializer.Deserialize(file_change.Item, file_change.Section.IndentLevel)
			if file_change.ChangeType == "Addition" {
				if strings.Contains(change_lines[0], "draft") {
					fmt.Print("Adding Draft PR: ", change_lines[3])
				} else {
					fmt.Print("Adding PR: ", change_lines[3])
				}
				fmt.Print(change_lines[2])
				doc.AddItemInSection(file_change.Section.Name, &file_change.Item)
			} else if file_change.ChangeType == "Replace" {
				doc.UpdateItemInSection(file_change.Section.Name, &file_change.Item)
			}
		}
		wg.Done()
	}
}

func NewManagerService(workflows []Workflow, release git_tools.DeployedVersion, oneoff bool) ManagerService {
	used_workflows := []Workflow{}
	for _, wf := range workflows {
		if strings.Contains(fmt.Sprintf("%T", wf), "ListMyPRsWorkflow") {
			// TODO: match the release getter with the repo
			fixed := wf.(ListMyPRsWorkflow)
			fixed.ReleasedVersion = release
			used_workflows = append(used_workflows, fixed)
		} else {
			used_workflows = append(used_workflows, wf)
		}
	}

	return ManagerService{
		Workflows:     used_workflows,
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
