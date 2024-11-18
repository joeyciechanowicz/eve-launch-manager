package main

import "github.com/charmbracelet/lipgloss"

var appStyle = lipgloss.NewStyle().MarginTop(1).MarginLeft(2)

// var appStyle = lipgloss.NewStyle().Padding(1, 2)

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(0, 1)

var statusMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
	Render

var highlightedTextStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

// var titleStyle = lipgloss.NewStyle().
// 	Bold(true).
// 	Foreground(lipgloss.Color("#FAFAFA")).
// 	Background(lipgloss.Color("#7D56F4")).
// 	PaddingTop(2).
// 	PaddingLeft(4)

// var statusMessageStyle = lipgloss.NewStyle().
// 	Foreground(lipgloss.Color("42"))

var errorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
