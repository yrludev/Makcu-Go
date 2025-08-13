package main

import (
	"fmt"
	"time"

	makcu "github.com/nullpkt/Makcu-Go"
)

func main() {
	makcu.Debug = false
	ComPort, _ := makcu.Find()

	MakcuConn, err := makcu.Connect(ComPort, 4000000)
	if err != nil || MakcuConn == nil {
		fmt.Printf("Error connecting: %v\n", err)
		fmt.Println("No MAKCU device found or failed to connect. Exiting gracefully. 🐱")
		return
	}

	time.Sleep(1 * time.Second)

	// these are just random values just for an example.
	MakcuConn.MoveMouseWithCurve(100, 100, 10, 70, 30)
	time.Sleep(100 * time.Millisecond)
	MakcuConn.MoveMouseWithCurve(-56, -200, 10, 89, 54)
	time.Sleep(2 * time.Second)
	MakcuConn.MoveMouseWithCurve(100, 100, 3)
	MakcuConn.Close()
}
