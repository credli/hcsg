package catalogs

import (
	"fmt"

	"github.com/credli/hcsg/base"
	"github.com/credli/hcsg/middleware"
	"github.com/credli/hcsg/models"
)

const (
	tmplCatalogList   base.TplName = "catalogs/list"
	tmplCatalogCreate base.TplName = "catalogs/create"
)

type CatalogCreateForm struct {
	Name        string `form:"name" binding:"Required"`
	Version     string `form:"version" binding:"Required"`
	Description string `form:"description"`
	Printable   bool   `form:"printable" binding:"Required"`
	Enabled     bool   `form:"enabled" binding:"Required"`
}

func List(ctx *middleware.Context) {
	ctx.Data["PageIsCatalogs"] = true

	catalogs, err := models.GetAllCatalogs()
	if err != nil {
		ctx.Handle(500, "GetAllCatalogs", err)
		return
	}
	ctx.Data["Catalogs"] = catalogs

	ctx.HTML(200, tmplCatalogList)
}

func Create(ctx *middleware.Context) {
	ctx.Data["PageIsCatalogs"] = true
	ctx.HTML(200, tmplCatalogCreate)
}

func CreatePost(ctx *middleware.Context, form CatalogCreateForm) {
	if ctx.HasError() {
		ctx.RenderWithErr("Something is wrong, check your entry then try again", tmplCatalogCreate, form)
		return
	}

	var user *models.User
	if user, ok := ctx.Data["SignedUser"].(*models.User); !ok || user == nil {
		ctx.Handle(500, "SignedUser", fmt.Errorf("Could not retrieve current user info"))
		return
	}

	c := &models.Catalog{
		Name:        form.Name,
		Version:     form.Version,
		Description: form.Description,
		Printable:   form.Printable,
		Enabled:     form.Enabled,
	}
	err := models.CreateCatalog(c, user)
	if err != nil {
		ctx.Handle(500, "CreateCatalog", err)
		return
	}

	ctx.Redirect("/catalogs", 302)
}
