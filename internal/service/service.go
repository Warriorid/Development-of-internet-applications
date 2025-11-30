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
	CalculationMaterial
	User
}

type Material interface{
	GetMaterials(materialTitle string) ([]model.Material, error)
	GetMaterial(id int) (model.Material, error)
	CreateMaterial(material *model.Material) (model.Material, error)
	UpdateMaterial(id int, material *model.Material) error
	DeleteMaterial(id int) error 
	AddMaterialToPit(userId, id int) error
	UploadMaterialImage(id int, file []byte, filename string) error
}

type Pits interface{
	GetDraftPitWithItemsCount(creatorID int) (int, int, error)
	GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculationListItem, error)
	GetPitWithMaterials(id, creatorId int) (model.PitsCalculationWithMaterials, error)
	UpdatePit(id int, input model.UpdatePitParam) error
	FormPit(id, creatorId int) error
	CompletePit(id, moderatorID int, status string) error
	DeletePit(id int) error
	IsPitOwner(pitID, userID int) (bool, error)
}

type CalculationMaterial interface{
	RemoveMaterialFromPit(calculationID, materialID int) error
	UpdateCalculationMaterial(calculationID, materialID int, slopeAngle int) error
 }

type User interface {
	Register(user *model.Users) (*model.Users, error)
	Login(username, password string) (*model.LoginResp, error)
	GetUserProfile(id int) (*model.Users, error)
	UpdateUserProfile(id int, username, password string) (*model.Users, error)
	Logout() error
}

func NewService(repo *repository.Repository, minio *minio.MinioClient) *Service {
	return &Service{
		Material: NewMaterialService(repo.Material, minio),
		Pits: NewPitsService(repo.Pits),
		CalculationMaterial: NewCalculationMaterialService(repo.CalculationMaterial),
		User: NewUserService(repo.User),
	}
}
