package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/configs"
	"net/http"
	"strings"
)

type IssueModel interface {
	GetIssues(token string, owner string, repo string, first int, order string, states string) (IssuesResp, error)
}

type issueModel struct{}

func NewIssueModel() IssueModel {
	return &issueModel{}
}

type (
	IssuesResp struct {
		Data IssueRepository `json:"data"`
	}

	IssueRepository struct {
		Repository Issues `json:"repository"`
	}

	Issues struct {
		CreatedAt     string     `json:"createdAt"`
		UpdatedAt     string     `json:"updatedAt"`
		Name          string     `json:"name"`
		NameWithOwner string     `json:"nameWithOwner"`
		Issues        IssueNodes `json:"issues"`
	}

	IssueNodes struct {
		Nodes []Issue `json:"nodes"`
	}

	Issue struct {
		ID        string       `json:"id"`
		CreatedAt string       `json:"createdAt"`
		UpdatedAt string       `json:"updatedAt"`
		State     string       `json:"state"`
		URL       string       `json:"url"`
		Title     string       `json:"title"`
		Number    int          `json:"number"`
		Body      string       `json:"body"`
		BodyHTML  string       `json:"bodyHTML"`
		Comments  CommentNodes `json:"comment"`
	}

	CommentNodes struct {
		Nodes []Comment `json:"nodes"`
	}

	Comment struct {
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
	}
)

func (im *issueModel) GetIssues(token string, owner string, repo string, first int, order string, states string) (IssuesResp, error) {
	issueResp := IssuesResp{}

	query := `
		query {
			repository(name: "%s", owner: "%s") {
				createdAt
				updatedAt
				name
				nameWithOwner
				issues(first: %d, states: %s, orderBy: { field: CREATED_AT, direction: %s }) {
					nodes {
						id
						createdAt
						updatedAt
						state
						url
						title
						number
						body
						bodyHTML
					}
				}
			}
		}
	`

	query = fmt.Sprintf(query, repo, owner, first, states, order)
	val := map[string]string{"query": query}

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal GraphQL query json:", err)
		return issueResp, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return issueResp, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Context-type", "application/json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error do client:", err)
		return issueResp, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return issueResp, err
	}

	json.Unmarshal(body, &issueResp)

	return issueResp, nil
}
