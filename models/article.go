package models

type Article struct {
	Id         int      `gorm:"primarykey;column:id"  json:"id"`
	Title      string   `gorm:"column:title" json:"title"`
	Content    string   `gorm:"column:content" json:"content"`
	Overview   string   `gorm:"column:overview" json:"overview"`
	CreateTime string   `gorm:"column:create" json:"create"`
	CategoryId int      `json:"categoryid"`
	Category   Category `gorm:"foreignKey:CategoryId"  json:"category"`
	Tags       []*Tag   `gorm:"many2many:article_tags;ForeignKey:id" json:"tags"`
}
