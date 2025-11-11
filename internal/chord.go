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
	if strings.HasPrefix(chord, "N.C") {
		return chord
	}
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
		{regexp.MustCompile(`(?i)aug`), "+"},
		{regexp.MustCompile(`(?i)halfdim`), "ø"},
		{regexp.MustCompile(`(?i)dim`), "°"},
		//{regexp.MustCompile(`(?i)m`), "ₘ"},// This char request the font "Hiragino Kaku Gothic ProN"
	}

	for _, r := range replacements {
		chord = r.pattern.ReplaceAllString(chord, r.replace)
	}

	// Superscript numbers (extensions)
	numbers := regexp.MustCompile(`([♯♭]?[0-9A-G])([^/]*)(/[♯♭]?[A-G0-7].*)?`).FindStringSubmatch(chord)
	if len(numbers) == 4 {
		var sb strings.Builder
		sb.WriteString(numbers[1])
		for _, c := range numbers[2] {
			if sup, ok := superscriptMap[c]; ok {
				sb.WriteRune(sup)
			} else {
				sb.WriteRune(c)
			}
		}
		sb.WriteString(numbers[3])
		chord = sb.String()
	}

	replacements = []struct {
		pattern *regexp.Regexp
		replace string
	}{
		{regexp.MustCompile(`%`), "<span class=\"text-xs\">%</span>"},
		{regexp.MustCompile(`(?i)min|m`), "<small>m</small>"},
		{regexp.MustCompile(`(?i)min7|m7`), "<small>m7</small>"},
		{regexp.MustCompile(`(?i)/([♯♭]?[A-G1-7].*)`), "<span class=\"over\">/$1</span>"},
		{regexp.MustCompile(`(?i)(\(.+\))`), "<small>$1</small>"},
		{regexp.MustCompile(`(?i)ø`), "<sup>ø</sup>"},
		{regexp.MustCompile(`(?i)°`), "<sup>o</sup>"},
	}

	for _, r := range replacements {
		chord = r.pattern.ReplaceAllString(chord, r.replace)
	}

	return chord
}
