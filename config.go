package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/pelletier/go-toml/v2"
	"gtdbot/workflows"
)

// Define your classes
type ClassA struct{}

type ClassB struct{}

type Config struct{
	Repos []string
	Workflows  []workflows.Workflow
}

func load_config() Config {
	// Load TOML config

	var intermediate_config struct {
		Repos []string
		Workflows []map[string]string
	}
	home_dir, err := os.UserHomeDir()
	the_bytes, err := os.ReadFile(filepath.Join(home_dir, ".config/gtdbot.toml"))
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(the_bytes, &intermediate_config)
	if err != nil {
		panic(err)
	}
	fmt.Println("config: ")
	fmt.Println(intermediate_config)

	//MatchWorkflows(intermediate_config.ClassName)

	return Config{}
}


func MatchWorkflows(workflow_names []string) []workflows.Workflow {
	for _, name := range workflow_names {
		fmt.Println("workflow name", name)
	}
	return []workflows.Workflow{}
}

//func BuildSyncReviewRequests()
