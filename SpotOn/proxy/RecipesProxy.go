package proxy

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Recipe struct {
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

type Nutrient struct {
	Name                string  `json:"name"`
	Amount              float64 `json:"amount"`
	Unit                string  `json:"unit"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds"`
}

type Property struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Flavonoid struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Ingredient struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Amount    float64    `json:"amount"`
	Unit      string     `json:"unit"`
	Nutrients []Nutrient `json:"nutrients"`
}

type CaloricBreakdown struct {
	PercentProtein float64 `json:"percentProtein"`
	PercentFat     float64 `json:"percentFat"`
	PercentCarbs   float64 `json:"percentCarbs"`
}

type WeightPerServing struct {
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type RecipeDetails struct {
	Nutrients        []Nutrient       `json:"nutrients"`
	Properties       []Property       `json:"properties"`
	Flavonoids       []Flavonoid      `json:"flavonoids"`
	Ingredients      []Ingredient     `json:"ingredients"`
	CaloricBreakdown CaloricBreakdown `json:"caloricBreakdown"`
	WeightPerServing WeightPerServing `json:"weightPerServing"`
}

type RecipesProxy struct {
}

func (rp *RecipesProxy) GetRecipesByIngredients(url string) (*[]Recipe, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var recipes []Recipe
	err = json.Unmarshal(body, &recipes)

	if err != nil {
		return nil, err
	}

	return &recipes, nil
}

func (rp *RecipesProxy) GetNutritionByRecipeId(id int) (*RecipeDetails, error) {
	url := "https://api.spoonacular.com/recipes/" + strconv.Itoa(id) + "/nutritionWidget.json?apiKey=8da80267f2bc4e3e81762e459bc4590d"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var recipeDetails RecipeDetails
	err = json.Unmarshal(body, &recipeDetails)

	if err != nil {
		return nil, err
	}

	return &recipeDetails, nil
}
