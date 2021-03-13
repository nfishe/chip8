package wait

import (
	"context"
	"fmt"
	"log"
	"time"
)

func Until(f func()) {
	UntilWithContext(context.Background(), func(context.Context) { f() })
}

func UntilWithContext(ctx context.Context, f func(context.Context)) {
	JitterUntilWithContext(ctx, f)
}

func JitterUntilWithContext(ctx context.Context, f func(context.Context)) {
	JitterUntil(ctx, func() { f(ctx) })
}

func JitterUntil(ctx context.Context, f func()) {
	BackoffUntil(ctx, time.Second, f)
}

func BackoffUntil(ctx context.Context, duration time.Duration, f func()) {
	t := time.NewTicker(duration)
	defer t.Stop()

	var tc <-chan time.Time
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(fmt.Sprintf("%#v", err))
				}
			}()

			f()
		}()

		select {
		case <-ctx.Done():
			return
		case <-tc:
		}
	}
}
