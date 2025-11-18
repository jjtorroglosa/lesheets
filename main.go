//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"embed"
	"flag"
	"fmt"
	"lesheets/internal"
	"lesheets/internal/cmds"
	"lesheets/internal/svg"
	"log"
	"os"
)

//go:embed build/*.css build/*.js build/abc2svg.woff2 build/*.wasm
var staticsFS embed.FS

//go:embed internal/svg/abc2svg/user.js internal/svg/abc2svg/tosvg.js vendorjs/abc2svg-1.cjs
var Abc2svg embed.FS

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
	printSong := flag.Bool("p", false, "Print song (only for the html command)")
	printTokens := flag.Bool("t", false, "Print tokens (only for the html command)")

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
		err := internal.RenderIndex(*outputDir, files)
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
		cmds.ServeCommand(*outputDir, 8008)
	case "watch":
		cleanup := svg.LoadJsRuntime(Abc2svg)
		defer cleanup()
		cmds.WatchCommand(staticsFS, dev, *outputDir, files, 8008)
	case "json":
		cmds.JsonCommand(files, *outputDir)
	case "html":
		cleanup := svg.LoadJsRuntime(Abc2svg)
		defer cleanup()
		cmds.HtmlCommand(staticsFS, files, *printTokens, *printSong, *outputDir)
	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}
