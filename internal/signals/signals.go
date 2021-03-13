package signals

import (
	"context"
	"os"
	"os/signal"
	"time"
)

var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func NewContext() context.Context {
	return &signalContext{stopCh: SetupSignalHandler()}
}

type signalContext struct {
	stopCh <-chan struct{}
}

// Deadline implements context.Context
func (scc *signalContext) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context.Context
func (scc *signalContext) Done() <-chan struct{} {
	return scc.stopCh
}

// Err implements context.Context. If the underlying stop channel is closed, Err
// always returns context.Canceled, and nil otherwise.
func (scc *signalContext) Err() error {
	select {
	case _, ok := <-scc.Done():
		if !ok {
			// TODO: revisit this behavior when Deadline() implementation is changed
			return context.Canceled
		}
	default:
	}
	return nil
}

// Value implements context.Context
func (scc *signalContext) Value(key interface{}) interface{} {
	return nil
}
