package internal

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func (s *Song) PrintFrontmatter(sb *strings.Builder) error {
	if len(s.FrontMatter) == 0 {
		return nil
	}
	sb.WriteString("---\n")
	yml, err := yaml.Marshal(s.FrontMatter)
	if err != nil {
		return err
	}
	sb.Write(yml)

	sb.WriteString("---\n")

	return nil
}

func (song *Song) PrintSections(sb *strings.Builder) error {
	for _, s := range song.Sections {
		if s.Name != "" {
			sb.WriteString("\n")
			if s.Break {
				sb.WriteString("#- ")
			} else {
				sb.WriteString("# ")
			}
			sb.WriteString(s.Name)
			sb.WriteString("\n\n")
		}
		for _, l := range s.Lines {
			l.PrintBarsLine(sb)
			sb.WriteString("\n")
		}
	}
	return nil
}

func (line *Line) PrintBarsLine(sb *strings.Builder) {
	if line.MultilineBacktick.Value != "" {
		sb.WriteString("```\n")
		sb.WriteString(line.MultilineBacktick.Value)
		sb.WriteString("```\n")
	}

	for i, b := range line.Bars {
		isLastOfLine := i >= len(line.Bars)-1
		var next *Bar
		if !isLastOfLine {
			next = &line.Bars[i+1]
		}
		b.PrintBar(sb, next)
	}
}

func (bar *Bar) PrintBar(sb *strings.Builder, next *Bar) {
	if bar.RepeatStart {
		sb.WriteString("||: ")
	}
	if bar.BarNote != "" {
		sb.WriteString(`"`)
		sb.WriteString(bar.BarNote)
		sb.WriteString(`" `)
	}

	if bar.Backtick.Value != "" {
		sb.WriteString("`")
		sb.WriteString(bar.Backtick.Value)
		sb.WriteString("`")
	} else {
		for _, c := range bar.Chords {
			if c.Annotation.Value != "" {
				sb.WriteString("!")
				sb.WriteString(c.Annotation.Value)
				sb.WriteString("!")
			}
			sb.WriteString(c.Value)
		}
	}

	if bar.RepeatEnd {
		sb.WriteString(" :||")
	} else if bar.DoubleBarEnd {
		sb.WriteString(" ||")
	} else if next != nil {
		if !next.RepeatStart {
			sb.WriteString(" | ")
		}
	}
}

func (s *Song) PrintLesheet() string {
	sb := &strings.Builder{}

	err := s.PrintFrontmatter(sb)
	if err != nil {
		return ""
	}

	err = s.PrintSections(sb)
	if err != nil {
		return ""
	}

	return sb.String()
}
