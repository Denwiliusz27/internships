package controllers

import (
	"SpotOn/models"
	"SpotOn/proxy"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"strings"
)

type RecipesController struct {
}

// GetRecipesByIngredients returns list of recipes for provided ingredients
func (rC *RecipesController) GetRecipesByIngredients(ingredientsList []string, numberOfRecipes int) ([]models.Recipe, error) {
	recipesProxy := proxy.RecipesProxy{}
	var recipes []models.Recipe
	var newRecipe models.Recipe
	db := ConnectToDatabase()

	db_recipes := getRecipesByIngredientsFromDatabase(ingredientsList, db)

	// if recipes are in database
	if len(db_recipes) >= numberOfRecipes {
		for nr := 0; nr < numberOfRecipes; nr++ {
			recipes = append(recipes, parseRecipeDBToRecipe(db_recipes[nr]))
		}

		return recipes, nil

	} else { // if recipes are not in database
		url := createIngredientsUrlLink(ingredientsList, numberOfRecipes)

		recipesExtended, err := recipesProxy.GetRecipesInfoByIngredients(url)
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

// ConnectToDatabase connects to database
func ConnectToDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  false,
			})})

	if err != nil {
		panic("Failed to connect with database")
	}

	err = db.AutoMigrate(&models.RecipeDB{}, &models.PresentIngredient{}, &models.MissingIngredient{})
	if err != nil {
		return nil
	}

	return db
}

// createIngredientsUrlLink creates url link for getting recipes from ingredients list
func createIngredientsUrlLink(ingredientsList []string, numberOfRecipes int) string {
	url := "https://api.spoonacular.com/recipes/findByIngredients?apiKey=f661c070cf4f4ce480a75ff371a12b92&ranking=1&ingredients="

	for _, ingredient := range ingredientsList {
		url += ingredient + ","
	}

	url = strings.TrimSuffix(url, ",")
	url = url + "&number=" + strconv.Itoa(numberOfRecipes)

	return url
}

// presentIngredientsExualsIngredients checks if present ingredients list is same as ingredients list
func presentIngredientsExualsIngredients(ingredientsList []string, presentIngredients []models.PresentIngredient) bool {
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

// getRecipesByIngredientsFromDatabase returns list of recipes for ingredients list
func getRecipesByIngredientsFromDatabase(ingredientsList []string, db *gorm.DB) []models.RecipeDB {
	var db_recipes []models.RecipeDB
	var final_recipes []models.RecipeDB

	query := db.Table("recipes").
		Preload("PresentIngredients").
		Preload("MissingIngredients").
		Joins("INNER JOIN present_ingredients ON present_ingredients.recipe_id = recipes.id").
		Where("present_ingredients.name LIKE ?", "%"+ingredientsList[0]+"%")

	for i := 1; i < len(ingredientsList); i++ {
		query = query.Or("present_ingredients.name LIKE ?", "%"+ingredientsList[i]+"%")
	}

	query = query.Distinct().Find(&db_recipes)

	if len(db_recipes) > 0 {
		for _, recipe := range db_recipes {
			if presentIngredientsExualsIngredients(ingredientsList, recipe.PresentIngredients) {
				final_recipes = append(final_recipes, recipe)
			}
		}

		return final_recipes
	}

	return nil
}

// parseRecipeDBToRecipe creates Recipe object from RecipeDB object from database
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

// parseRecipeExtendedToRecipe creates Recipe object from RecipeExtended object
func parseRecipeExtendedToRecipe(recipeExtended models.RecipeExtended) models.Recipe {
	var newRecipe models.Recipe
	recipesProxy := proxy.RecipesProxy{}

	recipeNutrientsInfo, err := recipesProxy.GetNutritionInfoByRecipeId(recipeExtended.ID)
	if err != nil {
		return models.Recipe{}
	}

	newRecipe.Name = recipeExtended.Title

	for _, nutrient := range recipeNutrientsInfo.Nutrients {
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

// addRecipeToDatabase adds recipe and corresponding ingredients to database
func addRecipeToDatabase(recipe models.Recipe, recipeExtended models.RecipeExtended, db *gorm.DB) {
	var db_recipe models.RecipeDB
	var exist bool

	err := db.Model(models.RecipeDB{}).
		Where("name = ? ", recipe.Name).
		Find(&exist).
		Error

	if err != nil {
		return
	}

	if !exist {
		tx := db.Begin()

		db_recipe = models.RecipeDB{
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
}
