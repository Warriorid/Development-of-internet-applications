// pits.go
package repository

import (
	"DIA/internal/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PitsPostgres struct {
	db *gorm.DB
}

func NewPitsPostgres(db *gorm.DB) *PitsPostgres {
	return &PitsPostgres{db: db}
}

func (r *PitsPostgres) GetDraftPitWithItemsCount(creatorID int) (int, int, error) {
	var pit model.PitsCalculation
	var itemsCount int64

	err := r.db.Where("status = ? AND creator_id = ?", "draft", creatorID).First(&pit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	err = r.db.Model(&model.CalculationMaterial{}).Where("calculation_id = ?", pit.ID).Count(&itemsCount).Error
	if err != nil {
		return 0, 0, err
	}

	return pit.ID, int(itemsCount), nil
}


func (r *PitsPostgres) GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculation, error) {
	var pits []model.PitsCalculation

	query := r.db.Where("status != ? AND status != ?", "draft", "deleted")

	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}
	if startDate != nil {
		query = query.Where("formed_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("formed_at <= ?", endDate)
	}
	err := query.Order("formed_at DESC").Find(&pits).Error
	return pits, err
}

func (r *PitsPostgres) GetPitWithMaterials(id, creatorId int) (model.PitsCalculation, []model.Material, error) {
    var pit model.PitsCalculation
    materials := make([]model.Material, 0)
    
    err := r.db.Where("id = ? AND creator_id = ? AND status != 'deleted'", id, creatorId).
        First(&pit).Error
    if err != nil {
        return model.PitsCalculation{}, materials, err
    }

    err = r.db.Joins("JOIN calculation_materials ON calculation_materials.material_id = materials.id").
        Where("calculation_materials.calculation_id = ?", pit.ID).
        Find(&materials).Error
    if err != nil {
        return model.PitsCalculation{}, materials, err
    }

    return pit, materials, nil
}


func (r *PitsPostgres) UpdatePit(id int, input model.UpdatePitParam) error {
    var count int64
    err := r.db.Model(&model.PitsCalculation{}).Where("id = ? AND status != 'deleted'", id).Count(&count).Error
    if err != nil {
        return err
    }
    if count == 0 {
        return gorm.ErrRecordNotFound
    }
	err = r.db.Model(&model.PitsCalculation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"pit_length": input.PitLength,
		"pit_width":  input.PitWidth,
		"pit_depth":  input.PitDepth,
	}).Error
    return err
}

func (r *PitsPostgres) FormPit(id int) error {
	return r.db.Model(&model.PitsCalculation{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    "formed",
			"formed_at": time.Now(),
		}).Error
}
func (r *PitsPostgres) ValidatePitForForming(id, creatorId int) error {
	var pit model.PitsCalculation
	var materialsCount int64
	var materialsWithNullSlope int64

	err := r.db.Where("id = ? AND status = 'draft' AND creator_id = ?", id, creatorId).First(&pit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("pit not found")  
		}
		return err
	}
	err = r.db.Model(&model.CalculationMaterial{}).
		Where("calculation_id = ?", id).
		Count(&materialsCount).Error
	if err != nil {
		return err
	}
	if materialsCount == 0 {
		return fmt.Errorf("pit must have at least one material")
	}
	err = r.db.Model(&model.CalculationMaterial{}).
		Where("calculation_id = ? AND slope_angle IS NULL", id).
		Count(&materialsWithNullSlope).Error
	if err != nil {
		return err
	}

	if materialsWithNullSlope > 0 {
		return fmt.Errorf("all materials must have slope angle specified")
	}
	if pit.PitLength <= 0.1 || pit.PitWidth <= 0.1 || pit.PitDepth <= 0.1 {
		return fmt.Errorf("pit dimensions must be greater than default values")
	}
	return nil
}

func (r *PitsPostgres) CompletePit(id, moderatorID int, status string) error {
	var pit model.PitsCalculation
	err := r.db.Where("id = ? AND status = 'formed'", id).First(&pit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("pit not found")
		}
		return err
	}
	if status == "completed" {
		_, err = r.CalculatePitVolume(id)
		if err != nil {
			return err
		}
	}

	result := r.db.Model(&model.PitsCalculation{}).Where("id = ? AND status = 'formed'", id).
		Updates(map[string]interface{}{
			"status":       status,
			"moderator_id": moderatorID,
			"completed_at": time.Now(),
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("pit cannot be %s", status)
	}
	
	return nil
}

func (r *PitsPostgres) CalculatePitVolume(pitID int) (float64, error) {
	var calculationMaterials []model.CalculationMaterial
	var totalVolume float64
	err := r.db.Preload("Material").
		Where("calculation_id = ?", pitID).
		Find(&calculationMaterials).Error
	if err != nil {
		return 0, err
	}
	var pit model.PitsCalculation
	err = r.db.Where("id = ?", pitID).First(&pit).Error
	if err != nil {
		return 0, err
	}

	for _, cm := range calculationMaterials {
		volume, err := model.CalculateExcavationVolume(
			pit.PitLength,
			pit.PitWidth,
			pit.PitDepth,
			float64(cm.SlopeAngle),
			cm.Material.Coefficient,
		)
		if err != nil {
			return 0, err
		}
		
		err = r.db.Model(&model.CalculationMaterial{}).
			Where("calculation_id = ? AND material_id = ?", pitID, cm.MaterialID).
			Update("volume_result", volume).Error
		if err != nil {
			return 0, err
		}
		
		totalVolume += volume
	}

	return totalVolume, nil
}

func (r *PitsPostgres) DeletePit(id int) error {
	var pit model.PitsCalculation
	err := r.db.Where("id = ? AND status != 'deleted'", id).First(&pit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("pit not found")
		}
		return err
	}
	result := r.db.Model(&model.PitsCalculation{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":    "deleted",
			"formed_at": time.Now(),
			"completed_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("pit cannot be deleted")
	}
	return nil
}