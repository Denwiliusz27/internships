package models

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

type RecipeNutrientInfo struct {
	Nutrients        []Nutrient       `json:"nutrients"`
	Properties       []Property       `json:"properties"`
	Flavonoids       []Flavonoid      `json:"flavonoids"`
	Ingredients      []Ingredient     `json:"ingredients"`
	CaloricBreakdown CaloricBreakdown `json:"caloricBreakdown"`
	WeightPerServing WeightPerServing `json:"weightPerServing"`
}
