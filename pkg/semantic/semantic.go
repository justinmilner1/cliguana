package semantic

import (
	"fmt"

	"cliguana/config"
	"cliguana/pkg/http/greptile"
	"cliguana/pkg/util"
)

// handleQuery handles the query command by sending the query to the Greptile API and displaying the results
func HandleQuery(cfg *config.Config, semanticQuery string, repoPath string) error {
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

	// Send the query request to the Greptile API
	response, err := greptile.SendQueryRepoRequest(cfg, repository, remote, branch, semanticQuery)
	if err != nil {
		return fmt.Errorf("error querying repository: %v", err)
	}

	// Display the response
	fmt.Println("Query Response:", response)
	return nil
}

// handleSearch handles the search command by sending the search query to the Greptile API and displaying the results
func HandleSearch(cfg *config.Config, searchQuery string, repoPath string) error {
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

	// Send the search request to the Greptile API
	response, err := greptile.SendSearchRepoRequest(cfg, repository, remote, branch, searchQuery)
	if err != nil {
		return fmt.Errorf("error searching repository: %v", err)
	}

	// Display the response
	fmt.Println("Search Response:", response)
	return nil
}
