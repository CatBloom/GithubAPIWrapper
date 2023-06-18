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

type RepoModels interface {
	GetReposByToken(token string, first int, order string) (RepoResp, error)
}

type repoModels struct{}

func NewRepoModels() RepoModels {
	return &repoModels{}
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
		Nodes []Repo `json:"nodes"`
	}

	Repo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}
)

func (rm *repoModels) GetReposByToken(token string, first int, order string) (RepoResp, error) {
	repoResp := RepoResp{}

	query := `
		query {
			viewer {
		  		repositories(first: %d, orderBy: { field: CREATED_AT, direction: %s }) {
					nodes {
			  			name
			  			description
			  			url
			  			createdAt
			  			updatedAt
					}
		  		}
			}
		}
	`

	query = fmt.Sprintf(query, first, order)
	val := map[string]string{"query": query}

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
