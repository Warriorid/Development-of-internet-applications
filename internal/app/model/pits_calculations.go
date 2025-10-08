package model

import "time"

type PitsCalculation struct {
	ID          int        `gorm:"primaryKey;column:id;autoIncrement"`
	Status      string     `gorm:"column:status;type:varchar(20);not null;check:status IN ('draft','deleted','formed','completed','rejected')"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:now()"`
	FormedAt    *time.Time `gorm:"column:formed_at"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
	CreatorID   int        `gorm:"column:creator_id;not null"`
	ModeratorID *int       `gorm:"column:moderator_id;null"`

	PitLength *float64 `gorm:"column:pit_length;null"`
	PitWidth  *float64 `gorm:"column:pit_width;null"`
	PitDepth  *float64 `gorm:"column:pit_depth;null"`

	Creator   Users                 `gorm:"foreignKey:CreatorID;references:ID;constraint:OnDelete:RESTRICT"`
	Moderator Users                 `gorm:"foreignKey:ModeratorID;references:ID;constraint:OnDelete:RESTRICT"`
	Materials []CalculationMaterial `gorm:"foreignKey:CalculationID;references:ID"`
}
