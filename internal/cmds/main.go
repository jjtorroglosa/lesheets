package cmds

import (
	"encoding/json"
	"fmt"
	"lesheets/internal"
	"log"
)

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
