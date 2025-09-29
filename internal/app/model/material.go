package model

type Material struct {
    ID          int     `gorm:"primaryKey;column:id"`
    Title       string  `gorm:"column:title"`
    Coefficient float64 `gorm:"column:coefficient"`
    ImageURL    string  `gorm:"column:image_url"`
    Description string  `gorm:"column:description"`
    IsDeleted   bool    `gorm:"column:is_deleted"`
}

