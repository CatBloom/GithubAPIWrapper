package models

import (
	"encoding/json"
	"io"
	"log"
	"main/configs"
	"net/http"
	"strings"
)

type UserModels interface {
	GetUserByToken(token string) (UserResp, error)
}

type userModels struct{}

func NewUserModels() UserModels {
	return &userModels{}
}

type (
	UserResp struct {
		Data UserViewer `json:"data"`
	}

	UserViewer struct {
		Viewer User `json:"viewer"`
	}

	User struct {
		Login     string `json:"login"`
		Name      string `json:"name"`
		URL       string `json:"url"`
		AvatarUrl string `json:"avatarUrl"`
	}
)

func (um *userModels) GetUserByToken(token string) (UserResp, error) {
	userResp := UserResp{}

	query := `
		query {
			viewer{
		  		login
		  		name  
		  		url    
		  		avatarUrl  
			}
	  	}
	`

	val := map[string]string{"query": query}

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal GraphQL query json:", err)
		return userResp, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubGraphQLEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return userResp, err
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
		return userResp, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return userResp, err
	}

	json.Unmarshal(body, &userResp)

	return userResp, nil
}
