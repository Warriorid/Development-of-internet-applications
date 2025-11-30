package model


type CalculationMaterial struct {
    CalculationID int     `gorm:"primaryKey;column:calculation_id"`
    MaterialID    int     `gorm:"primaryKey;column:material_id"`
    SlopeAngle    *int     `gorm:"column:slope_angle;null;check:slope_angle>=0 AND slope_angle<=45"`
    VolumeResult  *float64 `gorm:"column:volume_result;null"`
    Material    Material        `gorm:"foreignKey:MaterialID;references:ID;constraint:OnDelete:RESTRICT"`
    Calculation PitsCalculation `gorm:"foreignKey:CalculationID;references:ID;constraint:OnDelete:RESTRICT"`
}


