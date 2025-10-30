package internal

import (
	"regexp"
	"strings"
)

// superscriptMap maps digits to Unicode superscripts
var superscriptMap = map[rune]rune{
	'0': '⁰',
	'1': '¹',
	'2': '²',
	'3': '³',
	'4': '⁴',
	'5': '⁵',
	'6': '⁶',
	'7': '⁷',
	'8': '⁸',
	'9': '⁹',
}

// FormatChord formats a chord symbol like "F#min11" to "F♯m¹¹"
func FormatChord(chord string) string {
	// Replace sharps and flats
	chord = strings.ReplaceAll(chord, "#", "♯")
	chord = strings.ReplaceAll(chord, "b", "♭")

	// Replace common chord aliases
	replacements := []struct {
		pattern *regexp.Regexp
		replace string
	}{
		{regexp.MustCompile(`(?i)maj7`), "△7"},
		{regexp.MustCompile(`(?i)sus`), "ˢᵘˢ"},
		{regexp.MustCompile(`(?i)maj9`), "△9"},
		{regexp.MustCompile(`(?i)maj`), "△"},
		{regexp.MustCompile(`(?i)min`), "ₘ"},
		{regexp.MustCompile(`(?i)aug`), "+"},
		{regexp.MustCompile(`(?i)halfdim`), "ø"},
		{regexp.MustCompile(`(?i)dim`), "°"},
		{regexp.MustCompile(`(?i)m`), "ₘ"},
	}

	for _, r := range replacements {
		chord = r.pattern.ReplaceAllString(chord, r.replace)
	}

	// Superscript numbers (extensions)
	numbers := regexp.MustCompile(`([♯♭]?[0-9A-G])(.*)?`).FindStringSubmatch(chord)
	if len(numbers) == 3 {
		var sb strings.Builder
		sb.WriteString(numbers[1])
		for _, c := range numbers[2] {
			if sup, ok := superscriptMap[c]; ok {
				sb.WriteRune(sup)
			} else {
				sb.WriteRune(c)
			}
		}
		chord = sb.String()
	}

	return chord
}
