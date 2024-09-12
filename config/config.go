package config

import (
	"os"
)

type RepoConfig struct {
	RepoPath string `json:"repo_path"`
	Status   string `json:"status"`
}

type Config struct {
	AutouploadRepos []RepoConfig
	AutouploadDirs  []string
	BaseURL         string
	AuthToken       string
	GithubToken     string
	ConfigFile      string
}

func DefaultConfig() *Config {
	authToken := os.Getenv("GREPTILE_AUTH_TOKEN")
	if authToken == "" {
		print("Failed to get greptile token. See readme for setup")
		authToken = "Bearer <token>"
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		print("Failed to get github token. See readme for setup")
		githubToken = "Bearer <token>"
	}

	return &Config{
		AutouploadRepos: []RepoConfig{},
		BaseURL:         "https://api.greptile.com/v2/repositories",
		AuthToken:       authToken,
		GithubToken:     githubToken,
		ConfigFile:      "~/.cliguana/autoupload_repos.json",
	}
}

func LoadConfig() (*Config, error) {
	config := DefaultConfig()
	return config, nil
}

// Helper to expand the `~` to the user's home directory
func expandPath(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return path
}
