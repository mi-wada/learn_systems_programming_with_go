package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func singnalctx() {
	ctx := context.Background()

	siguser1ctx, cancelSiguser1ctx := signal.NotifyContext(ctx, syscall.SIGUSR1)
	defer cancelSiguser1ctx()

	siguser2ctx, cancelSiguser2ctx := signal.NotifyContext(ctx, syscall.SIGUSR2)
	defer cancelSiguser2ctx()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
		}
	}()

	select {
	case <-siguser1ctx.Done():
		fmt.Println("Received SIGUSR1")
	case <-siguser2ctx.Done():
		fmt.Println("Received SIGUSR2")
	}
}

func singnalChan() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1)

	go func() {
		for {
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
			time.Sleep(1 * time.Second)
		}
	}()

	for sig := range signals {
		fmt.Printf("sig.String(): %s\n", sig.String())
		fmt.Println("Received SIGUSR1")
	}
}

func main() {
	singnalChan()
}
