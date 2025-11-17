package cmds

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"lesheets/internal"
	"lesheets/internal/logger"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func HtmlCommand(staticsFS embed.FS, files []string, printTokens bool, printSong bool, outputDir string) {
	for _, inputFile := range files {
		if err := extractEmbeddedStatics(staticsFS, outputDir); err != nil {
			log.Fatalf("error extracting statics: %v", err)
		}

		parser, song, err := internal.ParseSongFromFile(inputFile)
		if err != nil {
			log.Fatalf("error parsing song: %v", err)
		}

		if printTokens {
			lexer := parser.Lexer
			lexer.PrintTokens()
		}

		if printSong {
			song.PrintSong()
		}

		if err := render(false, inputFile, outputDir); err != nil {
			log.Printf("Error rendering file %s: %v\n", inputFile, err)
		}
	}
}

func render(dev bool, inputFile string, outputDir string) error {
	waitForFile(inputFile)
	defer logger.LogElapsedTime("WholeRender:" + inputFile)()
	outputFilename := strings.TrimSuffix(inputFile, ".nns") + ".html"
	if err := os.MkdirAll("output/"+filepath.Dir(outputFilename), 0755); err != nil {
		return fmt.Errorf("failed to create outupt dir: %w", err)
	}

	sourceCode, err := internal.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file %s: %w", inputFile, err)
	}

	song, err := internal.ParseSongFromStringWithFileName(inputFile, sourceCode)
	if err != nil {
		// Write the error to the output html file
		if err2 := os.WriteFile(outputDir+"/"+outputFilename, []byte(internal.RenderError(err)), 0644); err2 != nil {
			return errors.Join(err, err2)
		}
		return err
	} else {
		fmt.Printf("Rendering %s to %s\n", inputFile, outputDir+"/"+outputFilename)
		err = internal.WriteSongHtmlToFile(dev, sourceCode, song, outputDir+"/"+outputFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

func waitForFile(file string) {
	// Wait until file exists (up to 10 seconds)
	timeout := time.After(3 * time.Second)
	tick := time.Tick(200 * time.Millisecond)

	for {
		select {
		case <-timeout:
			fmt.Printf("timed out waiting for file: %s\n", file)
			return
		case <-tick:
			if _, err := os.Stat(file); err == nil {
				return
			}
		}
	}
}

func extractEmbeddedStatics(staticsFS embed.FS, outputDir string) error {
	extensions := []string{".js", ".css", ".wasm", ".woff2", ".wasm.gz"}
	// Walk through embedded FS and write any .js files to disk
	err := fs.WalkDir(staticsFS, "build", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		found := false
		for _, i := range extensions {
			if !d.IsDir() && filepath.Ext(path) == i {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
		data, err := staticsFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}
		destPath := filepath.Join(outputDir, filepath.Base(path))
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}
		log.Printf("Extracted %s -> %s", path, destPath)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
