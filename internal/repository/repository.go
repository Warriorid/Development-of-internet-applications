package repository

import (
	"DIA/internal/model"
	"time"

	"gorm.io/gorm"
)

type Repository struct{
	Material
   Pits
}

type Material interface{
   GetMaterials(materialTitle string)([]model.Material, error)
   GetMaterial(id int) (model.Material, error)
   CreateMaterial(material *model.Material) (model.Material, error)
   UpdateMaterial(id int, material *model.Material) error
   DeleteMaterial(id int) error 
   AddMaterialToPit(calculationId, materialId int) error
   UpdateMaterialImage(id int, imageURL string) error

   GetOrCreateDraftPit() (int, error)

   
}

type Pits interface{
   GetDraftPitWithItemsCount(creatorID int) (int, int, error)
   GetPits(statusFilter string, startDate, endDate *time.Time) ([]model.PitsCalculation, error)
   GetPitWithMaterials(id, creatorId int) (model.PitsCalculation, []model.Material, error)
   UpdatePit(id int, input model.UpdatePitParam) error
   FormPit(id int) error
   ValidatePitForForming(id, creatorId int) error
   CompletePit(id, moderatorID int, status string) error
   DeletePit(id int) error
}



func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
      Material: NewMaterialPostgres(db),
      Pits: NewPitsPostgres(db),
   }
}
