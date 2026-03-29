package app

import (
	"github.com/haochend413/scut/internal/app/context"
	"github.com/haochend413/scut/internal/app/shortcut"
	"github.com/haochend413/scut/internal/db"
	"github.com/haochend413/scut/internal/models"
)

type App struct {
	ContextMgr  *context.ContextMgr
	ShortcutMgr *shortcut.ShortcutMgr
	DB          *db.DB
}

// initialize hepler
func (a *App) init() {
	a.ContextMgr.FetchContext()
	data := a.DB.FetchAll()

	// load data from DB
	a.ShortcutMgr.RefreshFromDB(data)
	// fetch context
	a.ShortcutMgr.SetCWD(a.ContextMgr.DisplayWD())

	// right now, shortcut and wd are in shortcutMgr, while history is in cotextmgr.
}

// close helper
func (a *App) OnClose() {
	// sync data to database
	a.DB.LoadAll(a.ShortcutMgr.ExportValues())
	// return
}

// constructor
func NewApp(dbConn *db.DB) *App {
	app := &App{
		DB:          dbConn,
		ContextMgr:  &context.ContextMgr{},
		ShortcutMgr: &shortcut.ShortcutMgr{},
	}
	app.init()
	return app
}

// APIs
func (a *App) DisplayCWDShortcuts() []models.Shortcut {
	// filter to values (convert []*models.Shortcut to []models.Shortcut)
	ptrs := a.ShortcutMgr.DisplayCWDShortcuts()
	values := make([]models.Shortcut, 0, len(ptrs))
	for _, p := range ptrs {
		if p != nil {
			values = append(values, *p)
		}
	}
	return values
}

func (a *App) AddShortcut(sc models.Shortcut) {
	a.ShortcutMgr.AddShortcut(sc)
}
