package cmds

import (
	"fmt"
	"log"
	"net/http"
)

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
