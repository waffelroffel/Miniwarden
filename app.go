package main

import "github.com/getlantern/systray"

type App struct {
	User    *systray.MenuItem
	Search  *systray.MenuItem
	Sync    *systray.MenuItem
	SignIn  *systray.MenuItem
	SignOut *systray.MenuItem
	Quit    *systray.MenuItem
}

func (app *App) Init() {
	systray.SetIcon(iconData)
	systray.SetTitle("Miniwarden")
	systray.SetTooltip("Miniwarden")

	app.User = systray.AddMenuItem(session.UserEmail, "")
	systray.AddSeparator()
	app.Search = systray.AddMenuItem("Search", "Search through vault")
	app.Sync = systray.AddMenuItem("Sync", "Sync vault state")
	systray.AddSeparator()
	app.SignIn = systray.AddMenuItem("Sign in", "Sign in to account")
	app.SignOut = systray.AddMenuItem("Sign out", "Sign out of account")
	app.Quit = systray.AddMenuItem("Quit", "Quit service")

	app.User.Disable()

	if session.UserEmail == "" {
		app.SetSignedOut()
		session.Clear()
	} else {
		app.SetSignedIn()
		session.LoadAllEntries()
	}
}

func (app *App) SetSignedIn() {
	app.User.SetTitle(session.UserEmail)
	app.Search.Show()
	app.Sync.Show()
	app.SignIn.Hide()
	app.SignOut.Show()
}

func (app *App) SetSignedOut() {
	app.User.SetTitle("Not signed in")
	app.Search.Hide()
	app.Sync.Hide()
	app.SignIn.Show()
	app.SignOut.Hide()
}
