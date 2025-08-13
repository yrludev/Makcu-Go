package makcu

//newest pdate
// 🐱 Imports
import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// 🐱 Debug flag
var Debug bool = false

// 🐱 Logger instance
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func DebugPrint(s string, a ...interface{}) {
	if Debug {
		logger.Debug(fmt.Sprintf("🐱 "+s, a...))
	}
}

// 🐱🐱🐱 Cat debug function! 🐱🐱🐱

func InfoPrint(s string, a ...interface{}) {
	logger.Info(fmt.Sprintf("🐱🟢 "+s, a...))
}

// 🐱🐱🐱 Cat info function! 🐱🐱🐱

func ErrorPrint(s string, a ...interface{}) {
	logger.Error(fmt.Sprintf("🐱🔴 "+s, a...))
}

// 🐱 Converts UTF16 buffer to string
func utf16ToString(buf []uint16) string {
	for i, v := range buf {
		if v == 0 {
			return syscall.UTF16ToString(buf[:i])
		}
	}
	return syscall.UTF16ToString(buf)
}

// 🐱🐱🐱 Cat string conversion! 🐱🐱🐱

// 🐱 Gets device info
func GetDeviceInfo(h, devInfo unsafe.Pointer, getProp *syscall.LazyProc, propertyCode uint32) string {
	buf := make([]uint16, 512)
	getProp.Call(uintptr(h), uintptr(devInfo), uintptr(propertyCode), 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)*2), 0)
	return utf16ToString(buf)
}

// 🐱🐱🐱 Cat device info! 🐱🐱🐱

// 🐱 Gets port name from registry
func GetPortName(hDevInfo uintptr, devInfoData *byte) (string, error) {
	setupapi := syscall.NewLazyDLL("setupapi.dll")
	setupDiOpenDevRegKey := setupapi.NewProc("SetupDiOpenDevRegKey")

	regKey, _, err := setupDiOpenDevRegKey.Call(hDevInfo, uintptr(unsafe.Pointer(devInfoData)), 0x00000001, 0, 0x00000001, 0x20019)
	if regKey == 0 || regKey == ^uintptr(0) {
		return "", fmt.Errorf("GetPortName: failed to open device registry key: %w", err)
	}

	defer syscall.RegCloseKey(syscall.Handle(regKey))

	key := registry.Key(regKey)
	portName, _, err := key.GetStringValue("PortName")
	if err != nil {
		return "", fmt.Errorf("GetPortName: failed to get 'PortName' value: %w", err)
	}

	return portName, nil
}

// 🐱🐱🐱 Cat port name! 🐱🐱🐱

// 🐱 Device property constants
const (
	DeviceDescription = 0x0 // Device Description (ex: USB-Enhanced-SERIAL CH343)
	HardwareID        = 0x1 // Hardware ID (ex: USB\VID_1A86&PID_55D3&REV_0445)
	DeviceName        = 0xC // Friendly name of the device (ex: USB-Enhanced-SERIAL CH343 (COM3) )
)

// Find searches for the MAKCU device by default name or VID/PID and returns the COM port.
func Find() (string, error) {
	setupapi := syscall.NewLazyDLL("setupapi.dll")
	getClassDevs := setupapi.NewProc("SetupDiGetClassDevsW")
	enumDeviceInfo := setupapi.NewProc("SetupDiEnumDeviceInfo")
	getDeviceProperty := setupapi.NewProc("SetupDiGetDeviceRegistryPropertyW")
	destroyDeviceList := setupapi.NewProc("SetupDiDestroyDeviceInfoList")

	guid := windows.GUID{0x4d36e978, 0xe325, 0x11ce, [8]byte{0xbf, 0xc1, 0x08, 0x00, 0x2b, 0xe1, 0x03, 0x18}}
	h, _, _ := getClassDevs.Call(uintptr(unsafe.Pointer(&guid)), 0, 0, uintptr(0x2))
	if h == 0 || h == ^uintptr(0) {
		DebugPrint("Failed to get device list\n")
		return "", fmt.Errorf("Find: failed to get device list")
	}

	defer func() {
		ret, _, _ := destroyDeviceList.Call(h)
		if ret == 0 {
			ErrorPrint("Failed to destroy device info list handle")
		}
	}()

	// 🐱🐱🐱 Cat device search! 🐱🐱🐱

	for index := 0; ; index++ {
		var devInfo struct {
			cbSize    uint32
			ClassGuid windows.GUID
			DevInst   uint32
			Reserved  uintptr
		}

		devInfo.cbSize = uint32(unsafe.Sizeof(devInfo))

		ok, _, _ := enumDeviceInfo.Call(h, uintptr(index), uintptr(unsafe.Pointer(&devInfo)))
		if ok == 0 {
			break
		}

		description := GetDeviceInfo(unsafe.Pointer(h), unsafe.Pointer(&devInfo), getDeviceProperty, DeviceDescription)
		hwid := GetDeviceInfo(unsafe.Pointer(h), unsafe.Pointer(&devInfo), getDeviceProperty, HardwareID)
		deviceNameStr := GetDeviceInfo(unsafe.Pointer(h), unsafe.Pointer(&devInfo), getDeviceProperty, DeviceName)

		if deviceNameStr == "" || description == "" || hwid == "" {
			continue
		}

		if strings.Contains(deviceNameStr, "USB-Enhanced-SERIAL CH343") || strings.Contains(hwid, "VID_1A86&PID_55D3") {
			DebugPrint("--------\n")
			DebugPrint("Name: %s\n", deviceNameStr)
			DebugPrint("Description: %s\n", description)
			DebugPrint("Hardware Info: %s\n", hwid)

			port := regexp.MustCompile(`COM\d+`).FindString(deviceNameStr) //creds to yrlu for this idea lawl
			if port != "" {
				DebugPrint("Port Name: %s\n", port)
				DebugPrint("--------\n")
				return port, nil
			}

			// Try to get port from registry if not found in name
			port, err := GetPortName(h, (*byte)(unsafe.Pointer(&devInfo)))
			if err != nil {
				DebugPrint("Failed to get port name: %v\n", err)
				return "", fmt.Errorf("Find: failed to get port name from registry: %w", err)
			}

			DebugPrint("Port Name: %s\n", port)
			DebugPrint("--------\n")

			if strings.Contains(port, "COM") {
				return port, nil
			}

			return "", nil
		}
	}

	fmt.Println("Failed to locate MAKCU!")
	return "", fmt.Errorf("Find: device not found")
}

// Sets the timeout settings for the COM port
func SetTimeouts(handle windows.Handle) error {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setCommTimeouts := kernel32.NewProc("SetCommTimeouts")

	var timeouts windows.CommTimeouts
	timeouts.ReadIntervalTimeout = 50         // Time to wait for a byte to arrive
	timeouts.ReadTotalTimeoutMultiplier = 10  // Time to wait before a read operation is finished
	timeouts.ReadTotalTimeoutConstant = 500   // Timeout in milliseconds for the entire read operation
	timeouts.WriteTotalTimeoutMultiplier = 10 // Timeout for writing
	timeouts.WriteTotalTimeoutConstant = 500  // Timeout in milliseconds for the entire write operation

	ret, _, err := setCommTimeouts.Call(uintptr(handle), uintptr(unsafe.Pointer(&timeouts)))
	if ret == 0 {
		return fmt.Errorf("SetTimeouts: SetCommTimeouts failed: %w", err)
	}

	return nil
}

// 🐱🐱🐱 Cat timeouts! 🐱🐱🐱

// 🐱 Handle for MAKCU device
type MakcuHandle struct {
	Port   string
	handle windows.Handle
	dcb    windows.DCB
}

// Make a connection to the COM port where our MAKCU was found.
func Connect(portName string, baudRate uint32) (*MakcuHandle, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	openPort := kernel32.NewProc("CreateFileW")
	setCommState := kernel32.NewProc("SetCommState")

	if !strings.HasPrefix(portName, `\\.\`) {
		portName = `\\.\` + portName
	}

	path, err := windows.UTF16PtrFromString(portName)
	if err != nil {
		return nil, fmt.Errorf("Connect: failed to convert port name to UTF16: %w", err)
	}

	handle, _, err := openPort.Call(uintptr(unsafe.Pointer(path)), syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, 0, 3, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if handle == uintptr(syscall.InvalidHandle) {
		return nil, fmt.Errorf("Connect: failed to open port: %w", err)
	}

	portHandle := windows.Handle(handle)
	// set the settings for the serial communications
	dcbOpts := &windows.DCB{}
	dcbOpts.DCBlength = uint32(unsafe.Sizeof(*dcbOpts))
	dcbOpts.BaudRate = baudRate
	dcbOpts.Flags = 0
	dcbOpts.ByteSize = 8
	dcbOpts.Parity = 0
	dcbOpts.StopBits = 1

	dcbOpts.Flags |= 0x00000400
	dcbOpts.Flags |= 0x00000800

	ret, _, err := setCommState.Call(uintptr(portHandle), uintptr(unsafe.Pointer(dcbOpts)))
	if ret == 0 {
		return nil, fmt.Errorf("Connect: failed to set communication state: %w", err)
	}

	err = SetTimeouts(portHandle)
	if err != nil {
		return nil, fmt.Errorf("Connect: failed to set timeouts: %w", err)
	}

	CleanPort := strings.TrimPrefix(portName, `\\.\`)
	DebugPrint("Successfully Connected to MAKCU! {Port %s | Baud Rate %d}\n", CleanPort, baudRate)

	return &MakcuHandle{
		Port:   CleanPort,
		handle: portHandle,
		dcb:    *dcbOpts,
	}, nil
}

// 🐱🐱🐱 Cat connect! 🐱🐱🐱

// Close the connection to the MAKCU
func (m *MakcuHandle) Close() error {
	if m == nil {
		return fmt.Errorf("Close: MakcuHandle is nil (no device connected)")
	}

	err := windows.CloseHandle(m.handle)
	if err != nil {
		return fmt.Errorf("Close: failed to close handle: %w", err)
	}

	return nil
}

// 🐱🐱🐱 Cat close! 🐱🐱🐱

// Sends the bytes needed to change the Baud Rate of the MAKCU to 4m and then returns a new Connection object with the new baud rate
// Note: This is NOT a permanent change and will reset back to the default 115200 baud rate after the MAKCU powers off and then back on again.
func ChangeBaudRate(m *MakcuHandle) (*MakcuHandle, error) {
	if m == nil {
		return nil, fmt.Errorf("ChangeBaudRate: MakcuHandle is nil (no device connected)")
	}

	n, err := m.Write([]byte{0xDE, 0xAD, 0x05, 0x00, 0xA5, 0x00, 0x09, 0x3D, 0x00})
	if err != nil {
		// Always try to close the handle on error
		_ = m.Close()
		return nil, fmt.Errorf("ChangeBaudRate: write error: %w", err)
	}

	if n != 9 {
		_ = m.Close()
		return nil, fmt.Errorf("ChangeBaudRate: wrong number of bytes written (got %d, want 9)", n)
	}

	if err := m.Close(); err != nil {
		ErrorPrint("ChangeBaudRate: failed to close old connection: %v", err)
		// Continue, but log the error
	}

	NewConn, err := Connect(m.Port, 4000000)
	if err != nil {
		return nil, fmt.Errorf("ChangeBaudRate: connect error: %w", err)
	}

	time.Sleep(1 * time.Second)

	_, err = NewConn.Write([]byte("km.version()\r"))
	if err != nil {
		_ = NewConn.Close()
		return nil, fmt.Errorf("ChangeBaudRate: write error after reconnect: %w", err)
	}

	ReadBuf := make([]byte, 32)
	n, err = NewConn.Read(ReadBuf)
	if err != nil {
		_ = NewConn.Close()
		return nil, fmt.Errorf("ChangeBaudRate: read error after reconnect: %w", err)
	}

	if !strings.Contains(string(ReadBuf[:n]), "MAKCU") {
		_ = NewConn.Close()
		return nil, fmt.Errorf("ChangeBaudRate: did not receive expected response, got: %q", string(ReadBuf[:n]))
	}

	time.Sleep(1 * time.Second)

	DebugPrint("Successfully Changed Baud Rate To %d!\n", 4000000)

	return NewConn, nil
}

// 🐱🐱🐱 Cat baud rate! 🐱🐱🐱

// Sends the given bytes to the MAKCU and returns the number of bytes written.
func (m *MakcuHandle) Write(data []byte) (int, error) {
	if m == nil {
		return -1, fmt.Errorf("Write: MakcuHandle is nil (no device connected)")
	}

	var bytesWritten uint32
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	writeFile := kernel32.NewProc("WriteFile")

	var overlapped windows.Overlapped
	overlapped.Offset = 0
	overlapped.OffsetHigh = 0

	DebugPrint("Sending %s\r\n", data[:])

	ret, _, err := writeFile.Call(uintptr(m.handle), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(unsafe.Pointer(&bytesWritten)), uintptr(unsafe.Pointer(&overlapped)))
	if ret == 0 {
		return -1, fmt.Errorf("Write: error writing to port: %w", err)
	}

	return int(bytesWritten), nil
}

// 🐱🐱🐱 Cat write! 🐱🐱🐱

// Reads data from the MAKCU and saves it to a given buffer then returns the number of bytes read.
func (m *MakcuHandle) Read(buffer []byte) (int, error) {
	if m == nil {
		return -1, fmt.Errorf("Write: MakcuHandle is nil (no device connected)")
	}

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	readFile := kernel32.NewProc("ReadFile")

	var bytesRead uint32
	ret, _, err := readFile.Call(uintptr(m.handle), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)), uintptr(unsafe.Pointer(&bytesRead)), 0)
	if ret == 0 {
		return -1, fmt.Errorf("Read: error reading from port: %w", err)
	}

	return int(bytesRead), nil
}

// 🐱🐱🐱 Cat read! 🐱🐱🐱

// 🐱 Mouse left down
func (m *MakcuHandle) LeftDown() error {
	if m == nil {
		return fmt.Errorf("LeftDown: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.left(1)\r"))
	if err != nil {
		DebugPrint("Failed to press mouse: Write Error: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat left down! 🐱🐱🐱

// 🐱 Mouse left up
func (m *MakcuHandle) LeftUp() error {
	if m == nil {
		return fmt.Errorf("LeftUp: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.left(0)\r"))
	if err != nil {
		DebugPrint("Failed to release mouse: Write Error: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat left up! 🐱🐱🐱

// 🐱 Mouse left click
func (m *MakcuHandle) LeftClick() error {
	if m == nil {
		return fmt.Errorf("LeftClick: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.left(1)\r km.left(0)\r"))
	if err != nil {
		DebugPrint("Failed to click mouse: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat left click! 🐱🐱🐱

// 🐱 Mouse right down
func (m *MakcuHandle) RightDown() error {
	if m == nil {
		return fmt.Errorf("RightDown: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.right(1)\r"))
	if err != nil {
		DebugPrint("Failed to press mouse: Write Error: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat right down! 🐱🐱🐱

// 🐱 Mouse right up
func (m *MakcuHandle) RightUp() error {
	if m == nil {
		return fmt.Errorf("RightUp: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.right(0)\r"))
	if err != nil {
		DebugPrint("Failed to release mouse: Write Error: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat right up! 🐱🐱🐱

// 🐱 Mouse right click
func (m *MakcuHandle) RightClick() error {
	if m == nil {
		return fmt.Errorf("RightClick: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.right(1)\r km.right(0)\r"))
	if err != nil {
		DebugPrint("Failed to right click mouse: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat right click! 🐱🐱🐱

// 🐱 Mouse middle down
func (m *MakcuHandle) MiddleDown() error {
	if m == nil {
		return fmt.Errorf("MiddleDown: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.middle(1)\r"))
	if err != nil {
		DebugPrint("Failed to press middle mouse button: Write Error: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat middle down! 🐱🐱🐱

// 🐱 Mouse middle up
func (m *MakcuHandle) MiddleUp() error {
	if m == nil {
		return fmt.Errorf("MiddleUp: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.middle(0)\r"))
	if err != nil {
		DebugPrint("Failed to release middle mouse button: Write Error: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat middle up 🖕! 🐱🐱🐱

// 🐱 Mouse middle click
func (m *MakcuHandle) MiddleClick() error {
	if m == nil {
		return fmt.Errorf("MiddleClick: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte("km.middle(1)\r km.middle(0)\r"))
	if err != nil {
		DebugPrint("Failed to middle click mouse: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat middle click! 🐱🐱🐱

// 🐱 Mouse button constants
const (
	MOUSE_BUTTON_LEFT   = 1
	MOUSE_BUTTON_RIGHT  = 2
	MOUSE_BUTTON_MIDDLE = 3
)

// 🐱 Clicks a mouse button
func (m *MakcuHandle) Click(i int, delay time.Duration) error {
	if m == nil {
		return fmt.Errorf("Click: MakcuHandle is nil (no device connected)")
	}

	// Basically, we create a function pointer which is just basically a variable that stores a function for us.
	// Then we can use that variable to call the function later on. :()
	type mouseAction func() error
	var down, up mouseAction

	switch i {
	case MOUSE_BUTTON_LEFT:
		down, up = m.LeftDown, m.LeftUp
	case MOUSE_BUTTON_RIGHT:
		down, up = m.RightDown, m.RightUp
	case MOUSE_BUTTON_MIDDLE:
		down, up = m.MiddleDown, m.MiddleUp
	default:
		return fmt.Errorf("invalid mouse button: %d", i)
	}

	if err := down(); err != nil {
		return err
	}

	if delay > 0 {
		time.Sleep(delay)
	}

	if err := up(); err != nil {
		return err
	}

	return nil
}

func (m *MakcuHandle) ClickMouse() error {
	if m == nil {
		return fmt.Errorf("ClickMouse: MakcuHandle is nil (no device connected)")
	}

	return m.Click(MOUSE_BUTTON_LEFT, 0)
}

// 🐱🐱🐱 Cat click! 🐱🐱🐱

// 🐱 Scrolls the mouse
func (m *MakcuHandle) ScrollMouse(amount int) error {
	if m == nil {
		return fmt.Errorf("ScrollMouse: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte(fmt.Sprintf("km.wheel(%d)\r", amount)))
	if err != nil {
		DebugPrint("Failed to scroll mouse: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat scroll! 🐱🐱🐱

// 🐱 Moves the mouse
func (m *MakcuHandle) MoveMouse(x, y int) error {
	if m == nil {
		return fmt.Errorf("MoveMouse: MakcuHandle is nil (no device connected)")
	}

	_, err := m.Write([]byte(fmt.Sprintf("km.move(%d, %d)\r", x, y)))
	if err != nil {
		DebugPrint("Failed to move mouse: Write Error: %v", err)
		return err
	}

	return nil
}

// 🐱🐱🐱 Cat move! 🐱🐱🐱

// use a curve with the built in curve functionality from MAKCU... i THINK this is only on fw v3+ ??? idk don't care to fact check it rn either :)
// "It is common sense that the higher the number of the third parameter, the smoother the curve will be fitted" - from MAKCU/km box docs
func (m *MakcuHandle) MoveMouseWithCurve(x, y int, params ...int) error {
	if m == nil {
		return fmt.Errorf("MoveMouseWithCurve: MakcuHandle is nil (no device connected)")
	}

	var cmd string
	switch len(params) {
	case 0:
		cmd = fmt.Sprintf("km.move(%d, %d)\r", x, y)
	case 1:
		cmd = fmt.Sprintf("km.move(%d, %d, %d)\r", x, y, params[0])
	case 3:
		cmd = fmt.Sprintf("km.move(%d, %d, %d, %d, %d)\r", x, y, params[0], params[1], params[2])
	default:
		DebugPrint("Invalid number of parameters")
		return fmt.Errorf("invalid number of parameters")
	}

	_, err := m.Write([]byte(cmd))
	if err != nil {
		DebugPrint("Failed to move mouse with curve: Write Error: %v", err)
		return err
	}

	return nil
}
