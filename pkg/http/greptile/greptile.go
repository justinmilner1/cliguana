package greptile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"cliguana/config"
	"cliguana/pkg/util"
)

// UploadRequest represents the structure of the payload to be sent to the API
type UploadRequest struct {
	Remote     string `json:"remote"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Reload     bool   `json:"reload"`
	Notify     bool   `json:"notify"`
}

// RepositoryInfo represents the response structure from the repository API
type RepositoryInfo struct {
	Repository      string   `json:"repository"`
	Remote          string   `json:"remote"`
	Branch          string   `json:"branch"`
	Private         bool     `json:"private"`
	Status          string   `json:"status"`
	FilesProcessed  int      `json:"filesProcessed"`
	NumFiles        int      `json:"numFiles"`
	SampleQuestions []string `json:"sampleQuestions"`
	Sha             string   `json:"sha"`
}

// Custom HTTP client with a timeout
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func SendIndexRequest(cfg *config.Config, repository string, remoteType string, branch string) error {
	uploadRequest := UploadRequest{
		Remote:     remoteType,
		Repository: repository,
		Branch:     branch,
		Reload:     false,
		Notify:     true,
	}

	payloadBytes, err := json.Marshal(uploadRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal upload request: %v", err)
	}

	req, err := http.NewRequest("POST", cfg.BaseURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+cfg.AuthToken)
	req.Header.Add("X-GitHub-Token", cfg.GithubToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make the request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		fmt.Println("Upload triggered successfully!")
	} else {
		return fmt.Errorf("failed to trigger upload, status: %s, response: %s", res.Status, string(body))
	}

	return nil
}

// SendGetInfoRequest sends a request to get repository information from the Greptile API
func SendGetInfoRequest(cfg *config.Config, repository string, remote string, branch string) (RepositoryInfo, error) {
	var repoInfo RepositoryInfo

	// Extract the repository name from the remote URL
	repoName := util.ExtractRepoName(remote)
	if repoName == "" {
		return repoInfo, fmt.Errorf("invalid remote URL: %s", remote)
	}

	// Format the repositoryId as remote:branch:owner/repository
	repositoryId := fmt.Sprintf("%s:%s:%s", util.GetRemoteType(remote), branch, repoName)

	// URL-encode the repositoryId
	url := fmt.Sprintf("%s/%s", cfg.BaseURL, url.PathEscape(repositoryId))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return repoInfo, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+cfg.AuthToken)

	res, err := httpClient.Do(req)
	if err != nil {
		return repoInfo, fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return repoInfo, fmt.Errorf("failed to read response: %v", err)
	}

	if res.StatusCode != 200 {
		return repoInfo, fmt.Errorf("received non-200 response: %s, response: %s", res.Status, string(body))
	}

	err = json.Unmarshal(body, &repoInfo)
	if err != nil {
		return repoInfo, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return repoInfo, nil
}

// SendQueryRepoRequest sends a semantic query request to the Greptile API
func SendQueryRepoRequest(cfg *config.Config, repository string, remote string, branch string, query string) (string, error) {
	url := "https://api.greptile.com/v2/query" // Corrected endpoint URL
	println("repository:", repository)
	println("remote:", remote)
	println("branch:", branch)
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{
				"id":      "<string>",
				"content": query,
				"role":    "user",
			},
		},
		"repositories": []map[string]string{
			{
				"remote":     util.GetRemoteType(remote),
				"branch":     branch,
				"repository": repository,
			},
		},
		"sessionId": "<session-id>", // Replace with actual session ID if needed
		"stream":    true,
		"genius":    true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal query request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+cfg.AuthToken)
	req.Header.Add("X-GitHub-Token", cfg.GithubToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make the request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return string(body), nil
	} else {
		return "", fmt.Errorf("failed to query repository, status: %s, response: %s", res.Status, string(body))
	}
}

// SendSearchRepoRequest sends a search query request to the Greptile API
func SendSearchRepoRequest(cfg *config.Config, repository string, remote string, branch string, query string) (string, error) {
	url := "https://api.greptile.com/v2/search"

	// Identify the remote type
	remoteType := util.GetRemoteType(remote)
	if remoteType == "" {
		return "", fmt.Errorf("invalid remote URL: %s", remote)
	}

	payload := map[string]interface{}{
		"query": query,
		"repositories": []map[string]string{
			{
				"remote":     remote,
				"branch":     branch,
				"repository": repository,
				"type":       remoteType, // Include the remote type
			},
		},
		"sessionId": "<session-id>", // Replace with actual session ID if needed
		"stream":    true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+cfg.AuthToken)
	req.Header.Add("X-GitHub-Token", cfg.GithubToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make the request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return string(body), nil
	} else {
		return "", fmt.Errorf("failed to search repository, status: %s, response: %s", res.Status, string(body))
	}
}
