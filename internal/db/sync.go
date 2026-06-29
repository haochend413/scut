package db

import "github.com/haochend413/scut/internal/models"

func (d *DB) InsertShortcut(sc models.Shortcut) (models.Shortcut, error) {
	result := d.Conn.Create(&sc)
	return sc, result.Error
}

func (d *DB) FetchAll() []models.Shortcut {
	var shortcuts []models.Shortcut
	d.Conn.Find(&shortcuts)
	return shortcuts
}

func (d *DB) DeleteShortcut(id uint) error {
	return d.Conn.Delete(&models.Shortcut{}, id).Error
}

func (d *DB) UpdateShortcut(id uint, newCommand string) error {
	return d.Conn.Model(&models.Shortcut{}).Where("id = ?", id).Update("command", newCommand).Error
}
