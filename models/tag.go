package models

type Tag struct {
	Id       int       `gorm:"primarykey" json:"id"`
	Name     string    `json:"name"`
	Articles []Article `gorm:"many2many:article_tags;" json:"articles"`
}
