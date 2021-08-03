package main

import (
	"syscall"
	"unsafe"
)

var (
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	getWindowText       = user32.NewProc("GetWindowTextW")
	getWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	// getClassName        = user32.NewProc("GetClassNameW")
)

func GetWindowTextLength(hwnd uintptr) int {
	ret, _, _ := getWindowTextLength.Call(hwnd)
	return int(ret)
}

func GetWindowText(hwnd uintptr) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	getWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

/*
func GetClassName(hwnd uintptr) string {
	n := make([]uint16, 256)
	p := &n[0]
	r0, _, _ := syscall.Syscall(getClassName.Addr(), 3, hwnd, uintptr(unsafe.Pointer(p)), uintptr(len(n)))
	if r0 == 0 {
		return ""
	}
	return syscall.UTF16ToString(n)
}
*/

func GetActiveWindow() string {
	hwnd, _, err := getForegroundWindow.Call()
	if hwnd == 0 {
		fatal(err)
	}
	return GetWindowText(hwnd)
}
