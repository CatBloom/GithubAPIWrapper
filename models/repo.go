package models

import (
	"encoding/json"
	"io"
	"log"
	"main/configs"
	"net/http"
	"strings"
)

type RepoModel interface {
	GetReposByToken(token string, first int, order string, after string) (RepoResp, error)
}

type repoModel struct{}

func NewRepoModel() RepoModel {
	return &repoModel{}
}

type (
	RepoResp struct {
		Data RepoViewer `json:"data"`
	}

	RepoViewer struct {
		Viewer RepoRepositories `json:"viewer"`
	}

	RepoRepositories struct {
		Repositories RepoNodes `json:"repositories"`
	}

	RepoNodes struct {
		Nodes    []Repo       `json:"nodes"`
		PageInfo RepoPageInfo `json:"pageInfo"`
	}

	RepoPageInfo struct {
		EndCutsor   string `json:"endCursor"`
		HasNextPage bool   `json:"hasNextPage"`
	}

	Repo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}
)

func (rm *repoModel) GetReposByToken(token string, first int, order string, after string) (RepoResp, error) {
	repoResp := RepoResp{}

	query := `
		query($first: Int!, $orderBy: RepositoryOrder!, $after: String) {
			viewer {
		  		repositories(first: $first, orderBy: $orderBy, after: $after) {
					nodes {
			 			name
			  			description
			  			url
			  			createdAt
			  			updatedAt
					}
					pageInfo {
			  			endCursor
			  			hasNextPage
			  		}
		  		}
			}
	  	}
	`

	val := map[string]string{"query": query, "variables": rm.makeVariables(first, order, after)}

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal GraphQL query json:", err)
		return repoResp, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return repoResp, err
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
		return repoResp, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return repoResp, err
	}

	json.Unmarshal(body, &repoResp)

	return repoResp, nil
}

func (rm *repoModel) makeVariables(first int, order string, after string) string {
	variables := make(map[string]interface{})

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

	if after != "" {
		variables["after"] = after
	}

	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		log.Println("Error make variables:", err)
	}

	return string(jsonVariables)
}
