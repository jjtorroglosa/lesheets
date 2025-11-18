package utils

import (
	"sync"
	"testing"
	"time"
)

func TestDebouncer(t *testing.T) {
	t.Run("it should call the function only once after the delay", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		var counter int
		debounce := NewDebouncer(100 * time.Millisecond)
		increment := func() {
			counter++
			wg.Done()
		}

		debounce("test", increment)
		debounce("test", increment)
		debounce("test", increment)

		wg.Wait()

		if counter != 1 {
			t.Errorf("Expected counter to be 1, but got %d", counter)
		}
	})

	t.Run("it should not call the function if another call is made within the delay", func(t *testing.T) {
		var counter int
		debounce := NewDebouncer(100 * time.Millisecond)
		increment := func() {
			counter++
		}

		debounce("test", increment)
		time.Sleep(50 * time.Millisecond)
		debounce("test", increment)
		time.Sleep(150 * time.Millisecond)

		if counter != 1 {
			t.Errorf("Expected counter to be 1, but got %d", counter)
		}
	})

	t.Run("it should call the function for a different key without waiting", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		var counter1 int
		var counter2 int
		debounce := NewDebouncer(100 * time.Millisecond)
		increment1 := func() {
			counter1++
			wg.Done()
		}
		increment2 := func() {
			counter2++
			wg.Done()
		}

		debounce("test1", increment1)
		debounce("test2", increment2)

		wg.Wait()

		if counter1 != 1 {
			t.Errorf("Expected counter1 to be 1, but got %d", counter1)
		}
		if counter2 != 1 {
			t.Errorf("Expected counter2 to be 1, but got %d", counter2)
		}
	})

	t.Run("it should call the function again after the delay has passed", func(t *testing.T) {
		var counter int
		debounce := NewDebouncer(100 * time.Millisecond)
		increment := func() {
			counter++
		}

		debounce("test", increment)
		time.Sleep(150 * time.Millisecond)
		debounce("test", increment)
		time.Sleep(150 * time.Millisecond)

		if counter != 2 {
			t.Errorf("Expected counter to be 2, but got %d", counter)
		}
	})
}
