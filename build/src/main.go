package main

import (
	"fmt"
	"maw/context"
	"os"
	"os/signal"
	"syscall"
)

// app entrypoint
func main() {
	appCtx := context.Build()
	signal.Notify(appCtx.Signal, os.Interrupt, syscall.SIGTERM)
	fmt.Println(appCtx.App.Run())
}
