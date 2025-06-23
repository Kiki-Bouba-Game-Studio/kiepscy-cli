package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"encoding/json"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	currentSeason *Season
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

func (m *model) selectPreviousSeason() {
	if m.currentSeason != nil {
		for i, season := range m.seasons {
			if season.Name == m.currentSeason.Name {
				m.list.Select(i)
				break
			}
		}
	}
}

func (m *model) createSeasonsList() {
	items := make([]list.Item, len(m.seasons))
	for i, season := range m.seasons {
		items[i] = season
	}
	m.list.SetItems(items)
	m.list.ResetFilter()
	m.list.Title = "Seasons"
	m.selectPreviousSeason()
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter", " "),
				key.WithHelp("enter/space", "select"),
			),
			key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "search all"),
			),
		}
	}
	m.state = "seasons"
}

func (m *model) createEpisodesList(season Season) {
	m.currentSeason = &season
	items := make([]list.Item, len(season.Episodes))
	for i, episode := range season.Episodes {
		items[i] = episode
	}
	m.list.SetItems(items)
	m.list.ResetFilter()
	m.list.Title = season.Name
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter", " "),
				key.WithHelp("enter/space", "select"),
			),
			key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "back"),
			),
		}
	}
	m.state = "episodes"
}

func (m *model) searchGlobally() {
	var allEpisodes []list.Item
	for _, season := range m.seasons {
		for _, episode := range season.Episodes {
			allEpisodes = append(allEpisodes, episode)
		}
	}
	m.list.SetItems(allEpisodes)
	m.list.ResetFilter()
	m.list.Title = "Global Search"
	m.state = "global_search"
}

func playVideo(url string, title string) {

	cmd := exec.Command("mpv", url, "--force-media-title="+title)
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
	m := model{
		state:   "seasons",
		seasons: seasons,
		list:    l,
	}
	m.createSeasonsList()
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "s":
			if m.state == "seasons" {
				m.searchGlobally()
				return m, nil
			}
		case "enter", " ":
			switch m.state {
			case "seasons":
				if selected, ok := m.list.SelectedItem().(Season); ok {
					m.createEpisodesList(selected)
				}
			case "episodes":
				if selected, ok := m.list.SelectedItem().(Episode); ok {
					playVideo(selected.Url(), selected.Title_)
				}
			case "global_search":
				if selected, ok := m.list.SelectedItem().(Episode); ok {
					playVideo(selected.Url(), selected.Title_)
				}
			}
		case "backspace":
			if m.state == "episodes" || m.state == "global_search" {
				m.createSeasonsList()
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
