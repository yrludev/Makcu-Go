package main

import (
	"fmt"
	"time"

	makcu "github.com/nullpkt/Makcu-Go"
	"golang.org/x/sys/windows"
)

var Enabled bool = false

func listenForKeyPress() {
	moduser32 := windows.NewLazySystemDLL("user32.dll")
	procGetAsyncKeyState := moduser32.NewProc("GetAsyncKeyState")
	for {
		ret, _, _ := procGetAsyncKeyState.Call(uintptr(0x20))
		if ret&0x8000 != 0 {
			Enabled = !Enabled
			if Enabled {
				fmt.Println("Autoclicker started...")
			} else {
				fmt.Println("Autoclicker stopped...")
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func Autoclicker(conn *makcu.MakcuHandle) {
	for {
		if Enabled {
			conn.ClickMouse()
			time.Sleep(2 * time.Millisecond)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	ComPort, _ := makcu.Find()
	MakcuConn, err := makcu.Connect(ComPort, 115200)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}

	time.Sleep(1 * time.Second)
	MakcuConn, _ = makcu.ChangeBaudRate(MakcuConn)

	go Autoclicker(MakcuConn)

	fmt.Println("Press 'Space' to toggle autoclicker on/off.")

	listenForKeyPress()

}
