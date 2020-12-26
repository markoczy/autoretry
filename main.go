package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func main() {
	timeoutFlag := flag.Duration("timeout", 1*time.Hour, "timeout for command retry")
	nologFlag := flag.Bool("nolog", false, "disable all logging")
	maxFlag := flag.Int("max", -1, "max retries '-1' for infinite")
	flag.Parse()

	timeout := *timeoutFlag
	nolog := *nologFlag
	max := *maxFlag
	args := flag.Args()

	retries := 0
	for max != -1 && retries < max {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		out, err := exec.CommandContext(ctx, args[0], args[1:]...).Output()
		if err == nil {
			fmt.Println(string(out))
			return
		}
		if !nolog {
			log.Println("Autoretry - ERROR:", err)
		}
		retries++
	}
}
