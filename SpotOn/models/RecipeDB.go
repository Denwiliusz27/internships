package models

type RecipeDB struct {
	ID                 uint
	Name               string
	Calories           string
	Proteins           string
	Carbs              string
	PresentIngredients []PresentIngredient `gorm:"foreignKey:RecipeID"`
	MissingIngredients []MissingIngredient `gorm:"foreignKey:RecipeID"`
}

type PresentIngredient struct {
	ID       uint
	RecipeID uint
	Name     string
}

type MissingIngredient struct {
	ID       uint
	RecipeID uint
	Name     string
}

func (r *RecipeDB) TableName() string {
	return "recipes"
}
