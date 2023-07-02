package types

type User struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	AvatarUrl string `json:"avatarUrl"`
}

// response user
type (
	UserRes struct {
		Data UserViewer `json:"data"`
	}

	UserViewer struct {
		Viewer User `json:"viewer"`
	}
)
