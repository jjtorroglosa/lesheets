package timer

import (
	"log"
	"time"
)

func LogElapsedTime(name string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		log.Printf("Timer:%s:%fms", name, elapsed.Seconds()*1000)
	}
}
