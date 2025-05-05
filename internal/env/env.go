// Package env manages environment
package env

const (
	clientEnvName  = "AI_CLIENT"
	githubToken    = "GITHUB_TOKEN"
	fileExtensions = "FILE_EXTENSIONS"
)

type EnvironmentManager interface {
	Client() string
	GithubToken() string
	ShouldProcessFile(fileName string) bool
}
