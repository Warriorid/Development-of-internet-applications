package model

import "time"

type PitsCalculation struct {
    ID          int        `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Status      string     `gorm:"column:status;type:varchar(20);not null;check:status IN ('draft','deleted','formed','completed','rejected')" json:"status"`
    CreatedAt   time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`
    FormedAt    time.Time `gorm:"column:formed_at;null" json:"formed_at"`   
    CompletedAt time.Time `gorm:"column:completed_at;null" json:"completed_at"`
    CreatorID   int        `gorm:"column:creator_id;not null" json:"creator_id"`
    ModeratorID *int       `gorm:"column:moderator_id;null" json:"moderator_id"`  
    
    PitLength   *float64 `gorm:"column:pit_length;null" json:"pit_length"`
    PitWidth    *float64 `gorm:"column:pit_width;null" json:"pit_width"`
    PitDepth    *float64 `gorm:"column:pit_depth;null" json:"pit_depth"`
}

type PitsCalculationWithMaterial struct {
    ID          int        `json:"id"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    FormedAt    time.Time `json:"formed_at"`   
    CompletedAt time.Time `json:"completed_at"`
    CreatorID   int        `json:"creator_id"`
    ModeratorID *int       `json:"moderator_id"`  
    
    PitLength   *float64 `json:"pit_length"`
    PitWidth    *float64 `json:"pit_width"`
    PitDepth    *float64 `json:"pit_depth"`

    Materials []Material

}

type UpdatePitParam struct {
    PitLength   *float64 `json:"pit_length"`
    PitWidth    *float64 `json:"pit_width"`
    PitDepth    *float64 `json:"pit_depth"`
}


func (p *PitsCalculation) ToPitsCalculationWithMaterial(materials []Material) PitsCalculationWithMaterial {
    result := PitsCalculationWithMaterial{
        ID:          p.ID,
        Status:      p.Status,
        CreatedAt:   p.CreatedAt,
        FormedAt:    p.FormedAt,
        CompletedAt: p.CompletedAt,
        CreatorID:   p.CreatorID,
        Materials:   materials,
    }
    if p.ModeratorID != nil {
        result.ModeratorID = p.ModeratorID
    }
    if p.PitLength != nil {
        result.PitLength = p.PitLength
    }
    if p.PitWidth != nil {
        result.PitWidth = p.PitWidth
    }
    if p.PitDepth != nil {
        result.PitDepth = p.PitDepth
    }
    
    return result
}

