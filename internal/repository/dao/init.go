package dao

import "gorm.io/gorm"

func InitUserTables(db *gorm.DB) error {
	// 严格来说，这个不是优秀实践
	// 但是为了简化代码，我们在这里使用 AutoMigrate
	return db.AutoMigrate(&User{})
}
