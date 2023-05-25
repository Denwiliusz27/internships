package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"))

	if err != nil {
		panic("Failed to connect with database")
	}

	err = db.AutoMigrate(&RecipeDB{}, &PresentIngredient{}, &MissingIngredient{})
	if err != nil {
		return nil
	}

	return db
}
