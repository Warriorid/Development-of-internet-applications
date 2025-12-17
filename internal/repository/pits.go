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
    err := r.db.Where("creator_id = ? AND status = 'draft'", creatorID).First(&pit).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return 0, 0, nil
        }
        return 0, 0, err
    }
    
    var count int64
    err = r.db.Model(&model.CalculationMaterial{}).
        Where("calculation_id = ?", pit.ID).
        Count(&count).Error
    if err != nil {
        return 0, 0, err
    }
    
    return pit.ID, int(count), nil
}


func (r *PitsPostgres) GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculationListItem, error) {
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
    if err != nil {
        return nil, err
    }
    result := make([]model.PitsCalculationListItem, 0, len(pits))
    for _, pit := range pits {
        calculatedCount, err := r.GetCalculatedMaterialsCount(pit.ID)
        if err != nil {
            return nil, err
        }

        var pitVolume *float64 = nil
        if pit.Status == "completed" {
            totalVolume, err := r.GetTotalPitVolume(pit.ID)
            if err == nil && totalVolume > 0 {
                pitVolume = &totalVolume
            }
        }
        
        result = append(result, pit.ToPitsCalculationListItem(calculatedCount, pitVolume))
    }
    return result, nil
}

func (r *PitsPostgres) GetTotalPitVolume(pitID int) (float64, error) {
    var totalVolume float64
    err := r.db.Model(&model.CalculationMaterial{}).
        Where("calculation_id = ? AND volume_result IS NOT NULL", pitID).
        Select("COALESCE(SUM(volume_result), 0)").
        Scan(&totalVolume).Error
        
    return totalVolume, err
}

func (r *PitsPostgres) GetCalculatedMaterialsCount(pitID int) (int, error) {
    var count int64
    
    err := r.db.Model(&model.CalculationMaterial{}).
        Where("calculation_id = ? AND volume_result IS NOT NULL", pitID).
        Count(&count).Error
        
    return int(count), err
}

func (r *PitsPostgres) GetPitWithMaterials(id, creatorId int) (model.PitsCalculation, []model.MaterialWithCalculationData, error) {
    var pit model.PitsCalculation
    materials := make([]model.MaterialWithCalculationData, 0)

    err := r.db.Where("id = ? AND creator_id = ? AND status != 'deleted'", id, creatorId).
        First(&pit).Error
    if err != nil {
        return model.PitsCalculation{}, materials, err
    }
    err = r.db.Table("materials").
        Select(`
            materials.*, 
            calculation_materials.slope_angle,
            calculation_materials.volume_result
        `).
        Joins("JOIN calculation_materials ON calculation_materials.material_id = materials.id").
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
			"completed_at": nil,
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
	return nil
}

func (r *PitsPostgres) CalculatePitVolume(pitID int, calculateFunc func(length, width, depth, angle, coefficient float64) (float64, error)) (float64, error) {
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
	length := 0.0
	if pit.PitLength != nil {
		length = *pit.PitLength
	}
	
	width := 0.0
	if pit.PitWidth != nil {
		width = *pit.PitWidth
	}
	
	depth := 0.0
	if pit.PitDepth != nil {
		depth = *pit.PitDepth
	}

	for _, cm := range calculationMaterials {
		slopeAngle := 0
		if cm.SlopeAngle != nil {
			slopeAngle = *cm.SlopeAngle
		}

		volume, err := calculateFunc(
			length, 
			width,   
			depth,     
			float64(slopeAngle), 
			cm.Material.Coefficient,
		)
		if err != nil {
			return 0, err
		}
		
		volumePtr := &volume
		err = r.db.Model(&model.CalculationMaterial{}).
			Where("calculation_id = ? AND material_id = ?", pitID, cm.MaterialID).
			Update("volume_result", volumePtr).Error
		if err != nil {
			return 0, err
		}
		
		totalVolume += volume
	}

	return totalVolume, nil
}

func (r *PitsPostgres) DeletePit(id int) error {
    result := r.db.Model(&model.PitsCalculation{}).Where("id = ?", id).
        Updates(map[string]interface{}{
            "status":       "deleted",
        })
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("pit not found or already deleted")
    }
    
    return nil
}

func (r *PitsPostgres) GetPitByID(id int) (*model.PitsCalculation, error) {
    var pit model.PitsCalculation
    err := r.db.Where("id = ?", id).First(&pit).Error
    if err != nil {
        return nil, err
    }
    return &pit, nil
}


func (r *PitsPostgres) GetPitWithMaterialsForAsync(id int) (*model.PitsCalculation, []model.MaterialWithCalculationData, error) {
    var pit model.PitsCalculation
    err := r.db.Where("id = ? AND status = 'formed'", id).First(&pit).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil, fmt.Errorf("pit not found")
        }
        return nil, nil, err
    }
    
    if pit.PitLength == nil || pit.PitWidth == nil || pit.PitDepth == nil {
        return nil, nil, fmt.Errorf("pit dimensions not set")
    }
    
    materials := make([]model.MaterialWithCalculationData, 0)
    err = r.db.Table("materials").
        Select(`
            materials.id,
            materials.title,
            materials.description,
            materials.coefficient,
            materials.image_url,
            calculation_materials.slope_angle,
            calculation_materials.volume_result
        `).
        Joins("JOIN calculation_materials ON calculation_materials.material_id = materials.id").
        Where("calculation_materials.calculation_id = ?", pit.ID).
        Find(&materials).Error
    if err != nil {
        return nil, nil, err
    }
    
    if len(materials) == 0 {
        return nil, nil, fmt.Errorf("no materials found for calculation")
    }
    
    return &pit, materials, nil
}

func (r *PitsPostgres) UpdateMaterialVolumeResult(calculationID, materialID int, volumeResult float64) error {
    return r.db.Model(&model.CalculationMaterial{}).
        Where("calculation_id = ? AND material_id = ?", calculationID, materialID).
        Update("volume_result", volumeResult).Error
}

func (r *PitsPostgres) UpdatePitStatus(id, moderatorID int, status string, completedAt *time.Time) error {
    updates := map[string]interface{}{
        "status": status,
        "moderator_id": moderatorID,
    }
    
    if completedAt != nil {
        updates["completed_at"] = completedAt
    }
    
    result := r.db.Model(&model.PitsCalculation{}).Where("id = ?", id).Updates(updates)
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("pit not found")
    }
    
    return nil
}