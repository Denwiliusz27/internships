package main

import (
	"SpotOn/proxy"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var (
	ingredientsRaw  string
	numberOfRecipes int
	ingredientsList []string

	rootCmd = &cobra.Command{
		Use:   "Parsing arguments",
		Short: "Get ingredients list and number of recipes",
		RunE: func(cmd *cobra.Command, args []string) error {

			if ingredientsRaw == "" && numberOfRecipes <= 0 {
				return fmt.Errorf("please provide a list of ingredients using the --ingredients flag and positive " +
					"number of recipes using the --numberOfRecipes flag")
			}

			if ingredientsRaw == "" {
				return fmt.Errorf("please provide a list of ingredients using the --ingredients flag")
			}

			if numberOfRecipes <= 0 {
				return fmt.Errorf("please provide positive number of recipes using the --numberOfRecipes flag")
			}

			ingredientsList = strings.Split(ingredientsRaw, ",")
			return nil
		},
	}
)

// parse ingredints list and number of recipes
func getIngredientsAndRecipesNumber() {
	rootCmd.Flags().StringVar(&ingredientsRaw, "ingredients", "", "Comma-separated list of ingredients")
	rootCmd.Flags().IntVar(&numberOfRecipes, "numberOfRecipes", 0, "Maximum number of recipes")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// creates
func createIngredientsUrlLink() string {
	url := "https://api.spoonacular.com/recipes/findByIngredients?apiKey=8da80267f2bc4e3e81762e459bc4590d&ingredients="

	for i := 0; i < len(ingredientsList); i++ {
		url += ingredientsList[i]

		if i != len(ingredientsList)-1 {
			url += ","
		}
	}

	url = url + "&number=" + strconv.Itoa(numberOfRecipes)

	return url
}

func main() {
	recipesProxy := proxy.RecipesProxy{}

	getIngredientsAndRecipesNumber()

	for i := 0; i < len(ingredientsList); i++ {
		println(ingredientsList[i])
	}
	println(numberOfRecipes)

	url := createIngredientsUrlLink()

	recipes, err := recipesProxy.GetRecipesByIngredients(url)
	if err != nil {
		return
	}

	for _, recipe := range *recipes {
		println("----------------------------")
		println("Name: " + recipe.Title)
		print("Present ingredients: ")

		for _, presentIngredient := range recipe.UsedIngredients {
			print(presentIngredient.Name + ", ")
		}
		println()

		print("Missing ingredients: ")

		for _, missedIngredient := range recipe.MissedIngredients {
			print(missedIngredient.Name + ", ")
		}
		println()

		recipeDetails, err := recipesProxy.GetNutritionByRecipeId(recipe.ID)

		if err != nil {
			return
		}

		for _, nutrient := range recipeDetails.Nutrients {
			if nutrient.Name == "Carbohydrates" {
				fmt.Printf("Carbs: %.2f %s\n", nutrient.Amount, nutrient.Unit)
			} else if nutrient.Name == "Protein" {
				fmt.Printf("Proteins: %.2f %s\n", nutrient.Amount, nutrient.Unit)
			} else if nutrient.Name == "Calories" {
				fmt.Printf("Calories: %.2f %s\n", nutrient.Amount, nutrient.Unit)
			}
		}
	}
}
