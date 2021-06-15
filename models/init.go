package models

import (
	"blog/core"
	"fmt"
	"os"
)

var (
	// Dbms database manager system
	Dbms *DBManager
)

// init  database
func init() {
	typ := core.Conf.Db.Type
	name := core.Conf.Db.Name
	password := core.Conf.Db.Password
	address := core.Conf.Db.Address
	sqlConnection := fmt.Sprintf(`%v:%v@tcp(%v)/mysql?parseTime=true&loc=Local`, name, password, address)
	Dbms = NewDBManager(typ, sqlConnection, true)
	if err := Dbms.ReloadDbConnect(); err != nil {
		fmt.Printf("open db [%v][%v] err[%v]", "mysql", err)
		os.Exit(-1)
		return
	}
	db := Dbms.Db
	db.AutoMigrate(&Article{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Tag{})
}
