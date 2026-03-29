package models

import "gorm.io/gorm"

// maybe we can wait a bit to give a thought on this.
type Shortcut struct {
	gorm.Model    // This contains ID.
	WorkDirectory string
	Command       string // a list of commands to execute
}
