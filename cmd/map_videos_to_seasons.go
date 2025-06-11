package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Episode struct {
	Number string `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url,omitempty"`
}

type Season struct {
	Title    string    `json:"title"`
	Episodes []Episode `json:"episodes"`
}

func readJSON(filename string) []Season {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Cannot read file %s: %v", filename, err)
	}
	var seasons []Season
	if err := json.Unmarshal(data, &seasons); err != nil {
		log.Fatalf("Cannot parse JSON file %s: %v", filename, err)
	}
	return seasons
}

func main() {
	seasons := readJSON("database/seasons.json")
	seasons_in_correct_order := readJSON("database/episodes-in-seasons.json")

	urlMap := make(map[string]string)
	for _, season := range seasons {
		for _, episode := range season.Episodes {
			key := episode.Title
			fmt.Println(episode.URL + " " + key)

			urlMap[key] = episode.URL
		}
	}

	for seasonIndex := range seasons_in_correct_order {
		for episodeIndex := range seasons_in_correct_order[seasonIndex].Episodes {
			episode := &seasons_in_correct_order[seasonIndex].Episodes[episodeIndex]
			key := episode.Title
			fmt.Println(key, urlMap[key])
			if url, found := urlMap[key]; found {
				episode.URL = url
			}
		}
	}

	result, err := json.MarshalIndent(seasons_in_correct_order, "", "  ")
	if err != nil {
		log.Fatalf("Cannot serialize JSON: %v", err)
	}
	if err := os.WriteFile("merged.json", result, 0644); err != nil {
		log.Fatalf("Cannot write to file merged.json: %v", err)
	}

	fmt.Println("🚀 File merged.json has been saved!")
}
