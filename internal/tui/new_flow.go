// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/monamaret/biblio/internal/scaffold"
	"github.com/monamaret/biblio/internal/theme"
)

// NewFlowResult carries the answers collected by the guided `new` flow —
// the TUI fallback `biblio-cli new` falls back to whenever its flags are
// omitted (contracts/cli-commands.md `biblio-cli new`).
type NewFlowResult struct {
	SiteType string
	Domain   string
	Theme    string
}

type newFlowStep int

const (
	stepSiteType newFlowStep = iota
	stepDomain
	stepTheme
)

type newFlowModel struct {
	step       newFlowStep
	typeIndex  int
	domain     string
	themes     []theme.Preset
	themeIndex int
	result     NewFlowResult
	err        error
}

func newNewFlowModel() newFlowModel {
	return newFlowModel{step: stepSiteType}
}

func (m newFlowModel) Init() tea.Cmd { return nil }

func (m newFlowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case "ctrl+c", "esc":
		m.err = fmt.Errorf("new: guided flow cancelled")
		return m, tea.Quit
	}

	switch m.step {
	case stepSiteType:
		return m.updateSiteType(keyMsg)
	case stepDomain:
		return m.updateDomain(keyMsg)
	case stepTheme:
		return m.updateTheme(keyMsg)
	}
	return m, nil
}

func (m newFlowModel) updateSiteType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.typeIndex > 0 {
			m.typeIndex--
		}
	case "down", "j":
		if m.typeIndex < len(scaffold.SiteTypes)-1 {
			m.typeIndex++
		}
	case "enter":
		m.result.SiteType = scaffold.SiteTypes[m.typeIndex]
		m.step = stepDomain
	}
	return m, nil
}

func (m newFlowModel) updateDomain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		domain := strings.TrimSpace(m.domain)
		if domain == "" {
			return m, nil
		}
		m.result.Domain = domain
		// Theme choices are filtered to the chosen site type — the contract's
		// promise that `--theme` (and this prompt) only ever offer presets
		// valid for `--type` (contracts/cli-commands.md `biblio-cli new`).
		m.themes = theme.PresetsForType(m.result.SiteType)
		m.themeIndex = 0
		m.step = stepTheme
	case tea.KeyBackspace:
		if len(m.domain) > 0 {
			m.domain = m.domain[:len(m.domain)-1]
		}
	case tea.KeyRunes:
		m.domain += string(msg.Runes)
	}
	return m, nil
}

func (m newFlowModel) updateTheme(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.themeIndex > 0 {
			m.themeIndex--
		}
	case "down", "j":
		if m.themeIndex < len(m.themes)-1 {
			m.themeIndex++
		}
	case "enter":
		if len(m.themes) == 0 {
			m.err = fmt.Errorf("new: no theme presets are defined for site type %q", m.result.SiteType)
			return m, tea.Quit
		}
		m.result.Theme = m.themes[m.themeIndex].ID
		return m, tea.Quit
	}
	return m, nil
}

func (m newFlowModel) View() string {
	var b strings.Builder
	switch m.step {
	case stepSiteType:
		b.WriteString(RenderTitle("Choose a site type") + "\n\n")
		for i, t := range scaffold.SiteTypes {
			b.WriteString(RenderOption(t, i == m.typeIndex) + "\n")
		}
		b.WriteString("\n" + RenderHelp("↑/↓ to choose, enter to confirm"))
	case stepDomain:
		b.WriteString(RenderTitle("Enter the site's domain") + "\n\n")
		b.WriteString(PromptStyle.Render("> "+m.domain) + "\n")
		b.WriteString("\n" + RenderHelp("type the domain, enter to confirm"))
	case stepTheme:
		b.WriteString(RenderTitle(fmt.Sprintf("Choose a theme for %s", m.result.SiteType)) + "\n\n")
		if len(m.themes) == 0 {
			b.WriteString(RenderError(fmt.Sprintf("no theme presets are defined for site type %q", m.result.SiteType)) + "\n")
		}
		for i, p := range m.themes {
			b.WriteString(RenderOption(p.Name, i == m.themeIndex) + "\n")
		}
		b.WriteString("\n" + RenderHelp("↑/↓ to choose, enter to confirm"))
	}
	if m.err != nil {
		b.WriteString("\n" + RenderError(m.err.Error()))
	}
	return b.String()
}

// RunNewFlow drives the interactive site-type/domain/theme prompts for
// `biblio-cli new`, offering only the theme presets valid for the chosen
// site type, and returns the collected answers (contracts/cli-commands.md
// `biblio-cli new` — flags are the scripting path, the guided flow is the
// default one per Principle III).
func RunNewFlow() (NewFlowResult, error) {
	finalModel, err := tea.NewProgram(newNewFlowModel()).Run()
	if err != nil {
		return NewFlowResult{}, fmt.Errorf("new: running guided flow: %w", err)
	}

	m := finalModel.(newFlowModel)
	if m.err != nil {
		return NewFlowResult{}, m.err
	}
	return m.result, nil
}
