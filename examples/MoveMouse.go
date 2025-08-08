package main

import (
	"fmt"
	"math"
	"time"

	"github.com/nullpkt/Makcu-Go"
)

func main() {
	ComPort, _ := makcu.Find()

	makcuConn, err := makcu.Connect(ComPort, 115200)
	if err != nil {
		fmt.Printf("Error connecting: %v", err)
	}

	time.Sleep(1 * time.Second)
	makcuConn, _ = makcu.ChangeBaudRate(makcuConn)

	time.Sleep(5 * time.Second)
	fmt.Printf("\033[2J\033[HMoving mouse in a circle...\n")
	time.Sleep(2 * time.Second)

	for i := 0; i < 5; i++ {
		a := float64(2560) / float64(1440)
		for i := 0; i < 50; i++ {
			t := 2 * math.Pi * float64(i) / float64(50)
			x := int(float64(0) + float64(25)*math.Cos(t))
			y := int(float64(0) + float64(25)*a*math.Sin(t))
			makcuConn.MoveMouse(x, y)
			time.Sleep(10 * time.Millisecond)
		}
	}

	time.Sleep(2 * time.Second)
	fmt.Printf("\033[2J\033[HScrolling mouse...\n")

	for i := 0; i < 5; i++ {
		makcuConn.ScrollMouse(-i)
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)

	for i := 0; i < 5; i++ {
		makcuConn.ScrollMouse(i)
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("Done")
	time.Sleep(50 * time.Second)
	makcuConn.Close()

}
