package recipes

import(

	"database/sql"
	"fmt"

)

// The Recipe entity is used to marshall/unmarshall JSON.
type Recipe struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	PrepTime   float32 `json:"preptime"`
	Difficulty int     `json:"difficulty"`
	Vegetarian bool    `json:"vegetarian"`
}


// UpdateRecipe is used to modify a specific recipe.
func (r *Recipe) UpdateRecipe(db *sql.DB) (res sql.Result, err error) {
	fmt.Printf("Entro aqui malditasea")
	res, err = db.Exec("UPDATE recipes SET name=$1, preptime=$2, difficulty=$3, vegetarian=$4 WHERE id=$5",
		r.Name, r.PrepTime, r.Difficulty, r.Vegetarian, r.ID)
	return res, err
}

// DeleteRecipe is used to delete a specific recipe.
func (r *Recipe) DeleteRecipe(db *sql.DB) (res sql.Result, err error) {
	res, err = db.Exec("DELETE FROM recipes WHERE id=$1", r.ID)
	return res, err
}

// CreateRecipe is used to create a single recipe.
func (r *Recipe) CreateRecipe(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO recipes(id, name, preptime, difficulty, vegetarian) VALUES($1, $2, $3, $4, $5) RETURNING id",
		r.ID, r.Name, r.PrepTime, r.Difficulty, r.Vegetarian).Scan(&r.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetRecipes returns a collection of known recipes.
func GetRecipes(db *sql.DB) ([]Recipe, error) {
	rows, err := db.Query(
		"SELECT id, name, preptime, difficulty, vegetarian FROM recipes")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	recipes := []Recipe{}
	for rows.Next() {
		var r Recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.PrepTime, &r.Difficulty, &r.Vegetarian); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}


func GetRecipesByName(db *sql.DB, name string) ([]Recipe, error) {

	query := "SELECT id, name, preptime, difficulty, vegetarian FROM recipes WHERE name LIKE $1"
	rows, err := db.Query(
		query, fmt.Sprintf("%%%s%%", name))

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	recipes := []Recipe{}
	for rows.Next() {
		var r Recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.PrepTime, &r.Difficulty, &r.Vegetarian); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}
