package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/discohaus/go-rcon/pkg/rcon"
)

// Minecraft Color Mapping
var mcColors = map[rune]string{
	'0': "#000000", // black
	'1': "#0000AA", // dark_blue
	'2': "#00AA00", // dark_green
	'3': "#00AAAA", // dark_aqua
	'4': "#AA0000", // dark_red
	'5': "#AA00AA", // dark_purple
	'6': "#FFAA00", // gold
	'7': "#AAAAAA", // gray
	'8': "#555555", // dark_gray
	'9': "#5555FF", // blue
	'a': "#55FF55", // green
	'b': "#55FFFF", // aqua
	'c': "#FF5555", // red
	'd': "#FF55FF", // light_purple
	'e': "#FFFF55", // yellow
	'f': "#FFFFFF", // white
}

func parseMinecraftCodes(input string) string {
	// remove Windows Carriage Returns (\r)
	input = strings.ReplaceAll(input, "\r", "")

	lines := strings.Split(input, "\n")
	var parsedLines []string

	for _, line := range lines {
		if !strings.Contains(line, "§") {
			parsedLines = append(parsedLines, line)
			continue
		}

		var builder strings.Builder
		currentStyle := lipgloss.NewStyle()

		parts := strings.Split(line, "§")
		builder.WriteString(parts[0])

		for _, part := range parts[1:] {
			if len(part) == 0 {
				continue
			}

			code := rune(strings.ToLower(string(part[0]))[0])
			text := part[1:]

			if hex, exists := mcColors[code]; exists {
				currentStyle = currentStyle.Foreground(lipgloss.Color(hex))
			} else {
				switch code {
				case 'l':
					currentStyle = currentStyle.Bold(true)
				case 'o':
					currentStyle = currentStyle.Italic(true)
				case 'n':
					currentStyle = currentStyle.Underline(true)
				case 'm':
					currentStyle = currentStyle.Strikethrough(true)
				case 'r':
					currentStyle = lipgloss.NewStyle()
				}
			}

			builder.WriteString(currentStyle.Render(text))
		}
		parsedLines = append(parsedLines, builder.String())
	}

	return strings.Join(parsedLines, "\n")
}

type model struct {
	viewport     viewport.Model
	textarea     textarea.Model
	messages     []string
	history      []string
	historyIdx   int
	senderStyle  lipgloss.Style
	headerStyle  lipgloss.Style
	headerHeight int
	err          error
	rconClient   *rcon.Client
}

func initialModel(rconClient *rcon.Client) model {
	ta := textarea.New()
	ta.Placeholder = "Enter command... (Esc to unfocus, /exit to quit)"
	ta.Focus()
	ta.Prompt = "❯ "
	ta.CharLimit = 280
	ta.SetWidth(80)
	ta.SetHeight(1)
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)

	return model{
		textarea:     ta,
		messages:     []string{},
		history:      []string{},
		historyIdx:   -1,
		viewport:     vp,
		headerHeight: 2,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		headerStyle:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")),
		rconClient:   rconClient,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			m.textarea.Blur()
			return m, nil

		case tea.KeyUp:
			if m.textarea.Focused() && len(m.history) > 0 {
				if m.historyIdx < len(m.history)-1 {
					m.historyIdx++
					m.textarea.SetValue(m.history[len(m.history)-1-m.historyIdx])
				}
				return m, nil
			}

		case tea.KeyDown:
			if m.textarea.Focused() && len(m.history) > 0 {
				if m.historyIdx > 0 {
					m.historyIdx--
					m.textarea.SetValue(m.history[len(m.history)-1-m.historyIdx])
				} else if m.historyIdx == 0 {
					m.historyIdx = -1
					m.textarea.Reset()
				}
				return m, nil
			}

		case tea.KeyEnter:
			if !m.textarea.Focused() {
				m.textarea.Focus()
				return m, textarea.Blink
			}

			input := strings.TrimSpace(m.textarea.Value())
			if input == "" {
				return m, nil
			}
			if input == "/exit" {
				return m, tea.Quit
			}

			m.history = append(m.history, input)
			m.historyIdx = -1

			parsedInput := parseMinecraftCodes(fmt.Sprintf("  §a❯§r %s", input))
			m.messages = append(m.messages, parsedInput)

			output, err := m.rconClient.Send(input)
			if err != nil {
				m.err = err
				output = fmt.Sprintf("§cFehler: %v", err)
			}

			cleanOutput := strings.ReplaceAll(output, "\t", "    ")

			parsedResponse := parseMinecraftCodes(cleanOutput)
			m.messages = append(m.messages, parsedResponse)

			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.textarea.Height() - m.headerHeight - 1
		m.textarea.SetWidth(msg.Width)
	}

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", m.viewport.Width))
	headerText := m.headerStyle.Render("RCON Client by github.com/discohaus")

	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		headerText,
		separator,
		m.viewport.View(),
		separator,
		m.textarea.View(),
	)
}
