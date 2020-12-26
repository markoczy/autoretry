package main

import (
	"context"
	"flag"
	"log"
	"os/exec"
	"time"
)

func main() {
	timeoutFlag := flag.Duration("timeout", 1*time.Hour, "timeout for command retry")
	nologFlag := flag.Bool("nolog", false, "disable all logging")
	flag.Parse()

	timeout := *timeoutFlag
	nolog := *nologFlag
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	args := flag.Args()
	for {
		err := exec.CommandContext(ctx, args[0], args[1:]...).Run()
		if err == nil {
			return
		}
		if !nolog {
			log.Println("Autoretry - ERROR:", err)
		}
	}
}
