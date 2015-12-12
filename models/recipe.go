package models

type Recipe struct {
	ID                string
	Formula           string
	Date              string
	Superseeded       string
	NewFormula        string
	Product           string
	Category          string
	SubCategory       string
	OldRmName         string
	OldRmCode         string
	Percentage        float64
	Component         string
	NewRmCode         string
	NewClassification string
	NewGenericName    string
	NewTradeName      string
	CostAssumption    string
	Remarks           string
}

func IsRecipeExist(recipeID string) (bool, error) {
	if len(recipeID) == 0 {
		return false, nil
	}
	db, err := GetDb()
	if err != nil {
		return false, err
	}
	rows, err := db.Query("SELECT * FROM Recipes WHERE ID = ?;", recipeID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func GetRecipeByID(recipeID string) (*Recipe, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT ID, formula, date, superseeded, new_formula, product, category, sub_category,
		old_rm_name, old_rm_code, percentage, component, new_rm_code, new_classification, new_generic_name,
		new_trade_name, cost_assumption, remarks
		FROM Recipes
		WHERE ID = ?`, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := new(Recipe)
		err = rows.Scan(&r.ID, &r.Formula, &r.Date, &r.Superseeded, &r.NewFormula, &r.Product, &r.Category, &r.SubCategory,
			&r.OldRmName, &r.OldRmCode, &r.Percentage, &r.Component, &r.NewRmCode, &r.NewClassification, &r.NewGenericName,
			&r.NewTradeName, &r.CostAssumption, &r.Remarks)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrRecipeNotExist{recipeID, ""}
}

func GetAllRecipes() ([]*Recipe, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT ID, formula, date, superseeded, new_formula, product, category, sub_category,
		old_rm_name, old_rm_code, percentage, component, new_rm_code, new_classification, new_generic_name,
		new_trade_name, cost_assumption, remarks
		FROM Recipes ORDER BY formula`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allRecipes = make([]*Recipe, 0)
	for rows.Next() {
		r := new(Recipe)
		err = rows.Scan(&r.ID, &r.Formula, &r.Date, &r.Superseeded, &r.NewFormula, &r.Product, &r.Category, &r.SubCategory,
			&r.OldRmName, &r.OldRmCode, &r.Percentage, &r.Component, &r.NewRmCode, &r.NewClassification, &r.NewGenericName,
			&r.NewTradeName, &r.CostAssumption, &r.Remarks)
		if err != nil {
			return nil, err
		}
		allRecipes = append(allRecipes, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allRecipes, nil
}
