package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func performInitialLoad() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		profileManager, err := NewProfileManager()
		if err != nil {
			return initLoadCompleteMsg{
				err: err,
			}
		}
		return initLoadCompleteMsg{
			profileManager: profileManager,
		}
	})
}

func loadProfile(profileManager *ProfileManager, profile string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := profileManager.SwitchProfile(profile)
		return switchedProfileMsg{
			err: err,
		}
	})
}

func createProfile(profileManager *ProfileManager, profileName, baseProfile string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := profileManager.CreateProfile(profileName, baseProfile)
		return createdProfileMsg{
			err: err,
		}
	})
}
