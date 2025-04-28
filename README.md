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

the basic format is root level config for "repos"
and then a list of tables called [[Workflows]] configuring each workflow.


Each workflow entry can take the fields:
```
WorkflowType: str
Name: str
Owner: str
Filters: list[str]
OrgFileName: str
SectionTitle: str
```

The WorkflowType is one of the following strings:
SyncReviewRequestsWorkflow
SingleRepoSyncReviewRequestsWorkflow
ListMyPRsWorkflow
ProjectListWorkflow

Single Repo Sync workflow takes an additional paramter, Repo.
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

[[Workflows]]
WorkflowType = "SyncReviewRequestsWorkflow"
Name = "List Open PRs"
Owner = "C-Hipple"
Filters = ["FilterNotDraft"]
OrgFileName = "reviews.org"
SectionTitle = "Open PRs"

[[Workflows]]
WorkflowType = "ListMyPRsWorkflow"
Name = "List Closed PRs"
Owner = "C-Hipple"
OrgFileName = "reviews.org"
SectionTitle = "Closed PRs"
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
WorkflowType = "SyncReviewRequestsWorkflow"
Name = "List Open PRs"
JiraEpic = "BOARD-123" # the epic key
Owner = "C-Hipple"
Repo = "diff-lsp"
OrgFileName = "reviews.org"
SectionTitle = "Diff LSP Upgrade Project"
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


That string will then be put into the title line of the PR via the org-serializer
