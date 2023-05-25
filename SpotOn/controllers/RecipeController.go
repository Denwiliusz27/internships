package controllers

import (
	"SpotOn/models"
	"SpotOn/proxy"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type RecipesController struct {
}

func ConnectToDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"))

	if err != nil {
		panic("Failed to connect with database")
	}

	err = db.AutoMigrate(&models.RecipeDB{}, &models.PresentIngredient{}, &models.MissingIngredient{})
	if err != nil {
		return nil
	}

	return db
}

func ifListContainsSubstring(ingredientsList []string, presentIngredients []models.PresentIngredient) bool {
	presentnNr := 0

	for _, presentIngredient := range presentIngredients {
		for _, ingredient := range ingredientsList {
			if strings.Contains(presentIngredient.Name, ingredient) {
				presentnNr += 1
				break
			}
		}
	}

	if presentnNr == len(presentIngredients) {
		nr := 0
		for _, ingredient := range ingredientsList {
			for _, presentIngredient := range presentIngredients {
				if strings.Contains(presentIngredient.Name, ingredient) {
					nr += 1
					break
				}
			}
		}

		if nr == len(ingredientsList) {
			return true
		}
	}

	return false
}

func getRecipesFromDatabase(ingredientsList []string, db *gorm.DB) []models.RecipeDB {
	var db_recipes []models.RecipeDB

	query := db.Table("recipes").
		Preload("PresentIngredients").
		Preload("MissingIngredients").
		Joins("INNER JOIN present_ingredients ON present_ingredients.recipe_id = recipes.id").
		Where("present_ingredients.name LIKE ?", "%"+ingredientsList[0]+"%")

	for i := 1; i < len(ingredientsList); i++ {
		query = query.Or("present_ingredients.name LIKE ?", "%"+ingredientsList[i]+"%")
	}

	query = query.Distinct().Find(&db_recipes)

	var final_recipes []models.RecipeDB

	if len(db_recipes) > 0 {
		println("przed: " + strconv.Itoa(len(db_recipes)))

		for _, recipe := range db_recipes {
			if ifListContainsSubstring(ingredientsList, recipe.PresentIngredients) {
				final_recipes = append(final_recipes, recipe)
			}
		}

		println("Po: " + strconv.Itoa(len(final_recipes)))

		return final_recipes
	}

	return nil
}

func parseRecipeDBToRecipe(recipeDB models.RecipeDB) models.Recipe {
	newRecipe := models.Recipe{}
	newRecipe.Name = recipeDB.Name
	newRecipe.Carbs = recipeDB.Carbs
	newRecipe.Proteins = recipeDB.Proteins
	newRecipe.Calories = recipeDB.Calories

	for _, presentRecipe := range recipeDB.PresentIngredients {
		newRecipe.PresentIngredients = append(newRecipe.PresentIngredients, presentRecipe.Name)
	}

	for _, missingRecipe := range recipeDB.MissingIngredients {
		newRecipe.MissingIngredients = append(newRecipe.MissingIngredients, missingRecipe.Name)
	}

	return newRecipe
}

func parseRecipeExtendedToRecipe(recipeExtended models.RecipeExtended) models.Recipe {
	var newRecipe models.Recipe

	recipesProxy := proxy.RecipesProxy{}
	newRecipe.Name = recipeExtended.Title
	recipeDetails, err := recipesProxy.GetNutritionByRecipeId(recipeExtended.ID)

	if err != nil {
		return models.Recipe{}
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

	for _, presentIngredient := range recipeExtended.UsedIngredients {
		newRecipe.PresentIngredients = append(newRecipe.PresentIngredients, presentIngredient.Name)
	}

	for _, missedIngredient := range recipeExtended.MissedIngredients {
		newRecipe.MissingIngredients = append(newRecipe.MissingIngredients, missedIngredient.Name)
	}

	return newRecipe
}

func addRecipeToDatabase(recipe models.Recipe, recipeExtended models.RecipeExtended, db *gorm.DB) {
	tx := db.Begin()

	db_recipe := models.RecipeDB{
		Name:     recipe.Name,
		Carbs:    recipe.Carbs,
		Proteins: recipe.Proteins,
		Calories: recipe.Calories,
	}

	tx.Create(&db_recipe)
	for _, presentIngredient := range recipe.PresentIngredients {
		ingredient := models.PresentIngredient{
			RecipeID: db_recipe.ID,
			Name:     presentIngredient,
		}
		tx.Create(&ingredient)
	}

	for _, missedIngredient := range recipeExtended.MissedIngredients {
		ingredient := models.MissingIngredient{
			RecipeID: db_recipe.ID,
			Name:     missedIngredient.Name,
		}
		tx.Create(&ingredient)
	}

	tx.Commit()
}

func (rC *RecipesController) GetRecipesByIngredients(ingredientsList []string, numberOfRecipes int) ([]models.Recipe, error) {
	recipesProxy := proxy.RecipesProxy{}
	var recipes []models.Recipe
	var newRecipe models.Recipe

	db := ConnectToDatabase()
	db_recipes := getRecipesFromDatabase(ingredientsList, db)

	if len(db_recipes) >= numberOfRecipes {
		println("-- biore z bazy --")

		for nr := 0; nr < numberOfRecipes; nr++ {
			recipes = append(recipes, parseRecipeDBToRecipe(db_recipes[nr]))
		}

		return recipes, nil

	} else {
		println("--- pobieram nowe -- ")

		url := createIngredientsUrlLink(ingredientsList, numberOfRecipes)

		recipesExtended, err := recipesProxy.GetRecipesByIngredients(url)
		if err != nil {
			return nil, err
		}

		for _, recipe := range *recipesExtended {
			newRecipe = parseRecipeExtendedToRecipe(recipe)
			recipes = append(recipes, newRecipe)
			addRecipeToDatabase(newRecipe, recipe, db)
		}

		return recipes, nil
	}
}

// creates url for getting recipes from ingredients list
func createIngredientsUrlLink(ingredientsList []string, numberOfRecipes int) string {
	url := "https://api.spoonacular.com/recipes/findByIngredients?apiKey=f661c070cf4f4ce480a75ff371a12b92&ingredients="

	for _, ingredient := range ingredientsList {
		url += ingredient + ","
	}

	url = strings.TrimSuffix(url, ",")
	url = url + "&number=" + strconv.Itoa(numberOfRecipes)

	return url
}
