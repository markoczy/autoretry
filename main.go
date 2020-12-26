package main

import (
	"context"
	"flag"
	"os/exec"
	"time"
)

func main() {
	timeoutFlag := flag.Duration("timeout", 1*time.Hour, "timeout for command retry")
	flag.Parse()

	timeout := *timeoutFlag
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	args := flag.Args()
	for {
		err := exec.CommandContext(ctx, args[0], args[1:]...).Run()
		if err == nil {
			return
		}
	}
}
