package model


type Material struct {
    ID          int     `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Title       string  `gorm:"column:title;type:varchar(100);not null;unique" json:"title"`
    Coefficient float64 `gorm:"column:coefficient;not null;check:coefficient>0" json:"coefficient"`
    ImageURL    string  `gorm:"column:image_url;type:varchar(255);null" json:"image_url"`
    Description string  `gorm:"column:description;type:text;not null" json:"description"`
    IsDeleted   bool    `gorm:"column:is_deleted;not null;default:false" json:"is_deleted"`
}