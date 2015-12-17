package categories

import (
	"fmt"

	"github.com/credli/hcsg/base"
	"github.com/credli/hcsg/middleware"
	"github.com/credli/hcsg/models"
)

const (
	tmplCategoriesList base.TplName = "categories/list"
)

type CategoryCreateStruct struct {
	Name     string `form:"name" binding:"Required"`
	ParentID string `form:"parentId"`
	Order    int    `form:"order"`

	// DisplayAttributes
	Color       string `form:"color"`
	Description string `form:"description"`
}

func List(ctx *middleware.Context) {
	catalogId := ctx.Params(":catalogId")
	if len(catalogId) == 0 {
		ctx.Redirect("/catalogs")
	}

	catalog, err := models.GetCatalogByID(catalogId)
	if err != nil {
		ctx.Handle(500, "GetCatalogByID", err)
		return
	} else if catalog == nil {
		ctx.Handle(404, "GetCatalogByID", fmt.Errorf("Catalog with id '%s' was not found", catalogId))
		return
	}

	ctx.Data["Catalog"] = catalog
	ctx.HTML(200, tmplCategoriesList)
}

func Create(ctx *middleware.Context, form CategoryCreateForm) {
	if ctx.HasError() {
		ctx.RenderWithErr("Something is wrong, check your entry then try again", tmplCategoriesList, form)
		return
	}

	ctx.HTML(200, tmplCategoriesList)
}
