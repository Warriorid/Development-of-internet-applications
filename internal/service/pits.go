package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
)

const (
    AsyncServiceURL   = "http://localhost:8000/api/calculate/"
    AuthToken         = "8ba1f1945d2a7c6e"
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
    fmt.Printf("[Go Service] CompletePit called: id=%d, moderatorID=%d, status=%s\n", 
        id, moderatorID, status)
    
    pit, materials, err := s.repo.GetPitWithMaterialsForAsync(id)
    if err != nil {
        fmt.Printf("[Go Service] Error getting pit: %v\n", err)
        return err
    }
    
    fmt.Printf("[Go Service] Pit status: %s, materials count: %d\n", pit.Status, len(materials))
    if pit.Status != "formed" {
        return fmt.Errorf("can only complete or reject pits with formed status")
    }
    
    completedAt := time.Now()
    if status == "rejected" {
        fmt.Printf("[Go Service] Rejecting pit %d\n", id)
        return s.repo.UpdatePitStatus(id, moderatorID, status, &completedAt)
    }
    
    fmt.Printf("[Go Service] Completing pit %d\n", id)
    if err := s.repo.UpdatePitStatus(id, moderatorID, "completed", &completedAt); err != nil {
        fmt.Printf("[Go Service] Error updating status: %v\n", err)
        return err
    }
    
    for _, material := range materials {
        if material.SlopeAngle == nil {
            fmt.Printf("[Go Service] Skipping material %d - no slope angle\n", material.ID)
            continue 
        }
        
        fmt.Printf("[Go Service] Starting async calculation for material %d (slope=%d, coeff=%f)\n", 
            material.ID, *material.SlopeAngle, material.Coefficient)
        go s.startAsyncCalculation(
            pit.ID,
            material.ID,
            *pit.PitLength,
            *pit.PitWidth,
            *pit.PitDepth,
            *material.SlopeAngle,
            material.Coefficient,
        )
    }
    
    fmt.Printf("[Go Service] CompletePit finished successfully\n")
    return nil
}

func (s *PitsService) startAsyncCalculation(calculationID, materialID int, length, width, depth float64, 
                                          slopeAngle int, coefficient float64) {
    fmt.Printf("[Go Service] startAsyncCalculation: calculation=%d, material=%d\n", 
        calculationID, materialID)
    
    request := model.AsyncCalculationRequest{
        CalculationID: calculationID,
        MaterialID:    materialID,
        PitLength:     length,
        PitWidth:      width,
        PitDepth:      depth,
        SlopeAngle:    slopeAngle,
        Coefficient:   coefficient,
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        fmt.Printf("[Go Service] Error marshaling: %v\n", err)
        return
    }
    
    fmt.Printf("[Go Service] Sending request to Django: %s\n", AsyncServiceURL)
    resp, err := http.Post(AsyncServiceURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("[Go Service] Error calling Django: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("[Go Service] Django response status: %d, body: %s\n", resp.StatusCode, string(body))
}

func (s *PitsService) CompleteAsyncCalculation(result model.AsyncCalculationResult) error {
    fmt.Printf("[Go Service] CompleteAsyncCalculation received: calculation=%d, material=%d, volume=%f, token=%s\n",
        result.CalculationID, result.MaterialID, result.VolumeResult, result.Token)
    
    if result.Token != AuthToken {
        fmt.Printf("[Go Service] Invalid token! Received: %s, Expected: %s\n", result.Token, AuthToken)
        return fmt.Errorf("unauthorized: invalid token")
    }
    
    fmt.Printf("[Go Service] Token valid. Updating volume result in DB for calculation=%d, material=%d\n",
        result.CalculationID, result.MaterialID)
    
    err := s.repo.UpdateMaterialVolumeResult(result.CalculationID, result.MaterialID, result.VolumeResult)
    if err != nil {
        fmt.Printf("[Go Service] Error updating DB: %v\n", err)
        return err
    }
    
    fmt.Printf("[Go Service] Volume result updated successfully\n")
    return nil
}

func (s *PitsService) CanUserDeletePit(pitID, userID, userRole int) (bool, error) {
    pit, err := s.repo.GetPitByID(pitID)
    if err != nil {
        return false, err
    }
    if userRole == 1 {
        return pit.Status == "draft", nil
    }
    
    return pit.CreatorID == userID && pit.Status == "draft", nil
}

func (s *PitsService) DeletePit(id int) error {
    pit, err := s.repo.GetPitByID(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("pit not found")
        }
        return err
    }
    
    if pit.Status != "draft" {
        return fmt.Errorf("only draft pits can be deleted")
    }
    
    err = s.repo.DeletePit(id)
    if err != nil {
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

func (s *PitsService) GetPitStatus(pitID int) (string, error) {
	pit, err := s.repo.GetPitByID(pitID)
	if err != nil {
		return "", err
	}
	return pit.Status, nil
}