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
	//DB *gorm.DB
}

func (rC *RecipesController) GetRecipesByIngredients(ingredientsList []string, numberOfRecipes int) ([]models.Recipe, error) {
	recipesProxy := proxy.RecipesProxy{}
	var recipes []models.Recipe
	var newRecipe models.Recipe

	db, err := gorm.Open(sqlite.Open("test.db"))

	if err != nil {
		panic("Failed to connect with database")
	}

	err = db.AutoMigrate(&models.RecipeDB{}, &models.PresentIngredient{}, &models.MissingIngredient{})
	if err != nil {
		return nil, err
	}

	var db_recipes []models.RecipeDB

	query := db.Table("recipes").
		Preload("PresentIngredients").
		Preload("MissingIngredients").
		Joins("INNER JOIN present_ingredients ON present_ingredients.recipe_id = recipes.id").
		Where("present_ingredients.name LIKE ?", ingredientsList[0])

	// Append additional LIKE conditions for each ingredient
	for i := 1; i < len(ingredientsList); i++ {
		query = query.Or("present_ingredients.name LIKE ?", ingredientsList[i])
	}

	query = query.Limit(numberOfRecipes)

	// Execute the query
	result := query.Find(&db_recipes)

	//db_recipes := []models.RecipeDB{}
	//db.Preload("PresentIngredients", "MissingIngredients").
	//	Where("present_ingredients.name LIKE ?", "%"+ingredientsList[0]+"%")
	//for i := 1; i < len(ingredientsList); i++ {
	//	db.Or("present_ingredients.name LIKE ?", "%"+ingredientsList[i]+"%")
	//}
	//result := db.Find(&db_recipes)

	if result.RowsAffected == int64(numberOfRecipes) {
		println("cos mam! : " + db_recipes[0].Name + ", " + db_recipes[0].Carbs + ", " + db_recipes[0].Proteins + ", " + db_recipes[0].Calories + ", ")

		for _, ing := range db_recipes[0].PresentIngredients {
			println("pres: " + ing.Name)
		}

		for _, ing := range db_recipes[0].MissingIngredients {
			println("miss: " + ing.Name)
		}

		return nil, nil

	} else {

		println("elo")

		url := createIngredientsUrlLink(ingredientsList, numberOfRecipes)

		recipesExtended, err := recipesProxy.GetRecipesByIngredients(url)
		if err != nil {
			return nil, err
		}

		for _, recipe := range *recipesExtended {
			newRecipe = models.Recipe{}
			newRecipe.Name = recipe.Title

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

			tx := db.Begin()

			db_recipe := models.RecipeDB{
				Name:     newRecipe.Name,
				Carbs:    newRecipe.Carbs,
				Proteins: newRecipe.Proteins,
				Calories: newRecipe.Calories,
			}

			tx.Create(&db_recipe)

			for _, presentIngredient := range recipe.UsedIngredients {
				newRecipe.PresentIngredients = append(newRecipe.PresentIngredients, presentIngredient.Name)

				ingredient := models.PresentIngredient{
					RecipeID: db_recipe.ID,
					Name:     presentIngredient.Name,
				}
				tx.Create(&ingredient)
			}

			for _, missedIngredient := range recipe.MissedIngredients {
				newRecipe.MissingIngredients = append(newRecipe.MissingIngredients, missedIngredient.Name)

				ingredient := models.MissingIngredient{
					RecipeID: db_recipe.ID,
					Name:     missedIngredient.Name,
				}
				tx.Create(&ingredient)
			}

			tx.Commit()
			recipes = append(recipes, newRecipe)

		}

		return recipes, nil
	}
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
