package shortcut

import (
	"github.com/haochend413/scut/internal/models"
)

type ShortcutMgr struct {
	WorkDirectory string
	Shortcuts     []*models.Shortcut
	CWDShortcuts  []*models.Shortcut
}

func (sm *ShortcutMgr) RefreshFromDB(shortcuts []models.Shortcut) {
	for i := range shortcuts {
		sm.Shortcuts = append(sm.Shortcuts, &shortcuts[i])
	}
}

func (sm *ShortcutMgr) SetCWD(wd string) {
	sm.WorkDirectory = wd
	sm.UpdateCWDShortcuts()
}

func (sm *ShortcutMgr) UpdateCWDShortcuts() {
	sm.CWDShortcuts = []*models.Shortcut{}
	for _, sc := range sm.Shortcuts {
		if sc.WorkDirectory == sm.WorkDirectory {
			sm.CWDShortcuts = append(sm.CWDShortcuts, sc)
		}
	}
}

func (sm *ShortcutMgr) DisplayCWDShortcuts() []*models.Shortcut {
	return sm.CWDShortcuts
}

func (sm *ShortcutMgr) GetSelectedShortCut(cursor int) *models.Shortcut {
	if cursor < 0 || cursor >= len(sm.CWDShortcuts) {
		return nil
	}
	return sm.CWDShortcuts[cursor]
}

func (sm *ShortcutMgr) HasShortcut(cmd, wd string) bool {
	for _, sc := range sm.Shortcuts {
		if sc.WorkDirectory == wd && sc.Command == cmd {
			return true
		}
	}
	return false
}

func (sm *ShortcutMgr) AddShortcut(sc models.Shortcut) {
	sm.Shortcuts = append(sm.Shortcuts, &sc)
	sm.UpdateCWDShortcuts()
}

func (sm *ShortcutMgr) DeleteShortcut(id uint) {
	for index, sc := range sm.Shortcuts {
		if sc.ID == id {
			sm.Shortcuts = append(sm.Shortcuts[:index], sm.Shortcuts[index+1:]...)
			sm.UpdateCWDShortcuts()
			return
		}
	}
}

func (sm *ShortcutMgr) UpdateShortcut(id uint, newCommand string) {
	for _, sc := range sm.Shortcuts {
		if sc.ID == id {
			sc.Command = newCommand
			return
		}
	}
}
