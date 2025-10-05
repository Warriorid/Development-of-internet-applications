package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"DIA/pkg"
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

type Pits interface{}

func NewService(repo *repository.Repository, minio *pkg.MinioClient) *Service {
	return &Service{
		Material: NewMaterialService(repo.Material, minio),
		Pits: NewPitsService(repo.Pits),
	}
}