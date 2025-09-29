package model


type CalculationMaterial struct {
    CalculationID int     `gorm:"primaryKey;column:calculation_id"`
    MaterialID    int     `gorm:"primaryKey;column:material_id"`
    SlopeAngle    int     `gorm:"null;column:slope_angle"`
    VolumeResult  float64 `gorm:"null;column:volume_result"`
    Material    Material        `gorm:"foreignKey:MaterialID;references:ID"`
    Calculation PitsCalculation `gorm:"foreignKey:CalculationID;references:ID"`
}


