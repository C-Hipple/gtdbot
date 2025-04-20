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

Depending on whether it's in the current build or not, this workflow has additional handling on checking if a PR is released or in a release-candidate tag.

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
