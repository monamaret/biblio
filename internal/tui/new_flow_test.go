package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func press(t *testing.T, m newFlowModel, keys ...tea.KeyMsg) newFlowModel {
	t.Helper()
	for _, k := range keys {
		next, _ := m.Update(k)
		m = next.(newFlowModel)
	}
	return m
}

var (
	down  = tea.KeyMsg{Type: tea.KeyDown}
	enter = tea.KeyMsg{Type: tea.KeyEnter}
)

// The guided `new` flow offers all four site types, quarterly after docs,
// and quarterly's theme step offers quarterly-default
// (specs/004-quarterly-site-type FR-001; contracts/cli.md `new`).
func TestNewFlow_OffersQuarterlyAndItsPreset(t *testing.T) {
	m := newNewFlowModel()
	view := m.View()
	last := -1
	for _, siteType := range []string{"profile", "log", "docs", "quarterly"} {
		i := strings.Index(view, siteType)
		if i < 0 {
			t.Fatalf("site type step does not offer %q:\n%s", siteType, view)
		}
		if i < last {
			t.Errorf("%q is listed out of order:\n%s", siteType, view)
		}
		last = i
	}

	m = press(t, m, down, down, down, enter)
	if m.result.SiteType != "quarterly" {
		t.Fatalf("SiteType = %q after choosing the fourth option, want quarterly", m.result.SiteType)
	}
	m = press(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("example.com")}, enter)
	if !strings.Contains(m.View(), "Quarterly Default") {
		t.Errorf("theme step does not offer Quarterly Default:\n%s", m.View())
	}
	m = press(t, m, enter)
	if m.result.Theme != "quarterly-default" || m.err != nil {
		t.Errorf("result = %+v, err = %v; want theme quarterly-default", m.result, m.err)
	}
}
