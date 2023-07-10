package models

import (
	"encoding/json"
	"io"
	"log"
	"main/configs"
	"main/types"
	"net/http"
	"strings"
)

type IssueModel interface {
	GetIssues(string, types.IssuesReq) (types.IssuesRes, error)
	GetIssue(string, types.IssueReq) (types.IssueRes, error)
	CreateIssue(string, types.IssueCreateReq) (types.IssueCreateRes, error)
}

type issueModel struct{}

func NewIssueModel() IssueModel {
	return &issueModel{}
}

func (im *issueModel) GetIssues(token string, issuesReq types.IssuesReq) (types.IssuesRes, error) {
	issuesRes := types.IssuesRes{}

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

	val := map[string]string{"query": query, "variables": im.makeVariables(issuesReq)}

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal GraphQL query json:", err)
		return issuesRes, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return issuesRes, err
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
		return issuesRes, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return issuesRes, err
	}

	json.Unmarshal(body, &issuesRes)

	return issuesRes, nil
}

func (im *issueModel) GetIssue(token string, issueReq types.IssueReq) (types.IssueRes, error) {
	issueRes := types.IssueRes{}

	query := `
		query($name: String!, $owner: String!, $number: Int!, $last: Int!) {
			repository(name: $name, owner: $owner) {
		  		issue(number: $number) {
					id
					createdAt
					updatedAt
					state
					url
					title
					number
					body
					bodyHTML
					comments(last: $last) {
			  			nodes {
							id
							createdAt
							updatedAt
							author {
								login
								url
								avatarUrl
							}
							authorAssociation
							body
							bodyHTML
			  			}
					}
		  		}
			}
	  	}
	`

	val := map[string]string{"query": query, "variables": im.makeVariables(issueReq)}

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal GraphQL query json:", err)
		return issueRes, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return issueRes, err
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
		return issueRes, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return issueRes, err
	}

	json.Unmarshal(body, &issueRes)

	return issueRes, nil
}

func (im *issueModel) CreateIssue(token string, issueCreateReq types.IssueCreateReq) (types.IssueCreateRes, error) {
	issueCreateRes := types.IssueCreateRes{}

	query := `
		mutation createIssue($input: CreateIssueInput!) {
			createIssue(input: $input) {
		  		issue {
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
	`

	val := map[string]string{"query": query, "variables": im.makeVariables(issueCreateReq)}

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal GraphQL query json:", err)
		return issueCreateRes, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return issueCreateRes, err
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
		return issueCreateRes, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return issueCreateRes, err
	}

	json.Unmarshal(body, &issueCreateRes)

	return issueCreateRes, nil
}

func (rm *issueModel) makeVariables(i interface{}) string {
	variables := make(map[string]interface{})

	// IssuesReq型の場合
	if obj, ok := i.(types.IssuesReq); ok {
		variables["owner"] = obj.Owner
		variables["name"] = obj.Repo
		variables["first"] = obj.First
		if obj.Order != "" {
			variables["orderBy"] = map[string]string{
				"field":     "CREATED_AT",
				"direction": obj.Order,
			}
		} else {
			variables["orderBy"] = map[string]string{
				"field":     "CREATED_AT",
				"direction": "DESC",
			}
		}
		if obj.States != "" {
			variables["states"] = obj.States
		} else {
			variables["states"] = "OPEN"
		}
		if obj.After != "" {
			variables["after"] = obj.After
		}
	}

	// IssueReq型の場合
	if obj, ok := i.(types.IssueReq); ok {
		variables["owner"] = obj.Owner
		variables["name"] = obj.Repo
		variables["number"] = obj.Number
		variables["last"] = 100 // コメントは最新の100件を取得
	}

	// IssueCreateReq型の場合
	if obj, ok := i.(types.IssueCreateReq); ok {
		input := make(map[string]interface{})
		input["repositoryId"] = obj.RepoID
		input["title"] = obj.Title
		input["body"] = obj.Body
		if len(obj.LabelIds) != 0 {
			input["labelIds"] = obj.LabelIds
		}
		variables["input"] = input
	}

	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		log.Println("Error make variables:", err)
	}

	return string(jsonVariables)
}
