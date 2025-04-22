package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type JiraPullRequestIdentifier struct {
	URL    string `json:"url"`
	Status string `json:"status"` // Note this is the status in JIRA, not the github PR's status.
}

type DevDetails struct {
	PullRequests []JiraPullRequestIdentifier `json:"pullRequests"`
}

type DevStatusResponse struct {
	Detail []DevDetails `json:"detail"`
}

type JiraSearchResponse struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	ID string `json:"id"`
}

func getDevURL(issueID string) string {
	return fmt.Sprintf("https://multimediallc.atlassian.net/rest/dev-status/1.0/issue/details?issueId=%s&applicationType=github&dataType=pullrequest", issueID)
}

func getAuth() (string, string) {
	token := os.Getenv("JIRA_API_TOKEN")
	jiraEmail := os.Getenv("JIRA_API_EMAIL")
	return jiraEmail, token
}

func GetProjectPRKeys(epicKey string, repo_name string) []int {
	/// Get all of teh cards under a JIRA epic

	fmt.Printf("Searching for project shas for project: `%s`\n", epicKey)
	searchURL := "https://multimediallc.atlassian.net/rest/api/3/search"

	jiraEmail, token := getAuth()

	params := url.Values{}
	params.Add("jql", fmt.Sprintf("Parent = %s", epicKey))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []int{}
	}
	req.URL.RawQuery = params.Encode()

	req.SetBasicAuth(jiraEmail, token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return []int{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error getting JIRA search!")
		// Read the response body to get the error message.
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body) // We ignore the error here.
		fmt.Println(string(body[:n]))
		return []int{}
	}

	var data JiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return []int{}
	}

	return processIssues(data, repo_name)

}

func processIssues(data JiraSearchResponse, target_repo string) []int {
	// this function right now only works for a single repo.
	// Returns a list of the PR numbers.

	var (
		PRNumbers []int
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	errChan := make(chan error, len(data.Issues))   // Buffered channel for errors
	resultsChan := make(chan int, len(data.Issues)) // Buffered channel for merge SHAs

	for _, issue := range data.Issues {
		wg.Add(1)
		go func(issue Issue) {
			defer wg.Done()

			pr, err := getPRLinkForIssue(issue.ID)
			if pr == nil {
				errChan <- fmt.Errorf("pr is nil %s", issue.ID)
				return
			}

			if err != nil {
				errChan <- fmt.Errorf("error getting PR link for issue %s: %w", issue.ID, err)
				return
			}

			if pr.URL == "" {
				errChan <- fmt.Errorf("Err getting PR link for issue %s: URL is empty", issue.ID)
				return
			}

			// fmt.Println("checking the url: " + pr.URL)
			prNumber := strings.Split(pr.URL, "/")
			repo := prNumber[len(prNumber)-3]
			if repo != target_repo {
				errChan <- fmt.Errorf("Issue PR is for a separate repo: %s", repo)
				return
			}
			prNum := prNumber[len(prNumber)-1]
			num, err := strconv.Atoi(prNum)
			if err != nil {
				errChan <- fmt.Errorf("Failed to convert prNum %s to int", prNum)
			}
			resultsChan <- num
		}(issue)
	}

	wg.Wait()
	close(errChan)
	close(resultsChan)

	// Collect errors
	// for err := range errChan {
	//	fmt.Println(err)
	// }

	// Collect merge SHAs
	for PRNumber := range resultsChan {
		mu.Lock()
		PRNumbers = append(PRNumbers, PRNumber)
		mu.Unlock()
	}

	return PRNumbers
}

func getPRLinkForIssue(issueID string) (*JiraPullRequestIdentifier, error) {
	/// Get first the PRs (Jira calls them dev-status) for an issue
	jiraEmail, token := getAuth()
	devURL := getDevURL(issueID)

	req, err := http.NewRequest("GET", devURL, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(jiraEmail, token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var data DevStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Detail) == 0 || len(data.Detail[0].PullRequests) == 0 {
		fmt.Printf("No PR for issue: %s\n", issueID)
		return nil, nil // Indicate no PR found without an error
	}

	pr := data.Detail[0].PullRequests[0]
	// fmt.Printf("URL: %s\n", pr.URL)
	// fmt.Printf("Status: %s\n", pr.Status)

	return &pr, nil
}
