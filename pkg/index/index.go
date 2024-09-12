package index

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cliguana/config"
	"cliguana/pkg/http/greptile"
	"cliguana/pkg/util"
)

// Check if a repository is already in the autoupload list
func isRepoInAutoupload(cfg *config.Config, repoPath string) bool {
	for _, repo := range cfg.AutouploadRepos {
		if repo.RepoPath == repoPath {
			return true
		}
	}
	return false
}

// Add a directory to the autoupload list
func AddRepoToAutoupload(cfg *config.Config, repoPath string) error {
	// Check if the directory is already in the autoupload list
	for _, dir := range cfg.AutouploadDirs {
		if dir == repoPath {
			fmt.Println("Directory is already flagged for autoupload.")
			return nil
		}
	}

	// Add new directory
	cfg.AutouploadDirs = append(cfg.AutouploadDirs, repoPath)

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("error saving configuration: %v", err)
	}

	fmt.Println("Directory added for autoupload:", repoPath)
	return nil
}

// Delete a directory from the autoupload list
func DeleteRepoFromAutoupload(cfg *config.Config, repoPath string) error {
	newDirs := []string{}
	found := false
	for _, dir := range cfg.AutouploadDirs {
		if dir != repoPath {
			newDirs = append(newDirs, dir)
		} else {
			found = true
		}
	}

	if !found {
		fmt.Println("Directory not found in autoupload list.")
		return nil
	}

	cfg.AutouploadDirs = newDirs
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("error saving configuration: %v", err)
	}

	fmt.Println("Directory removed from autoupload:", repoPath)
	return nil
}

// Trigger an API call to upload the repository
func TriggerUploadAPI(cfg *config.Config, repoPath string) error {
	// Get remote URL
	remote := util.GetRemoteUrl(repoPath)
	if remote == "" {
		return fmt.Errorf("failed to get remote URL")
	}

	// Get current branch
	branch := util.GetCurrentBranch(repoPath)
	if branch == "" {
		return fmt.Errorf("failed to get current branch")
	}

	// Extract repository name from remote URL
	repository := util.ExtractRepoName(remote)
	if repository == "" {
		return fmt.Errorf("invalid remote URL: %s", remote)
	}

	// Check if the directory is a valid Git repository
	if _, err := os.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
		return fmt.Errorf("the directory %s is not a valid Git repository", repoPath)
	}

	// Handle detached HEAD state
	if branch == "HEAD" {
		var branchCmd = exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
		branchOutput, err := branchCmd.Output()
		if err != nil {
			fmt.Printf("Error getting current commit hash: %v\n", err)
			return fmt.Errorf("failed to get current commit hash: %v", err)
		}
		branch = strings.TrimSpace(string(branchOutput))
		fmt.Printf("Current commit hash: %s\n", branch)
	}

	// Map the remote URL to the appropriate remote type
	remoteType := util.GetRemoteType(remote)
	if remoteType == "" {
		return fmt.Errorf("invalid remote URL: %s", remote)
	}

	return greptile.SendIndexRequest(cfg, repository, remoteType, branch)
}

// Simulate an API call for deletion
func TriggerDeleteAPI(cfg *config.Config, repoPath string) error {
	// Simulated API call for deleting the repository
	fmt.Printf("Triggering delete API for repository: %s\n", repoPath)
	fmt.Print("Delete not implemented yet\n")
	return nil
}

// Wrap the `git clone` command
func GitCloneAndUpload(cfg *config.Config, repoURL string, repoPath string) error {
	if repoPath == "." {
		repoPath = util.ExtractRepoName(repoURL)
	}

	cmd := exec.Command("git", "clone", repoURL, repoPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Cloning the repository...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone the repository: %v", err)
	}

	fmt.Println("Repository cloned successfully.")

	return TriggerUploadAPI(cfg, repoPath)
}
