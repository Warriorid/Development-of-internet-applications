package repository

import (
	"DIA/internal/app/model"
	"fmt"

	"github.com/sirupsen/logrus"
)


func (r *Repository) GetMaterials() ([]model.Material, error) {
	var materials []model.Material
	err := r.db.Find(&materials).Error
	if err != nil {
		return nil, err
	}
	if len(materials) == 0 {
		return nil, fmt.Errorf("массив материалов пустой")
	}
	return materials, nil

}

func (r *Repository) GetMaterial(materialId int) (model.Material, error) {
    query := "SELECT id, title, coefficient, image_url, description, is_deleted FROM materials WHERE id = $1"
    var material model.Material
    err := r.db.Raw(query, materialId).Row().Scan(
        &material.ID,
        &material.Title,
        &material.Coefficient,
        &material.ImageURL,
        &material.Description,
        &material.IsDeleted,
    )
    if err != nil {
        return model.Material{}, err
    }
    return material, nil
}

func (r *Repository) GetMaterialsByTitle(title string) ([]model.Material, error) {
    var materials []model.Material
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}

func (r *Repository) AddMaterialToPit(materialID int) error {
    pitID, err := r.GetOrCreateDraftPit()
    if err != nil {
        return fmt.Errorf("failed to get or create draft pit: %v", err)
    }
    
    var existingMaterial model.CalculationMaterial
    err = r.db.Where("calculation_id = ? AND material_id = ?", pitID, materialID).
        First(&existingMaterial).Error

    if err == nil {
        logrus.Println("material already in pit")
        return nil
    }
    
    calculationMaterial := model.CalculationMaterial{
        CalculationID: pitID,
        MaterialID:    materialID,
    }
        
    err = r.db.Create(&calculationMaterial).Error
    if err != nil {
        return err
    }
        
    return nil
}