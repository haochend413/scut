package app

import (
	"fmt"
	"strings"

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

func (a *App) init() {
	a.ContextMgr.FetchContext()
	data := a.DB.FetchAll()
	a.ShortcutMgr.RefreshFromDB(data)
	a.ShortcutMgr.SetCWD(a.ContextMgr.DisplayWD())
}

func (a *App) OnClose() {}

func NewApp(dbConn *db.DB) *App {
	app := &App{
		DB:          dbConn,
		ContextMgr:  &context.ContextMgr{},
		ShortcutMgr: &shortcut.ShortcutMgr{},
	}
	app.init()
	return app
}

func (a *App) GetCWD() string {
	return a.ContextMgr.DisplayWD()
}

// GetHistory returns recent shell commands, most-recent-first, deduped and with scut calls filtered out.
func (a *App) GetHistory() []string {
	all := a.ContextMgr.CurrentCtx.History
	seen := make(map[string]bool)
	result := make([]string, 0, len(all))
	for i := len(all) - 1; i >= 0; i-- {
		cmd := all[i]
		if strings.HasPrefix(cmd, "scut") || seen[cmd] {
			continue
		}
		seen[cmd] = true
		result = append(result, cmd)
	}
	return result
}

func (a *App) DisplayCWDShortcuts() []models.Shortcut {
	ptrs := a.ShortcutMgr.DisplayCWDShortcuts()
	values := make([]models.Shortcut, 0, len(ptrs))
	for _, p := range ptrs {
		if p != nil {
			values = append(values, *p)
		}
	}
	return values
}

func (a *App) AddShortcut(sc models.Shortcut) error {
	if a.ShortcutMgr.HasShortcut(sc.Command, sc.WorkDirectory) {
		return fmt.Errorf("shortcut already exists")
	}
	inserted, err := a.DB.InsertShortcut(sc)
	if err != nil {
		return err
	}
	a.ShortcutMgr.AddShortcut(inserted)
	return nil
}

func (a *App) DeleteShortcut(id uint) {
	a.ShortcutMgr.DeleteShortcut(id)
	_ = a.DB.DeleteShortcut(id)
}

func (a *App) UpdateShortcut(id uint, newCommand string) error {
	if err := a.DB.UpdateShortcut(id, newCommand); err != nil {
		return err
	}
	a.ShortcutMgr.UpdateShortcut(id, newCommand)
	return nil
}

func (a *App) DuplicateShortcut(sc models.Shortcut) (models.Shortcut, error) {
	inserted, err := a.DB.InsertShortcut(sc)
	if err != nil {
		return models.Shortcut{}, err
	}
	a.ShortcutMgr.AddShortcut(inserted)
	return inserted, nil
}
