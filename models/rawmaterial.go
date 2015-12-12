package models

type RawMaterial struct {
	ID              string
	OldName         string
	StockCount      int64
	Cost            int64
	InventoryName   string
	RmCode          string
	Classification  string
	NewGenericName  string
	NewTradeName    string
	NewRmcc         string
	NewSku          string
	CostAssumption  string
	TdsHyperlink    string
	ReplacementCost float64
	LastPurCost     float64
}

func IsRawMaterialExist(rmID string) (bool, error) {
	if len(rmID) == 0 {
		return false, nil
	}
	db, err := GetDb()
	if err != nil {
		return false, err
	}
	rows, err := db.Query("SELECT * FROM RawMaterials WHERE ID = ?;", rmID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func GetRawMaterialByID(rmID string) (*RawMaterial, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT ID, old_name, stock_count, cost, inventory_name, rm_code,
		classification, new_generic_name, new_trade_name, new_rmcc, new_sku, cost_assumption,
		tds_hyperlink, replacement_cost, last_pur_cost
		FROM RawMaterials
		WHERE ID = ?`, rmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rm := new(RawMaterial)
		err = rows.Scan(&rm.ID, &rm.OldName, &rm.StockCount, &rm.Cost, &rm.InventoryName,
			&rm.RmCode, &rm.Classification, &rm.NewGenericName, &rm.NewTradeName, &rm.NewRmcc, &rm.NewSku,
			&rm.CostAssumption, &rm.TdsHyperlink, &rm.ReplacementCost, &rm.LastPurCost)
		if err != nil {
			return nil, err
		}
		return rm, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrRawMaterialNotExist{rmID, ""}
}

func GetAllRawMaterials() ([]*RawMaterial, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT ID, old_name, stock_count, cost, inventory_name, rm_code,
		classification, new_generic_name, new_trade_name, new_rmcc, new_sku, cost_assumption,
		tds_hyperlink, replacement_cost, last_pur_cost
		FROM RawMaterials ORDER BY new_trade_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allRMs = make([]*RawMaterial, 0)
	for rows.Next() {
		rm := new(RawMaterial)
		err = rows.Scan(&rm.ID, &rm.OldName, &rm.StockCount, &rm.Cost, &rm.InventoryName,
			&rm.RmCode, &rm.Classification, &rm.NewGenericName, &rm.NewTradeName, &rm.NewRmcc, &rm.NewSku,
			&rm.CostAssumption, &rm.TdsHyperlink, &rm.ReplacementCost, &rm.LastPurCost)
		if err != nil {
			return nil, err
		}
		allRMs = append(allRMs, rm)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allRMs, nil
}
