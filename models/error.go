package models

import (
	"errors"
	"fmt"
)

var (
	ErrNameEmpty = errors.New("Name provided is empty")
)

type ErrNameReserved struct {
	Name string
}

func IsErrNameReserved(err error) bool {
	_, ok := err.(ErrNameReserved)
	return ok
}

func (e ErrNameReserved) Error() string {
	return fmt.Sprintf("Name '%s' is reserved by the system and cannot be used", e.Name)
}

type ErrNamePatternNotAllowed struct {
	Pattern string
}

func IsErrNamePatternNotAllowed(err error) bool {
	_, ok := err.(ErrNamePatternNotAllowed)
	return ok
}

func (e ErrNamePatternNotAllowed) Error() string {
	return fmt.Sprintf("Name pattern '%s' is not allowed", e.Pattern)
}

// Product error

type ErrProductNotExist struct {
	ID   string
	Name string
}

func IsErrProductNotExist(err error) bool {
	_, ok := err.(ErrProductNotExist)
	return ok
}

func (e ErrProductNotExist) Error() string {
	return fmt.Sprintf("Product doesn't exist [ID = %s, Name = %s]", e.ID, e.Name)
}

// Catalog error

type ErrCatalogNotExist struct {
	ID   string
	Name string
}

func IsErrCatalogNotExist(err error) bool {
	_, ok := err.(ErrCatalogNotExist)
	return ok
}

func (e ErrCatalogNotExist) Error() string {
	return fmt.Sprintf("Catalog doesn't exist [ID = %s, Name = %s]", e.ID, e.Name)
}

type ErrCatalogAlreadyExist struct {
	ID   string
	Name string
}

func IsErrCatalogAlreadyExist(err error) bool {
	_, ok := err.(ErrCatalogAlreadyExist)
	return ok
}

func (e ErrCatalogAlreadyExist) Error() string {
	return fmt.Sprintf("Catalog already exists [ID = %s, Name = %s]", e.ID, e.Name)
}

// RawMaterial error

type ErrRawMaterialNotExist struct {
	ID   string
	Name string
}

func IsErrRawMaterialNotExist(err error) bool {
	_, ok := err.(ErrRawMaterialNotExist)
	return ok
}

func (e ErrRawMaterialNotExist) Error() string {
	return fmt.Sprintf("Raw Material doesn't exist [ID = %s, Name = %s]", e.ID, e.Name)
}

// Recipe error

type ErrRecipeNotExist struct {
	ID      string
	Formula string
}

func IsErrRecipeNotExist(err error) bool {
	_, ok := err.(ErrRecipeNotExist)
	return ok
}

func (e ErrRecipeNotExist) Error() string {
	return fmt.Sprintf("Recipe doesn't exist [ID = %s, Formula = %s]", e.ID, e.Formula)
}

// User error

type ErrUserNotExist struct {
	UserID   string
	UserName string
}

func IsErrUserNotExist(err error) bool {
	_, ok := err.(ErrUserNotExist)
	return ok
}

func (e ErrUserNotExist) Error() string {
	return fmt.Sprintf("User doesn't exist [UserID = %s, UserName = %s]", e.UserID, e.UserName)
}
