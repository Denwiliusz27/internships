package models

type RecipeDB struct {
	ID                 uint
	Name               string
	Calories           string
	Proteins           string
	Carbs              string
	PresentIngredients []PresentIngredient `gorm:"foreignKey:RecipeID"`
	MissingIngredients []MissingIngredient `gorm:"foreignKey:RecipeID"`
}

type PresentIngredient struct {
	ID       uint
	RecipeID uint
	Name     string
}

type MissingIngredient struct {
	ID       uint
	RecipeID uint
	Name     string
}

//type RecipeDB struct {
//	gorm.Model
//	Name               string
//	Carbs              string
//	Proteins           string
//	Calories           string
//	PresentIngredients []IngredientDB `gorm:"many2many:present_ingredients;"`
//	MissingIngredients []IngredientDB `gorm:"many2many:missing_ingredients;"`
//}

//type IngredientDB struct {
//	gorm.Model
//	Name   string
//	Recipe []RecipeDB `gorm:"many2many:recipe_ingredient;"`
//}
//
//type PresentIngredient struct {
//	gorm.Model
//	IngredientId uint `gorm:"primaryKey" column:"ingredient_id"`
//	RecipeId     uint `gorm:"primaryKey" column:"recipe_id"`
//}
//
//type MissingIngredient struct {
//	gorm.Model
//	IngredientId uint `gorm:"primaryKey" column:"ingredient_id"`
//	RecipeId     uint `gorm:"primaryKey" column:"recipe_id"`
//}
//
//func (r *RecipeDB) TableName() string {
//	return "recipe"
//}
//
//func (i *IngredientDB) TableName() string {
//	return "ingredient"
//}
//
//func (pi *PresentIngredient) TableName() string {
//	return "p_ingredient"
//}
//
//func (mi *MissingIngredient) TableName() string {
//	return "m_ingredient"
//}

func (r *RecipeDB) TableName() string {
	return "recipes"
}
