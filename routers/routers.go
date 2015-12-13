package routers

import (
	"log"
	"net/url"

	"github.com/credli/hcsg/auth"
	"github.com/credli/hcsg/middleware"
	"github.com/credli/hcsg/models"
	"github.com/credli/hcsg/settings"

	"gopkg.in/macaron.v1"
)

const (
	tmplHome  = "home"
	tmplLogin = "login"

	//Catalogs
	tmplCatalogList = "catalogs/list"
)

type LoginForm struct {
	Username string `form:"name" binding:"Required"`
	Password string `form:"password" binding:"Required"`
}

func NewServices() {
	settings.NewServices()
}

func GlobalInit() {
	settings.NewContext()

	log.Printf("Custom path: %s\n", settings.CustomPath)
	log.Printf("Log path: %s\n", settings.LogRootPath)

	models.LoadConfigs()
	NewServices()
	if models.EnableSQLite3 {
		log.Println("SQLite Supported")
	}
	if models.EnableODBC {
		log.Println("ODBC Supported")
	}
	switch settings.Cfg.Section("").Key("RUN_MODE").String() {
	case "prod":
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
		settings.ProdMode = true
	}
}

func Home(ctx *middleware.Context) {
	ctx.Data["PageIsHome"] = true
	uname := ctx.GetCookie("uname")
	if len(uname) != 0 {
		ctx.Redirect("/login", 302)
		return
	}
	ctx.HTML(200, tmplHome)
}

func Login(ctx *middleware.Context) {
	ctx.Data["Title"] = "Login"

	success, err := middleware.AutoSignIn(ctx)
	if err != nil {
		ctx.Handle(500, "AutoSignIn", err)
		return
	}

	if success {
		if redirectTo, _ := url.QueryUnescape(ctx.GetCookie("redirect_to")); len(redirectTo) > 0 {
			ctx.SetCookie("redirect_to", "", -1, "/")
			ctx.Redirect(redirectTo)
			return
		}
		ctx.Redirect("/")
		return
	}

	ctx.HTML(200, tmplLogin)
}

func LoginPost(ctx *middleware.Context, form auth.LoginForm) {
	ctx.Data["Title"] = "Login"

}

func Logout(ctx *middleware.Context) {
	ctx.SetCookie(settings.CookieUserName, "", -1, "/")
	ctx.Redirect(settings.AppSubURL)
}

func CatalogList(ctx *middleware.Context) {
	ctx.Data["PageIsCatalogs"] = true
	ctx.HTML(200, tmplCatalogList)
}

func NotFound(ctx *middleware.Context) {
	ctx.Data["Title"] = "Page Not Found"
	ctx.Handle(404, "Page Not Found", nil)
}
