package middleware

import (
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"

	"github.com/credli/hcsg/auth"
	"github.com/credli/hcsg/base"
	"github.com/credli/hcsg/models"
	"github.com/credli/hcsg/settings"
)

type Context struct {
	*macaron.Context
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	User     *models.User
	SignedIn bool
}

func (ctx *Context) Handle(status int, title string, err error) {
	if err != nil {
		log.Printf("ERROR: %s - %v\n", title, err)
		if macaron.Env != macaron.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}

	switch status {
	case 404:
		ctx.Data["Title"] = "Page Not Found"
	case 500:
		ctx.Data["Title"] = "Internal Server Error"
	}
	ctx.HTML(status, base.TplName(fmt.Sprintf("status/%d", status)))
}

func (ctx *Context) HandleText(status int, title string) {
	if (status/100 == 4) || (status/100 == 5) {
		log.Printf("ERROR: %s\n", title)
	}
	ctx.PlainText(status, []byte(title))
}

func (ctx *Context) HasError() bool {
	has, ok := ctx.Data["HasError"]
	if !ok {
		return false
	}
	ctx.Flash.ErrorMsg = ctx.Data["ErrorMsg"].(string)
	ctx.Data["Flash"] = ctx.Flash
	return has.(bool)
}

func (ctx *Context) HasValue(name string) bool {
	_, ok := ctx.Data[name]
	return ok
}

func (ctx *Context) GetErrorMsg() string {
	return ctx.Data["ErrorMsg"].(string)
}

// RenderWithErr used for page has form validation but need to prompt error to users.
func (ctx *Context) RenderWithErr(msg string, tpl base.TplName, form interface{}) {
	if form != nil {
		auth.AssignForm(form, ctx.Data)
	}
	ctx.Flash.ErrorMsg = msg
	ctx.Data["Flash"] = ctx.Flash
	ctx.HTML(200, tpl)
}

func (ctx *Context) HTML(status int, name base.TplName, data ...interface{}) {
	ctx.Context.HTML(status, string(name), data...)
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash, x csrf.CSRF) {
		ctx := &Context{
			Context: c,
			csrf:    x,
			Flash:   f,
			Session: sess,
		}

		ctx.Data["Link"] = settings.AppSubURL + strings.TrimSuffix(ctx.Req.URL.Path, "/")

		ctx.Data["PageStartTime"] = time.Now()

		// Get user from session if logged in
		ctx.User, _ = auth.SignedInUser(ctx.Context, ctx.Session)

		if ctx.User != nil {
			ctx.SignedIn = true
			ctx.Data["SignedIn"] = ctx.SignedIn
			ctx.Data["SignedUser"] = ctx.User
			ctx.Data["SingedUserID"] = ctx.User.UserID
			ctx.Data["SignedUserName"] = ctx.User.UserName
			ctx.Data["SignedDisplayName"] = ctx.User.DisplayName()
			ctx.Data["IsAdmin"] = ctx.User.IsAdmin
		} else {
			ctx.Data["SignedUserID"] = ""
			ctx.Data["SignedUserName"] = ""
		}

		// If request sends files, parse them here otherwise the Query() can't be parsed and the CsrfToken will be invalid.
		if ctx.Req.Method == "POST" && strings.Contains(ctx.Req.Header.Get("Content-Type"), "multipart/form-data") {
			if err := ctx.Req.ParseMultipartForm(settings.AttachmentMaxSize << 20); err != nil && !strings.Contains(err.Error(), "EOF") { // 32MB max size
				ctx.Handle(500, "ParseMultipartForm", err)
				return
			}
		}

		ctx.Data["CsrfToken"] = x.GetToken()
		ctx.Data["CsrfTokenHtml"] = template.HTML(`<input type="hidden" name="_csrf" value"` + x.GetToken() + `">`)

		c.Map(ctx)
	}
}
