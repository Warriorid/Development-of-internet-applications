// internal/model/pits_calculations.go
package model

import "time"

type PitsCalculation struct {
    ID          int        `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Status      string     `gorm:"column:status;type:varchar(20);not null;check:status IN ('draft','deleted','formed','completed','rejected')" json:"status"`
    CreatedAt   time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`
    FormedAt    *time.Time `gorm:"column:formed_at;null" json:"formed_at"`   
    CompletedAt *time.Time `gorm:"column:completed_at;null" json:"completed_at"`
    CreatorID   int        `gorm:"column:creator_id;not null" json:"creator_id"`
    ModeratorID *int       `gorm:"column:moderator_id;null" json:"moderator_id"`  
    
    PitLength   *float64 `gorm:"column:pit_length;null" json:"pit_length"`
    PitWidth    *float64 `gorm:"column:pit_width;null" json:"pit_width"`
    PitDepth    *float64 `gorm:"column:pit_depth;null" json:"pit_depth"`
}

type MaterialWithCalculationData struct {
    Material
    SlopeAngle   *int     `json:"slope_angle"`
    VolumeResult *float64 `json:"volume_result"`
}

type PitsCalculationWithMaterials struct {
    ID          int           `json:"id"`
    Status      string        `json:"status"`
    CreatedAt   time.Time     `json:"created_at"`
    FormedAt    *time.Time    `json:"formed_at"`   
    CompletedAt *time.Time    `json:"completed_at"`
    CreatorID   int           `json:"creator_id"`
    ModeratorID *int          `json:"moderator_id"`  
    
    PitLength   *float64 `json:"pit_length"`
    PitWidth    *float64 `json:"pit_width"`
    PitDepth    *float64 `json:"pit_depth"`

    Materials []MaterialWithCalculationData `json:"materials"`
}

type UpdatePitParam struct {
    PitLength   *float64 `json:"pit_length"`
    PitWidth    *float64 `json:"pit_width"`
    PitDepth    *float64 `json:"pit_depth"`
}

func (p *PitsCalculation) ToPitsCalculationWithMaterials(materials []MaterialWithCalculationData) PitsCalculationWithMaterials {
    result := PitsCalculationWithMaterials{
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

// internal/model/pits_calculations.go

type PitsCalculationListItem struct {
    ID          int        `json:"id"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    FormedAt    *time.Time `json:"formed_at"`   
    CompletedAt *time.Time `json:"completed_at"`
    CreatorID   int        `json:"creator_id"`
    ModeratorID *int       `json:"moderator_id"`  
    
    PitLength   *float64 `json:"pit_length"`
    PitWidth    *float64 `json:"pit_width"`
    PitDepth    *float64 `json:"pit_depth"`
    PitVolume   *float64 `json:"pit_volume"`

    CalculatedMaterialsCount int `json:"calculated_materials_count"`
}

func (p *PitsCalculation) ToPitsCalculationListItem(calculatedCount int, pitVolume *float64) PitsCalculationListItem {
    result := PitsCalculationListItem{
        ID:          p.ID,
        Status:      p.Status,
        CreatedAt:   p.CreatedAt,
        FormedAt:    p.FormedAt,
        CompletedAt: p.CompletedAt,
        CreatorID:   p.CreatorID,
        CalculatedMaterialsCount: calculatedCount,
        PitVolume:   pitVolume,
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

type AsyncCalculationResult struct {
    CalculationID int      `json:"calculation_id"`
    MaterialID    int      `json:"material_id"`
    VolumeResult  float64  `json:"volume_result"`
    Token         string   `json:"token"`
}

type AsyncCalculationRequest struct {
    CalculationID  int      `json:"calculation_id"`
    MaterialID     int      `json:"material_id"`
    PitLength      float64  `json:"pit_length"`
    PitWidth       float64  `json:"pit_width"`
    PitDepth       float64  `json:"pit_depth"`
    SlopeAngle     int      `json:"slope_angle"`
    Coefficient    float64  `json:"coefficient"`
}

