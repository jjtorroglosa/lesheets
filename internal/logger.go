package internal

func Println(args ...any) (n int) {
	println(args)
	return 0
}

func Printf(format string, a ...any) (n int) {
	for _, a := range a {
		println(a)
	}
	println(format)
	return 0
}

func Fatalf(format string, a ...any) (n int) {
	println(a)
	panic(format)
}
