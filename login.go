package main

import (
	"strings"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

var (
	glwW = 270
	glwH = 170
	glwX = screenX - glwW
	glwY = 40
)

type GowardLoginWindow struct {
	*walk.MainWindow
	un *walk.LineEdit
	pw *walk.LineEdit
	si *walk.PushButton
	o  *walk.Label
}

func (glw *GowardLoginWindow) un_OnKeyDown(key walk.Key) {
	if key.String() == "Return" {
		glw.pw.SetFocus()
	}
}

func (glw *GowardLoginWindow) pw_OnKeyDown(key walk.Key) {
	if key.String() == "Return" {
		glw.Login()
	}
}

func (glw *GowardLoginWindow) Login() {
	glw.si.SetEnabled(false)
	out, err := cmdLogin(glw.un.Text(), glw.pw.Text())
	if err != nil {
		if strings.Index(out.String(), "You are already logged in as") == 0 {
			glw.Unlock()
		} else {
			glw.o.SetText(out.String())
			glw.si.SetEnabled(true)
		}
		return
	}
	session.UserEmail = glw.un.Text()
	arr := strings.Split(out.String(), " ")
	session.Key = arr[len(arr)-1]
	glw.Dispose()
}

func (glw *GowardLoginWindow) Unlock() {
	out, err := cmdUnlock(glw.pw.Text())
	if err != nil {
		glw.o.SetText(out.String())
		glw.si.SetEnabled(true)
		return
	}
	session.UserEmail = glw.un.Text()
	arr := strings.Split(out.String(), " ")
	session.Key = arr[len(arr)-1]
	glw.Dispose() // change to show-hide
}

func (glw *GowardLoginWindow) Start() {
	err := declarative.MainWindow{
		AssignTo: &glw.MainWindow,
		Title:    "MW Login",
		Layout:   declarative.VBox{},
		Bounds:   declarative.Rectangle{X: glwX, Y: glwY, Width: glwW, Height: glwH},
		Children: []declarative.Widget{
			declarative.Label{Text: "Email"},
			declarative.LineEdit{
				AssignTo:  &glw.un,
				OnKeyDown: glw.un_OnKeyDown,
			},
			declarative.Label{Text: "Password"},
			declarative.LineEdit{
				AssignTo:     &glw.pw,
				PasswordMode: true,
				OnKeyDown:    glw.pw_OnKeyDown,
			},
			declarative.PushButton{AssignTo: &glw.si, Text: "Sign in", OnClicked: glw.Login},
			declarative.Label{AssignTo: &glw.o},
		},
	}.Create()
	fatal(err)

	if session.UserEmail != "" {
		glw.un.SetText(session.UserEmail)
		glw.un.SetEnabled(false)
		glw.pw.SetFocus()
	}
	glw.SetIcon(icon)
	setDefaultStyle(glw.Handle())
	glw.Run()
	glw.Show()
}
