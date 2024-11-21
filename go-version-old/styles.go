package main

import "github.com/charmbracelet/lipgloss"

var docStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2)

var highlightedTextStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4)

var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("42"))

var errorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

var profileNameText = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))
