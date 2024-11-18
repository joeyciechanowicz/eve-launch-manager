package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
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
	baseProfileList         list.Model
	currentScreen           screen
	dump                    io.Writer
	err                     error
	isEveRunning            bool
	mainList                list.Model
	newProfileName          string
	profileList             list.Model
	profileManager          *ProfileManager
	profileNameInput        textinput.Model
	spinner                 spinner.Model
	switchToSelectedProfile string
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Message types
type backupCompleteMsg struct {
	err      error
	filename string
}
type switchedProfileMsg struct {
	err error
}
type initLoadCompleteMsg struct {
	profileManager *ProfileManager
	err            error
}
type createdProfileMsg struct {
	err error
}
type checkProcessMsg struct {
	running bool
}

func initialModel() model {
	// Initialize lists with dimensions
	mainList := list.New([]list.Item{
		item{title: "Load a profile", desc: "Load an existing profile"},
		item{title: "Create a profile", desc: "Create a new profile"},
		item{title: "Backup", desc: "Backup your profiles"},
	}, list.NewDefaultDelegate(), 0, 0)
	mainList.Title = "EVE Launcher Manager"
	mainList.SetShowFilter(false)
	mainList.SetShowStatusBar(true)
	mainList.SetFilteringEnabled(false)
	mainList.SetShowPagination(false)
	mainList.StatusMessageLifetime = time.Second * 3

	profileList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	profileList.Title = "Select Profile"
	profileList.SetShowTitle(true)
	profileList.SetShowFilter(false)
	profileList.SetShowHelp(true)

	baseProfileList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	baseProfileList.Title = "Select Base Profile"
	baseProfileList.SetShowTitle(true)
	baseProfileList.SetShowFilter(false)
	baseProfileList.SetShowHelp(true)

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter profile name"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 20

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Hamburger
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		currentScreen:    initLoadingScreen, // Start with loading screen
		mainList:         mainList,
		profileList:      profileList,
		baseProfileList:  baseProfileList,
		profileNameInput: ti,
		spinner:          s,
	}
}

func (m model) Init() tea.Cmd {
	// Start with both the spinner and the initial loading timer
	return tea.Batch(
		m.spinner.Tick,
		performInitialLoad(),
		checkProcess(),
	)
}

func (m *model) UpdateLists() {
	profiles := []list.Item{}
	baseProfiles := []list.Item{}

	baseProfiles = append(baseProfiles, item{title: "None", desc: "Create an empty profile"})
	for _, profileName := range m.profileManager.Config.Profiles {
		var desc string
		if profileName == m.profileManager.Config.ActiveProfile {
			desc = "Active"
		}

		profiles = append(profiles,
			item{title: profileName, desc: desc},
		)
		baseProfiles = append(baseProfiles,
			item{title: profileName, desc: desc},
		)
	}

	m.mainList.Title = fmt.Sprintf("EVE Launcher Manager (%s)", m.profileManager.Config.ActiveProfile)
	m.baseProfileList.SetItems(baseProfiles)
	m.profileList.SetItems(profiles)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(spinner.TickMsg); !ok && m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.mainList.SetSize(msg.Width-h, (msg.Height - v))
		m.profileList.SetSize(msg.Width-h, (msg.Height - v))
		m.baseProfileList.SetSize(msg.Width-h, (msg.Height - v))

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.currentScreen == createProfileScreen {
				// Don't allow escape during operations
				return m, nil
			}
		case "esc":
			if m.currentScreen == backupScreen ||
				m.currentScreen == loadingProfileScreen ||
				m.currentScreen == initLoadingScreen {
				// Don't allow escape during operations
				return m, nil
			}
			m.currentScreen = mainScreen
			m.err = nil

			return m, m.mainList.NewStatusMessage("")
		}

	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)

		if m.currentScreen == mainScreen {
			m.mainList.Update(msg)
		}
		return m, spinnerCmd

	case checkProcessMsg:
		m.isEveRunning = msg.running
		return m, checkProcess()

	case initLoadCompleteMsg:
		init := initLoadCompleteMsg(msg)
		if init.err != nil {
			m.err = init.err
		} else {
			m.profileManager = init.profileManager
			m.currentScreen = mainScreen

			m.UpdateLists()
		}
		return m, nil

	case backupCompleteMsg:
		m.currentScreen = mainScreen

		if msg.err != nil {
			return m, m.mainList.NewStatusMessage(msg.err.Error())
		}

		return m, m.mainList.NewStatusMessage(fmt.Sprintf("Created %s", msg.filename))

	case switchedProfileMsg:
		m.currentScreen = mainScreen

		m.UpdateLists()
		if msg.err != nil {
			return m, m.mainList.NewStatusMessage(msg.err.Error())
		}

		m.mainList.Title = fmt.Sprintf("EVE Launcher Manager (%s)", m.profileManager.Config.ActiveProfile)

		return m, tea.Batch(m.mainList.NewStatusMessage(fmt.Sprintf("Switched to profile %s", m.switchToSelectedProfile)))

	case createdProfileMsg:
		m.currentScreen = mainScreen

		m.UpdateLists()
		if msg.err != nil {
			return m, m.mainList.NewStatusMessage(msg.err.Error())
		}

		return m, tea.Batch(m.mainList.NewStatusMessage(fmt.Sprintf("Created profile %s", m.newProfileName)))
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
			m.switchToSelectedProfile = selected.FilterValue()
			m.currentScreen = loadingProfileScreen
			return m, tea.Batch(
				m.spinner.Tick,
				loadProfile(m.profileManager, m.switchToSelectedProfile),
			)
		}

	case createProfileScreen:
		var tiCmd tea.Cmd
		m.profileNameInput, tiCmd = m.profileNameInput.Update(msg)
		cmd = tiCmd

		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "enter":
				if isValidProfileName(m.profileNameInput.Value(), m) {
					m.newProfileName = m.profileNameInput.Value()
					m.baseProfileList.Title = fmt.Sprintf("Select Base Profile for '%s'", m.newProfileName)
					m.currentScreen = selectBaseProfileScreen
				} else {
					m.err = fmt.Errorf("profile name must be unique and contain only alphanumeric characters")
				}
			}
		}

	case selectBaseProfileScreen:
		m.baseProfileList, cmd = m.baseProfileList.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
			selected := m.baseProfileList.SelectedItem()

			if selected.FilterValue() == "None" {
				return m, createProfile(m.profileManager, m.newProfileName, "")
			} else {
				return m, createProfile(m.profileManager, m.newProfileName, selected.FilterValue())
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	if m.isEveRunning {
		return docStyle.Render(fmt.Sprintf("EVE Launcher is running (%s), please close %s", profileNameText.Render(m.profileManager.Config.ActiveProfile), m.spinner.View()))
	}

	var sections []string = []string{}

	switch m.currentScreen {
	case initLoadingScreen:
		if m.err != nil {
			sections = append(sections, fmt.Sprintf("Error initializing: %s", m.err.Error()))
		} else {
			sections = append(sections, highlightedTextStyle.Render(fmt.Sprintf("Initializing... %s", m.spinner.View())))
		}

	case mainScreen:
		sections = append(sections, m.mainList.View())

	case loadProfileScreen:
		sections = append(sections, m.profileList.View())

	case createProfileScreen:

		sections = append(sections, highlightedTextStyle.Render("Enter profile name:"))
		sections = append(sections, m.profileNameInput.View())
		if m.err != nil {
			sections = append(sections, errorTextStyle.Render("Error: "+m.err.Error()))
		}

	case selectBaseProfileScreen:
		sections = append(sections, m.baseProfileList.View())

	case backupScreen:
		sections = append(sections, highlightedTextStyle.Render("Backing up profiles..."+m.spinner.View()+" "))

	case loadingProfileScreen:

		sections = append(sections, highlightedTextStyle.Render(fmt.Sprintf("Loading profile %s... %s", m.switchToSelectedProfile, m.spinner.View())))
	}

	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func isValidProfileName(name string, m model) bool {
	if name == "None" {
		return false
	}

	for _, profileName := range m.profileManager.Config.Profiles {
		if name == profileName {
			return false
		}
	}

	match, _ := regexp.MatchString("^[a-zA-Z0-9\\-_]+$", name)
	return match
}

func main() {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}
	m := initialModel()
	m.dump = dump
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
