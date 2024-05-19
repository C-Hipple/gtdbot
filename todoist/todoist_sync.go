package todoist


import (
	"net/http"
)

/*
Project structure

   {
    "id": "2203306141",
    "name": "Shopping List",
    "comment_count": 0,
    "color": "charcoal",
    "is_shared": false,
    "order": 1,
    "is_favorite": false,
    "is_inbox_project": false,
    "is_team_inbox": false,
    "view_style": "list",
    "url": "https://todoist.com/showProject?id=2203306141",
    "parent_id": null
}
*/

type Project struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CommentCount int `json:"comment_count"`
	Color string `json:"color"`
	IsShared bool `json:"is_shared"`
	Order int `json:"order"`
	IsFavorite bool `json:"is_favorite"`
	IsInboxProject bool `json:"is_inbox_project"`
	IsTeamInbox bool `json:"is_team_inbox"`
	ViewStyle string `json:"view_style"`
	Url string `json:"url"`
	ParentId string `json:"parent_id"`
}

/*
Task structure

   {
        "creator_id": "2671355",
        "created_at": "2019-12-11T22:36:50.000000Z",
        "assignee_id": "2671362",
        "assigner_id": "2671355",
        "comment_count": 10,
        "is_completed": false,
        "content": "Buy Milk",
        "description": "",
        "due": {
            "date": "2016-09-01",
            "is_recurring": false,
            "datetime": "2016-09-01T12:00:00.000000Z",
            "string": "tomorrow at 12",
            "timezone": "Europe/Moscow"
        },
        "duration": null,
        "id": "2995104339",
        "labels": ["Food", "Shopping"],
        "order": 1,
        "priority": 1,
        "project_id": "2203306141",
        "section_id": "7025",
        "parent_id": "2995104589",
        "url": "https://todoist.com/showTask?id=2995104339"
    },
 */
type Task struct {
	CreatorId string `json:"creator_id"`
	CreatedAt string `json:"created_at"`
	AssigneeId string `json:"assignee_id"`
	AssignerId string `json:"assigner_id"`
	CommentCount int `json:"comment_count"`
	IsCompleted bool `json:"is_completed"`
	Content string `json:"content"`
	Description string `json:"description"`
	Due struct {
		Date string `json:"date"`
		IsRecurring bool `json:"is_recurring"`
		Datetime string `json:"datetime"`
		String string `json:"string"`
		Timezone string `json:"timezone"`
	} `json:"due"`
	Duration string `json:"duration"`
	Id string `json:"id"`
	Labels []string `json:"labels"`
	Order int `json:"order"`
	Priority int `json:"priority"`
	ProjectId string `json:"project_id"`
	SectionId string `json:"section_id"`
	ParentId string `json:"parent_id"`
	Url string `json:"url"`
}

type TodoistService struct {
	client *http.Client
}

func New(client *http.Client) (*TodoistService, error) {
	return &TodoistService{client: client}, nil
}

func (s *TodoistService) GetProjects() ([]Project, error) {


