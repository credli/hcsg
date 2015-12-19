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
	tmplCatalogEdit   base.TplName = "catalogs/edit"
	tmplCatalogDelete base.TplName = "catalogs/delete"
)

type CatalogForm struct {
	ID          string `form:"id"`
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

func CreatePost(ctx *middleware.Context, form CatalogForm) {
	ctx.Data["PageIsCatalogs"] = true

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

	ctx.Flash.Success(fmt.Sprintf("Catalog %s was created", c.Name))
	ctx.Redirect("/catalogs", 302)
}

func Edit(ctx *middleware.Context) {
	ctx.Data["PageIsCatalogs"] = true

	id := ctx.Params(":catalogId")
	catalog, err := models.GetCatalogByID(id)
	if err != nil {
		ctx.Handle(500, "GetCatalogByID", err)
		return
	}
	ctx.Data["Catalog"] = catalog

	ctx.HTML(200, tmplCatalogEdit)
}

func EditPost(ctx *middleware.Context, form CatalogForm) {
	ctx.Data["PageIsCatalogs"] = true

	if ctx.HasError() {
		ctx.RenderWithErr("Something is wrong, check your entry then try again", tmplCatalogEdit, form)
		return
	}

	if len(form.ID) == 0 {
		ctx.Handle(500, "form.ID", fmt.Errorf("Something went wrong, please try editing it again."))
		return
	}

	c := &models.Catalog{
		ID:          form.ID,
		Name:        form.Name,
		Version:     form.Version,
		Description: form.Description,
		Printable:   form.Printable,
		Enabled:     form.Enabled,
	}
	err := models.UpdateCatalog(c)
	if err != nil {
		ctx.Handle(500, "UpdateCatalog", err)
		return
	}

	ctx.Flash.Success(fmt.Sprintf("Changes to catalog %s were saved", c.Name))
	ctx.Redirect(fmt.Sprintf("/catalogs/%s", c.ID), 302)
}

func DisablePost(ctx *middleware.Context) {
	catalogId := ctx.Params(":catalogId")
	if len(catalogId) == 0 {
		ctx.Handle(500, "catalogId", fmt.Errorf("catalogId is empty"))
		return
	}

	err := models.DisableCatalog(catalogId)
	if err != nil {
		ctx.Handle(500, "EnableCatalog", err)
		return
	}

	ctx.Redirect("/catalogs", 302)
}

func EnablePost(ctx *middleware.Context) {
	catalogId := ctx.Params(":catalogId")
	if len(catalogId) == 0 {
		ctx.Handle(500, "catalogId", fmt.Errorf("catalogId is empty"))
		return
	}

	err := models.EnableCatalog(catalogId)
	if err != nil {
		ctx.Handle(500, "EnableCatalog", err)
		return
	}

	ctx.Redirect("/catalogs", 302)
}

func Delete(ctx *middleware.Context) {
	catalogId := ctx.Params(":catalogId")
	if len(catalogId) == 0 {
		ctx.Handle(500, "catalogId", fmt.Errorf("catalogId is empty"))
		return
	}

	c, err := models.GetCatalogByID(catalogId)
	if err != nil {
		ctx.Handle(500, "GetCatalogByID", err)
		return
	}

	ctx.Data["Catalog"] = c
	ctx.HTML(200, tmplCatalogDelete)
}

func DeletePost(ctx *middleware.Context) {
	catalogId := ctx.Params(":catalogId")
	if len(catalogId) == 0 {
		ctx.Handle(500, "catalogId", fmt.Errorf("catalogId is empty"))
		return
	}

	c, err := models.GetCatalogByID(catalogId)
	if err != nil {
		ctx.Handle(500, "GetCatalogByID", err)
		return
	}
	name := c.Name

	err = models.DeleteCatalog(catalogId)
	if err != nil {
		ctx.Handle(500, "DeleteCatalog", err)
		return
	}

	ctx.Flash.Error(fmt.Sprintf("Catalog %s was deleted", name))
	ctx.Redirect("/catalogs", 302)
}
