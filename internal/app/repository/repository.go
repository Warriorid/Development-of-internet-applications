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

type Cart struct {
	ID int
	CartItems []Material
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

func (r *Repository)GetCarts()([]Cart, error){
	carts := []Cart {
		{
			ID: 1,
			CartItems: []Material{
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
	if len(carts) == 0 {
		return []Cart{}, fmt.Errorf("Не найдено ни одной корзины")
	}
	return carts, nil
}

func (r *Repository) GetCart(id int)(Cart, int, error){
	
	carts, err := r.GetCarts()
	if err != nil {
		return Cart{}, 0, err
	}
	for _, cart := range carts {
		if cart.ID == id {
			cartSize := len(cart.CartItems)
			return cart, cartSize, nil
		}
	}
	
    return Cart{}, 0, fmt.Errorf("корзина пуста")
}