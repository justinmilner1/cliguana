package util

import (
	"fmt"
	"os/exec"
	"strings"
)

// Helper function to extract repository name from remote URL
func ExtractRepoName(remote string) string {
	// Example remote URLs:
	// HTTPS: https://github.com/owner/repo.git
	// SSH: git@github.com:owner/repo.git

	var repo string
	if strings.HasPrefix(remote, "http://") || strings.HasPrefix(remote, "https://") {
		// Handle HTTP(S) URLs
		parts := strings.Split(remote, "/")
		if len(parts) >= 2 {
			repo = parts[len(parts)-2] + "/" + strings.TrimSuffix(parts[len(parts)-1], ".git")
		}
	} else if strings.HasPrefix(remote, "git@") {
		// Handle SSH URLs
		parts := strings.Split(remote, ":")
		if len(parts) >= 2 {
			repo = parts[1]
			repo = strings.TrimSuffix(repo, ".git")
		}
	}

	return repo
}

// Helper function to map remote URL to remote type
func GetRemoteType(remote string) string {
	if strings.Contains(remote, "github.com") {
		return "github"
	} else if strings.Contains(remote, "gitlab.com") {
		return "gitlab"
	} else if strings.Contains(remote, "azure.com") {
		return "azure"
	}
	return ""
}

func GetRemoteUrl(absPath string) string {
	// Get remote URL
	remoteCmd := exec.Command("git", "-C", absPath, "remote", "get-url", "origin")
	remoteOutput, err := remoteCmd.Output()
	if err != nil {
		fmt.Println("Failed to get remote URL:", err)
		return ""
	}
	remote := strings.TrimSpace(string(remoteOutput))
	return remote
}

func GetCurrentBranch(absPath string) string {
	branchCmd := exec.Command("git", "-C", absPath, "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, err := branchCmd.Output()
	var branch = "master" // Default to 'master' if the current branch cannot be determined
	if err != nil {
		fmt.Printf("Error getting current branch: %v\n", err)
	} else {
		branch = strings.TrimSpace(string(branchOutput))
	}
	return branch
}
