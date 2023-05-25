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

func Connect() *gorm.DB {
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

	//print("List: ")
	//for _, i := range ingredientsList {
	//print(i + ",")
	//}

	//print("\nPresentIngrd: ")

	for _, presentIngredient := range presentIngredients {
		//print(presentIngredient.Name + ",")
		for _, ingredient := range ingredientsList {
			if strings.Contains(presentIngredient.Name, ingredient) {
				presentnNr += 1
				break
			}
		}
	}

	println()

	if presentnNr < len(presentIngredients) {
		return false
	} else {
		nr := 0
		for _, ingredient := range ingredientsList {
			for _, presentIngredient := range presentIngredients {
				if strings.Contains(presentIngredient.Name, ingredient) {
					nr += 1
					break
				}
			}
		}

		if nr < len(ingredientsList) {
			return false
		} else {
			return true
		}
	}
}

func searchDatabase(ingredientsList []string, numberOfRecipes int, db *gorm.DB) []models.RecipeDB {
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

	for _, recipe := range db_recipes {
		println("----------------------------")
		println("Name: " + recipe.Name)
		print("Present ingredients: ")

		for _, presentIngredient := range recipe.PresentIngredients {
			print(presentIngredient.Name + ", ")
		}
		println()

		print("Missing ingredients: ")

		for _, missedIngredient := range recipe.MissingIngredients {
			print(missedIngredient.Name + ", ")
		}
		println()

		println("Proteins: " + recipe.Proteins)
		println("Calories: " + recipe.Calories)
		fmt.Printf("Carbs: %s\n\n", recipe.Carbs)
	}

	if len(db_recipes) > 0 {
		println("-- PRZED: " + strconv.Itoa(len(db_recipes)))
		for _, recipe := range db_recipes {
			println("*****")

			if ifListContainsSubstring(ingredientsList, recipe.PresentIngredients) {
				final_recipes = append(final_recipes, recipe)
			}
		}

		println("-- PO: " + strconv.Itoa(len(final_recipes)))

		return final_recipes
	}

	return db_recipes
}

func (rC *RecipesController) GetRecipesByIngredients(ingredientsList []string, numberOfRecipes int) ([]models.Recipe, error) {
	recipesProxy := proxy.RecipesProxy{}
	var recipes []models.Recipe
	var newRecipe models.Recipe

	db := Connect()

	db_recipes := searchDatabase(ingredientsList, numberOfRecipes, db)

	if len(db_recipes) == numberOfRecipes {
		println("cos mam! : " + db_recipes[0].Name + ", " + db_recipes[0].Carbs + ", " + db_recipes[0].Proteins + ", " + db_recipes[0].Calories + ", ")

		for _, ing := range db_recipes[0].PresentIngredients {
			println("pres: " + ing.Name)
		}

		for _, ing := range db_recipes[0].MissingIngredients {
			println("miss: " + ing.Name)
		}

		return nil, nil

	} else {

		println("--- pobieram nowe -- ")

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
				println("NUTRIENT: " + nutrient.Name)
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
	url := "https://api.spoonacular.com/recipes/findByIngredients?apiKey=f661c070cf4f4ce480a75ff371a12b92&ingredients="

	for _, ingredient := range ingredientsList {
		url += ingredient + ","
	}

	url = strings.TrimSuffix(url, ",")
	url = url + "&number=" + strconv.Itoa(numberOfRecipes)

	return url
}
