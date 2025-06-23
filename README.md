# kiepscy-cli

![Go](https://img.shields.io/badge/Go-1.22-blue.svg?style=flat-square&logo=go)
![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-Terminal%20UI-purple.svg?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)

## W tym kraju nie ma pracy dla ludzi z moim wykształceniem 🍺

`kiepscy-cli` is a small, terminal-based program written in Go, utilizing the [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss) libraries. It allows you to play any episode of the iconic Polish TV series **"Świat według Kiepskich"** directly from your terminal!

![Peek 2025-06-23 20-57](https://github.com/user-attachments/assets/aa55ac92-5187-42bf-9495-5c455e4a9425)

## Panie! Afera jest! 💥 (Features)

* **Season Browse:** An intuitive list of all available seasons.
* **Episode Listing:** After selecting a season, you'll see a list of all episodes within that season.
* **Fuzzy Finding:** Quickly search for seasons and episodes by typing partial names.
* **Terminal Playback:** Select an episode, and it will be played using the external `mpv` media player.

## Wie pan, to nie są tanie rzeczy! 💰 (Dependencies)

For the project to function correctly, you'll need the following dependencies:

### Go Modules:

* `github.com/charmbracelet/bubbletea`
* `github.com/charmbracelet/bubbles`
* `github.com/charmbracelet/lipgloss`

### External Video Player:

* `mpv` - A lightweight and versatile media player. It must be installed on your system and accessible within your `PATH` environment variable.

    **`mpv` Installation (examples):**
    * **macOS:** `brew install mpv`
    * **Linux (Debian/Ubuntu):** `sudo apt-get install mpv`
    * **Windows:** Download the installer from `mpv.io` or use a package manager (e.g., `scoop install mpv` / `choco install mpv`).

## No i gitara! 🎸 (Installation and Usage)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Kiki-Bouba-Game-Studio/kiepscy-cli.git # Change to your repository URL
    cd kiepscy-cli
    ```

2.  **Build the application:**
    ```bash
    go build
    ```
    This will compile the project and create an executable file named `kiepscy-cli` (or `kiepscy-cli.exe` on Windows) in the current directory.

3.  **Run the application:**
    ```bash
    ./kiepscy-cli # On Linux/macOS
    .\kiepscy-cli.exe # On Windows
    ```
    You can also run the application directly without prior compilation:
    ```bash
    go run main.go
    ```

## Co jest, kurde? Noc jest, kurde! No idę, kurde! 🌛 (Controls)

The program's interface is intuitive and keyboard-driven:

* **`↑` (Up Arrow) / `↓` (Down Arrow):** Navigate through the list (seasons or episodes).
* **`Enter` / `Space`:**
    * When on the season list: Selects a season and moves to its episode list.
    * When on the episode list: Plays the selected episode using `mpv`.
* **`Backspace` / `Esc`:**
    * When on the episode list: Go back to the season list.
    * When on the season list: Exit the application.
* **Any text input:** Typing activates fuzzy finding, filtering the current list.
* **`Ctrl+C`:** Immediately exit the application.

## Albo jutro się spotkamy wiesz gdzie. U adwokata! 💼 (License)

This project is licensed under the MIT License. See the `LICENSE` file for more details.
