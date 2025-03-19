package dao

import "gorm.io/gorm"

type Product struct {
	gorm.Model
}

func InitTables(db *gorm.DB) error {
	err1 := db.AutoMigrate(&User{})
	if err1 != nil {
		return err1
	}
	err2 := db.AutoMigrate(&Product{})
	if err2 != nil {
		return err2
	}
	return nil
}
