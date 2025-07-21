package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// This struct implements all possible values a workflow can define, then they're written as-needed.
type RawWorkflow struct {
	WorkflowType        string
	Name                string
	Owner               string
	Repo                string
	Repos               []string
	JiraEpic            string
	Filters             []string
	OrgFileName         string
	SectionTitle        string
	PRState             string
	ReleaseCheckCommand string
	Prune               string
}

// Define your classes
type Config struct {
	Repos         []string
	RawWorkflows  []RawWorkflow
	SleepDuration time.Duration
	OrgFileDir    string
	JiraDomain    string
}

var C Config

func init() {

	var intermediate_config struct {
		Repos         []string
		JiraDomain    string
		SleepDuration int64
		Workflows     []RawWorkflow
		OrgFileDir    string
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
	parsed_sleep_duration := time.Duration(1) * time.Minute
	if intermediate_config.SleepDuration == 0 {
		parsed_sleep_duration = time.Duration(intermediate_config.SleepDuration) * time.Minute
	}

	if intermediate_config.OrgFileDir == "" {
		intermediate_config.OrgFileDir = "~/"
	}

	C = Config{
		Repos:         intermediate_config.Repos,
		RawWorkflows:  intermediate_config.Workflows,
		SleepDuration: parsed_sleep_duration,
		OrgFileDir:    intermediate_config.OrgFileDir,
		JiraDomain:    intermediate_config.JiraDomain,
	}
}
