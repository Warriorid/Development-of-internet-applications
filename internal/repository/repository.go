package repository

import (
	"DIA/internal/model"

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

type Pits interface{}



func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
      Material: NewMaterialPostgres(db),
      Pits: NewPitsPostgres(db),
   }
}
