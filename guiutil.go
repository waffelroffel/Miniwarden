package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

var (
	screenX = int(win.GetSystemMetrics(win.SM_CXSCREEN))
	screenY = int(win.GetSystemMetrics(win.SM_CYSCREEN))
	icon    = getIconFromResource()
)

func getIconFromResource() *walk.Icon {
	icon, _ := walk.NewIconFromResourceId(2)
	/*
		if err != nil {
			// show err
		}
	*/
	return icon
}

func setDefaultStyle(hwnd win.HWND) {
	win.SetWindowLong(hwnd, win.GWL_STYLE, win.GetWindowLong(hwnd, win.GWL_STYLE) & ^win.WS_MINIMIZEBOX & ^win.WS_MAXIMIZEBOX & ^win.WS_SIZEBOX)
}
