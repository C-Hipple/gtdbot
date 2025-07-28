# gtdbot
Once-Simple bot to help with my gtd workflows

# Installation


Build from source, it's go, just do it.
Binaries not provided.



```bash
export GTDBOT_GITHUB_TOKEN="Github Token"


go build
go run
```

# Configuration

gtdbot works from a toml config expected at the path `~/config/gtdbot.toml`.

the basic format is root level config for general fields

and then a list of tables called [[Workflows]] configuring each workflow.

The general fields are:
-
```
Repos: list[str]
SleepDuration: int (in minutes, optional, default=1 minute)
OrgFileDir: str
```

OrgFileDir will default to "~/" if it's not defined.


Each workflow entry can take the fields:
```
WorkflowType: str
Name: str
Owner: str
Filters: list[str]
OrgFileName: str
SectionTitle: str
ReleaseCommandCheck: str
Prune: string
IncludeDiff: bool
```

The WorkflowType is one of the following strings:
SyncReviewRequestsWorkflow
SingleRepoSyncReviewRequestsWorkflow
ListMyPRsWorkflow
ProjectListWorkflow

Prune tells the workflow runner whether or not to remove PRs from the section if they're no longer relevant.  The default behavior is to do nothing, and the options are:
Delete: Removes the item from the section.
Archive: Tags the items with :ARCHIVE: so that org functions can clean them up
Keep: Leave existing items in the section untouched.

IncludeDiff is a boolean that will include the diff of the PR in the org file entry.

### Workflow specific configurations
Single Repo Sync workflow takes an additional parameter, Repo.
```
Repo: str
```

ListMyPRsWorkflow takes the additional parameter PRState, which is passed through to the github API when filtering for PRs.
```
PRState: str [open/closed/nil]
```



An Example complete config file is below

```toml

Repos = [
    "C-Hipple/gtdbot",
    "C-Hipple/diff-lsp",
    "C-Hipple/diff-lsp.el",
]
SleepDuration = 5
OrgFileDir = "~/gtd/"

[[Workflows]]
WorkflowType = "SyncReviewRequestsWorkflow"
Name = "List Open PRs"
Owner = "C-Hipple"
Filters = ["FilterNotDraft"]
OrgFileName = "reviews.org"
SectionTitle = "Open PRs"
Prune = "Archive"

[[Workflows]]
WorkflowType = "ListMyPRsWorkflow"
Name = "List Closed PRs"
Owner = "C-Hipple"
OrgFileName = "reviews.org"
SectionTitle = "Closed PRs"
Prune = "Delete"
```

## Filters

Each workflow can the available filters:

*   `FilterMyReviewRequested`
*   `FilterNotDraft`
*   `FilterIsDraft`
*   `FilterMyTeamRequested`
*   `FilterNotMyPRs`


## JIRA Integration

The `ProjectListWorkflow` pulls information from Jira to build a realtime list of all PRs which are linked to children cards of the Jira epic given in the config.

Each workflow is tied to a single github repository, if you want multiple repos per project, create two workflows and have them use the same SectionTitle.

```bash
export JIRA_API_TOKEN="Jira API Token"
export JIRA_AIP_EMAIL="your email with your jira account"
```

```toml
JiraDomain="https://your-company.atlassain.net"

[[Workflows]]
WorkflowType = "ProjectListWorkflow"
Name = "Project - Example"
Owner = "C-Hipple"
Repo = "diff-lsp"
OrgFileName = "reviews.org"
SectionTitle = "Diff LSP Upgrade Project"
JiraEpic = "BOARD-123" # the epic key
```


## Release Checking

Often for work-workflows, it's very important to know when your particular PR is not just merged, but released to production, or in a release client.

You can configure a release check command which is run when PRs are added to the org file or updated.  GTDBOT will call-out to that program and expected a single string in response for

example. If we have a program on our PATH variable named release-check, you should call it like this:

```
$ release-check C-Hipple gtdbot abcdef
released

$ release-check C-Hipple gtdbot hijklm
release-client

$ release-check C-Hipple gtdbot nopqrs
merged
```


That string will then be put into the title line of the PR via the org-serializer.

## Emacs integration

This project ships with `gtdbot.el` for running and configuring this in emacs seamlessly.

### Installation

#### Spacemacs
```elisp
   ;; in dotspacemacs-additional-packages
   (gtdbot :location (recipe
                      :fetcher github
                      :repo "C-Hipple/gtdbot"
                      :files ("*.el")))
```

### Keybinds


You'll likely want to bind run-gtdbot-oneoff and/or run-gtdbot-service.

By default this package sets (if you use evil mode) `,r l` and `, r s` for those two commands.

If you don't use evil mode, you'll have to pick your own keybinds.

In org-agenda mode, this package adds a new command "R" which allows for a quick review (filtered by day/week/month/sprint) of completed items.

## Org-mode Review Notes

The default value for the files searched by the review functionality is:

```elisp
(setq gtdbot-org-agenda-files '("~/gtd/inbox.org" "~/gtd/gtd.org" "~/gtd/notes.org" "~/gtd/next_actions.org" "~/gtd/reviews.org"))
```

You can set this variable to wherever you keep your org files
