package main

import "strings"

// progressBar renders an ASCII progress bar
// e.g. [████████████░░░░░░░░░░░░░░░░░░]
func progressBar(value, max float64, width int) string {
	if max == 0 {
		return "[" + strings.Repeat("░", width) + "]"
	}
	pct := value / max
	if pct > 1 {
		pct = 1
	}
	filled := int(pct * float64(width))
	empty  := width - filled

	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
	return bar
}
