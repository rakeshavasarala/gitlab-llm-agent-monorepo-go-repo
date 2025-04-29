package shared

import (
	"log"
	"os"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

func InitGitLabClient() *gitlab.Client {
	token := os.Getenv("GITLAB_TOKEN")
	url := os.Getenv("GITLAB_URL")
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func GetOpenMergeRequests(client *gitlab.Client, projectID string) ([]*gitlab.MergeRequest, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.String("opened"),
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}
	mrs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, opts)
	return mrs, err
}

func HasRequiredLabels(mr *gitlab.MergeRequest, required []string) bool {
	labelSet := map[string]bool{}
	for _, l := range mr.Labels {
		labelSet[strings.ToLower(l)] = true
	}
	for _, r := range required {
		if !labelSet[strings.ToLower(r)] {
			return false
		}
	}
	return true
}

func ApproveMR(client *gitlab.Client, projectID string, mrID int) error {
	_, _, err := client.MergeRequestApprovals.ApproveMergeRequest(projectID, mrID, &gitlab.ApproveMergeRequestOptions{})
	return err
}

func AddLabel(client *gitlab.Client, projectID string, mrID int, label string) {
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrID, nil)
	if err != nil {
		log.Printf("Failed to get MR !%d: %v", mrID, err)
		return
	}
	found := false
	for _, l := range mr.Labels {
		if l == label {
			found = true
			break
		}
	}
	if !found {
		mr.Labels = append(mr.Labels, label)
		_, _, err = client.MergeRequests.UpdateMergeRequest(projectID, mrID, &gitlab.UpdateMergeRequestOptions{
			Labels: &mr.Labels,
		})
		if err != nil {
			log.Printf("Failed to update labels on MR !%d: %v", mrID, err)
		}
	}
}
