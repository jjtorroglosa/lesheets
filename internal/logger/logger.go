package logger

import (
	"log"
	"strconv"
	"time"
)

var IsProd = false

func Println(args ...any) {
	log.Println(args...)
	// for _, i := range args {
	// 	print(i)
	// 	print("  ")
	// }
	// println()
}

func Printf(format string, a ...any) {
	log.Printf(format, a...)
	// print(format)
	// print(": ")
	// Println(a...)
}

func LogElapsedTime(name string) func() {
	if IsProd {
		return func() {}
	}

	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		Println("Timer:" + name + ":" + strconv.FormatFloat(elapsed.Seconds()*1000.0, 'f', 2, 64) + "ms")
	}
}
