package database

import (
	"engine/app/model"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		model.Detail{},
		model.Rule{},
		model.Score{},
	)
}
