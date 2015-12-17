package routers

import (
	"log"
	"net/url"

	"github.com/credli/hcsg/auth"
	"github.com/credli/hcsg/base"
	"github.com/credli/hcsg/middleware"
	"github.com/credli/hcsg/models"
	"github.com/credli/hcsg/settings"

	"gopkg.in/macaron.v1"
)

const (
	tmplHome  base.TplName = "home"
	tmplLogin base.TplName = "user/login"
)

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

	// Perform an auto-sign in if 'remember me' was selected
	success, err := middleware.AutoSignIn(ctx)
	if err != nil {
		ctx.Handle(500, "AutoSignIn", err)
		return
	}

	if success {
		if redirectTo, _ := url.QueryUnescape(ctx.GetCookie("redirect_to")); len(redirectTo) > 0 {
			ctx.SetCookie("redirect_to", "", -1, settings.AppSubURL)
			ctx.Redirect(redirectTo)
		} else {
			ctx.Redirect(settings.AppSubURL + "/")
		}
		return
	}

	ctx.HTML(200, tmplLogin)
}

func LoginPost(ctx *middleware.Context, form auth.LoginForm) {
	if ctx.HasError() {
		ctx.HTML(200, tmplLogin)
		ctx.RenderWithErr("", tmplLogin, form)
		return
	}
	log.Printf("form: %s %s %s", form.Username, form.Password, form.Remember)

	u, err := models.UserSignIn(form.Username, form.Password)
	log.Printf("Logged in: %v (err: %v)", u, err)
	if err != nil {
		if models.IsErrUserNotExist(err) {
			ctx.RenderWithErr("Incorrect username or password, please try again", tmplLogin, &form)
		} else {
			ctx.Handle(500, "UserLogin", err)
		}
		return
	}

	if form.Remember {
		days := 86400 * settings.LoginRememberDays
		ctx.SetCookie(settings.CookieUserName, u.UserName, days, settings.AppSubURL)
		// FIXME: should always require PasswordSalt when encoding to cookie, regardless of PasswordFormat
		if u.PasswordFormat == "0" {
			ctx.SetSuperSecureCookie(base.EncodeMD5(u.Password), settings.CookieRememberName, u.UserName, settings.AppSubURL)
		} else if u.PasswordFormat == "1" {
			ctx.SetSuperSecureCookie(base.EncodeMD5(u.PasswordSalt+u.Password), settings.CookieRememberName, u.UserName, days, settings.AppSubURL)
		}
	}

	ctx.Session.Set("user", u)
	ctx.Session.Set("uid", u.UserID)
	ctx.Session.Set("uname", u.UserName)
	if redirectTo, _ := url.QueryUnescape(ctx.GetCookie("redirect_to")); len(redirectTo) > 0 {
		ctx.SetCookie("redirect_to", "", -1, settings.AppSubURL)
		ctx.Redirect(redirectTo)
		return
	}

	ctx.Redirect(settings.AppSubURL + "/")
}

func Logout(ctx *middleware.Context) {
	ctx.Session.Delete("user")
	ctx.Session.Delete("uid")
	ctx.Session.Delete("uname")
	ctx.SetCookie(settings.CookieUserName, "", -1, settings.AppSubURL)
	ctx.SetCookie(settings.CookieRememberName, "", -1, settings.AppSubURL)
	ctx.Redirect(settings.AppSubURL + "/")
}

func NotFound(ctx *middleware.Context) {
	ctx.Data["Title"] = "Page Not Found"
	ctx.Handle(404, "Page Not Found", nil)
}
