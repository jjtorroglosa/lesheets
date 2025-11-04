package svg

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"os/exec"
	"time"
)

func AbcToHtml(measure string, abcInput string) template.HTML {
	abc := `
%%topspace 0
%%pagewidth 700px
%%scale 1.1
%%topmargin      0px
%%botmargin      0px
%%leftmargin     0px
%%rightmargin    0px
%%titlespace     0px
%%measurebox 1
` + abcInput
	return template.HTML(AbcToSvg(abc))
}

func InlineAbcToHtml(measure string, abcInput string) template.HTML {
	abc := `
%%topspace 0
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
L:` + measure + `
K:none clef=none stafflines=0 stem=up
%%voicemap all2A
` + abcInput
	return template.HTML(AbcToSvg(abc))
}

func AbcToSvg(abcInput string) string {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		log.Printf("SVG rendering took: %dms", duration.Milliseconds())
	}()
	// Example ABC notation

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "input-*.abc")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Write ABC input to the temp file
	if _, err := tmpFile.WriteString(abcInput); err != nil {
		log.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close() // Close so the script can read it

	// Run the abc script with the temp file as argument
	cmd := exec.Command("abc", tmpFile.Name())

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
