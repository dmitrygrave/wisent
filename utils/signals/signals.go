package signals

import (
	"os"
	"os/signal"
	"sync"
)

// interruptEvents is a list of interrupt callbacks
var interruptEvents struct {
	lock  sync.Mutex
	funcs []func()
}

// AppendInterrupt appends an interrupt to the callback array
func AppendInterrupt(fn func()) {
	interruptEvents.lock.Lock()
	interruptEvents.funcs = append(interruptEvents.funcs, fn)
	interruptEvents.lock.Unlock()
}

// HandleInterrupt loops through the callbacks and calls them
func HandleInterrupt() {
	interruptEvents.lock.Lock()
	for _, fn := range interruptEvents.funcs {
		fn()
	}
	interruptEvents.lock.Unlock()
}

// HandleInterrupts handles unix interupts
func HandleInterrupts() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	go func() {
		HandleInterrupt()
	}()
}
