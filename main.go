package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Author struct {
	Name string `json:"name"`
}

type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	Author  Author `json:"author"`
}

type Actor struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	URL   string `json:"url"`
}

type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Payload struct {
	Ref     string   `json:"ref"`
	Head    string   `json:"head"`
	Size    int      `json:"size"`
	Commits []Commit `json:"commits"`
}

type GitHubEvent struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Actor     Actor   `json:"actor"`
	Repo      Repo    `json:"repo"`
	Payload   Payload `json:"payload"`
	Public    bool    `json:"public"`
	CreatedAt string  `json:"created_at"`
}

func main() {
	// Commandline stuff
	if len(os.Args) < 2 {
		fmt.Println("Please provide a Github username")
		os.Exit(1)
	}

	userName := os.Args[1]
	fmt.Printf("Looking up activity for: %s\n\n", userName)

	apiURL := "https://api.github.com/users/" + userName + "/events"

	// GET request with error handling
	res, err := http.Get(apiURL)
	if err != nil {
		os.Exit(1)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d\n", res.StatusCode)
	}
	if err != nil {
		os.Exit(1)
	}

	var events []GitHubEvent
	json.Unmarshal(body, &events)

	// Looping over the JSON structs for output
	for _, event := range events {
		if event.Type == "PushEvent" { // Getting a suggestion to switch to tagged switching but I can't get it to work
			fmt.Printf(userName+" commited to: %s\n\n", event.Repo.Name)
		} else if event.Type == "WatchEvent" {
			fmt.Printf(userName+" starred: %s\n", event.Repo.Name)
			fmt.Println("------------")
		}
		// fmt.Printf("Event Type: %s\n", event.Type)

		if len(event.Payload.Commits) > 0 {
			fmt.Printf("Commit Message:\n %s\n", event.Payload.Commits[0].Message)
			fmt.Println("------------")
		}
	}
}
