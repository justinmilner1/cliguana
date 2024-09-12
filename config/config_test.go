package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary config file
func createTempConfigFile(t *testing.T, content []byte) string {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

// Test loading a valid config file
func TestLoadConfig_ValidFile(t *testing.T) {
	// Mock environment variables
	os.Setenv("GREPTILE_AUTH_TOKEN", "Bearer valid_token")
	os.Setenv("GITHUB_TOKEN", "Bearer github_token")
	defer os.Unsetenv("GREPTILE_AUTH_TOKEN")
	defer os.Unsetenv("GITHUB_TOKEN")

	configContent := []byte(`{
		"AutouploadRepos": [{"repo_path": "/path/to/repo", "status": "enabled"}],
		"BaseURL": "https://api.greptile.com/v2/repositories",
		"AuthToken": "Bearer valid_token",
		"GithubToken": "Bearer github_token",
		"ConfigFile": "config.json"
	}`)

	configFilePath := createTempConfigFile(t, configContent)
	defer os.Remove(configFilePath)

	// Override the default config file path
	cfg := DefaultConfig()
	cfg.ConfigFile = configFilePath

	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedConfig.BaseURL != "https://api.greptile.com/v2/repositories" {
		t.Errorf("Expected BaseURL to be 'https://api.greptile.com/v2/repositories', got '%s'", loadedConfig.BaseURL)
	}
	if loadedConfig.AuthToken != "Bearer valid_token" {
		t.Errorf("Expected AuthToken to be 'Bearer valid_token', got '%s'", loadedConfig.AuthToken)
	}
	if loadedConfig.GithubToken != "Bearer github_token" {
		t.Errorf("Expected GithubToken to be 'Bearer github_token', got '%s'", loadedConfig.GithubToken)
	}
	if len(loadedConfig.AutouploadRepos) != 1 || loadedConfig.AutouploadRepos[0].RepoPath != "/path/to/repo" {
		t.Errorf("Expected AutouploadRepos to contain '/path/to/repo', got '%v'", loadedConfig.AutouploadRepos)
	}
}

// Test loading a non-existent config file
func TestLoadConfig_NonExistentFile(t *testing.T) {
	// Mock environment variables
	os.Setenv("GREPTILE_AUTH_TOKEN", "Bearer valid_token")
	os.Setenv("GITHUB_TOKEN", "Bearer github_token")
	defer os.Unsetenv("GREPTILE_AUTH_TOKEN")
	defer os.Unsetenv("GITHUB_TOKEN")

	// Override the default config file path to a non-existent file
	cfg := DefaultConfig()
	cfg.ConfigFile = filepath.Join(os.TempDir(), "non_existent_config.json")

	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loadedConfig.AutouploadRepos) != 0 {
		t.Errorf("Expected AutouploadRepos to be empty, got '%v'", loadedConfig.AutouploadRepos)
	}
}

// Test loading an invalid config file
func TestLoadConfig_InvalidFile(t *testing.T) {
	// Mock environment variables
	os.Setenv("GREPTILE_AUTH_TOKEN", "Bearer valid_token")
	os.Setenv("GITHUB_TOKEN", "Bearer github_token")
	defer os.Unsetenv("GREPTILE_AUTH_TOKEN")
	defer os.Unsetenv("GITHUB_TOKEN")

	configContent := []byte(`invalid json content`)

	configFilePath := createTempConfigFile(t, configContent)
	defer os.Remove(configFilePath)

	// Override the default config file path
	cfg := DefaultConfig()
	cfg.ConfigFile = configFilePath

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected error when loading invalid config file, got nil")
	}
}
