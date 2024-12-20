package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"encoding/json"
	"log"
)

// The data struct for the decoded data
// Notice that all fields must be exportable!
type Episode struct {
	Number     string `json:"number"`
	Url_       string `json:"url"`
	Title_     string `json:"title"`
	Duration_  string `json:"duration"`
	Thumbnail_ string `json:"thumbnail"`
}

type Season struct {
	Name     string    `json:"title"`
	Episodes []Episode `json:"episodes"`
}

func (e Episode) Title() string       { return e.Title_ }
func (e Episode) Description() string { return e.Url_ }
func (e Episode) Url() string         { return e.Url_ }
func (e Episode) Thumbnail() string   { return e.Thumbnail_ }
func (e Episode) Duration() string    { return e.Duration_ }
func (e Episode) FilterValue() string { return e.Title_ }

func (s Season) Title() string       { return s.Name }
func (s Season) Description() string { return fmt.Sprintf("%d episodes", len(s.Episodes)) }
func (s Season) FilterValue() string { return s.Name }

type model struct {
	list          list.Model
	seasons       []Season
	episodes      []Episode
	state         string
	currentSeason int
}

func getSeasonsFromJSON() []Season {
	content, err := os.ReadFile("./utils/seasons_with_videos.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []Season
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal()", err)
	}

	// seasons := make([]Season, len(payload))

	// for i, episodeList := range payload {
	// 	seasons[i] = Season{
	// 		Name:     fmt.Sprintf("Season %d", i+1),
	// 		Episodes: episodeList,
	// 	}
	// }

	// return seasons
	return payload
}

func playVideo(url string) {
	cmd := exec.Command("mpv", url)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) Init() tea.Cmd {
	return nil
}

func initializeModel() model {
	seasons := getSeasonsFromJSON()

	items := make([]list.Item, len(seasons))
	for i, season := range seasons {
		items[i] = season
	}

	l := list.New(items, list.NewDefaultDelegate(), 30, 15)
	l.Title = "Seasons"

	return model{
		state:   "seasons",
		seasons: seasons,
		list:    l,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "KeyEnter", " ":
			if m.state == "seasons" {
				selected, ok := m.list.SelectedItem().(Season)
				if ok {
					items := make([]list.Item, len(selected.Episodes))
					for i, episode := range selected.Episodes {
						items[i] = episode
					}
					episodeList := list.New(items, list.NewDefaultDelegate(), 30, 15)
					episodeList.Title = selected.Name

					m.state = "episodes"
					m.list = episodeList
					m.currentSeason = m.list.Index()
				}
			} else if m.state == "episodes" {
				selected, ok := m.list.SelectedItem().(Episode)
				if ok {
					playVideo(selected.Url())
				}
			}
		case "backspace":
			if m.state == "episodes" {
				items := make([]list.Item, len(m.seasons))
				fmt.Println(len(m.seasons))
				for i, season := range m.seasons {
					items[i] = season
				}
				seasonList := list.New(items, list.NewDefaultDelegate(), 30, 15)
				seasonList.Title = "Seasons"

				m.state = "seasons"
				m.list = seasonList
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	m := initializeModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
