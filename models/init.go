package models

import (
	ogorm "end/plugins/gorm"
	"fmt"
	"os"
)

// DBMS   database manager object
type DBMS struct {
	Gcfg *ogorm.Manager
}

var (
	// Dbms database manager system
	Dbms *DBMS
)

// InitDB init database
func InitDB() {
	dbmCfg := ogorm.NewManager("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql?parseTime=true", true)
	if err := dbmCfg.ReloadDbConnect(); err != nil {
		fmt.Printf("open db [%v][%v] err[%v]", "mysql", err)
		os.Exit(-1)
		return
	}

	Dbms = &DBMS{
		Gcfg: dbmCfg,
	}

	db := Dbms.Gcfg.GetDB()

	db.AutoMigrate(&Article{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Tag{})
}
