package shared

import (
	"log"
	"os"
	"strings"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// InitGitLabClient initializes the GitLab API client.
func InitGitLabClient() *gitlab.Client {
	token := os.Getenv("GITLAB_TOKEN")
	url := os.Getenv("GITLAB_URL")
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// GetOpenMergeRequests returns all open MRs in the given project.
func GetOpenMergeRequests(client *gitlab.Client, projectID string) ([]*gitlab.MergeRequest, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State: gitlab.Ptr("opened"),
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}
	mrsBasic, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, opts)
	if err != nil {
		return nil, err
	}

	var mrs []*gitlab.MergeRequest
	for _, mrBasic := range mrsBasic {
		mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrBasic.IID, nil)
		if err != nil {
			return nil, err
		}
		mrs = append(mrs, mr)
	}
	return mrs, nil
}

// GetChangedFiles fetches the list of modified file paths in a MR.
func GetChangedFiles(client *gitlab.Client, projectID string, mrID int) ([]string, error) {
	diffs, _, err := client.MergeRequests.ListMergeRequestDiffs(projectID, mrID, nil)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, diff := range diffs {
		if diff.NewPath != "" {
			files = append(files, diff.NewPath)
		}
		if diff.OldPath != "" {
			files = append(files, diff.OldPath)
		}
	}
	// Remove duplicates
	uniqueFiles := make(map[string]struct{})
	for _, file := range files {
		uniqueFiles[file] = struct{}{}
	}
	files = files[:0]
	for file := range uniqueFiles {
		files = append(files, file)
	}
	return files, nil
}

// ClassifyFiles groups file paths into categories for the prompt.
func ClassifyFiles(files []string) map[string][]string {
	classified := map[string][]string{
		"docs":          {},
		"infrastructure": {},
		"source_code":   {},
		"tests":         {},
		"other":         {},
	}

	for _, file := range files {
		switch {
		case strings.HasSuffix(file, ".md") || strings.HasSuffix(file, ".rst"):
			classified["docs"] = append(classified["docs"], file)
		case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") || strings.HasSuffix(file, ".tf") || strings.Contains(file, "helm") || strings.Contains(file, "infra"):
			classified["infrastructure"] = append(classified["infrastructure"], file)
		case strings.HasSuffix(file, ".go") || strings.HasSuffix(file, ".py") || strings.HasSuffix(file, ".ts") || strings.HasSuffix(file, ".js") || strings.HasSuffix(file, ".java") || strings.HasSuffix(file, ".rb"):
			classified["source_code"] = append(classified["source_code"], file)
		case strings.Contains(file, "test") || strings.Contains(file, "spec"):
			classified["tests"] = append(classified["tests"], file)
		default:
			classified["other"] = append(classified["other"], file)
		}
	}
	return classified
}

// AddLabel adds a GitLab label to a MR if not already present.
func AddLabel(client *gitlab.Client, projectID string, mrID int, label string) {
	mr, _, err := client.MergeRequests.GetMergeRequest(projectID, mrID, nil)
	if err != nil {
		log.Printf("Failed to get MR !%d: %v", mrID, err)
		return
	}
	// Avoid duplicate labels
	for _, l := range mr.Labels {
		if l == label {
			return
		}
	}
	mr.Labels = append(mr.Labels, label)
	_, _, err = client.MergeRequests.UpdateMergeRequest(projectID, mrID, &gitlab.UpdateMergeRequestOptions{
		Labels: &mr.Labels,
	})
	if err != nil {
		log.Printf("Failed to update labels on MR !%d: %v", mrID, err)
	}
}

// PostUniqueComment creates or updates a comment containing the marker tag.
func PostUniqueComment(client *gitlab.Client, projectID string, mrID int, body string, marker string) {
	notes, _, _ := client.Notes.ListMergeRequestNotes(projectID, mrID, nil)
	for _, note := range notes {
		if strings.Contains(note.Body, "<!-- "+marker+" -->") {
			// Update existing comment
			note.Body = body
			_, _, err := client.Notes.UpdateMergeRequestNote(projectID, mrID, note.ID, &gitlab.UpdateMergeRequestNoteOptions{Body: &body})
			if err != nil {
				log.Printf("Failed to update comment for MR !%d: %v", mrID, err)
			}
			return
		}
	}
	// Create new comment
	_, _, err := client.Notes.CreateMergeRequestNote(projectID, mrID, &gitlab.CreateMergeRequestNoteOptions{Body: &body})
	if err != nil {
		log.Printf("Failed to create comment for MR !%d: %v", mrID, err)
	}
}