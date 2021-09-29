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

func main() {
	fatal(hderr)
	session.FetchUserEmail()
	warning(session.LoadSessionKey())

	systray.Run(onReady, nil)
}

func onReady() {
	app.Init()
	hotkeys.Init()
	search := hotkeys.Register(Ctrl+Shift, 'A')

	go func() {
		for {
			select {
			case <-search:
				MakeGMW(&session.entries).Start(0)
			case <-app.Search.ClickedCh:
				MakeGMW(&session.entries).Start(1)
			case <-app.Sync.ClickedCh:
				session.FetchAllEntries()
			case <-app.SignIn.ClickedCh:
				session.FetchAllEntries()
				if session.UserEmail != "" {
					app.SetSignedIn()
				}
			case <-app.SignOut.ClickedCh:
				fatal(cmdLogout())
				app.SetSignedOut()
				session.Clear()
			case <-app.Quit.ClickedCh:
				systray.Quit()
			}
		}
	}()
	go hotkeys.Listen()
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func warning(err error) {
	if err != nil {
		log.Println(err)
	}
}
