package configs

var githubGraphQLEndPoint = "https://api.github.com/graphql"
var githubAccessTokenEndPoint = "https://github.com/login/oauth/access_token"

var localHost = "http://localhost:5173"
var host = ""

func GetGithubGraphQLEndPoint() string {
	return githubGraphQLEndPoint
}

func GetGithubAccessTokenEndPoint() string {
	return githubAccessTokenEndPoint
}

func GetLocalHost() string {
	return localHost
}

func GetHost() string {
	return host
}
