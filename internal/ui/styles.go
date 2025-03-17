package ui

import "github.com/charmbracelet/lipgloss"

// Estilos globales
var (
    TitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FFA500"))

    SuccessStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#00FF00"))

    ErrorStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FF0000"))
	
	UseSuccess = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00"))

	UseError = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF0000"))
    WarningStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FFA500"))
)
