package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type PitsService struct {
	repo repository.Pits
}

func NewPitsService(repo repository.Pits) *PitsService {
	return &PitsService{repo: repo}
}

func (s *PitsService) GetDraftPitWithItemsCount(creatorID int) (int, int, error) {
    return s.repo.GetDraftPitWithItemsCount(creatorID)
}

func (s *PitsService) GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculationListItem, error) {
	return s.repo.GetPits(statusFilter, startDate, endDate)
}

func (s *PitsService) GetPitWithMaterials(id, creatorId int) (model.PitsCalculationWithMaterials, error) {
    pit, materials, err := s.repo.GetPitWithMaterials(id, creatorId)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return model.PitsCalculationWithMaterials{}, fmt.Errorf("not found")
        }
        return model.PitsCalculationWithMaterials{}, err
    }

    return pit.ToPitsCalculationWithMaterials(materials), nil
}

func (s *PitsService) UpdatePit(id int, input model.UpdatePitParam) error {
	err := s.repo.UpdatePit(id, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("not found")
		}
		return err
	}
	return err
}

func (s *PitsService) FormPit(id, creatorId int) error {
	err := s.repo.ValidatePitForForming(id, creatorId)
	if err != nil {
		return err
	}
	return s.repo.FormPit(id)
}

func (s *PitsService) CompletePit(id, moderatorID int, status string) error {
	err := s.repo.CompletePit(id, moderatorID, status, s.calculateExcavationVolume)
	if err != nil {
		if err.Error() == "pit not found" {
			return fmt.Errorf("not found")
		}
		return err
	}
	return nil
}

func (s *PitsService) calculateExcavationVolume(length, width, depth, angle, coefficient float64) (float64, error) {
	mainVolume := length * width * depth

	if angle == 0 {
		return mainVolume * coefficient, nil
	}

	angleRad := angle * math.Pi / 180
	tan := math.Tan(angleRad)
	
	topLength := length + 2*depth*tan
	topWidth := width + 2*depth*tan
	
	bottomArea := length * width     
	topArea := topLength * topWidth
	
	pyramidVolume := (depth / 3) * (bottomArea + topArea + math.Sqrt(bottomArea*topArea))
	totalVolume := pyramidVolume * coefficient
	
	return math.Round(totalVolume * 100) / 100, nil
}

func (s *PitsService) DeletePit(id int) error {
	err := s.repo.DeletePit(id)
	if err != nil {
		if err.Error() == "pit not found" {
			return fmt.Errorf("not found")
		}
		return err
	}
	return nil
}

func (s *PitsService) IsPitOwner(pitID, userID int) (bool, error) {
    pit, err := s.repo.GetPitByID(pitID)
    if err != nil {
        return false, err
    }
    return pit.CreatorID == userID, nil
}