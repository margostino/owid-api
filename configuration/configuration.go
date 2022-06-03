package configuration

import "github.com/margostino/owid-api/common"

type Configuration struct {
	MetadataPath      string `yaml:"metadata_path"`
	MetadataBaseUrl   string `yaml:"metadata_base_url"`
	SchemaFile        string `yaml:"schema_file"`
	GithubAccessToken string `yaml:"github_access_token"`
}

func GetConfig() *Configuration {
	var configuration Configuration
	common.UnmarshalYaml("./configuration.yml", &configuration)
	return &configuration
}
