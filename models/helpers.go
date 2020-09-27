package models

import (
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/postgres" // Not directly used, but needed to help gorm communicate with postgres
)

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return errNotFound
	}
	return err
}

func last(db *gorm.DB, dst interface{}) error {
	err := db.Last(dst).Error
	if err == gorm.ErrRecordNotFound {
		return errNotFound
	}
	return err
}