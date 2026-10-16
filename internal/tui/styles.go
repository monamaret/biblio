// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package tui holds the Lip Gloss styles and Bubble Tea prompt helpers shared
// across every guided flow (the `new` flow first, in v0.01; later flows reuse
// the same look rather than redefining it).
package tui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	PromptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	SelectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	ErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	HelpStyle     = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("241"))
)

// RenderTitle renders a guided-flow step heading (e.g. "Choose a site type").
func RenderTitle(s string) string {
	return TitleStyle.Render(s)
}

// RenderError renders an inline validation/error message below a prompt.
func RenderError(s string) string {
	return ErrorStyle.Render(s)
}

// RenderHelp renders dim help text such as keybinding hints.
func RenderHelp(s string) string {
	return HelpStyle.Render(s)
}

// RenderOption renders one selectable choice in a list prompt, highlighted
// when it is the current selection.
func RenderOption(label string, selected bool) string {
	if selected {
		return SelectedStyle.Render("> " + label)
	}
	return PromptStyle.Render("  " + label)
}
