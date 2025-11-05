package svg

import (
	"bytes"
	"html/template"
	"log"
	"nasheets/internal/timer"
	"os/exec"
	"strings"
)

func AbcToHtml(sourceFile string, defaultLength string, abcInput string) template.HTML {
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
	return template.HTML(AbcToSvg(sourceFile, abc))
}

func InlineAbcToHtml(sourceFile string, defaultLength string, abcInput string) template.HTML {
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
%%map all2A * print=A
X:1
M:none
L:` + defaultLength + `
K:none clef=none stafflines=0 stem=up
%%voicemap all2A
` + abcInput
	return template.HTML(AbcToSvg(sourceFile, abc))
}

func AbcToSvg(sourceFile string, abcInput string) string {
	defer timer.LogElapsedTime("RenderSvg")()
	if true {
		res, err := RenderAbcToSvg(sourceFile, abcInput)
		if err != nil {
			log.Fatalf("error rendering abc to svg: %v", err)
		}
		return res

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
			log.Fatalf("Error running abc script: %v\nStderr: %s", err, stderr.String())
		}

		// Get SVG output
		return out.String()
	}
}
