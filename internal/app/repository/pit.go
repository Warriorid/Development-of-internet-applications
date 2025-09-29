package repository

import (
	"DIA/internal/app/model"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetDraftPit(id int) (model.PitsCalculation, error) {
    var pit model.PitsCalculation
    creatorID := 2
    
    err := r.db.Preload("Creator").Preload("Moderator").Preload("Materials").Preload("Materials.Material").
        Where("id = ? AND status = ? AND creator_id = ?", id, "draft", creatorID).
        First(&pit).Error
    if err != nil {
        return model.PitsCalculation{}, err
    }
    return pit, nil
}

func (r *Repository) GetPitCount() int64 {
    var count int64
    creatorID := 2

    err := r.db.Model(&model.PitsCalculation{}).
        Joins("JOIN calculation_materials cm ON pits_calculations.id = cm.calculation_id").
        Where("pits_calculations.creator_id = ? AND pits_calculations.status = ?", creatorID, "draft").
        Count(&count).Error

    if err != nil {
        logrus.Println("Error counting materials in draft pit:", err)
        return 0
    }
    return count
}

func (r *Repository) GetDraftPitID() (int, error) {
    var pit model.PitsCalculation
    creatorID := 2 

    err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
        Select("id").
        First(&pit).Error
        
    if err != nil {
        return 0, err
    }
    return pit.ID, nil
}

func (r *Repository) DeletePit(pitId int) error {
    query := "UPDATE pits_calculations SET status = 'deleted' WHERE id = $1 AND status = 'draft'"
    result := r.db.Exec(query, pitId)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("pit not found or already deleted")
    }
    return nil
}


func (r *Repository) GetOrCreateDraftPit() (int, error) {
    creatorID := 2
    
    var pit model.PitsCalculation
    err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
        Select("id").
        First(&pit).Error
        
    if err != nil {
        newPit := model.PitsCalculation{
            CreatorID: creatorID,
            Status:    "draft",
        }
        
        err = r.db.Create(&newPit).Error
        if err != nil {
            return 0, err
        }
        return newPit.ID, nil
    }
    
    return pit.ID, nil
}