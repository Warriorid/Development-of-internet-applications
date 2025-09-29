package model

import "time"

type PitsCalculation struct {
    ID          int        `gorm:"primaryKey;column:id"`
    Status      string     `gorm:"column:status"`
    CreatedAt   time.Time  `gorm:"column:created_at"`
    FormedAt    *time.Time `gorm:"column:formed_at"`   
    CompletedAt *time.Time `gorm:"column:completed_at"`
    CreatorID   int        `gorm:"column:creator_id"`
    ModeratorID *int       `gorm:"column:moderator_id"`  
    
    PitLength   float64 `gorm:"null;column:pit_length"`
    PitWidth    float64 `gorm:"null;column:pit_width"`
    PitDepth    float64 `gorm:"null;column:pit_depth"`
 
    Creator     Users                   `gorm:"foreignKey:CreatorID"`
    Moderator   Users                   `gorm:"foreignKey:ModeratorID"`
    Materials   []CalculationMaterial   `gorm:"foreignKey:CalculationID"`
}
