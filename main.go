package main

import (
	"fmt"
	"os"
	"os/exec"
	_ "embed"

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

//go:embed database/seasons.json
var content []byte
func getSeasonsFromJSON() []Season {
	var payload []Season
	err := json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal()", err)
	}

	return payload
}

func playVideo(url string) {

	cmd := exec.Command("mpv", url)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())

		logfile, err := os.CreateTemp("", "kiepscy-cli-mpv-*.log")
		if err != nil {
			log.Fatal(err)
		}

		logfile.Write(stdout)

		if err := logfile.Close(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Logs can be found in ", logfile.Name())
		return
	}
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
		case "enter", " ":
			if m.state == "seasons" {
				selected, ok := m.list.SelectedItem().(Season)
				if ok {
					items := make([]list.Item, len(selected.Episodes))
					for i, episode := range selected.Episodes {
						items[i] = episode
					}
					episodeList := list.New(items, list.NewDefaultDelegate(), m.list.Width(), m.list.Height())
					episodeList.Title = selected.Name

					m.state = "episodes"
					m.list = episodeList
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
				seasonList := list.New(items, list.NewDefaultDelegate(), m.list.Width(), m.list.Height())
				seasonList.Title = "Seasons"

				m.state = "seasons"
				m.list = seasonList
				return m, nil
			}
		case "esc":
			if m.state == "episodes" {
				items := make([]list.Item, len(m.seasons))
				for i, season := range m.seasons {
					items[i] = season
				}
				seasonList := list.New(items, list.NewDefaultDelegate(), m.list.Width(), m.list.Height())
				seasonList.Title = "Seasons"
				m.state = "seasons"
				m.list = seasonList
				return m, nil
			} else if m.state == "seasons" {
				return m, tea.Quit
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
