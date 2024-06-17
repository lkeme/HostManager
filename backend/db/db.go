package db

import (
	"HostManager/backend/model"
	"HostManager/backend/pkg/wxSQLite3"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net/url"
	"sync"
	"time"
)

var once sync.Once

var globalDB *gorm.DB

// GetSqlDB 获取数据库连接
func GetSqlDB() (sqlDB *sql.DB, err1 error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Database panic failed: ", err)
			err1 = errors.New(fmt.Sprintf("%v", err))
		}
	}()
	sqlDB, err1 = globalDB.DB()
	return
}

// GetGlobalDB 获取全局数据库连接
func GetGlobalDB() *gorm.DB {
	return globalDB
}

// GetState 获取目前数据库状态参数
func GetState(gormDB *gorm.DB) *sql.DBStats {
	db, err := gormDB.DB()
	if err != nil {
		return nil
	}
	state := db.Stats()
	return &state
}

// GetHealthStatus 检查链接是否健康
func GetHealthStatus() (err error) {
	sqlDB, err := GetSqlDB()
	if err != nil {
		log.Println("Database fetch failed: ", err.Error())
		return
	}
	// 验证与数据库的连接是否仍然存在
	err = sqlDB.Ping()
	if err != nil {
		log.Println("Database ping failed: ", err.Error())
		return
	}
	err = globalDB.Raw(`SELECT 1`).Error
	if err != nil {
		log.Println("Database select failed: ", err.Error())
		return
	}
	return
}

// Connect create a new database connection
func Connect(dbPath, key string) (err error) {
	db, err := openDatabase(getDialector(dbPath, key), getDBConfig())
	//
	if err != nil {
		log.Println("Database connection failed: ", err.Error())
		// 处理密码错误
		if err.Error() == "file is not a database" {
			err = errors.New("密码错误")
		}
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)               // 最大空闲连接数
		sqlDB.SetMaxOpenConns(100)              // 最多可容纳的连接数
		sqlDB.SetConnMaxLifetime(time.Hour * 4) // 连接最大复用时间，
		sqlDB.SetConnMaxIdleTime(time.Hour * 2) // 连接最大空闲时间
		//err = sqlDB.Ping()
		//if err != nil {
		//	log.Fatalf("数据库连接失败ping:%v", err)
		//}
		globalDB = db
	}
	return
}

// ChangePassword 修改密码
func ChangePassword(dbName, oldKey, newKey string) (err error) {
	err = wxSQLite3.ResetDBKey(dbName, oldKey, newKey)
	if err != nil {
		log.Println("Database password change failed: ", err.Error())
	}
	_ = Close()
	return
}

func Migrate() (err error) {
	err = globalDB.AutoMigrate(
		&model.Folder{},
		&model.File{},
	)
	if err != nil {
		log.Println("Database migration failed: ", err.Error())
	}
	return
}

// Close closes the database.
// It returns ErrClosed if the database is already closed.
func Close() error {
	sqlDB, err := GetSqlDB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// IsEncrypted check if the database is encrypted
func IsEncrypted(dbPath string) (encrypted bool, err error) {
	//once.Do(func() {})
	encrypted, err = wxSQLite3.IsEncrypted(dbPath)
	if err != nil {
		log.Println("Database encryption check failed: ", err)
	}
	return
}

func getDBConfig() *gorm.Config {
	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "my_",
		},
		FullSaveAssociations: false,
		PrepareStmt:          true,
	}
}

func getDialector(dbPath, key string) gorm.Dialector {
	//dbPathWithDSN := fmt.Sprintf("%s?_pragma_key=%s&_pragma_cipher_page_size=4096", dbPath, key)
	dbPathWithDSN := fmt.Sprintf("%s?_db_key=%s", dbPath, url.QueryEscape(key))
	//dbPathWithDSN := fmt.Sprintf("%s", dbPath)
	return wxSQLite3.Open(dbPathWithDSN)
}

func openDatabase(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(dialector, config)
}
