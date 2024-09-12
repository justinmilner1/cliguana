package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

// Load the configuration from the JSON file
func LoadConfig() (*Config, error) {
	config := DefaultConfig()
	filePath := expandPath(config.ConfigFile)

	// Check if the file exists, create if not
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		config = DefaultConfig() // Reinitialize to ensure default values are set
		if err := SaveConfig(config); err != nil {
			return config, fmt.Errorf("could not save default config: %v", err)
		}

	} else {
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return config, fmt.Errorf("could not read config: %v", err)
		}

		err = json.Unmarshal(file, &config)
		if err != nil {
			return config, fmt.Errorf("could not unmarshal config: %v", err)
		}
	}

	return config, nil
}

// Save the configuration back to the JSON file
func SaveConfig(config *Config) error {
	filePath := expandPath(config.ConfigFile)
	dir := filepath.Dir(filePath)

	// Create the directory if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("could not create config directory: %v", err)
		}
	}

	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %v", err)
	}

	err = ioutil.WriteFile(filePath, file, 0644)
	if err != nil {
		return fmt.Errorf("could not write config: %v", err)
	}

	return nil
}

// Helper to expand the `~` to the user's home directory
func expandPath(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	return path
}
