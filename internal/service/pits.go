package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"fmt"
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
	pitId, materialsCount, err := s.repo.GetDraftPitWithItemsCount(creatorID)
	if pitId == 0 {
		return 0, 0, fmt.Errorf("pit not found")
	}
	return pitId, materialsCount, err
}

func (s *PitsService) GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculation, error) {
	return s.repo.GetPits(statusFilter, startDate, endDate)
}

func (s *PitsService) GetPitWithMaterials(id, creatorId int) (model.PitsCalculationWithMaterial, error) {
	pit, materials, err := s.repo.GetPitWithMaterials(id, creatorId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.PitsCalculationWithMaterial{}, fmt.Errorf("not found")
		}
		return model.PitsCalculationWithMaterial{}, err
	}

	return pit.ToPitsCalculationWithMaterial(materials), nil
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

	err := s.repo.CompletePit(id, moderatorID, status)
	if err != nil {
		if err.Error() == "pit not found" {
			return fmt.Errorf("not found")
		}
		return err
	}
	return nil
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