package models

import (
	"time"
)

type Catalog struct {
	ID          string
	AddedDate   time.Time
	Name        string
	Version     string
	Description string
	Enabled     bool
}

func IsCatalogExist(cid string) (bool, error) {
	if len(cid) == 0 {
		return false, nil
	}
	db, err := GetDb()
	if err != nil {
		return false, err
	}
	rows, err := db.Query("SELECT * FROM Catalogs WHERE ID = ?;", cid)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func GetCatalogByID(cid string) (*Catalog, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT ID, AddedDate, Name, Version, Description, Enabled
		FROM Catalogs WHERE ID = ?`, cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := new(Catalog)
		err = rows.Scan(&c.ID, &c.AddedDate, &c.Name, &c.Version, &c.Description, &c.Enabled)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nil, ErrCatalogNotExist{cid, ""}
}

func GetAllCatalogs() ([]*Catalog, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT ID, AddedDate, Name, Version, Description, Enabled
		FROM Catalogs ORDER BY AddedDate DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allCatalogs = make([]*Catalog, 0)
	for rows.Next() {
		c := new(Catalog)
		err = rows.Scan(&c.ID, &c.AddedDate, &c.Name, &c.Version, &c.Description, &c.Enabled)
		if err != nil {
			return nil, err
		}
		allCatalogs = append(allCatalogs, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allCatalogs, nil
}
