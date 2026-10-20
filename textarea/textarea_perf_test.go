package textarea

import (
	"strings"
	"testing"
)

// largePerfModel builds a textarea holding a large buffer: 3000 logical
// lines of ~200 runes each. See https://github.com/charmbracelet/bubbles/issues/1023.
func largePerfModel() Model {
	m := New()
	m.Focus()
	m.SetWidth(100)
	m.SetHeight(20)
	var sb strings.Builder
	for i := 0; i < 3000; i++ {
		sb.WriteString(strings.Repeat("lorem ipsum dolor sit amet ", 8))
		sb.WriteByte('\n')
	}
	m.SetValue(sb.String())
	// Warm the wrap cache so the benchmark measures render cost, not
	// one-off wrapping.
	_ = m.View()
	return m
}

// BenchmarkViewLargeBuffer reproduces charmbracelet/bubbles#1023: view()
// walks and styles the entire buffer on every frame, so render cost grows
// with buffer size. Guards the per-line render overhead in view().
func BenchmarkViewLargeBuffer(b *testing.B) {
	m := largePerfModel()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}

// BenchmarkTypeAndViewLargeBuffer measures the steady-state cost of typing
// into a large buffer and re-rendering: the "every update" path from #1023.
func BenchmarkTypeAndViewLargeBuffer(b *testing.B) {
	b.StopTimer()
	m := largePerfModel()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.InsertRune('x')
		_ = m.View()
	}
}
