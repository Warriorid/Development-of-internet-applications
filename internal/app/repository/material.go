package repository

import (
	"DIA/internal/app/model"
	"fmt"
	"strings"
)

func (r *Repository) GetMaterials() ([]model.Material, error) {
    materials := []model.Material{
        {
            ID: 1,
            Title: "Песок",
            Coefficient: "1.15",
            ImageURL: "http://localhost:9000/materials/sand.jpg",
			Description: "Рыхлый, несвязный грунт с низким коэффициентом разрыхления.",
        },
        {
            ID: 2,
            Title: "Глина",
            Coefficient: "1.35",
            ImageURL: "http://localhost:9000/materials/clay.jpg",
			Description: "Плотный, связный грунт с высоким коэффициентом разрыхления, обладающий высокой пластичностью и водонепроницаемостью.",
        },
        {
            ID: 3,
            Title: "Суглинка",
            Coefficient: "1.25",
            ImageURL: "http://localhost:9000/materials/loam.jpg",
			Description: "Сбалансированная смесь песка и глины с умеренным коэффициентом разрыхления, сочетающая дренажные свойства песка и связность глины.",
        },
		{
            ID: 4,
            Title: "Чернозем",
            Coefficient: "1.2",
            ImageURL: "http://localhost:9000/materials/black_earth.jpeg",
			Description: "Плодородный почвенный слой с высоким содержанием гумуса и органических веществ.",
			
        },
		
    }
    
    if len(materials) == 0 {
        return nil, fmt.Errorf("массив материалов пустой")
    }
    
    return materials, nil
}

func (r *Repository) GetMaterial(materialId int) (model.Material, error) {
	materials, err := r.GetMaterials()
	if err != nil {
		return model.Material{}, err
	}
	for _, material := range materials {
		if material.ID == materialId {
			return material, nil
		}
	}
	return model.Material{}, fmt.Errorf("Материал не найден")
}



func (r *Repository) GetMaterialsByTitle(title string) ([]model.Material, error) {
    materials, err := r.GetMaterials()
    if err != nil {
        return []model.Material{}, err
    }

    var result []model.Material
    for _, material := range materials {
        if strings.Contains(strings.ToLower(material.Title), strings.ToLower(title)) {
            result = append(result, material)
        }
    }

    return result, nil
}