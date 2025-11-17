package cmds

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatchCommand(staticsFS embed.FS, dev bool, outputDir string, files []string, port int) {
	if len(files) < 1 {
		log.Fatal("must specify at least one file to watch")
	}

	if err := extractEmbeddedStatics(staticsFS, outputDir); err != nil {
		log.Fatalf("error extracting statics: %v", err)
	}
	hub := NewSSEHub()

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
