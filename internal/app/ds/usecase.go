package ds

import (
	"gorm.io/gorm"
)

// type UseCase struct {
// 	gorm.Model
// 	ID          uint   `gorm:primarykey`
// 	name        string `gorm:notNull`
// 	url         string `gorm:notNull`
// 	description string `gorm:default:'-'`
// 	duration    uint   `gorm:default:'0'`
// 	IsDelete    bool   `grom:default:'0'`
// }

type UseCase struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(255)"`
	URL         string `gorm:"type:varchar(500)"`
	Description string `gorm:"type:text;default:'-'"`
	Consumption uint   `gorm:"default:0"`
	IsDelete    bool   `gorm:"default:false"`
}
