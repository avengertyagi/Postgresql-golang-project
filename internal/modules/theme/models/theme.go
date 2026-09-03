package models

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Theme struct {
	gorm.Model
	Name        string `json:"name"            gorm:"index;type:varchar(100);uniqueIndex;not null"`
	Slug        string `json:"slug"              gorm:"index;type:varchar(100);uniqueIndex;not null"`
	Description string `json:"description"      gorm:"type:text; nullable"`
	Image       string `json:"image"            gorm:"type:varchar(191); nullable"`
	Status      bool   `json:"status"          gorm:"index;default:true"`
}

func (t *Theme) BeforeCreate(tx *gorm.DB) (err error) {
	t.Slug = slug.Make(t.Name)
	return nil
}

func (t *Theme) BeforeUpdate(tx *gorm.DB) (err error) {
	t.Slug = slug.Make(t.Name)
	return nil
}
