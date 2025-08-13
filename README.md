# Golang library For Interacting With a makcu

Idek why I made this. I don't know who in their right mind would make anything for MAKCU in golang... but yk here it is for whoever wants it :)

This project’s been a really fun learning experience, and I’m still figuring stuff out as I go... So If you’ve got any feedback, suggestions, or ideas, feel free to reach out!

Fork it, change it, break it, improve it whatever you want idrc. If you end up making something you think I’ll like, reach out. If I find it cool or interesting, I might add it to the main repo or whatever you want and I'll make sure to give credit where credit is due.

---

# Code Examples

### basic usage:
```GO 
package main

import (
   "github.com/nullpkt/Makcu-Go"
)

func main(){
    MakcuPort, err := makcu.Find()
    if err != nil {
        fmt.Print("%v", err)
        return
    }
    
    MakcuConn, err := makcu.Connect(MakcuPort, 115200)
    if err != nil {
        fmt.Print("%v", err)
        return
    }
    
    MakcuConn.MoveMouse(100, 100)
    
    MakcuConn.Close()
}
```
### How to change Baud Rate:
Make a Connection to the MAKCU using 115200 Baud (default for makcu). Then run the ChangeBuadRate func to change the baud rate to 4m. 
```go
package main

import (
   "github.com/nullpkt/Makcu-Go"
)

func main(){
    MakcuPort, err := makcu.Find()
    
    MakcuConn, err := makcu.Connect(MakcuPort, 115200)
    if err != nil {
        fmt.Print("%v", err)
        return
    }
    
   MakcuConn, err = makcu.ChangeBaudRate(MakcuConn)
    if err != nil {
        fmt.Print("%v", err)
        return
    }

   MakcuConn.Close()
    
}
```
## Function documentation

### **Functions**

- **makcu.Debug**: A boolean flag to enable or disable debug printouts.
    ```go
    makcu.Debug = true
    ```
- **makcu.Find()**: Searches for and returns the COM port associated with the makcu.
    ``` go
    ComPort, err := makcu.Find()
    ```
- **makcu.Connect(port string, baudRate int)**: Establishes a connection to the makcu via the specified COM port and baud rate, returning a makcu instance.
  ```go
  MakcuConn, err := makcu.Connect("COM3", 115200)
  ```
- **makcu.ChangeBaudRate(MakcuConn *makcu)**: Changes the baud rate to 4m baud and returns a new makcu instance using the updated baud rate.
    ```go
    MakcuConn, err := makcu.ChangeBaudRate(MakcuConn)
    ```

### **MAKCU Commands**

- **MakcuConn.Write(data []byte)**:  Sends the provided data to the makcu.
    ```go
    n, err := MakcuConn.Write([]byte("Hello World!"))
    ```
- **MakcuConn.Read(buf []byte)**: Reads data from the makcu and stores it into the provided buffer.
    ```go
    buf := make([]byte, 32)
    n, err := MakcuConn.Read(buf)
    ```
- **MakcuConn.Close()**: Closes the current connection to the makcu.
    ```go
    err := MakcuConn.Close()
    ```
- **MakcuConn.LeftDown()**: Simulates pressing the left mouse button.
    ```go
    err := MakcuConn.LeftDown()
    ```
- **MakcuConn.LeftUp()**: Simulates releasing the left mouse button.
    ```go
    err := MakcuConn.LeftUp()
    ```
- **MakcuConn.LeftClick()**: Simulates a full left mouse click (press and release).
    ```go
    err := MakcuConn.LeftClick()
    ```
- **MakcuConn.RightDown()**: Simulates pressing the right mouse button.
    ```go
    err := MakcuConn.RightDown()
    ```
- **MakcuConn.RightUp()**: Simulates releasing the right mouse button.
    ```go
    err := MakcuConn.RightUp()
    ```
- **MakcuConn.RightClick()**: Simulates a full right mouse click (press and release).
    ```go
    err := MakcuConn.RightClick()
    ```
- **MakcuConn.MiddleDown()**: Simulates pressing the middle mouse button.
    ```go
    err := MakcuConn.MiddleDown()
    ```
- **MakcuConn.MiddleUp()**: Simulates releasing the middle mouse button.
    ```go
    err := MakcuConn.MiddleUp()
    ```
- **MakcuConn.MiddleClick()**: Simulates a full middle mouse click (press and release).
    ```go
    err := MakcuConn.MiddleClick()
    ```
- **MakcuConn.Click(i int, delay time.Duration)**: Simulates a mouse click with a given delay between the press and the release.
    ```go
    // MOUSE_BUTTON_LEFT = 1, MOUSE_BUTTON_RIGHT = 2, MOUSE_BUTTON_MIDDLE = 3
    err := MakcuConn.Click(makcu.MOUSE_BUTTON_LEFT, 1 *time.Second)
    err := MakcuConn.Click(makcu.MOUSE_BUTTON_RIGHT, 1 *time.Second)
    err := MakcuConn.Click(makcu.MOUSE_BUTTON_MIDDLE, 1 *time.Second)
    ```
- **MakcuConn.MoveMouse(x, y int)**: Moves the mouse cursor over (x, y) pixels.
    ```go
    err := MakcuConn.MoveMouse(100, 100)
    ```
- **MakcuConn.MoveMouseWithCurve(x, y int, ...int)**: Moves the mouse cursor along a curve.
    ```go
    // simple curve
    err := MakcuConn.MoveMouseWithCurve(100, 100, 10)
    // custom curve trijectory
    err := MakcuConn.MoveMouseWithCurve(100, 100, 10, 50, 60)
    ```
- **MakcuConn.ScrollMouse(amount int)**: Scrolls the mouse by the specified amount (positive for up, negative for down).
    ```go
    err := MakcuConn.ScrollMouse(6)
    fix v23
    ```

  
