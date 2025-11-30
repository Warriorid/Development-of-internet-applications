package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"DIA/pkg/minio"
	"bytes"
	"fmt"
)

type MaterialService struct {
    repo repository.Material
	imageSrorage minio.ImageStorage
}

func NewMaterialService(repo repository.Material, imageSrorage minio.ImageStorage) *MaterialService {
    return &MaterialService{
		repo: repo,
		imageSrorage: imageSrorage,
	}
}

func (s *MaterialService) GetMaterials(materialTitle string) ([]model.Material, error) {
	return s.repo.GetMaterials(materialTitle)
}

func (s *MaterialService) GetMaterial(id int) (model.Material, error) {
	material, err := s.repo.GetMaterial(id)
	if material.ID == 0 {
		return model.Material{}, fmt.Errorf("material not found")
	}
	return material, err
}

func (s *MaterialService) CreateMaterial(material *model.Material) (model.Material, error) {
	return s.repo.CreateMaterial(material)
}

func (s *MaterialService) UpdateMaterial(id int, material *model.Material) error {
    existingMaterial, err := s.repo.GetMaterial(id)
    if err != nil || existingMaterial.ID == 0 {
        return fmt.Errorf("material not found")
    }
    return s.repo.UpdateMaterial(id, material)
}

func (s *MaterialService) DeleteMaterial(id int) error {
	material, err := s.repo.GetMaterial(id)
	if material.ID == 0 {
		return fmt.Errorf("material not found: %s", err)
	}
	if material.ImageURL != "" {
        if err := s.imageSrorage.DeleteImage(material.Title); err != nil {
            fmt.Printf("Warning: failed to delete image from Minio: %v\n", err)
        }
    }
    
	return s.repo.DeleteMaterial(id)
}

func (s *MaterialService) AddMaterialToPit(userId, materialId int) error {
    pitId, err := s.repo.GetOrCreateDraftPit(userId)
    if err != nil {
        return err
    }
    return s.repo.AddMaterialToPit(pitId, materialId)
}

func (s *MaterialService) UploadMaterialImage(id int, file []byte, filename string) error {
    material, err := s.repo.GetMaterial(id)
    if material.ID == 0 {
        return fmt.Errorf("material not found")
    }
    
    reader := bytes.NewReader(file)
    imageURL, err := s.imageSrorage.UploadImage(material.Title, reader, int64(len(file)), filename)
    if err != nil {
        return err
    }
    
    return s.repo.UpdateMaterialImage(id, imageURL)
}