package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var InterruptHandlerChannel chan func() = make(chan func(), 10)

func InterruptHandler() {
	var interruptHandlers []func()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	func() {
		for {
			select {
			case handler := <-InterruptHandlerChannel:
				interruptHandlers = append(interruptHandlers, handler)
			case msg := <-c:
				close(InterruptHandlerChannel)
				fmt.Printf("%s\n", msg)
				return
			}
		}
	}()

	fmt.Println("Gracefully shutting down...")
	for _, handler := range interruptHandlers {
		handler()
	}
	fmt.Println("Server was successful shutdown.")
}
