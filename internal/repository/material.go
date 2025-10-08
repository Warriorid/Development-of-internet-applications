package repository

import (
	"DIA/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)


type MaterialPostgres struct {
    db *gorm.DB
}

func NewMaterialPostgres(db *gorm.DB) *MaterialPostgres {
    return &MaterialPostgres{db: db}
}

func (r *MaterialPostgres) GetMaterials(titleFilter string) ([]model.Material, error) {
	var materials []model.Material
	err := r.db.Where("title ILIKE ? AND is_deleted = ?", "%"+titleFilter+"%", false).Find(&materials).Error
	return materials, err
}

func (r *MaterialPostgres) GetMaterial(id int) (model.Material, error) {
    var material model.Material
    err := r.db.Where("id = ? AND is_deleted = ?", id, false).Find(&material).Error

    return material, err
}

func (r *MaterialPostgres) CreateMaterial(material *model.Material) (model.Material, error){
    var createdMaterial model.Material
    createdMaterial = *material
    
    err := r.db.Create(&createdMaterial).Error
    return createdMaterial, err
}

func (r *MaterialPostgres) UpdateMaterial(id int, material *model.Material) error {
	err := r.db.Where("id = ?", id).Updates(material).Error
	return err
}

func (r *MaterialPostgres) DeleteMaterial(id int) error {
	err := r.db.Model(&model.Material{}).Where("id = ?", id).Update("is_deleted", true).Error
	return err
}


func (r *MaterialPostgres) GetOrCreateDraftPit(userID int) (int, error) {
    var pit model.PitsCalculation
    err := r.db.Where("status = ? AND creator_id = ?", "draft", userID).First(&pit).Error
    
    if err == nil {
        return pit.ID, nil
    }
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        newPit := model.PitsCalculation{
            Status:    "draft",
            CreatorID: userID,
            CreatedAt: time.Now(),
            PitLength: 0.1,
            PitWidth:  0.1,
            PitDepth:  0.1,
        }
        err = r.db.Create(&newPit).Error
        if err != nil {
            return 0, err
        }
        return newPit.ID, nil
    }
    
    return 0, err
}

func (r *MaterialPostgres) AddMaterialToPit(calculationId, materialId int) error {
	calculationMaterial := model.CalculationMaterial{
		CalculationID: calculationId,
		MaterialID: materialId,
		SlopeAngle: 0,
	}
	err := r.db.Create(&calculationMaterial).Error
	return err
}

func (r *MaterialPostgres) UpdateMaterialImage(id int, imageURL string) error {
    return r.db.Model(&model.Material{}).Where("id = ?", id).Update("image_url", imageURL).Error
}