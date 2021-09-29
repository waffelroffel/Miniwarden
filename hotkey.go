package main

import (
	"runtime"

	"github.com/lxn/win"
)

const (
	Alt uintptr = 1 << iota
	Ctrl
	Shift
	Win
)

type HotKey struct {
	Id        uintptr
	Modifiers uintptr
	KeyCode   uintptr
	channel   chan bool
}

type HotKeys struct {
	hotkeys    map[uintptr]HotKey
	numHotkeys uintptr
}

func (hkm *HotKeys) Init() {
	hkm.hotkeys = make(map[uintptr]HotKey)
}

func (hkm *HotKeys) Register(mods uintptr, kc uintptr) chan bool {
	hkm.numHotkeys++
	err := registerHotKey(hkm.numHotkeys, mods, kc)
	fatal(err)

	hkm.hotkeys[hkm.numHotkeys] = HotKey{
		Id:        hkm.numHotkeys,
		Modifiers: mods,
		KeyCode:   kc,
		channel:   make(chan bool),
	}
	return hkm.hotkeys[hkm.numHotkeys].channel
}

func (hkm *HotKeys) Listen() {
	runtime.LockOSThread()
	msg := &win.MSG{}
	for {
		win.GetMessage(msg, 0, 0, 0)
		hkm.hotkeys[msg.WParam].channel <- true
	}
}
