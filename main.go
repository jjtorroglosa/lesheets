//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"nasheets/internal"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

//go:embed build/*.css
var templateFS embed.FS

func main() {
	outputDir := flag.String("d", "output", "Output dir")
	printSong := flag.Bool("p", false, "Print song")
	printTokens := flag.Bool("t", false, "Print tokens")

	// Parse CLI args
	flag.Parse()

	// Remaining non-flag arguments
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		internal.Fatalf("invalid args")
	}
	i := 0
	cmd := args[i]
	i++
	files := []string{}
	for ; i < len(args); i++ {
		files = append(files, args[i])
	}

	// Read the song file

	switch cmd {
	case "html":
		err := ExtractCSS(*outputDir)
		if err != nil {
			log.Fatalf("error extracting css: %v", err)
		}
	case "serve":
		runServe(*outputDir, 8008)
		return

	case "watch":
		for _, inputFile := range files {
			render(inputFile, *outputDir)
		}

		watch(*outputDir, 8008, func(f string) {
			render(f, *outputDir)
		}, files...)
		return
	}
	for _, inputFile := range files {
		data, err := os.ReadFile(inputFile)

		outputFilename := strings.TrimSuffix(inputFile, ".nns") + ".html"
		outputFilename = filepath.Base(outputFilename)

		if err != nil {
			internal.Fatalf("Failed to read file: %v", err)
		}

		song, err := internal.ParseSongFromString(string(data))
		if err != nil {
			internal.Fatalf("error parsing song: %v", err)
		}

		switch cmd {
		case "json":
			j, err := json.Marshal(song)
			if err != nil {
				log.Fatalf("Error marshalling json: %v", err)
			}
			fmt.Println(string(j))
		case "html":
			if *printTokens {
				lexer := internal.NewLexer(string(data))
				lexer.PrintTokens()
			}

			if *printSong {
				song.PrintSong()
			}

			render(inputFile, *outputDir)
		}
	}
}

func render(inputFile string, outputDir string) {
	// Wait until file exists (up to 10 seconds)
	timeout := time.After(3 * time.Second)
	tick := time.Tick(200 * time.Millisecond)

	for {
		select {
		case <-timeout:
			fmt.Printf("❌ Timed out waiting for file: %s\n", inputFile)
			return
		case <-tick:
			if _, err := os.Stat(inputFile); err == nil {
				// File exists — continue
				goto READ
			}
		}
	}

READ:
	data, err := os.ReadFile(inputFile)

	outputFilename := strings.TrimSuffix(inputFile, ".nns") + ".html"
	outputFilename = filepath.Base(outputFilename)

	if err != nil {
		internal.Fatalf("Failed to read file: %v", err)
	}

	song, err := internal.ParseSongFromString(string(data))
	if err != nil {
		internal.Fatalf("error parsing song: %v", err)
	}

	fmt.Printf("Rendering %s to %s\n", inputFile, outputDir+"/"+outputFilename)
	internal.RenderSongHTML(song, outputDir+"/"+outputFilename)
}

func ExtractCSS(outputDir string) error {
	// Walk through embedded FS and write any .css files to disk
	return fs.WalkDir(templateFS, "build", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".css" {
			data, err := templateFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read embedded file %s: %w", path, err)
			}
			destPath := filepath.Join(outputDir, filepath.Base(path))
			if err := os.WriteFile(destPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", destPath, err)
			}
			log.Printf("Extracted %s -> %s", path, destPath)
		}
		return nil
	})
}

func runServe(outputDir string, port int) {
	// Serve static files (HTML/CSS) from outputDir
	fs := http.FileServer(http.Dir(outputDir))

	http.Handle("/", fs)
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("🌐 Serving files from %s at http://localhost%s\n", outputDir, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func watch(outputDir string, port int, render func(f string), files ...string) {
	hub := internal.NewSSEHub()

	if len(files) < 1 {
		log.Fatal("must specify at least one file to watch")
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("creating a new watcher: %v", err)
	}
	defer w.Close()

	// Start listening for events.
	go fileLoop(w, files, func(f string) {
		hub.Broadcast("start")
		render(f)
		hub.Broadcast("reload")
	})

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

	// Add all files from the commandline.
	for _, p := range files {
		st, err := os.Lstat(p)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if st.IsDir() {
			log.Fatalf("%q is a directory, not a file", p)
		}

		// Watch the directory, not the file itself.
		err = w.Add(filepath.Dir(p))
		if err != nil {
			log.Fatalf("%q: %s", p, err)
		}
	}

	log.Println("ready; press ^C to exit")
	<-make(chan struct{}) // Block forever
}

func fileLoop(w *fsnotify.Watcher, files []string, render func(f string)) {
	i := 0
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

			// Ignore files we're not interested in. Can use a
			// map[string]struct{} if you have a lot of files, but for just a
			// few files simply looping over a slice is faster.
			var found bool
			for _, f := range files {
				if event.Op != fsnotify.Chmod && f == event.Name {
					found = true
					render(f)
				}
			}
			if !found {
				continue
			}

			// Just print the event nicely aligned, and keep track how many
			// events we've seen.
			i++
			log.Printf("%3d %s\n", i, event.Op.String())
		}
	}
}
