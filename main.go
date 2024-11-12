package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	initLoadingScreen screen = iota
	mainScreen
	loadProfileScreen
	createProfileScreen
	selectBaseProfileScreen
	backupScreen
	loadingProfileScreen
)

type model struct {
	profileManager   ProfileManager
	currentScreen    screen
	mainList         list.Model
	profileList      list.Model
	baseProfileList  list.Model
	profileNameInput textinput.Model
	profileName      string
	err              error
	spinner          spinner.Model
	statusMsg        string
	width            int
	height           int
	selectedProfile  string
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Message types
type backupCompleteMsg struct{}
type profileLoadedMsg struct{}
type initLoadCompleteMsg struct {
	profileManager ProfileManager
	err            error
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func initialModel() model {
	// Main menu items
	mainItems := []list.Item{
		item{title: "Load a profile", desc: "Load an existing profile"},
		item{title: "Create a profile", desc: "Create a new profile"},
		item{title: "Backup", desc: "Backup your profiles"},
	}

	// Profile list items
	profileItems := []list.Item{
		// item{title: "Profile1", desc: "Profile 1 description"},
		// item{title: "Profile2", desc: "Profile 2 description"},
		// item{title: "Profile3", desc: "Profile 3 description"},
	}

	// Set default list styles
	defaultListStyles := list.DefaultStyles()
	defaultListStyles.Title = defaultListStyles.Title.
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	// Initialize lists with dimensions
	mainList := list.New(mainItems, list.NewDefaultDelegate(), 0, 0)
	mainList.Title = "Main Menu"
	mainList.SetShowTitle(true)
	mainList.SetShowFilter(false)
	mainList.SetShowHelp(true)
	mainList.Styles = defaultListStyles

	profileList := list.New(profileItems, list.NewDefaultDelegate(), 0, 0)
	profileList.Title = "Select Profile"
	profileList.SetShowTitle(true)
	profileList.SetShowFilter(false)
	profileList.SetShowHelp(true)
	profileList.Styles = defaultListStyles

	baseProfileList := list.New(profileItems, list.NewDefaultDelegate(), 0, 0)
	baseProfileList.Title = "Select Base Profile"
	baseProfileList.SetShowTitle(true)
	baseProfileList.SetShowFilter(false)
	baseProfileList.SetShowHelp(true)
	baseProfileList.Styles = defaultListStyles

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter profile name"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 20

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		currentScreen:    initLoadingScreen, // Start with loading screen
		mainList:         mainList,
		profileList:      profileList,
		baseProfileList:  baseProfileList,
		profileNameInput: ti,
		spinner:          s,
		width:            80,
		height:           20,
	}
}

// Simulate initial loading
func performInitialLoad() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		profileManager, err := NewProfileManager()
		if err != nil {
			return initLoadCompleteMsg{
				err: err,
			}
		}
		return initLoadCompleteMsg{
			profileManager: *profileManager,
		}
	})
}

// Simulate backup process
func performBackup() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		return backupCompleteMsg{}
	})
}

// Simulate profile loading process
func loadProfile() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		return profileLoadedMsg{}
	})
}

func (m model) Init() tea.Cmd {
	// Start with both the spinner and the initial loading timer
	return tea.Batch(
		m.spinner.Tick,
		performInitialLoad(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.mainList.SetSize(msg.Width-h, msg.Height-v)
		m.profileList.SetSize(msg.Width-h, msg.Height-v)
		m.baseProfileList.SetSize(msg.Width-h, msg.Height-v)

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.currentScreen == backupScreen ||
				m.currentScreen == loadingProfileScreen ||
				m.currentScreen == initLoadingScreen {
				// Don't allow escape during operations
				return m, nil
			}
			m.currentScreen = mainScreen
			m.err = nil
			m.statusMsg = ""
			return m, nil
		}

	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		return m, spinnerCmd

	case initLoadCompleteMsg:
		init := initLoadCompleteMsg(msg)
		if init.err != nil {
			m.err = init.err
		} else {
			m.profileManager = init.profileManager
			m.currentScreen = mainScreen
		}
		return m, nil

	case backupCompleteMsg:
		m.currentScreen = mainScreen
		m.statusMsg = "Backup completed successfully!"
		return m, nil

	case profileLoadedMsg:
		m.currentScreen = mainScreen
		m.statusMsg = fmt.Sprintf("Switched to profile %s", m.selectedProfile)
		return m, nil
	}

	switch m.currentScreen {
	case mainScreen:
		m.mainList, cmd = m.mainList.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
			switch m.mainList.SelectedItem().FilterValue() {
			case "Load a profile":
				m.currentScreen = loadProfileScreen
			case "Create a profile":
				m.currentScreen = createProfileScreen
			case "Backup":
				m.currentScreen = backupScreen
				return m, tea.Batch(
					m.spinner.Tick,
					performBackup(),
				)
			}
		}

	case loadProfileScreen:
		m.profileList, cmd = m.profileList.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
			selected := m.profileList.SelectedItem()
			m.selectedProfile = selected.FilterValue()
			m.currentScreen = loadingProfileScreen
			return m, tea.Batch(
				m.spinner.Tick,
				loadProfile(),
			)
		}

	case createProfileScreen:
		var tiCmd tea.Cmd
		m.profileNameInput, tiCmd = m.profileNameInput.Update(msg)
		cmd = tiCmd

		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "enter":
				if isValidProfileName(m.profileNameInput.Value()) {
					m.profileName = m.profileNameInput.Value()
					m.currentScreen = selectBaseProfileScreen
				} else {
					m.err = fmt.Errorf("profile name must contain only alphanumeric characters")
				}
			}
		}

	case selectBaseProfileScreen:
		m.baseProfileList, cmd = m.baseProfileList.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
			selected := m.baseProfileList.SelectedItem()
			fmt.Printf("Created profile '%s' based on '%s'\n", m.profileName, selected.FilterValue())
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	var sb strings.Builder

	// Always show status message at the top if it exists
	if m.statusMsg != "" {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")). // Green color
			Render(m.statusMsg) + "\n\n")
	}

	switch m.currentScreen {
	case initLoadingScreen:
		if m.err != nil {
			sb.WriteString(fmt.Sprintf("Error initializing: %s", m.err.Error()))
		} else {
			sb.WriteString(fmt.Sprintf("Initializing... %s", m.spinner.View()))
		}

	case mainScreen:
		return sb.String() + "\n" + m.mainList.View()

	case loadProfileScreen:
		return sb.String() + "\n" + m.profileList.View()

	case createProfileScreen:
		sb.WriteString("Enter profile name:\n")
		sb.WriteString(m.profileNameInput.View())
		if m.err != nil {
			sb.WriteString("Error: " + m.err.Error())
		}

	case selectBaseProfileScreen:
		sb.WriteString("Select base profile for '" + m.profileName + "':")
		sb.WriteString(m.baseProfileList.View())

	case backupScreen:
		sb.WriteString("Backing up profiles... ")
		sb.WriteString(m.spinner.View())

	case loadingProfileScreen:
		sb.WriteString(fmt.Sprintf("Loading profile %s... ", m.selectedProfile))
		sb.WriteString(m.spinner.View())
	}

	return docStyle.Render(sb.String())
}

func isValidProfileName(name string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9\\-_]+$", name)
	return match
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
