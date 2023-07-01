package types

type CommentNodes struct {
	Nodes []Comment `json:"nodes"`
}

type Comment struct {
	ID                string `json:"id"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	AuthorAssociation string `json:"authorAssociation"`
	Author            Author `json:"author"`
	Body              string `json:"body"`
	BodyHTML          string `json:"bodyHTML"`
}

type Author struct {
	Login string `json:"login"`
}

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}
