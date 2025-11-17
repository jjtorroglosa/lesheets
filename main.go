//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"lesheets/internal"
	"lesheets/internal/logger"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

//go:embed build/*.css build/*.js build/abc2svg.woff2 build/*.wasm
var staticsFS embed.FS

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <command> <file1> ... <fileN>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	fmt.Fprintf(os.Stderr, "  watch   Watch the input files for changes, rendering the html files for them in outdir dir\n")
	fmt.Fprintf(os.Stderr, "  serve   Run a server for the previously generated html files\n")
	fmt.Fprintf(os.Stderr, "  html    Render html files for all the files provided as arguments\n")
	fmt.Fprintf(os.Stderr, "  json    Print a json representation of the song\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	outputDir := flag.String("d", "output", "Output dir")
	printSong := flag.Bool("p", false, "Print song")
	printTokens := flag.Bool("t", false, "Print tokens")

	// Parse CLI args
	flag.Parse()

	// Remaining non-flag arguments
	args := flag.Args()
	if len(args) < 1 {
		usage()
		log.Fatalf("invalid args")
	}
	cmd := args[0]
	files := []string{}
	if len(args) > 1 {
		files = args[1:]
	}
	dev := cmd != "html"
	shouldRenderIndex := cmd == "html" || cmd == "watch"

	if shouldRenderIndex {
		err := internal.RenderIndex(files)
		if err != nil {
			log.Fatalf("error rendering list: %v", err)
		}
		err = internal.WriteEditorToHtmlFile(dev, "output/editor.html")
		if err != nil {
			log.Fatalf("Error rendering editor: %v", err)
		}
	}

	switch cmd {
	case "serve":
		ServeCommand(*outputDir, 8008)
	case "watch":
		WatchCommand(dev, *outputDir, files, 8008)
	case "json":
		JsonCommand(files, *outputDir)
	case "html":
		HtmlCommand(files, *printTokens, *printSong, *outputDir)
	}
}

func ServeCommand(outputDir string, port int) {
	// Serve previously generated files (HTML/CSS) from outputDir
	fs := http.FileServer(http.Dir(outputDir))

	http.Handle("/", fs)
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("🌐 Serving files from %s at http://localhost%s\n", outputDir, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func JsonCommand(files []string, outputDir string) {
	for _, inputFile := range files {
		_, song, err := internal.ParseSongFromFile(inputFile)
		if err != nil {
			log.Fatalf("error parsing song: %v", err)
		}
		j, err := json.Marshal(song)
		if err != nil {
			log.Fatalf("Error marshalling json: %v", err)
		}
		fmt.Println(string(j))
	}
}

func HtmlCommand(files []string, printTokens bool, printSong bool, outputDir string) {
	for _, inputFile := range files {
		if err := extractEmbeddedStatics(outputDir); err != nil {
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

func WatchCommand(dev bool, outputDir string, files []string, port int) {
	if len(files) < 1 {
		log.Fatal("must specify at least one file to watch")
	}

	if err := extractEmbeddedStatics(outputDir); err != nil {
		log.Fatalf("error extracting statics: %v", err)
	}
	hub := internal.NewSSEHub()

	onChange := func(f string) {
		hub.Broadcast("start")
		if err := render(dev, f, outputDir); err != nil {
			log.Printf("Error rendering: %v\n", err)
		}
		hub.Broadcast("reload")
	}

	// Trigger a first render of all the files
	for _, f := range files {
		onChange(f)
	}

	// Start listening for events.
	go watcherFileLoop(files, onChange)

	// Serve static files from the output directory
	fs := http.FileServer(http.Dir(outputDir))
	http.Handle("/", fs)
	http.Handle("/events", hub)
	addr := fmt.Sprintf(":%d", port)

	go func() {
		fmt.Printf("🌐 Serving at http://localhost%s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Println("ready; press ^C to exit")
	<-make(chan struct{}) // Block forever
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
		log.Printf("Rendering %s to %s\n", inputFile, outputDir+"/"+outputFilename)
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

func extractEmbeddedStatics(outputDir string) error {
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

func watcherFileLoop(files []string, onChange func(f string)) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("creating a new watcher: %v", err)
	}
	defer w.Close()

	// Watch all files from the commandline.
	for _, p := range files {
		st, err := os.Lstat(p)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if st.IsDir() {
			log.Fatalf("%q is a directory, not a file", p)
		}

		// Watch the parent directory, not the file itself.
		if err = w.Add(filepath.Dir(p)); err != nil {
			log.Fatalf("%q: %s", p, err)
		}
	}

	i := 0
	const debounceDelay = 200 * time.Millisecond

	// debounceTimers := make(map[string]*time.Timer)
	debounce := createDebounce(debounceDelay)
	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Printf("ERROR: %s\n", err)
		// Read from Events.
		case event, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			// Ignore files we're not interested in.
			var found bool
			for _, f := range files {
				if event.Op != fsnotify.Chmod && f == event.Name {
					found = true
				}
			}
			if !found {
				continue
			}

			// Print the event
			i++
			log.Printf("%3d %s\n", i, event.Op.String())
			// if timer, exists := debounceTimers[event.Name]; exists {
			// 	timer.Stop()
			// }

			// debounceTimers[event.Name] = time.AfterFunc(debounceDelay, func() {
			// 	onChange(event.Name)
			// })
			debounce(event.Name, func() {
				onChange(event.Name)
			})
		}
	}
}

func createDebounce(delay time.Duration) func(key string, fn func()) {
	debounceTimers := map[string]*time.Timer{}

	return func(key string, fn func()) {
		if timer, exists := debounceTimers[key]; exists {
			timer.Stop()
		}
		debounceTimers[key] = time.AfterFunc(delay, fn)
	}
}
