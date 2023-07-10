# GithubAPIWrapper with golang

アプリで使用する最低限のGithubAPIのWrapper


## Technologies Used
* Golang v1.20
* Gin v1.9.0
* GitHubAPI
## API
Local End Point  
```
http://localhost:8080/api
```

header with access token
### /viewer  

* GET /user  

* GET /repos  

param|required|default|description
|--|--|--|--|
first|true|| min 1 max 100
order||DESC| ASC or DESC

### /issue
* GET /list  

param|required|default|description
|--|--|--|--|
first|true|| min 1 max 100
owner|true|| owner name
repo|true|| repository name
order||DESC| ASC or DESC
states||OPEN| OPEN or CLOSE

* GET /

param|required|default|description
|--|--|--|--|
owner|true|| owner name
repo|true|| repository name
number|true|| issue number

* POST /

body|required|default|description
|--|--|--|--|
repoID|true|| repository id
title|true|| issue title
body|true|| issue body
labelIds||| issue label id array

## Local
### Run
1. .env.exsampleを.envに変更する
2. .envの項目を設定する
3. `docker compose up --build`
4. localhost:8080でアプリケーションが起動

### UnitTest
1. `docker exec -it github_api_container bash`
2. `gotest -v ./...`


