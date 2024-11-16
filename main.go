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
type Item struct {
	Url_       string `json:"url"`
	Title_     string `json:"title"`
	Duration_  string `json:"duration"`
	Thumbnail_ string `json:"thumbnail"`
}

func getVideosFromJSON() []list.Item {
	content, err := os.ReadFile("./videos.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload []Item
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal()", err)
	}

	items := []list.Item{}

	for _, video := range payload {
		items = append(items, video)
	}

	return items
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

func (i Item) Title() string       { return i.Title_ }
func (i Item) Description() string { return i.Url_ }
func (i Item) Url() string         { return i.Url_ }
func (i Item) Thumbnail() string   { return i.Thumbnail_ }
func (i Item) Duration() string    { return i.Duration_ }
func (i Item) FilterValue() string { return i.Title_ }

type model struct {
	list   list.Model
	videos []list.Item
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "Return", " ":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				playVideo(i.Url())
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

	videos := getVideosFromJSON()

	m := model{list: list.New(videos, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Sezon 1"
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
