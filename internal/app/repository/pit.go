package repository

import (
	"DIA/internal/app/model"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetPit(id int) (model.PitsCalculation, error) {
    var pit model.PitsCalculation
    creatorID := 2
    
    err := r.db.Preload("Creator").Preload("Moderator").Preload("Materials").Preload("Materials.Material").
        Where("id = ? AND creator_id = ? AND status != 'deleted'", id, creatorID).
        First(&pit).Error
    if err != nil {
        return model.PitsCalculation{}, err
    }
    return pit, nil
}

func (r *Repository) GetPitCount(id int) int64 {
    var count int64
    err := r.db.Model(&model.CalculationMaterial{}).
        Where("calculation_id = ?", id).
        Count(&count).Error
    if err != nil {
        logrus.Println("Error counting materials in pit:", err)
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
    query := "UPDATE pits_calculations SET status = 'deleted' WHERE id = $1"
    result := r.db.Exec(query, pitId)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("pit not found or already deleted")
    }
    return nil
}


func (r *Repository) CreateDraftPit() (int, error) {
    creatorID := 2
    
    newPit := model.PitsCalculation{
        CreatorID: creatorID,
        Status:    "draft",
        PitLength: nil,
        PitWidth:  nil,
        PitDepth:  nil,
    }
        
    err := r.db.Create(&newPit).Error
    if err != nil {
        return 0, err
    }
    return newPit.ID, nil
}