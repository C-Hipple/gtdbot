package main

import (
	"flag"
	"gtdbot/git_tools"
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
	one_off := flag.Bool("oneoff", false, "Pass oneoff to only run once")
	flag.Parse()
	config := LoadConfig()
	ms := get_manager(*one_off, &config)
	ms.Run()
}
