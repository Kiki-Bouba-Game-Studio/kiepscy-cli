package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Episode struct {
	Number string `json:"number"`
	Title  string `json:"title"`
}

type Season struct {
	Title    string    `json:"title"`
	Episodes []Episode `json:"episodes"`
}

type Video struct {
	Url_       string `json:"url"`
	Title_     string `json:"title"`
	Duration_  string `json:"duration"`
	Thumbnail_ string `json:"thumbnail"`
}

type EpisodeWithVideo struct {
	Number     string `json:"number"`
	Title      string `json:"title"`
	Url_       string `json:"url"`
	Duration_  string `json:"duration"`
	Thumbnail_ string `json:"thumbnail"`
}

type SeasonWithVideos struct {
	Title    string             `json:"title"`
	Episodes []EpisodeWithVideo `json:"episodes"`
}

func read_bytes(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("File read error: %v", err)
	}

	return bytes
}

func read_seasons_from_file() []Season {
	bytes := read_bytes("./seasons_data.json")

	var seasons []Season
	if err := json.Unmarshal(bytes, &seasons); err != nil {
		log.Fatalf("JSON parsing error: %v", err)
	}

	return seasons
}

func read_videos_from_file() []Video {
	bytes := read_bytes("../data.json")

	var videos []Video
	if err := json.Unmarshal(bytes, &videos); err != nil {
		log.Fatalf("JSON parsing error: %v", err)
	}

	return videos
}

func normalizeString(input string) string {
	replacer := strings.NewReplacer(
		"–", "-", // Zamiana półpauzy na myślnik
		"—", "-", // Zamiana pauzy na myślnik
		"„", `"`, // Zamiana cudzysłowów
		"”", `"`,
	)
	normalized := replacer.Replace(input)
	return strings.ToLower(strings.TrimSpace(normalized)) // Zamiana na małe litery i usunięcie spacji
}

func main() {

	seasons := read_seasons_from_file()
	videos := read_videos_from_file()

	var seasonsWithVideos []SeasonWithVideos
	for _, season := range seasons {
		var episodesWithVideo []EpisodeWithVideo
		for _, episode := range season.Episodes {
			var bestMatch Video
			highestScore := -1
			for _, video := range videos {
				score := fuzzy.RankMatch(normalizeString(video.Title_), normalizeString(episode.Title))
				// fmt.Printf("%s, %s, score: %d, highestScore: %d\n", video.Title_, episode.Title, score, highestScore)
				if score > highestScore {
					highestScore = score
					bestMatch = video
				}
				if video.Title_ == episode.Title {
					bestMatch = video
					break
				}
			}

			const matchThreshold = 0
			if highestScore < matchThreshold {
				bestMatch = Video{}
			}

			episodesWithVideo = append(episodesWithVideo, EpisodeWithVideo{
				Number:     episode.Number,
				Title:      episode.Title,
				Url_:       bestMatch.Url_,
				Duration_:  bestMatch.Duration_,
				Thumbnail_: bestMatch.Thumbnail_,
			})
		}
		seasonsWithVideos = append(seasonsWithVideos, SeasonWithVideos{
			Title:    season.Title,
			Episodes: episodesWithVideo,
		})
	}

	result, err := json.MarshalIndent(seasonsWithVideos, "", "  ")
	if err != nil {
		log.Fatalf("JSON serialization error: %v", err)
	}
	fmt.Printf(string(result))

}
