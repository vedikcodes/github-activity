package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Event struct for GitHub event
type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}

	username := os.Args[1]
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Failed to fetch activity: %s (status %d)\n", http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	if len(events) == 0 {
		fmt.Println("No recent activity found.")
		return
	}

	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			var payload struct {
				Commits []struct{} `json:"commits"`
			}
			_ = json.Unmarshal(event.Payload, &payload)
			fmt.Printf("- Pushed %d commit(s) to %s\n", len(payload.Commits), event.Repo.Name)

		case "IssuesEvent":
			fmt.Printf("- Opened a new issue in %s\n", event.Repo.Name)

		case "WatchEvent":
			fmt.Printf("- Starred %s\n", event.Repo.Name)

		case "ForkEvent":
			fmt.Printf("- Forked %s\n", event.Repo.Name)

		default:
			fmt.Printf("- %s in %s\n", event.Type, event.Repo.Name)
		}
	}
}
