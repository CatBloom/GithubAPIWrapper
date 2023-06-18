# GithubAPIWrapper with golang

アプリで使用する最低限のGithubAPIのWrapper


## Technologies Used
* Golang v1.20
* Gin v1.9.0
* GitHubAPI
## API
* api
    * user
        * get
    * repo
        * get -> first max 100, order ASC or DESC
## Local
### Run
1. .env.exsampleを.envに変更する
2. .envの項目を設定する
3. `docker compose up --build`
4. localhost:8080でアプリケーションが起動
`
### UnitTest
1. `docker exec -it github_api_container bash`
2. `gotest -v ./...`


