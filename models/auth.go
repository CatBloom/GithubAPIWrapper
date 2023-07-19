package models

import (
	"encoding/json"
	"io"
	"log"
	"main/configs"
	"main/types"
	"net/http"
	"os"
	"strings"
)

type AuthModel interface {
	GetToken(types.AuthReq) (types.AuthRes, error)
}

type authModel struct{}

func NewAuthModel() AuthModel {
	return &authModel{}
}

func (am *authModel) GetToken(authReq types.AuthReq) (types.AuthRes, error) {
	authRes := types.AuthRes{}

	val := make(map[string]string)

	val["client_id"] = os.Getenv("CLIENT_ID")
	val["client_secret"] = os.Getenv("CLIENT_SECRET")
	val["code"] = authReq.Code

	data, err := json.Marshal(val)
	if err != nil {
		log.Println("Error marshal query json:", err)
		return authRes, err
	}

	req, err := http.NewRequest(http.MethodPost, configs.GetGithubAccessTokenEndPoint(), strings.NewReader(string(data)))
	if err != nil {
		log.Println("Error request github api:", err)
		return authRes, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error do client:", err)
		return authRes, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		return authRes, err
	}

	json.Unmarshal(body, &authRes)

	return authRes, nil
}
