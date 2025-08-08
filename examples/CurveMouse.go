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
	if err != nil {
		fmt.Printf("Error connecting: %v", err)
	}

	time.Sleep(1 * time.Second)

	MakcuConn.MoveMouseWithCurve(100, 100, 70, 30)
	time.Sleep(100 * time.Millisecond)
	MakcuConn.MoveMouseWithCurve(-56, -200, 89, 54)
	time.Sleep(2 * time.Second)
	MakcuConn.MoveMouseWithCurve(100, 100, 3)
	MakcuConn.Close()
}
