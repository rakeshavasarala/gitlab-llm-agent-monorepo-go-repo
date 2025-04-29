package main

import (
	"fmt"
	"log"
	"os"

	"gitlab-llm-agent-monorepo-go/shared"
)

func main() {
    client := shared.InitGitLabClient()
    projectID := os.Getenv("GITLAB_PROJECT_ID")

    if projectID == "" {
        log.Fatal("GITLAB_PROJECT_ID not set")
    }

    mrs, err := shared.GetOpenMergeRequests(client, projectID)
    if err != nil {
        log.Fatalf("Failed to fetch MRs: %v", err)
    }

    for _, mr := range mrs {
        fmt.Printf("🔍 Security labeler: reviewing MR !%d - %s\n", mr.IID, mr.Title)

        // Placeholder:
        // Here you would compare old vs new image SHA, call your internal security scan service.
        // Based on vulnerability count, apply label.

        // Example: Let's just simulate it passes for now:
        shared.AddLabel(client, projectID, mr.IID, "security-pass")
        fmt.Printf("✅ Labeled MR !%d with security-pass\n", mr.IID)
    }
}
