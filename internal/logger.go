package internal

import "log"

func debugf(format string, args ...any) {
	log.Printf(format, args...)
}

func Println(args ...any) (n int) {
	log.Println(args...)
	return 0
}

func Printf(format string, a ...any) {
	log.Printf(format, a...)
	// for _, a := range a {
	// 	println(a)
	// }
	// println(format)
}

func Fatalf(format string, a ...any) {
	log.Fatalf(format, a...)
	// println(a)
	// panic(format)
}
