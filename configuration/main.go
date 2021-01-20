package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type RepositoryConfig struct {
	Name  string `yaml:name`
	Owner string `yaml:owner`
}

type Config struct {
	GithubToken    string             `yaml:"github_token"`
	GithubUsername string             `yaml:"github_username"`
	Repositories   []RepositoryConfig `yaml:"repositories"`
}

func Load(path string) (config Config, err error) {
	var d []byte
	d, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(d, &config)
	if err != nil {
		return
	}
	return
}
