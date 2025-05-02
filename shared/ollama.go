package shared

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type OllamaRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
}

type OllamaResponse struct {
    Response string `json:"response"`
}

func QueryOllama(prompt string) (string, error) {
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