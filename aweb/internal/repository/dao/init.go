package dao

import "gorm.io/gorm"

type Product struct {
	gorm.Model
}

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{},
		&AsyncSms{},
		&Job{},
	)
}
