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

type RepoModel interface {
	GetReposByToken(string, types.ReposReq) (types.ReposRes, error)
}

type repoModel struct{}

func NewRepoModel() RepoModel {
	return &repoModel{}
}

func (rm *repoModel) GetReposByToken(token string, reposReq types.ReposReq) (types.ReposRes, error) {
	repoRes := types.ReposRes{}

	query := `
		query($first: Int!, $orderBy: RepositoryOrder!, $after: String) {
			viewer {
		  		repositories(first: $first, orderBy: $orderBy, after: $after) {
					nodes {
						id
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

	val := map[string]string{"query": query, "variables": rm.makeVariables(reposReq)}

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal GraphQL query json:", err)
		return repoRes, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return repoRes, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error do client:", err)
		return repoRes, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return repoRes, err
	}

	json.Unmarshal(body, &repoRes)

	return repoRes, nil
}

func (rm *repoModel) makeVariables(i interface{}) string {
	variables := make(map[string]interface{})

	// ReposReq型の場合
	if obj, ok := i.(types.ReposReq); ok {
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

		if obj.After != "" {
			variables["after"] = obj.After
		}
	}

	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		log.Println("Error make variables:", err)
	}

	return string(jsonVariables)
}
