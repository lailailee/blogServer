package models

import "time"

type Article struct {
	Id          int       `gorm:"primarykey;column:id"  json:"id"`
	Title       string    `gorm:"column:title;UNIQUE" json:"title"`
	Content     string    `gorm:"column:content" json:"content"`
	Overview    string    `gorm:"column:overview" json:"overview"`
	CategoryId  int       `gorm:"column:categoryId" json:"categoryId"`
	SeriesId    int       `gorm:"column:seriesId" json:"seriesId"`
	ViewCount   int       `gorm:"column:viewCount" json:"viewCount"`
	SeriesIndex int       `gorm:"column:seriesIndex" json:"seriesIndex"`
	Series      Series    `gorm:"foreignKey:SeriesId"  json:"series"`
	Category    Category  `gorm:"foreignKey:CategoryId"  json:"category"`
	Tags        []Tag     `gorm:"many2many:article_tags;ForeignKey:id" json:"tags"`
	CreatedAt   time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

type Category struct {
	Id        int       `gorm:"primarykey;column:id"  json:"id"`
	Name      string    `gorm:"column:name;UNIQUE" json:"name"`
	Articles  []Article `gorm:"foreignKey:CategoryId;" json:"articles"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

type Series struct {
	Id        int       `gorm:"primarykey;column:id"  json:"id"`
	Name      string    `gorm:"column:name;UNIQUE" json:"name"`
	Articles  []Article `gorm:"foreignKey:SeriesId;" json:"articles"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

type Tag struct {
	Id        int       `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"column:name;UNIQUE" json:"name"`
	Articles  []Article `gorm:"many2many:article_tags;" json:"articles"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

type User struct {
	Id       int    `gorm:"primarykey;column:id"  json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Password string `gorm:"column:password" json:"password"`
}
