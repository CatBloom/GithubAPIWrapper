package models

import (
	"encoding/json"
	"io"
	"log"
	"main/configs"
	"net/http"
	"strings"
)

type IssueModel interface {
	GetIssues(token string, owner string, repo string, first int, order string, states string, after string) (IssuesResp, error)
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
		Nodes    []Issue       `json:"nodes"`
		PageInfo IssuePageInfo `json:"pageInfo"`
	}

	IssuePageInfo struct {
		EndCutsor   string `json:"endCursor"`
		HasNextPage bool   `json:"hasNextPage"`
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

func (im *issueModel) GetIssues(token string, owner string, repo string, first int, order string, states string, after string) (IssuesResp, error) {
	issueResp := IssuesResp{}

	query := `
		query($name: String! ,$owner: String!, $first: Int!, $orderBy: IssueOrder!, $states: [IssueState!]!, $after: String) {
			repository(name: $name, owner: $owner) {
				createdAt
				updatedAt
				name
				nameWithOwner
				issues(first: $first, states: $states, orderBy: $orderBy, after: $after) {
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
					pageInfo {
						endCursor
						hasNextPage
					}
				}
			}
		}
	`

	val := map[string]string{"query": query, "variables": im.makeVariables(owner, repo, first, order, states, after)}

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

func (rm *issueModel) makeVariables(owner string, repo string, first int, order string, states string, after string) string {
	variables := make(map[string]interface{})

	variables["owner"] = owner
	variables["name"] = repo
	variables["first"] = first

	if order != "" {
		variables["orderBy"] = map[string]string{
			"field":     "CREATED_AT",
			"direction": order,
		}
	} else {
		variables["orderBy"] = map[string]string{
			"field":     "CREATED_AT",
			"direction": "DESC",
		}
	}

	if states != "" {
		variables["states"] = states
	} else {
		variables["states"] = "OPEN"
	}

	if after != "" {
		variables["after"] = after
	}

	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		log.Println("Error make variables:", err)
	}

	return string(jsonVariables)
}
