package repository

import (
	"DIA/internal/app/model"
	"fmt"
)



func (r *Repository) GetPits() ([]model.Pit, error) {
    pits := []model.Pit{
        {
            ID:     1,
            Length: 10,
            Width:  10,
            Depth:  10,
            PitMaterials: []model.PitMaterial{
                {
                    PitID:      1,
                    MaterialID: 1,
                    SlopeAngle: 45,
                    Volume:     1150,
                    Material: model.Material{
                        ID:          1,
                        Title:       "Песок",
                        Coefficient: "1.15",
                        ImageURL:    "http://localhost:9000/materials/sand.jpg",
                        Description: "Рыхлый, несвязный грунт с низким коэффициентом разрыхления.",
                    },
                },
                {
                    PitID:      1,
                    MaterialID: 2,
                    SlopeAngle: 30,
                    Volume:     1350,
                    Material: model.Material{
                        ID:          2,
                        Title:       "Глина",
                        Coefficient: "1.35",
                        ImageURL:    "http://localhost:9000/materials/clay.jpg",
                        Description: "Плотный, связный грунт с высоким коэффициентом разрыхления...",
                    },
                },
            },
        },
    }
    
    if len(pits) == 0 {
        return nil, fmt.Errorf("Не найдено ни одного котлована")
    }
    return pits, nil
}

func (r *Repository) GetPit(id int) (model.Pit, int, error) {
    pits, err := r.GetPits()
    if err != nil {
        return model.Pit{}, 0, err
    }
    
    for _, pit := range pits {
        if pit.ID == id {
            pitsCount := len(pit.PitMaterials)
            return pit, pitsCount, nil
        }
    }
    
    return model.Pit{}, 0, fmt.Errorf("Котлован не найден")
}

