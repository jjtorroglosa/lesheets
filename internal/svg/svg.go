package svg

import (
	"bytes"
	"fmt"
	"html/template"
	"nasheets/internal/timer"
	"os/exec"
	"strings"
)

func AbcToHtml(sourceFile string, defaultLength string, abcInput string) (template.HTML, error) {
	abc := `
%%topspace 0
%%musicfont
%%pagewidth 700px
%%scale 1.1
%%topmargin      0px
%%botmargin      0px
%%leftmargin     0px
%%rightmargin    0px
%%titlespace     0px
` + abcInput
	svg, err := AbcToSvg(sourceFile, abc)
	if err != nil {
		return template.HTML(""), err
	}

	return template.HTML(svg), nil
}

func InlineAbcToHtml(sourceFile string, defaultLength string, abcInput string) (template.HTML, error) {
	abc := `
%%topspace 0
%%musicfont
%%pagewidth 300px
%%scale 1.1
%%topmargin      0px
%%botmargin      0px
%%leftmargin     20px
%%rightmargin    0px
%%titlespace     0px
%%map all2A * print=F
X:1
M:none
L:` + defaultLength + `
K:none clef=none stafflines=0 stem=up
%%voicemap all2A
` + abcInput
	svg, err := AbcToSvg(sourceFile, abc)
	if err != nil {
		return template.HTML(""), err
	}

	return template.HTML(svg), nil
}

func AbcToSvg(sourceFile string, abcInput string) (string, error) {
	defer timer.LogElapsedTime("RenderSvg")()
	if true {
		res, err := RenderAbcToSvg(sourceFile, abcInput)
		if err != nil {
			return "", fmt.Errorf("error rendering abc to svg: %w", err)
		}
		return res, nil

	} else {
		// Example ABC notation

		// Run the abc script with the temp file as argument
		cmd := exec.Command("dash", "/Users/jtorr/Downloads/abc2svg-trystdin/abcqjs", "tosvg.js", "-")

		cmd.Stdin = strings.NewReader(abcInput)

		// Capture stdout and stderr
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		// Execute the command
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("error running abc script: %w, stderr: %s", err, stderr.String())
		}

		// Get SVG output
		return out.String(), nil
	}
}
