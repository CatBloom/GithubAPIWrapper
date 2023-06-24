# GithubAPIWrapper with golang

アプリで使用する最低限のGithubAPIのWrapper


## Technologies Used
* Golang v1.20
* Gin v1.9.0
* GitHubAPI
## API
* api
    * viewer(with token)
        * user
            * get
        * repos 
            * get params -> first: min 1 max 100, order: ASC or DESC
    * issue(with token)
        * list 
            * get params -> owner: owner name, repo: repository name, first: min 1 max 100, order: ASC or DESC, states: OPEN or CLOSE

## Local
### Run
1. .env.exsampleを.envに変更する
2. .envの項目を設定する
3. `docker compose up --build`
4. localhost:8080でアプリケーションが起動

### UnitTest
1. `docker exec -it github_api_container bash`
2. `gotest -v ./...`


