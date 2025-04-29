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
        fmt.Printf("🔍 Checking MR: !%d - %s\n", mr.IID, mr.Title)

        if shared.HasRequiredLabels(mr, []string{"auto-approve", "security-pass"}) {
            err := shared.ApproveMR(client, projectID, mr.IID)
            if err != nil {
                log.Printf("❌ Could not approve MR !%d: %v", mr.IID, err)
            } else {
                fmt.Printf("✅ Approved MR: !%d\n", mr.IID)
            }
        } else {
            fmt.Println("⏭️  Skipping MR due to missing labels.")
        }
    }
}
