package workflows

import (
	"fmt"
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
		if file_change.ChangeType != "No Change" {
			doc := org.GetOrgDocument(file_change.Filename, org.BaseOrgSerializer{})
			change_lines := doc.Serializer.Deserialize(file_change.Item, file_change.Section.IndentLevel)
			if file_change.ChangeType == "Addition" {
				if strings.Contains(change_lines[0], "draft") {
					fmt.Print("Adding Draft PR: ", change_lines[3])
				} else {
					fmt.Print("Adding PR: ", change_lines[3])
				}
				fmt.Print(change_lines[2])
				doc.AddItemInSection(file_change.Section.Name, &file_change.Item)
			} else if file_change.ChangeType == "Update" {
				fmt.Print("Updating item in section: ", change_lines[3])
				doc.UpdateItemInSection(file_change.Section.Name, &file_change.Item)
			} else if file_change.ChangeType == "Delete" {
				fmt.Print("Removing item from section: ", change_lines[3])
				doc.DeleteItemInSection(file_change.Section.Name, &file_change.Item)
			}
		}
		wg.Done() // The add is done when we enqueue the FileChange in the channel
	}
}

func NewManagerService(workflows []Workflow, oneoff bool) ManagerService {
	used_workflows := []Workflow{}
	for _, wf := range workflows {
		if strings.Contains(fmt.Sprintf("%T", wf), "ListMyPRsWorkflow") {
			// TODO: match the release getter with the repo
			fixed := wf.(ListMyPRsWorkflow)
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

func (ms ManagerService) runWorkflow(workflow Workflow, workflow_chan chan FileChanges, file_change_wg *sync.WaitGroup) {
	// Helper which times the workflow run command.
	fmt.Println("Starting Workflow: ", workflow.GetName())
	start := time.Now()
	result, err := workflow.Run(workflow_chan, file_change_wg)
	duration := time.Since(start)
	if err != nil {
		fmt.Println("Errored in Workflow: ", workflow.GetName(), " After: ", duration, ": ", err)
	}
	fmt.Println("Finishing Workflow: ", workflow.GetName(), " Took: ", duration, ":", result.Report())
}

func (ms ManagerService) RunOnce(file_change_wg *sync.WaitGroup) {
	var wg sync.WaitGroup
	for _, workflow := range ms.Workflows {
		wg.Add(1)
		go func(workflow Workflow) {
			defer wg.Done()
			ms.runWorkflow(workflow, ms.workflow_chan, file_change_wg)
		}(workflow)
	}
	wg.Wait()
	println("Completed RunOnce Waitgroup")
}

func (ms ManagerService) Run() {
	fmt.Println("Starting Service: ")
	var listener_wg sync.WaitGroup
	listener_wg.Add(1)
	go ListenChanges(ms.workflow_chan, &listener_wg)
	if ms.oneoff {
		fmt.Println("Running Once")
		ms.RunOnce(&listener_wg)
	} else {
		cycle_count := 0
		for {
			fmt.Println("Cycle: ", cycle_count)
			ms.RunOnce(&listener_wg)
			time.Sleep(ms.sleep_time)
			cycle_count++
		}
	}
	listener_wg.Done()
	listener_wg.Wait()
	fmt.Println("Exiting Service")
}

func (ms *ManagerService) Initialize() {
	// Ensure all required sections exist.
	// Does this sync since GetSection has creation side effect
	for _, wf := range ms.Workflows {
		// Don't need to check release command here
		doc := org.GetOrgDocument(wf.GetOrgFilename(), org.BaseOrgSerializer{ReleaseCheckCommand: ""})
		doc.GetSection(wf.GetOrgSectionName())
	}
}
