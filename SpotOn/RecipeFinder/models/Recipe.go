package models

type Recipe struct {
	Name               string
	PresentIngredients []string
	MissingIngredients []string
	Carbs              string
	Proteins           string
	Calories           string
}
