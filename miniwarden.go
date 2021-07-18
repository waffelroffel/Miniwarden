package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
)

var (
	app            = App{}
	session        = Session{}
	hotkeys        = HotKeys{}
	homeDir, hderr = os.UserHomeDir()
	confDir        = filepath.Join(homeDir, ".bwgo")
	confFile       = filepath.Join(confDir, "bwgo_conf.txt")
)

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
	app.Sync = systray.AddMenuItem("Sync", "Update vault state")
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
		session.FetchAllEntries()
		session.SaveToDisk() // make optional
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

func main() {
	fatal(hderr)
	session.FetchUserEmail()
	session.LoadFromDisk()

	defer systray.Quit()
	systray.Run(onReady, onExit)
}

func onReady() {
	app.Init()
	hotkeys.Init()
	search := hotkeys.Register(Ctrl+Shift, 'A')

	go func() {
		for {
			select {
			case <-search:
				MakeGMW(&session.entries).Start(0) // return if not logged in
			case <-app.Search.ClickedCh:
				MakeGMW(&session.entries).Start(1) // return if not logged in
			case <-app.Sync.ClickedCh:
				session.FetchAllEntries()
				session.SaveToDisk() // make optional
			case <-app.SignIn.ClickedCh:
				session.FetchAllEntries()
				if session.UserEmail != "" {
					app.SetSignedIn()    // split into login, sync
					session.SaveToDisk() // make optional
				}
			case <-app.SignOut.ClickedCh:
				fatal(cmdLogout()) // add popup for errors
				app.SetSignedOut()
				session.Clear()
			case <-app.Quit.ClickedCh:
				systray.Quit()
			}
		}
	}()
	go hotkeys.Listen()
}

func onExit() {
	session.SaveToDisk()
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
