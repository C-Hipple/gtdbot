package main

import (
	"flag"
	"gtdbot/config"
	"gtdbot/logger"
	"gtdbot/org"
	"gtdbot/workflows"
	"log/slog"
)

func main() {
	log := logger.New()
	slog.SetDefault(log)
	slog.Info("Starting!")
	one_off := flag.Bool("oneoff", false, "Pass oneoff to only run once")
	parse := flag.Bool("parse", false, "Pass parse to only parse the review file for testing/debugging.")
	initOnly := flag.Bool("init", false, "Pass init to only only setup the org file.")
	flag.Parse()
	if *parse {
		doc := org.GetOrgDocument("reviews.org", org.BaseOrgSerializer{}, config.C.OrgFileDir)
		doc.PrintAll()
		return
	}

	workflows_list := workflows.MatchWorkflows(config.C.RawWorkflows, &config.C.Repos, config.C.JiraDomain)
	ms := workflows.NewManagerService(
		workflows_list,
		*one_off,
		config.C.SleepDuration,
	)
	ms.Initialize()
	if *initOnly {
		slog.Info("Finished Initilization, Exiting.")

		return
	}

	ms.Run(log)
}
