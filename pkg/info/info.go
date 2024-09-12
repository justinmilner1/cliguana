package info

import (
	"fmt"
	"time"

	"cliguana/config"
	"cliguana/pkg/http/greptile"
	"cliguana/pkg/util"
)

// Display a simple progress bar in the terminal
func displayProgressBar(percentage float64) {
	barWidth := 50                   // Width of the progress bar
	completed := int(percentage / 2) // Calculate the filled portion of the bar

	fmt.Print("\r[") // Carriage return to update the same line
	for i := 0; i < completed; i++ {
		fmt.Print("=")
	}
	for i := completed; i < barWidth; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("] %.2f%%", percentage) // Display percentage
}

// Function to fetch repository progress and calculate the percentage
func CheckProgress(cfg *config.Config, repoPath string) (int, int, error) {
	// Get remote URL
	remote := util.GetRemoteUrl(repoPath)
	if remote == "" {
		return 0, 0, fmt.Errorf("failed to get remote URL")
	}

	// Get current branch
	branch := util.GetCurrentBranch(repoPath)
	if branch == "" {
		return 0, 0, fmt.Errorf("failed to get current branch")
	}

	// Extract repository name from remote URL
	repository := util.ExtractRepoName(remote)
	if repository == "" {
		return 0, 0, fmt.Errorf("invalid remote URL: %s", remote)
	}

	repoInfo, err := greptile.SendGetInfoRequest(cfg, repository, remote, branch)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get repository info: %v", err)
	}

	// Return the number of files processed and the total number of files
	return repoInfo.FilesProcessed, repoInfo.NumFiles, nil
}

// Function to repeatedly check progress until completion
func MonitorProgress(cfg *config.Config, repoPath string) error {
	for {
		filesProcessed, numFiles, err := CheckProgress(cfg, repoPath)
		if err != nil {
			return err
		}

		// Calculate the progress percentage
		if numFiles == 0 {
			return fmt.Errorf("number of files is zero, cannot calculate progress")
		}
		progress := (float64(filesProcessed) / float64(numFiles)) * 100

		// Display the progress bar
		displayProgressBar(progress)

		// Check if the upload is complete
		if filesProcessed >= numFiles {
			fmt.Println("\nUpload complete!")
			break
		}

		// Wait for a while before checking again
		time.Sleep(4 * time.Second)
	}

	return nil
}
