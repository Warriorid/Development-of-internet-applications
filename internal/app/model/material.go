package model

import "math"

type Material struct {
    ID          int     `gorm:"primaryKey;column:id;autoIncrement"`
    Title       string  `gorm:"column:title;type:varchar(100);not null;unique"`
    Coefficient float64 `gorm:"column:coefficient;not null;check:coefficient>0"`
    ImageURL    string  `gorm:"column:image_url;type:varchar(255);null"`
    Description string  `gorm:"column:description;type:text;not null"`
    IsDeleted   bool    `gorm:"column:is_deleted;not null;default:false"`
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
	
	return totalVolume, nil
}