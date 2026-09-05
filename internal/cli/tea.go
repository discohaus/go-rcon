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

var mcColors = map[rune]string{
	'0': "#000000", '1': "#0000AA", '2': "#00AA00", '3': "#00AAAA", '4': "#AA0000", '5': "#AA00AA", '6': "#FFAA00", '7': "#AAAAAA",
	'8': "#555555", '9': "#5555FF", 'a': "#55FF55", 'b': "#55FFFF", 'c': "#FF5555", 'd': "#FF55FF", 'e': "#FFFF55", 'f': "#FFFFFF",
}

var placeHolderText = "rcon command ... | \"/\" for UI commands, Esc to unfocus"

func parseMinecraftCodes(input string) string {
	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")
	for lineIndex, line := range lines {
		if !strings.Contains(line, "§") {
			continue
		}
		var builder strings.Builder
		style := lipgloss.NewStyle()
		parts := strings.Split(line, "§")
		builder.WriteString(parts[0])
		for _, part := range parts[1:] {
			if part == "" {
				continue
			}
			code := rune(part[0])
			if hex, ok := mcColors[code]; ok {
				style = style.Foreground(lipgloss.Color(hex))
			} else {
				switch code {
				case 'l':
					style = style.Bold(true)
				case 'o':
					style = style.Italic(true)
				case 'n':
					style = style.Underline(true)
				case 'm':
					style = style.Strikethrough(true)
				case 'r':
					style = lipgloss.NewStyle()
				}
			}
			builder.WriteString(style.Render(part[1:]))
		}
		lines[lineIndex] = builder.String()
	}
	return strings.Join(lines, "\n")
}

type sendResult struct {
	output string
	err    error
}

type model struct {
	viewport     viewport.Model
	textarea     textarea.Model
	messages     []string
	history      []string
	historyIdx   int
	registry     commandRegistry
	executor     executor
	busy         bool
	status       string
	palette      []string
	paletteIndex int
}

func initialModel(client *rcon.Client) model { return newModel(rconExecutor{client: client}) }

func newModel(e executor) model {
	ta := textarea.New()
	ta.Placeholder = placeHolderText
	ta.Focus()
	ta.Prompt = "❯ "
	ta.CharLimit = 280
	ta.SetHeight(1)
	ta.ShowLineNumbers = false
	return model{textarea: ta, viewport: viewport.New(80, 20), historyIdx: -1, registry: newCommandRegistry(), executor: e, status: "connected"}
}

func (m model) Init() tea.Cmd { return textarea.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case sendResult:
		m.busy = false
		if message.err != nil {
			m.addMessage(fmt.Sprintf("Error: %v", message.err), lipgloss.NewStyle().Foreground(lipgloss.Color("9")))
			m.status = "error"
		} else {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
			if message.output == "" {
				message.output = "[no response]"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
			}
			m.addMessage(message.output, style)
			m.status = "ready"
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.resize(message.Width, message.Height)
	case tea.KeyMsg:
		if message.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if message.Type == tea.KeyUp && len(m.history) > 0 && m.textarea.Focused() {
			if m.historyIdx < len(m.history)-1 {
				m.historyIdx++
			}
			m.textarea.SetValue(m.history[len(m.history)-1-m.historyIdx])
			return m, nil
		}
		if message.Type == tea.KeyDown && len(m.history) > 0 && m.textarea.Focused() {
			if m.historyIdx > 0 {
				m.historyIdx--
				m.textarea.SetValue(m.history[len(m.history)-1-m.historyIdx])
			} else if m.historyIdx == 0 {
				m.historyIdx = -1
				m.textarea.Reset()
			}
			return m, nil
		}
		if m.palette != nil && (message.Type == tea.KeyUp || message.Type == tea.KeyDown) {
			if message.Type == tea.KeyUp && m.paletteIndex > 0 {
				m.paletteIndex--
			}
			if message.Type == tea.KeyDown && m.paletteIndex < len(m.palette)-1 {
				m.paletteIndex++
			}
			return m, nil
		}
		switch message.Type {
		case tea.KeyTab:
			if len(m.palette) > 0 {
				m.textarea.SetValue(m.palette[m.paletteIndex] + " ")
				m.palette = nil
			}
			return m, nil
		case tea.KeyEsc:
			m.palette = nil
			m.textarea.Blur()
			m.textarea.Placeholder = "Enter to type ..."
			return m, nil
		case tea.KeyEnter:
			if !m.textarea.Focused() {
				m.textarea.Focus()
				m.textarea.Placeholder = placeHolderText
				return m, textarea.Blink
			}
			return m.submit()
		}
	}
	var textareaCmd, viewportCmd tea.Cmd
	m.textarea, textareaCmd = m.textarea.Update(msg)
	value := strings.TrimSpace(m.textarea.Value())
	if strings.HasPrefix(value, "/") {
		m.palette = m.registry.suggestions(strings.Fields(value)[0])
		m.paletteIndex = 0
	} else {
		m.palette = nil
	}
	m.viewport, viewportCmd = m.viewport.Update(msg)
	return m, tea.Batch(textareaCmd, viewportCmd)
}

func (m *model) submit() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.textarea.Value())
	if input == "" || m.busy {
		return m, nil
	}
	m.history = append(m.history, input)
	m.historyIdx = -1
	m.addMessage("❯ "+input, lipgloss.NewStyle().Foreground(lipgloss.Color("14")))
	m.textarea.Reset()
	m.palette = nil
	parsed := parseInput(input)
	if parsed.isUI {
		item, ok := m.registry.find(parsed.command)
		if !ok {
			m.addMessage("Unknown UI command: "+parsed.command, lipgloss.NewStyle().Foreground(lipgloss.Color("9")))
			return m, nil
		}
		action := item.Handler(m)
		if action.message != "" {
			m.addMessage(action.message, lipgloss.NewStyle().Foreground(lipgloss.Color("252")))
		}
		if action.quit {
			return m, tea.Quit
		}
		return m, nil
	}
	m.busy = true
	m.status = "sending..."
	return m, func() tea.Msg {
		output, err := m.executor.Send(parsed.command)
		return sendResult{output: strings.ReplaceAll(output, "\t", "    "), err: err}
	}
}

func (m *model) addMessage(message string, style lipgloss.Style) {
	m.messages = append(m.messages, style.Render(parseMinecraftCodes(message)))
	m.refreshViewportContent()
	m.viewport.GotoBottom()
}

func (m *model) resize(width, height int) {
	m.viewport.Width = clampMinimum(1, width)
	m.viewport.Height = clampMinimum(1, height-m.textarea.Height()-4)
	m.textarea.SetWidth(clampMinimum(1, width))
	m.refreshViewportContent()
}

func (m *model) refreshViewportContent() {
	content := lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n"))
	m.viewport.SetContent(content)
}

func clampMinimum(minimum, value int) int {
	if value > minimum {
		return value
	}
	return minimum
}

func (m model) View() string {
	accent := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	statusColor := "10"
	if m.status == "error" {
		statusColor = "9"
	}
	status := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Render("● " + m.status)
	header := accent.Render("RCON Client by DiscoHaus") + "  " + status
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", clampMinimum(1, m.viewport.Width)))
	palette := ""
	if len(m.palette) > 0 {
		palette = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render(strings.Join(m.palette, "  "))
	}
	return strings.Join([]string{header, separator, m.viewport.View(), separator, m.textarea.View() + palette}, "\n")
}
