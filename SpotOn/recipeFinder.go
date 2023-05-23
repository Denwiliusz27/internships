package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	ingredientsRaw  string
	numberOfRecipes int
	ingredientsList []string

	rootCmd = &cobra.Command{
		Use:   "recipeFinder",
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

func getIngredientsAndRecipesNumber() {
	rootCmd.Flags().StringVar(&ingredientsRaw, "ingredients", "", "Comma-separated list of ingredients")
	rootCmd.Flags().IntVar(&numberOfRecipes, "numberOfRecipes", 0, "Maximum number of recipes")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	getIngredientsAndRecipesNumber()

	for i := 0; i < len(ingredientsList); i++ {
		println(ingredientsList[i])
	}
	println(numberOfRecipes)

}
