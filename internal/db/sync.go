package db

import "github.com/haochend413/scut/internal/models"

// This functions sync the current data into DB and returns the synced data.
func (d *DB) LoadAndFetchAll(scs []models.Shortcut) []models.Shortcut {
	// Upsert all provided shortcuts
	for _, sc := range scs {
		d.Conn.Save(&sc)
	}
	// // Collect all IDs from input
	// ids := make([]uint, 0, len(scs))
	// for _, sc := range scs {
	// 	ids = append(ids, sc.ID)
	// }
	// // Delete shortcuts not in the input list
	// d.Conn.Where("id NOT IN ?", ids).Delete(&models.Shortcut{})

	// Return the current list of shortcuts from DB
	var synced []models.Shortcut
	d.Conn.Find(&synced)
	return synced
}

// fetch all data from db
func (d *DB) FetchAll() []models.Shortcut {
	// Return the current list of shortcuts from DB
	var synced []models.Shortcut
	d.Conn.Find(&synced)
	return synced
}

// fetch all data from db
func (d *DB) LoadAll(scs []models.Shortcut) {
	for _, sc := range scs {
		d.Conn.Save(&sc)
	}
}
