package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"cliguana/config"
	"cliguana/pkg/index"
	"cliguana/pkg/info"
	"cliguana/pkg/semantic"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{Use: "cliguana"}

	// Helper function to get absolute path
	getAbsPath := func(repoPath string) (string, error) {
		absPath, err := filepath.Abs(repoPath)
		if err != nil {
			return "", fmt.Errorf("error getting absolute path: %v", err)
		}
		return absPath, nil
	}

	// `clone` command to wrap git clone and automatically upload after clone
	var cloneCmd = &cobra.Command{
		Use:   "clone [repo_url] [repo_path]",
		Short: "Clone a repository and automatically upload it for indexing",
		Long:  "Clone a repository from the given URL and automatically upload it to the Greptile API for indexing.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoURL := args[0]
			repoPath := "."
			if len(args) > 1 {
				repoPath = args[1]
			}

			if err := index.GitCloneAndUpload(cfg, repoURL, repoPath); err != nil {
				fmt.Println("Error during clone and upload:", err)
			}
		},
	}

	// `index` command to manually index a repository
	var monitorProgress bool
	var indexCmd = &cobra.Command{
		Use:   "index [repo_path]",
		Short: "Index a specific repository",
		Long:  "Send a repository to greptile for indexing",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoPath := "."
			if len(args) > 0 {
				repoPath = args[0]
			}
			absPath, err := getAbsPath(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := index.TriggerUploadAPI(cfg, absPath); err != nil {
				fmt.Println("Error during indexing:", err)
				return
			}
			if monitorProgress {
				//sleep for 4 seconds
				time.Sleep(4 * time.Second)
				if err := info.MonitorProgress(cfg, absPath); err != nil {
					fmt.Println("Error monitoring progress:", err)
				}
			}
		},
	}
	indexCmd.Flags().BoolVar(&monitorProgress, "monitor-progress", true, "Monitor the progress of the repository upload")

	// `unindex` command to manually index a repository
	var unindexCmd = &cobra.Command{
		Use:   "unindex [repo_path]",
		Short: "Unindex a specific repository",
		Long:  "Request a repo be unindex by greptile",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoPath := "."
			if len(args) > 0 {
				repoPath = args[0]
			}
			absPath, err := getAbsPath(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := index.TriggerDeleteAPI(cfg, absPath); err != nil {
				fmt.Println("Error during unindexing:", err)
			}
		},
	}

	// `check-progress` command to check upload progress
	var checkProgressCmd = &cobra.Command{
		Use:   "check-progress [repo_path]",
		Short: "Check the progress of the repository upload",
		Long:  "Check the progress of the repository upload by querying the Greptile API.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoPath := "."
			if len(args) > 0 {
				repoPath = args[0]
			}
			absPath, err := getAbsPath(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			filesProcessed, numFiles, err := info.CheckProgress(cfg, absPath)
			if err != nil {
				fmt.Println("Error checking progress:", err)
				return
			}

			// Calculate the progress percentage
			if numFiles == 0 {
				fmt.Println("Number of files is zero, cannot calculate progress")
				return
			}
			progress := (float64(filesProcessed) / float64(numFiles)) * 100

			// Output the progress information
			fmt.Printf("Files Processed: %d/%d (%.2f%%)\n", filesProcessed, numFiles, progress)
		},
	}

	// `monitor-progress` command to check upload progress
	var monitorProgressCmd = &cobra.Command{
		Use:   "monitor-progress [repo_path]",
		Short: "Monitor the progress of the repository upload",
		Long:  "Monitor the progress of the repository upload by querying the Greptile API.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoPath := "."
			if len(args) > 0 {
				repoPath = args[0]
			}
			absPath, err := getAbsPath(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			if err := info.MonitorProgress(cfg, absPath); err != nil {
				fmt.Println("Error monitoring progress:", err)
			}
		},
	}

	// `query` command to submit a semantic query
	var queryCmd = &cobra.Command{
		Use:   "query [semantic_query] [repo_path]",
		Short: "Submit a semantic query about the codebase",
		Long:  "Submit a natural language query about the codebase and get a natural language answer with a list of relevant code references (filepaths, line numbers, etc).",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			semanticQuery := args[0]
			repoPath := "."
			if len(args) > 1 {
				repoPath = args[1]
			}
			absPath, err := getAbsPath(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Call the function to handle the query (this function needs to be implemented)
			if err := semantic.HandleQuery(cfg, semanticQuery, absPath); err != nil {
				fmt.Println("Error during query:", err)
			}
		},
	}

	// `search` command to submit a search query
	var searchCmd = &cobra.Command{
		Use:   "search [search_query] [repo_path]",
		Short: "Submit a search query about the codebase",
		Long:  "Submit a natural language search query about the codebase and get a list of relevant code references (filepaths, line numbers, etc).",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			searchQuery := args[0]
			repoPath := "."
			if len(args) > 1 {
				repoPath = args[1]
			}
			absPath, err := getAbsPath(repoPath)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Call the function to handle the search
			if err := semantic.HandleSearch(cfg, searchQuery, absPath); err != nil {
				fmt.Println("Error during search:", err)
			}
		},
	}

	// `getEnabledDirectories` command to print the list of enabled directories
	var getEnabledDirsCmd = &cobra.Command{
		Use:   "autoindex-list",
		Short: "Print the list of enabled directories",
		Long:  "Print the list of directories that are enabled for autoupload.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfg.AutouploadDirs) == 0 {
				fmt.Println("No directories are enabled for autoupload.")
				return
			}

			fmt.Println("Enabled directories for autoupload:")
			for _, dir := range cfg.AutouploadDirs {
				fmt.Println(dir)
			}
		},
	}

	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(unindexCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(checkProgressCmd)
	rootCmd.AddCommand(monitorProgressCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(getEnabledDirsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
