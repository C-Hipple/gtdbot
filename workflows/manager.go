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
	var serialziedChannel = make(chan SerializedFileChange)
	go ApplyChanges(serialziedChannel, wg)
	for fileChange := range channel {
		fileChange.Log()

		if fileChange.ChangeType == "No Change" {
			wg.Done()
			continue
		}

		go func() {
			serialziedChannel <- fileChange.Deserialize()
		}()
	}
}

func ApplyChanges(channel chan SerializedFileChange, wg *sync.WaitGroup) {
	for deserializedChange := range channel {
		doc := org.GetOrgDocument(deserializedChange.FileChange.Filename, deserializedChange.FileChange.ItemSerializer)
		switch deserializedChange.FileChange.ChangeType {
		case "Addition":
			doc.AddDeserializedItemInSection(deserializedChange.FileChange.Section.Name, deserializedChange.Lines)
		case "Update", "Archive":
			doc.UpdateDeserializedItemInSection(deserializedChange.FileChange.Section.Name, &deserializedChange.FileChange.Item, deserializedChange.FileChange.ChangeType == "Archive", deserializedChange.Lines)
		case "Delete":
			doc.DeleteItemInSection(deserializedChange.FileChange.Section.Name, &deserializedChange.FileChange.Item)
		}
		wg.Done()

	}

}

func NewManagerService(workflows []Workflow, oneoff bool, sleep_time time.Duration) ManagerService {
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
		sleep_time:    sleep_time,
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
	fmt.Println("Completed RunOnce Waitgroup")
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
