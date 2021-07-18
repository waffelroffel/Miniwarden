package main

import (
	"runtime"
	"syscall"

	"github.com/lxn/win"
)

const (
	Alt uintptr = 1 << iota
	Ctrl
	Shift
	Win
)

var (
	user32         = syscall.NewLazyDLL("user32")
	rhk            = user32.NewProc("RegisterHotKey")
	hkid   uintptr = 0
)

type HotKey struct {
	Id        uintptr
	Modifiers uintptr
	KeyCode   uintptr
	channel   chan bool
}

type HotKeys struct {
	hotkeys map[uintptr]HotKey
}

func (hkm *HotKeys) Init() {
	hkm.hotkeys = make(map[uintptr]HotKey)
}

func (hkm *HotKeys) Register(mods uintptr, kc uintptr) chan bool {
	hkid++
	rhk.Call(0, hkid, mods, kc)
	hkm.hotkeys[hkid] = HotKey{
		Id:        hkid,
		Modifiers: mods,
		KeyCode:   kc,
		channel:   make(chan bool),
	}
	return hkm.hotkeys[hkid].channel
}

func (hkm *HotKeys) Listen() {
	runtime.LockOSThread()
	for {
		msg := &win.MSG{}
		win.GetMessage(msg, 0, 0, 0)
		hkm.hotkeys[msg.WParam].channel <- true
	}
}
