package utils

import "time"

func NewDebouncer(delay time.Duration) func(key string, fn func()) {
	debounceTimers := map[string]*time.Timer{}

	return func(key string, fn func()) {
		if timer, exists := debounceTimers[key]; exists {
			timer.Stop()
		}
		debounceTimers[key] = time.AfterFunc(delay, fn)
	}
}
