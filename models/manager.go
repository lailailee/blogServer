package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBManager db manager object
type DBManager struct {
	Db            *gorm.DB
	sqlConnection string
	dialect       string
	dbg           bool
}

// NewDBManager new a db manager object
func NewDBManager(dialect, sqlConnection string, dbg bool) *DBManager {
	v := DBManager{
		Db:            nil,
		dialect:       dialect,
		dbg:           dbg,
		sqlConnection: sqlConnection,
	}
	return &v
}

// ReloadDbConnect re connect database server
func (v *DBManager) ReloadDbConnect() error {
	var err error
	if v.Db != nil {
		v.Db = nil
	}
	if v.dialect == "sqlite" {
		v.Db, err = gorm.Open(
			sqlite.Open(v.sqlConnection), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
	} else if v.dialect == "mysql" {
		v.Db, err = gorm.Open(mysql.Open(v.sqlConnection), &gorm.Config{})
	}
	if err != nil {
		return err
	}
	return nil
}
