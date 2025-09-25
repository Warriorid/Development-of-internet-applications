package repository

import (
	"fmt"
	"strings"
)

type Repository struct{}


func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Material struct {
    ID          int
    Title       string
    Coefficient string
    ImageURL    string
	Description string
}

type Pits struct {
	ID int
	PitsItems []Material
}


func (r *Repository) GetMaterials() ([]Material, error) {
    materials := []Material{
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

func (r *Repository) GetMaterial(materialId int) (Material, error) {
	materials, err := r.GetMaterials()
	if err != nil {
		return Material{}, err
	}
	for _, material := range materials {
		if material.ID == materialId {
			return material, nil
		}
	}
	return Material{}, fmt.Errorf("Материал не найден")
}



func (r *Repository) GetMaterialsByTitle(title string) ([]Material, error) {
    materials, err := r.GetMaterials()
    if err != nil {
        return []Material{}, err
    }

    var result []Material
    for _, material := range materials {
        if strings.Contains(strings.ToLower(material.Title), strings.ToLower(title)) {
            result = append(result, material)
        }
    }

    return result, nil
}

func (r *Repository)GetPits()([]Pits, error){
	pits := []Pits {
		{
			ID: 1,
			PitsItems: []Material{
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
			},
		},
	}
	if len(pits) == 0 {
		return []Pits{}, fmt.Errorf("Не найдено ни одного котлована")
	}
	return pits, nil
}

func (r *Repository) GetPit(id int)(Pits, int, error){
	
	pits, err := r.GetPits()
	if err != nil {
		return Pits{}, 0, err
	}
	for _, pit := range pits {
		if pit.ID == id {
			pitsCount := len(pit.PitsItems)
			return pit, pitsCount, nil
		}
	}
	
    return Pits{}, 0, fmt.Errorf("Не добавлено ни одного котлована")
}