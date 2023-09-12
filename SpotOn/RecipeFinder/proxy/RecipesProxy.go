package proxy

import (
	"SpotOn/models"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type RecipesProxy struct {
}

func (rp *RecipesProxy) GetRecipesInfoByIngredients(url string) (*[]models.RecipeExtended, error) {
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

	var recipes []models.RecipeExtended
	err = json.Unmarshal(body, &recipes)

	if err != nil {
		return nil, err
	}

	return &recipes, nil
}

func (rp *RecipesProxy) GetNutritionInfoByRecipeId(id int) (*models.RecipeNutrientInfo, error) {
	url := "https://api.spoonacular.com/recipes/" + strconv.Itoa(id) + "/nutritionWidget.json?apiKey=f661c070cf4f4ce480a75ff371a12b92"

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

	var recipeDetails models.RecipeNutrientInfo
	err = json.Unmarshal(body, &recipeDetails)

	if err != nil {
		return nil, err
	}

	return &recipeDetails, nil
}
