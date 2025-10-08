package model

import "math"

type Material struct {
    ID          int     `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Title       string  `gorm:"column:title;type:varchar(100);not null;unique" json:"title"`
    Coefficient float64 `gorm:"column:coefficient;not null;check:coefficient>0" json:"coefficient"`
    ImageURL    string  `gorm:"column:image_url;type:varchar(255);null" json:"image_url"`
    Description string  `gorm:"column:description;type:text;not null" json:"description"`
    IsDeleted   bool    `gorm:"column:is_deleted;not null;default:false" json:"is_deleted"`
}


// V = K × [ (H/3) × (L×B + (L + 2H·tan(α))×(B + 2H·tan(α)) + √(L×B×(L + 2H·tan(α))×(B + 2H·tan(α))) ) ]

func CalculateExcavationVolume(length, width, depth, angle, coefficient float64) (float64, error) {
	
	mainVolume := length * width * depth

	if angle == 0 {
		return mainVolume * coefficient, nil
	}

	angleRad := angle * math.Pi / 180
	tan := math.Tan(angleRad)
	
	topLength := length + 2*depth*tan
	topWidth := width + 2*depth*tan
	
	bottomArea := length * width     
	topArea := topLength * topWidth
	
	pyramidVolume := (depth / 3) * (bottomArea + topArea + math.Sqrt(bottomArea*topArea))
	totalVolume := pyramidVolume * coefficient
	
	return math.Round(totalVolume * 100) / 100, nil
}