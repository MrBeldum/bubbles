package textinput

// Regression test for charmbracelet/bubbles#562: navigating suggestions with
// the arrow keys must not abort the program. With no matched suggestions
// (e.g. empty input), pressing up previously drove currentSuggestionIndex to
// -1, and any later read of the current suggestion panicked with
// "index out of range [-1]".

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestIssue562_NavigateSuggestionsWithNoMatches(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.Focus()
	m.SetSuggestions([]string{"foo", "bar"})

	// No input typed: no matched suggestions. Navigate as the reporter did.
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})

	// The app reads the current suggestion while rendering suggestion
	// navigation state. This must not panic and the index must stay valid.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CurrentSuggestion() panicked after navigating with no matches: %v", r)
		}
	}()
	if got := m.CurrentSuggestion(); got != "" {
		t.Fatalf("expected no current suggestion with no matches, got %q", got)
	}
	if idx := m.CurrentSuggestionIndex(); idx < 0 {
		t.Fatalf("CurrentSuggestionIndex() is negative after navigating with no matches: %d", idx)
	}
}

func TestIssue562_NavigationStillWorksWithMatches(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.Focus()
	m.SetSuggestions([]string{"foo", "foobar"})
	m.SetValue("foo")
	m.updateSuggestions()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if got := m.CurrentSuggestion(); got != "foobar" {
		t.Fatalf("expected suggestion %q after down, got %q", "foobar", got)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if got := m.CurrentSuggestion(); got != "foo" {
		t.Fatalf("expected wrap-around to %q, got %q", "foo", got)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if got := m.CurrentSuggestion(); got != "foobar" {
		t.Fatalf("expected wrap-around to %q, got %q", "foobar", got)
	}
}
