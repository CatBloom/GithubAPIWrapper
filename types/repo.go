package types

type Repo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// request repos
type ReposReq struct {
	First int    `form:"first" binding:"required,max=100,min=1"`
	Order string `form:"order" binding:"omitempty,oneof=ASC DESC"`
	After string `form:"after" binding:"omitempty"`
}

// response repos
type (
	ReposRes struct {
		Data ReposViewer `json:"data"`
	}

	ReposViewer struct {
		Viewer ReposRepositories `json:"viewer"`
	}

	ReposRepositories struct {
		Repositories ReposNodes `json:"repositories"`
	}

	ReposNodes struct {
		Nodes    []Repo   `json:"nodes"`
		PageInfo PageInfo `json:"pageInfo"`
	}
)
