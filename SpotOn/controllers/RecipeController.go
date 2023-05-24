package controllers

import (
	"SpotOn/models"
	"SpotOn/proxy"
	"fmt"
	"strconv"
	"strings"
)

type RecipesController struct {
}

func (rC *RecipesController) GetRecipesByIngredients(recipesList []string, numberOfRecipes int) ([]models.Recipe, error) {
	recipesProxy := proxy.RecipesProxy{}
	var recipes []models.Recipe
	var newRecipe models.Recipe

	url := createIngredientsUrlLink(recipesList, numberOfRecipes)

	recipesExtended, err := recipesProxy.GetRecipesByIngredients(url)
	if err != nil {
		return nil, err
	}

	for _, recipe := range *recipesExtended {
		newRecipe = models.Recipe{}
		newRecipe.Name = recipe.Title

		for _, presentIngredient := range recipe.UsedIngredients {
			newRecipe.PresentIngredients = append(newRecipe.PresentIngredients, presentIngredient.Name)
		}

		for _, missedIngredient := range recipe.MissedIngredients {
			newRecipe.MissingIngredients = append(newRecipe.MissingIngredients, missedIngredient.Name)
		}

		recipeDetails, err := recipesProxy.GetNutritionByRecipeId(recipe.ID)

		if err != nil {
			return nil, err
		}

		for _, nutrient := range recipeDetails.Nutrients {
			if nutrient.Name == "Carbohydrates" {
				newRecipe.Carbs = fmt.Sprintf("%.2f %s", nutrient.Amount, nutrient.Unit)
			} else if nutrient.Name == "Protein" {
				newRecipe.Proteins = fmt.Sprintf("%.2f %s", nutrient.Amount, nutrient.Unit)
			} else if nutrient.Name == "Calories" {
				newRecipe.Calories = fmt.Sprintf("%.2f %s", nutrient.Amount, nutrient.Unit)
			}
		}

		recipes = append(recipes, newRecipe)
	}

	return recipes, nil

}

// creates url for getting recipes from ingredients list
func createIngredientsUrlLink(ingredientsList []string, numberOfRecipes int) string {
	url := "https://api.spoonacular.com/recipes/findByIngredients?apiKey=8da80267f2bc4e3e81762e459bc4590d&ingredients="

	for _, ingredient := range ingredientsList {
		url += ingredient + ","
	}

	url = strings.TrimSuffix(url, ",")
	url = url + "&number=" + strconv.Itoa(numberOfRecipes)

	return url
}
