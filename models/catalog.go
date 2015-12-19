package models

import (
	"time"

	"github.com/credli/hcsg/base"
)

type Catalog struct {
	ID          string
	AddedBy     string
	AddedDate   time.Time
	Name        string
	Version     string
	Description string
	Printable   bool
	Enabled     bool
}

func IsCatalogExist(cid, cname string) (bool, error) {
	if len(cid) == 0 {
		return false, nil
	}
	db, err := GetDb()
	if err != nil {
		return false, err
	}
	rows, err := db.Query("SELECT * FROM Catalogs WHERE ID = ? OR Name = ?;", cid, cname)
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
	rows, err := db.Query(`SELECT ID, AddedBy, AddedDate, Name, Version, Description, Printable, Enabled
		FROM Catalogs WHERE ID = ?`, cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := new(Catalog)
		err = rows.Scan(&c.ID, &c.AddedBy, &c.AddedDate, &c.Name, &c.Version, &c.Description, &c.Printable, &c.Enabled)
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
	rows, err := db.Query(`SELECT ID, AddedBy, AddedDate, Name, Version, Description, Printable, Enabled
		FROM Catalogs ORDER BY AddedDate DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allCatalogs = make([]*Catalog, 0)
	for rows.Next() {
		c := new(Catalog)
		err = rows.Scan(&c.ID, &c.AddedBy, &c.AddedDate, &c.Name, &c.Version, &c.Description, &c.Printable, &c.Enabled)
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

func CreateCatalog(c *Catalog, u *User) (err error) {
	if err = IsUsableName(c.Name); err != nil {
		return err
	}

	exists, err := IsCatalogExist(c.ID, c.Name)
	if err != nil {
		return err
	} else if exists {
		return ErrCatalogAlreadyExist{c.ID, c.Name}
	}

	uuid, err := base.GenerateUUID()
	if err != nil {
		return err
	}

	c.ID = uuid
	c.AddedBy = c.Name
	c.AddedDate = time.Now()

	db, err := GetDb()
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO Catalogs (ID, AddedBy, AddedDate, Name, Version, Description, Printable, Enabled)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?);`, c.ID, c.AddedBy, c.AddedDate, c.Name, c.Version, c.Description, c.Printable, c.Enabled)
	return err
}

func UpdateCatalog(c *Catalog) (err error) {
	if err = IsUsableName(c.Name); err != nil {
		return err
	}

	if len(c.ID) == 0 {
		return ErrEntityNotPersisted
	}

	db, err := GetDb()
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE Catalogs SET Name = ?, Version = ?, Description = ?, Printable = ?, Enabled = ? WHERE ID = ?;`,
		c.Name, c.Version, c.Description, c.Printable, c.Enabled, c.ID)
	return err
}

func EnableCatalog(id string) error {
	return toggleCatalogEnabled(id, true)
}

func DisableCatalog(id string) error {
	return toggleCatalogEnabled(id, false)
}

func toggleCatalogEnabled(id string, enabled bool) (err error) {
	if len(id) == 0 {
		return ErrEntityNotPersisted
	}

	db, err := GetDb()
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE Catalogs SET Enabled = ? WHERE ID = ?;`, enabled, id)
	return err
}

func DeleteCatalog(id string) error {
	if len(id) == 0 {
		return ErrEntityNotPersisted
	}

	db, err := GetDb()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM Catalogs WHERE ID = ?;`, id)
	return err
}
