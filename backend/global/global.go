package global

import (
	"HostManager/backend/db"
	"gorm.io/gorm"
)

// DB 数据库实例
func DB() *gorm.DB {
	return db.GetGlobalDB()
}
