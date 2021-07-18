package main

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/sahilm/fuzzy"
)

var (
	gmwW    = 280
	gmwH    = 200
	gmwX    = screenX - gmwW
	gmwY    = 40
	st      = ""
	ranking = Entries{}
)

type GowardMainWindow struct {
	*walk.MainWindow
	fss   *walk.LineEdit
	res   *walk.ListBox
	model *Model
}

func (gmw *GowardMainWindow) res_Update() {
	gmw.res.SetModel(InitModel(&ranking))
}

func (gmw *GowardMainWindow) res_ToClipboard(i int) {
	if i < ranking.Len() {
		fatal(walk.Clipboard().SetText(session.Decrypt(ranking[i].Login.Password)))
	}
	st = gmw.fss.Text()
	gmw.Dispose()
}

func (gmw *GowardMainWindow) fss_OnTextChange() {
	fmt.Println(gmw.Height())
	if gmw.fss.Text() == "" {
		ranking = gmw.model.items
	} else {
		results := fuzzy.FindFrom(gmw.fss.Text(), gmw.model.items)
		ranking = make(Entries, results.Len())
		for i, r := range results {
			ranking[i] = gmw.model.items[r.Index]
		}
	}
	gmw.res_Update()
}

func (gmw *GowardMainWindow) fss_OnKeyDown(key walk.Key) {
	switch key.String() {
	case "Return":
		gmw.res_ToClipboard(0)
	}
	// TODO TBD
}

func (gmw *GowardMainWindow) res_OnItemActivated() {
	gmw.res_ToClipboard(gmw.res.CurrentIndex())
}

type Uri struct {
	Uri string
}

type Login struct {
	Uris     []Uri
	Username string
	Password string
}

type Entry struct {
	Name  string
	Login Login
	Type  uint8 // int
}

type Entries []Entry

func (es Entries) String(i int) string {
	return es[i].Name
}

func (es Entries) Len() int {
	return len(es)
}

type Model struct {
	walk.ListModelBase
	items Entries
}

func (m *Model) ItemCount() int {
	return len(m.items)
}

func (m *Model) Value(index int) interface{} {
	return m.items[index].Name + " - " + m.items[index].Login.Username
}

func InitModel(entries *Entries) *Model {
	return &Model{items: *entries}
}

func MakeGMW(entries *Entries) *GowardMainWindow {
	ranking = *entries
	return &GowardMainWindow{model: InitModel(entries)}
}

func (gmw *GowardMainWindow) Start(pos uint) {
	err := declarative.MainWindow{ // change to dialog
		AssignTo: &gmw.MainWindow,
		Title:    "MW Search",
		Layout:   declarative.VBox{MarginsZero: true},
		Bounds:   declarative.Rectangle{X: gmwX, Y: gmwY, Width: gmwW, Height: gmwH},
		Children: []declarative.Widget{
			declarative.LineEdit{
				AssignTo:      &gmw.fss,
				OnTextChanged: gmw.fss_OnTextChange,
				OnKeyDown:     gmw.fss_OnKeyDown,
			},
			declarative.ListBox{
				AssignTo:        &gmw.res,
				Model:           gmw.model,
				OnItemActivated: gmw.res_OnItemActivated,
			},
		},
	}.Create()
	fatal(err)

	if pos == 0 {
		cursor := win.POINT{}
		win.GetCursorPos(&cursor)
		bounds := walk.Rectangle{X: int(cursor.X), Y: int(cursor.Y), Width: gmwW, Height: gmwH}
		if int(cursor.X)+gmwW > screenX {
			bounds.X -= gmwW
		}
		if int(cursor.Y)+gmwH > screenY {
			bounds.Y -= gmwH
		}
		gmw.SetBounds(bounds)
	}
	gmw.fss.SetText(st)
	gmw.fss.SetTextSelection(0, len(st))
	gmw.fss_OnTextChange()
	gmw.SetIcon(icon)
	setDefaultStyle(gmw.Handle())
	gmw.Run()
	gmw.Show()
}
