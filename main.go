package main

import (
	"flag"
	"fmt"
	"gtdbot/git_tools"
	"gtdbot/org"
	"gtdbot/workflows"
)

func get_manager(one_off bool, config *Config) workflows.ManagerService {
	release, err := git_tools.GetDeployedVersion()
	if err != nil {
		panic(err)
	}
	return workflows.NewManagerService(
		config.Workflows,
		release,
		one_off,
	)
}

func main() {
	fmt.Println("Starting!")
	one_off := flag.Bool("oneoff", false, "Pass oneoff to only run once")
	parse := flag.Bool("parse", false, "Pass parse to only parse the review file for testing/debugging.")
	initOnly := flag.Bool("init", false, "Pass init to only only setup the org file.")
	flag.Parse()
	if *parse {
		doc := org.GetBaseOrgDocument("reviews.org")
		doc.PrintAll()
		return
	}
	config := LoadConfig()
	ms := get_manager(*one_off, &config)
	ms.Initialize()
	if *initOnly {
		fmt.Println("Finished Initilization, Exiting.")

		return
	}

	ms.Run()
}
