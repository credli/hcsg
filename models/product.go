package models

import "time"

type Product struct {
	StockItemID           string
	FormulationID         string
	FormulaRef            string
	CategoryID            string
	SeqNumber             int
	BrandName             string
	ProductRef            int
	FullDescription       string
	MainPurpose           string
	Standards             string
	Dosage                string
	Color                 string
	PrototypeNumber       int
	DateOfFirstProduction time.Time
	DateOfLastProduction  time.Time
	ImagePath             string
	ImageAltText          string
	Published             bool
	PublishedDate         time.Time
}

func IsProductExist(pid string) (bool, error) {
	if len(pid) == 0 {
		return false, nil
	}
	db, err := GetDb()
	if err != nil {
		return false, err
	}
	rows, err := db.Query("SELECT * FROM Products WHERE ID = ?;", pid)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func GetProductByID(pid string) (*Product, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT StockItemID, FormulationID, FormulaRef, CategoryID, SeqNumber, BrandName, ProductRef,
		FullDescription, MainPurpose, Standards, Dosage, Color, PrototypeNumber, DateOfFirstProduction, DateOfLastProduction,
		ImagePath, ImageAltText, Published, PublishedDate
		FROM Products
		WHERE ID = ?`, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := new(Product)
		err = rows.Scan(&p.StockItemID, &p.FormulationID, &p.FormulaRef, &p.CategoryID, &p.SeqNumber, &p.BrandName,
			&p.ProductRef, &p.FullDescription, &p.MainPurpose, &p.Standards, &p.Dosage, &p.Color, &p.PrototypeNumber,
			&p.DateOfFirstProduction, &p.DateOfLastProduction, &p.ImagePath, &p.ImageAltText, &p.Published, &p.PublishedDate)
		if err != nil {
			return nil, err
		}
		return p, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrProductNotExist{pid, ""}
}

func GetAllProducts() ([]*Product, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT StockItemID, FormulationID, FormulaRef, CategoryID, SeqNumber, BrandName, ProductRef,
		FullDescription, MainPurpose, Standards, Dosage, Color, PrototypeNumber, DateOfFirstProduction, DateOfLastProduction,
		ImagePath, ImageAltText, Published, PublishedDate
		FROM Products ORDER BY BrandName`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allProducts = make([]*Product, 0)
	for rows.Next() {
		p := new(Product)
		err = rows.Scan(&p.StockItemID, &p.FormulationID, &p.FormulaRef, &p.CategoryID, &p.SeqNumber, &p.BrandName,
			&p.ProductRef, &p.FullDescription, &p.MainPurpose, &p.Standards, &p.Dosage, &p.Color, &p.PrototypeNumber,
			&p.DateOfFirstProduction, &p.DateOfLastProduction, &p.ImagePath, &p.ImageAltText, &p.Published, &p.PublishedDate)
		if err != nil {
			return nil, err
		}
		allProducts = append(allProducts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allProducts, nil
}
