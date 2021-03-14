package main

import (
	"context"
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/nfishe/chip8"
	"github.com/nfishe/chip8/internal/signals"
	"github.com/nfishe/chip8/util/wait"
)

var (
	cpu *chip8.CPU

	worker = &workerImpl{
		workAvailable: make(chan interface{}),
	}
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(signals.NewContext())
	defer cancel()

	cpu = chip8.New()

	go gameLoop(ctx, func() {
		if err := cpu.Cycle(); err != nil {
			panic(err)
		}
	})

	for i := 0; i < 2; i++ {
		go wait.Until(runWorker)
	}

	<-ctx.Done()

	log.Println("Exiting...")
}

type workerImpl struct {
	workAvailable chan interface{}
}

func runWorker() {
	for worker.processNext() {
	}
}

func (*workerImpl) processNext() bool {
	return true
}

func gameLoop(ctx context.Context, fn func()) {
	runtime.LockOSThread()

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fn()
		}
	}
}
