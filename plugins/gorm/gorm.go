package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// M short of map[string]interface{}
type M map[string]interface{}

// Manager db manager obeject
type Manager struct {
	Db            *gorm.DB
	sqlConnection string
	dialect       string
	dbg           bool
}

// NewManager new a db manager object
func NewManager(dialect, sqlConnection string, dbg bool) *Manager {
	v := Manager{
		Db:            nil,
		dialect:       dialect,
		dbg:           dbg,
		sqlConnection: sqlConnection,
	}

	return &v
}

// ReloadDbConnect re connect database server
func (v *Manager) ReloadDbConnect() error {
	var err error

	if v.Db != nil {
		// v.db.Close()
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

	// if v.dbg {
	// 	v.db.Debug()
	// }
	return nil
}

// GetDB get db object
func (v *Manager) GetDB() *gorm.DB {
	if v.Db == nil {
		v.ReloadDbConnect()
	}

	return v.Db
}

// GetOne get one row data
func (v *Manager) GetOne(out interface{}, condDict M) (e error) {
	// fmt.Printf("%+v, %T\n", cond_dict, condDict)
	e = v.Db.Where(map[string]interface{}(condDict)).First(out).Error
	return
}

// GetFirst get fist row data
func (v *Manager) GetFirst(out interface{}) (e error) {
	e = v.Db.First(out).Error
	return
}

// GetFirst get fist row data
// func (v *Manager) GetFirst(out interface{}) (e error) {
// 	e = v.db.Model(out).Preload("Category").Preload("Tag").First(out).Error
// 	return
// }

// GetLast get fist row data
func (v *Manager) GetLast(out interface{}) (e error) {
	e = v.Db.Last(out).Error
	return
}

// FindByPage get fist row data
func (v *Manager) FindByPage(limit int, offset int, out interface{}) (e error) {
	e = v.Db.Limit(limit).Offset(offset).Find(out).Error
	return
}

// GetSomeByPage get some row data by condtion and limit , offset
func (v *Manager) GetSomeByPage(limit int, offset int, out interface{}, condDict M) (e error) {
	e = v.Db.Where(map[string]interface{}(condDict)).Order("created_at desc").Limit(limit).Offset(offset).Find(out).Error
	return
}

// GetSome get some row data by condtion
func (v *Manager) GetSome(out interface{}, condDict M) (e error) {
	e = v.Db.Where(map[string]interface{}(condDict)).Find(out).Error
	return
}

// GetAll get all row data
func (v *Manager) GetAll(out interface{}) (e error) {
	e = v.Db.Find(out).Error
	return
}

// Delete delete all row data by condtion
func (v *Manager) Delete(condDict M, modelObj interface{}) (e error) {
	e = v.Db.Where(map[string]interface{}(condDict)).Delete(modelObj).Error
	return
}

// Update update some row by condtion
func (v *Manager) Update(condDict M, updateInfo M, modelObj interface{}) (e error) {
	// e = v.db.Model(modelObj).Where(map[string]interface{}(condDict)).Updates(updateInfo).Error
	e = v.Db.Model(modelObj).Where(map[string]interface{}(condDict)).Updates(map[string]interface{}(updateInfo)).Error
	return
}

func (v *Manager) UpdateStruct(condDict M, updateInfo interface{}, modelObj interface{}) (e error) {
	// e = v.db.Model(modelObj).Where(map[string]interface{}(condDict)).Updates(updateInfo).Error
	e = v.Db.Model(modelObj).Where(map[string]interface{}(condDict)).Updates(updateInfo).Error
	return
}

// Save update all fileds
func (v *Manager) Save(modelObj interface{}) (e error) {
	e = v.Db.Save(modelObj).Error
	return
}

// Create update all fileds
func (v *Manager) Create(modelObj interface{}) (e error) {
	e = v.Db.Create(modelObj).Error
	return
}

// UpdateAll update all all row
func (v *Manager) UpdateAll(updateInfo M, modelObj interface{}) (e error) {
	e = v.Db.Model(modelObj).Updates(map[string]interface{}(updateInfo)).Error
	return
}

// Add add one row data
func (v *Manager) Add(obj interface{}) (e error) {
	e = v.Db.Create(obj).Error
	return
}

// GetCount query count
func (v *Manager) GetCount(modelObj interface{}, count *int64) (e error) {
	e = v.Db.Model(modelObj).Count(count).Error
	return
}
