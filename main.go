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
	if len(os.Args) < 2 {
		fmt.Println("Please provide a Github username")
		os.Exit(1)
	}

	userName := os.Args[1]
	fmt.Printf("Looking up activity for: %s\n", userName)

	apiURL := "https://api.github.com/users/" + userName + "/events"
	fmt.Println(apiURL)

	res, err := http.Get(apiURL)
	if err != nil {
		os.Exit(1)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		os.Exit(1)
	}

	var events []GitHubEvent
	json.Unmarshal(body, &events)

	for _, event := range events {
		fmt.Printf("Event Type: %s\n", event.Type)
		fmt.Printf("User: %s\n", event.Actor.Login)
		fmt.Printf("Repository: %s\n", event.Repo.Name)

		if len(event.Payload.Commits) > 0 {
			fmt.Printf("Commit Message: %s\n", event.Payload.Commits[0].Message)
		}

		fmt.Println("------------")
	}
}
