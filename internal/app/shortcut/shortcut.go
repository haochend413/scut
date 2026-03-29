package shortcut

import (
	"github.com/haochend413/scut/internal/models"
)

type ShortcutMgr struct {
	WorkDirectory string
	Shortcuts     []*models.Shortcut // This should be all the shortcuts. Since this is easier.
	CWDShortcuts  []*models.Shortcut
}

// load stored from db
func (sm *ShortcutMgr) RefreshFromDB(shortcuts []models.Shortcut) {
	// fetch from db.
	for i := range shortcuts {
		sm.Shortcuts = append(sm.Shortcuts, &shortcuts[i])
	}
}

// update workdir
func (sm *ShortcutMgr) SetCWD(wd string) {
	// fetch from db.
	sm.WorkDirectory = wd
	sm.UpdateCWDShortcuts()
}

func (sm *ShortcutMgr) UpdateCWDShortcuts() {
	sm.CWDShortcuts = []*models.Shortcut{}
	// loop through and get

	for _, sc := range sm.Shortcuts {
		if sc.WorkDirectory == sm.WorkDirectory {
			sm.CWDShortcuts = append(sm.CWDShortcuts, sc)
		}
	}
}

// it might be better to just use value here. maybe.
func (sm *ShortcutMgr) DisplayCWDShortcuts() []*models.Shortcut {
	return sm.CWDShortcuts
}

func (sm *ShortcutMgr) GetSelectedShortCut(cursor int) *models.Shortcut {
	if cursor < 0 || cursor >= len(sm.CWDShortcuts) {
		return nil
	}
	return sm.CWDShortcuts[cursor]
}

// add a shortcut
func (sm *ShortcutMgr) AddShortcut(sc models.Shortcut) {
	// add shortcut
	sm.Shortcuts = append(sm.Shortcuts, &sc)
	sm.UpdateCWDShortcuts()
}

// delete a shortcut
func (sm *ShortcutMgr) DeleteShortcut(id uint) {
	// remove shortcut
	for index, sc := range sm.Shortcuts {
		if sc.ID == id {
			//remove sc
			sm.Shortcuts = append(sm.Shortcuts[:index], sm.Shortcuts[index+1:]...)
			return
		}
	}
}

// for sync
func (sm *ShortcutMgr) ExportValues() []models.Shortcut {
	shortcut_val := []models.Shortcut{}
	for _, sc := range sm.Shortcuts {
		if sc != nil {
			shortcut_val = append(shortcut_val, *sc)
		}
	}
	return shortcut_val
}
