package main

import (
	"strings"

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
	ranking = []Entry{}
)

type GowardMainWindow struct {
	*walk.MainWindow
	fss   *walk.LineEdit
	res   *walk.ListBox
	model *Model
}

type Model struct {
	walk.ListModelBase
	items []Entry
}

func (m *Model) String(i int) string {
	return m.items[i].Name
}

func (m *Model) Value(i int) interface{} {
	return m.String(i) + " - " + m.items[i].Login.Username
}

func (m *Model) Len() int {
	return m.ItemCount()
}

func (m *Model) ItemCount() int {
	return len(m.items)
}

func (gmw *GowardMainWindow) findBestMatchingEntry() {
	winText, err := GetActiveWindow()
	if err != nil {
		warning(err)
		return
	}

	tokens := strings.Split(winText, " ")
	best := ""
	bestscore := 0
	for _, t := range tokens {
		if len(t) <= 1 {
			continue
		}
		res := fuzzy.FindFrom(t, gmw.model)
		if res.Len() == 0 {
			continue
		}
		if res[0].Score > bestscore {
			best = res[0].Str
			bestscore = res[0].Score
		}
	}
	gmw.fss.SetText(best)
	gmw.fss.SetTextSelection(0, len(best))
}

func (gmw *GowardMainWindow) res_Update() {
	fatal(gmw.res.SetModel(initModel(&ranking)))
}

func (gmw *GowardMainWindow) res_AutoType(i int) {
	gmw.Dispose()
	go AutoType(ranking[i].Login.Username, ranking[i].Login.Password)
}

func (gmw *GowardMainWindow) fss_OnTextChange() {
	if gmw.fss.Text() == "" {
		ranking = gmw.model.items
	} else {
		results := fuzzy.FindFrom(gmw.fss.Text(), gmw.model)
		ranking = make([]Entry, results.Len())
		for i, r := range results {
			ranking[i] = gmw.model.items[r.Index]
		}
	}
	gmw.res_Update()
}

func (gmw *GowardMainWindow) fss_OnKeyDown(key walk.Key) {
	switch key.String() {
	case "Return":
		gmw.res_AutoType(0)
	}
}

func (gmw *GowardMainWindow) res_OnItemActivated() {
	gmw.res_AutoType(gmw.res.CurrentIndex())
}

func initModel(entries *[]Entry) *Model {
	return &Model{items: *entries}
}

func MakeGMW(entries *[]Entry) *GowardMainWindow {
	ranking = *entries
	return &GowardMainWindow{model: initModel(entries)}
}

func (gmw *GowardMainWindow) Start(pos int) {
	err := declarative.MainWindow{
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
	gmw.findBestMatchingEntry()
	gmw.SetIcon(icon)
	setDefaultStyle(gmw.Handle())
	gmw.Run()
	gmw.Show()
}
