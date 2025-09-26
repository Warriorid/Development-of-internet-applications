package model

type Material struct {
    ID          int
    Title       string
    Coefficient string
    ImageURL    string
	Description string
}
type PitMaterial struct {
    PitID      int
    MaterialID int
    SlopeAngle int   
    Volume     float64
    Material   Material
}

type Pit struct {
    ID         int
    Length     float64
    Width      float64
    Depth      float64
    PitMaterials []PitMaterial
}




