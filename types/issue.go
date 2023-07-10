package types

type Issue struct {
	ID        string       `json:"id"`
	CreatedAt string       `json:"createdAt"`
	UpdatedAt string       `json:"updatedAt"`
	State     string       `json:"state"`
	URL       string       `json:"url"`
	Title     string       `json:"title"`
	Number    int          `json:"number"`
	Body      string       `json:"body"`
	BodyHTML  string       `json:"bodyHTML"`
	Comments  CommentNodes `json:"comments"`
}

// request issues
type IssuesReq struct {
	Owner  string `form:"owner" binding:"required"`
	Repo   string `form:"repo" binding:"required"`
	First  int    `form:"first" binding:"required,max=100,min=1"`
	Order  string `form:"order" binding:"omitempty,oneof=ASC DESC"`
	States string `form:"states" binding:"omitempty,oneof=OPEN CLOSE"`
	After  string `form:"after" binding:"omitempty"`
}

// request issue
type IssueReq struct {
	Owner  string `form:"owner" binding:"required"`
	Repo   string `form:"repo" binding:"required"`
	Number int    `form:"number" binding:"required,min=1"`
}

type IssueCreateReq struct {
	RepoID   string   `form:"repoId" binding:"required"`
	Title    string   `form:"title" binding:"required"`
	Body     string   `form:"body" binding:"required"`
	LabelIds []string `form:"labelIds" binding:"omitempty"`
}

// response issues
type (
	IssuesRes struct {
		Data IssuesRepository `json:"data"`
	}

	IssuesRepository struct {
		Repository Issues `json:"repository"`
	}

	Issues struct {
		CreatedAt     string      `json:"createdAt"`
		UpdatedAt     string      `json:"updatedAt"`
		Name          string      `json:"name"`
		NameWithOwner string      `json:"nameWithOwner"`
		Issues        IssuesNodes `json:"issues"`
	}

	IssuesNodes struct {
		Nodes    []Issue  `json:"nodes"`
		PageInfo PageInfo `json:"pageInfo"`
	}
)

// response issue
type (
	IssueRes struct {
		Data IssueRepository `json:"data"`
	}

	IssueRepository struct {
		Repository IssueNode `json:"repository"`
	}

	IssueNode struct {
		Issue Issue `json:"issue"`
	}
)

// response create issue
type (
	IssueCreateRes struct {
		Data IssueCreate `json:"data"`
	}

	IssueCreate struct {
		CreateIssue IssueNode `json:"createIssue"`
	}
)
