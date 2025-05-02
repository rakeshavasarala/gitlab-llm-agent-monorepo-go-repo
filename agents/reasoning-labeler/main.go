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
        prompt := fmt.Sprintf("Merge Request Title: %s\nMerge Request Description: %s\nShould we auto-approve? Answer only YES or NO.", mr.Title, mr.Description)
        response, err := shared.QueryOllama(prompt)
        if err != nil {
            log.Printf("LLM query failed for MR !%d: %v", mr.IID, err)
            continue
        }

        fmt.Printf("🤖 LLM Response for MR !%d: %s\n", mr.IID, response)

        if response != "" && (response == "YES" || response == "yes") {
            shared.AddLabel(client, projectID, mr.IID, "auto-approve")
            fmt.Printf("✅ Added label auto-approve to MR !%d\n", mr.IID)
        } else {
            shared.AddLabel(client, projectID, mr.IID, "manual-review")
            fmt.Printf("🛑 Added label manual-review to MR !%d\n", mr.IID)
        }
    }
}
