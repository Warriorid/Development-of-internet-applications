package repository

import (
	"DIA/internal/model"
	"fmt"

	"gorm.io/gorm"
)


type CalculationMaterialPostgres struct {
	db *gorm.DB
}

func NewCalculationMaterialPostgres(db *gorm.DB) *CalculationMaterialPostgres{
	return &CalculationMaterialPostgres{db: db}
}


func (r *CalculationMaterialPostgres) RemoveMaterialFromPit(calculationID, materialID int) error {
	var count int64
	err := r.db.Model(&model.CalculationMaterial{}).
		Where("calculation_id = ? AND material_id = ?", calculationID, materialID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("material not found in pit")
	}
	result := r.db.Where("calculation_id = ? AND material_id = ?", calculationID, materialID).
		Delete(&model.CalculationMaterial{})
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to remove material from pit")
	}
	return nil
}


func (r *CalculationMaterialPostgres) UpdateCalculationMaterial(calculationID, materialID int, slopeAngle int) error {
	var count int64
	err := r.db.Model(&model.CalculationMaterial{}).
		Where("calculation_id = ? AND material_id = ?", calculationID, materialID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("material not found in pit")
	}
	result := r.db.Model(&model.CalculationMaterial{}).
		Where("calculation_id = ? AND material_id = ?", calculationID, materialID).
		Update("slope_angle", slopeAngle)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update slope angle")
	}
	return nil
}