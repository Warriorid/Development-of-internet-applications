package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"DIA/pkg/minio"
	"time"
)

type Service struct{
	Material
	Pits
}

type Material interface{
	GetMaterials(materialTitle string) ([]model.Material, error)
	GetMaterial(id int) (model.Material, error)
	CreateMaterial(material *model.Material) (model.Material, error)
	UpdateMaterial(id int, material *model.Material) error
	DeleteMaterial(id int) error 
	AddMaterialToPit(id int) error
	UploadMaterialImage(id int, file []byte, filename string) error
}

type Pits interface{
	GetDraftPitWithItemsCount(creatorID int) (int, int, error)
	GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculation, error)
	GetPitWithMaterials(id, creatorId int) (model.PitsCalculationWithMaterial, error)
	UpdatePit(id int, input model.UpdatePitParam) error
	FormPit(id, creatorId int) error
	CompletePit(id, moderatorID int, status string) error
	DeletePit(id int) error
}

func NewService(repo *repository.Repository, minio *minio.MinioClient) *Service {
	return &Service{
		Material: NewMaterialService(repo.Material, minio),
		Pits: NewPitsService(repo.Pits),
	}
}