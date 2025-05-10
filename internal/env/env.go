package env

// EnvironmentManager interface for environment variables
type EnvironmentManager interface {
	Client() string
	FileExtensions() []string
	GithubToken() string
	AwsAccessKeyID() string
	AwsSecretAccessKey() string
	AwsRegion() string
	QReviewAPIEndpoint() string
	ShouldProcessFile(fileName string) bool
	ContextLines() int
}