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
	GetIssues(string, types.IssuesReq) (types.IssuesResp, error)
	GetIssue(string, types.IssueReq) (types.IssueResp, error)
}

type issueModel struct{}

func NewIssueModel() IssueModel {
	return &issueModel{}
}

func (im *issueModel) GetIssues(token string, issuesReq types.IssuesReq) (types.IssuesResp, error) {
	issuesResp := types.IssuesResp{}

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
		return issuesResp, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return issuesResp, err
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
		return issuesResp, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return issuesResp, err
	}

	json.Unmarshal(body, &issuesResp)

	return issuesResp, nil
}

func (im *issueModel) GetIssue(token string, issueReq types.IssueReq) (types.IssueResp, error) {
	issueResp := types.IssueResp{}

	query := `
		query($name: String!, $owner: String!, $number: Int!) {
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
					comments(first: 100) {
			  			nodes {
							id
							createdAt
							updatedAt
							author {
								login
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
	}

	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		log.Println("Error make variables:", err)
	}

	return string(jsonVariables)
}
