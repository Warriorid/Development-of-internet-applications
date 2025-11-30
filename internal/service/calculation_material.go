package service

import (
	"DIA/internal/repository"
	"fmt"
)

type CalculationMaterialService struct {
	repo repository.CalculationMaterial
}

func NewCalculationMaterialService(repo repository.CalculationMaterial) *CalculationMaterialService{
	return &CalculationMaterialService{repo: repo}
}

func (s *CalculationMaterialService) RemoveMaterialFromPit(calculationID, materialID int) error {
	return s.repo.RemoveMaterialFromPit(calculationID, materialID)
}

func (s *CalculationMaterialService) UpdateCalculationMaterial(calculationID, materialID int, slopeAngle int) error {
	if slopeAngle < 0 || slopeAngle > 45 {
		return fmt.Errorf("slope angle must be between 0 and 45 degrees")
	}
	return s.repo.UpdateCalculationMaterial(calculationID, materialID, slopeAngle)
}
