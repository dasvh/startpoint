package promptui

import (
	"startpoint/tui/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: colors from theme
var promptStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("#cdd6f44")).BorderStyle(lipgloss.RoundedBorder()).Padding(1)
var inputStyle = lipgloss.NewStyle().BorderForeground(lipgloss.Color("#cdd6f44")).BorderStyle(lipgloss.NormalBorder()).Padding(1).MarginTop(1)
var errInputStyle = inputStyle.Copy().BorderForeground(lipgloss.Color("#f38ba8"))
var descriptionStyle = lipgloss.NewStyle()
var helpStyle = lipgloss.NewStyle().Padding(1, 1, 0, 1)

type keyMap struct {
	Save key.Binding
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save, k.Quit}, // first column
	}
}

var keys = keyMap{
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "ok"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

type Model struct {
	nameInput textinput.Model
	context   PromptContext
	label     string
	keys      keyMap
	help      help.Model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newWidth := min(64, msg.Width-2)
		m.nameInput.Width = newWidth
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Cmd(func() tea.Msg {
				return PromptCancelledMsg{}
			})
		case tea.KeyEnter:
			return m, tea.Cmd(func() tea.Msg {
				return PromptAnsweredMsg{
					Context: m.context,
					Input:   m.nameInput.Value(),
				}
			})
		}
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)

	return m, cmd
}

func (m Model) View() string {

	helpView := helpStyle.Render(m.help.View(m.keys))

	inputViews := []string{}
	var descStyle = descriptionStyle.Width(m.nameInput.Width)
	inputViews = append(inputViews, descStyle.Render(m.label))

	var style = inputStyle.Width(m.nameInput.Width)
	if m.nameInput.Err != nil {
		style = errInputStyle.Width(m.nameInput.Width)
	}

	inputViews = append(inputViews, style.Render(m.nameInput.View()))
	inputViews = append(inputViews, helpView)

	return promptStyle.Render(lipgloss.JoinVertical(lipgloss.Left, inputViews...))

}

func New(context PromptContext, initialValue string, label string, validator func(s string) error, w int) Model {

	theme := styles.GetTheme()
	commonStyles := styles.GetCommonStyles(theme)

	nameInput := textinput.New()
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = min(64, w-2)
	nameInput.SetValue(initialValue)
	nameInput.Prompt = ""
	if validator != nil {
		// validator blocks writing on invalid input; there is a fix for this, but it is not released yet
		nameInput.Validate = validator
	}

	help := help.New()
	help.Styles.FullKey = commonStyles.HelpKeyStyle
	help.Styles.FullDesc = commonStyles.HelpDescStyle
	help.Styles.ShortKey = commonStyles.HelpKeyStyle
	help.Styles.ShortDesc = commonStyles.HelpDescStyle
	help.Styles.ShortSeparator = commonStyles.HelpSeparatorStyle
	help.Styles.FullSeparator = commonStyles.HelpSeparatorStyle
	help.ShortSeparator = "  "
	help.FullSeparator = "  "

	return Model{
		context:   context,
		nameInput: nameInput,
		label:     label,
		keys:      keys,
		help:      help,
	}
}

type PromptContext struct {
	Key        string
	Additional interface{}
}

type PromptAnsweredMsg struct {
	Context PromptContext
	Input   string
}

type PromptCancelledMsg struct{}
