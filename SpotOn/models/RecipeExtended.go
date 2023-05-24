package models

type RecipeExtended struct {
	ID                    int           `json:"id"`
	Image                 string        `json:"image"`
	ImageType             string        `json:"imageType"`
	Likes                 int           `json:"likes"`
	MissedIngredientCount int           `json:"missedIngredientCount"`
	MissedIngredients     []IngredientR `json:"missedIngredients"`
	Title                 string        `json:"title"`
	UnusedIngredients     []IngredientR `json:"unusedIngredients"`
	UsedIngredientCount   int           `json:"usedIngredientCount"`
	UsedIngredients       []IngredientR `json:"usedIngredients"`
}

type IngredientR struct {
	Aisle        string   `json:"aisle"`
	Amount       float64  `json:"amount"`
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Meta         []string `json:"meta"`
	Name         string   `json:"name"`
	Original     string   `json:"original"`
	OriginalName string   `json:"originalName"`
	Unit         string   `json:"unit"`
	UnitLong     string   `json:"unitLong"`
	UnitShort    string   `json:"unitShort"`
}
