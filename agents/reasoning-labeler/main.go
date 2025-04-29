package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab-llm-agent-monorepo-go/shared"
)

type OllamaRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
}

type OllamaResponse struct {
    Response string `json:"response"`
}

func queryOllama(prompt string) (string, error) {
    req := OllamaRequest{
        Model:  "mistral",
        Prompt: prompt,
        Stream: false,
    }
    body, _ := json.Marshal(req)

    resp, err := http.Post(os.Getenv("OLLAMA_URL")+"/api/generate", "application/json", bytes.NewReader(body))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result OllamaResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    return result.Response, err
}

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
        response, err := queryOllama(prompt)
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
