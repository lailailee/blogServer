package models

type Category struct {
	Id   int    `gorm:"primarykey;column:id"  json:"id"`
	Name string `gorm:"column:name" json:"name"`
}
