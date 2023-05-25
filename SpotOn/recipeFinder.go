package main

import (
	"SpotOn/controllers"
	"SpotOn/models"
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

// getIngredientsAndRecipesNumber sets ingredients list and number of recipes from arguments
func getIngredientsAndRecipesNumber() {
	rootCmd.Flags().StringVar(&ingredientsRaw, "ingredients", "", "Comma-separated list of ingredients")
	rootCmd.Flags().IntVar(&numberOfRecipes, "numberOfRecipes", 0, "Maximum number of recipes")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// displayAvailableRecipes displays recipes for provided ingredients
func displayAvailableRecipes(recipes []models.Recipe) {
	println("\nThere are available recipes to prepare from provided ingredients:\n")

	for i, recipe := range recipes {
		println("-------------- " + strconv.Itoa(i+1) + " --------------")
		println("Name: " + recipe.Name)
		print("Present ingredients: ")

		for _, presentIngredient := range recipe.PresentIngredients {
			print(presentIngredient + ", ")
		}

		print("\nMissing ingredients: ")

		for _, missedIngredient := range recipe.MissingIngredients {
			print(missedIngredient + ", ")
		}

		println("\nProteins: " + recipe.Proteins)
		println("Calories: " + recipe.Calories)
		fmt.Printf("Carbs: %s\n\n", recipe.Carbs)
	}
}

func main() {
	recipeController := controllers.RecipesController{}

	getIngredientsAndRecipesNumber()

	recipes, err := recipeController.GetRecipesByIngredients(ingredientsList, numberOfRecipes)
	if err != nil {
		return
	}

	displayAvailableRecipes(recipes)
}
