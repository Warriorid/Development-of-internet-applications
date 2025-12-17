package repository

import (
	"DIA/internal/model"
	"time"

	"gorm.io/gorm"
)
type Repository struct{
	Material
   Pits
   CalculationMaterial
   User
}
type Material interface{
   GetMaterials(materialTitle string)([]model.Material, error)
   GetMaterial(id int) (model.Material, error)
   CreateMaterial(material *model.Material) (model.Material, error)
   UpdateMaterial(id int, material *model.Material) error
   DeleteMaterial(id int) error 
   AddMaterialToPit(calculationId, materialId int) error
   UpdateMaterialImage(id int, imageURL string) error

   GetOrCreateDraftPit(userId int) (int, error)

   
}
type Pits interface{
   GetDraftPitWithItemsCount(creatorID int) (int, int, error)
   GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculationListItem, error)
   GetPitWithMaterials(id, creatorId int) (model.PitsCalculation, []model.MaterialWithCalculationData, error)
   UpdatePit(id int, input model.UpdatePitParam) error
   FormPit(id int) error
   ValidatePitForForming(id, creatorId int) error
   // CompletePit(id, moderatorID int, status string, calculateFunc func(length, width, depth, angle, coefficient float64) (float64, error)) error
   DeletePit(id int) error
   GetPitByID(id int) (*model.PitsCalculation, error)
   GetPitWithMaterialsForAsync(id int) (*model.PitsCalculation, []model.MaterialWithCalculationData, error)
   UpdateMaterialVolumeResult(calculationID, materialID int, volumeResult float64) error
   UpdatePitStatus(id, moderatorID int, status string, completedAt *time.Time) error
   GetCalculatedMaterialsCount(pitID int) (int, error)
   GetTotalPitVolume(pitID int) (float64, error)
}
type CalculationMaterial interface{
   RemoveMaterialFromPit(calculationID, materialID int) error
   UpdateCalculationMaterial(calculationID, materialID int, slopeAngle int) error
}
type User interface {
   CreateUser(user *model.Users) error
   GetUserByID(id int) (*model.Users, error)
   UpdateUser(id int, user *model.Users) error
   GetUserByUsername(username string) (*model.Users, error)
   Authentication(username, password string) (*model.Users, error)
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
      Material: NewMaterialPostgres(db),
      Pits: NewPitsPostgres(db),
      CalculationMaterial: NewCalculationMaterialPostgres(db),
      User: NewUserPostgres(db),
   }
}
