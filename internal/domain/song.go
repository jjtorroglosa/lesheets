package domain

import (
	"encoding/json"
	"errors"
	"lesheets/internal/logger"
)

type Song struct {
	FrontMatter map[string]string `json:"front_matter"`
	Sections    []Section         `json:"sections"`
}

type Section struct {
	Name  string `json:"name"`
	Lines []Line `json:"lines"`
	Break bool   `json:"break"`
}

type Line struct {
	Bars              []Bar             `json:"bars"`
	MultilineBacktick MultilineBacktick `json:"multiline_backtick"`
}

func (song *Song) PrintSong() {
	logger.Println("Frontmatter:")
	for k, v := range song.FrontMatter {
		logger.Printf("%s: %s\n", k, v)
	}
	i := 1
	for _, sec := range song.Sections {
		logger.Printf("Section: %s\n", sec.Name)
		for _, line := range sec.Lines {
			if line.MultilineBacktick.Value != "" {
				logger.Printf("MultilineBacktick: %s", line.MultilineBacktick.Value)
			} else {
				for _, bar := range line.Bars {
					logger.Printf("  Bar %d (%s) '%s': ", i+1, bar.Type, bar.BarNote)
					for _, t := range bar.Chords {
						logger.Printf("Chord (%s): %s", t.Annotation.Value, t.Value)
					}
					i++
				}
			}
			logger.Printf("\n")
		}
	}
}

func (song *Song) ToJson() (string, error) {
	j, err := json.MarshalIndent(song, "", "  ")
	if err != nil {
		return "", errors.New("error marshalling json: " + err.Error())
	}
	return string(j), nil
}

func (s *Song) DefaultLength() string {
	if s == nil {
		return "1/16"
	}
	defaultLength, ok := s.FrontMatter["L"]
	if !ok || defaultLength == "" {
		return "1/16"
	}
	return defaultLength
}

func (section *Section) IsEmpty() bool {
	return len(section.Lines) == 0 && section.Name == ""
}
