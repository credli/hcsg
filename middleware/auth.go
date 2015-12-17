package middleware

import (
	"fmt"
	"log"
	"net/url"

	"gopkg.in/macaron.v1"

	"github.com/credli/hcsg/base"
	"github.com/credli/hcsg/models"
	"github.com/credli/hcsg/settings"
)

type ToggleOptions struct {
	SignInRequired  bool
	SignOutRequired bool
	AdminRequired   bool
	DisableCsrf     bool
}

func AutoSignIn(ctx *Context) (bool, error) {
	if !models.Connected {
		return false, nil
	}

	uname := ctx.GetCookie(settings.CookieUserName)
	if len(uname) == 0 {
		return false, nil
	}

	succeeded := false
	defer func() {
		if !succeeded {
			log.Println("TRACE: auto-login cookie cleared: %s", uname)
			ctx.SetCookie(settings.CookieUserName, "", -1, settings.AppSubURL)
			ctx.SetCookie(settings.CookieRememberName, "", -1, settings.AppSubURL)
		}
	}()

	if ctx.Session.Get("uid") != nil {
		return true, nil
	}

	u, err := models.GetUserByName(uname)
	if err != nil {
		if !models.IsErrUserNotExist(err) {
			return false, fmt.Errorf("GetUserByName: %v", err)
		}
		return false, nil
	}

	var hash string
	switch u.PasswordFormat {
	case "0":
		hash, _ = ctx.GetSuperSecureCookie(base.EncodeMD5(u.Password), settings.CookieRememberName)
	case "1":
		hash, _ = ctx.GetSuperSecureCookie(base.EncodeMD5(u.PasswordSalt+u.Password), settings.CookieRememberName)
	}
	if hash != u.UserName {
		return false, nil
	}
	succeeded = true
	ctx.Session.Set("uid", u.UserID)
	ctx.Session.Set("uname", u.UserName)

	return succeeded, nil
}

func Toggle(options *ToggleOptions) macaron.Handler {
	return func(ctx *Context) {
		if options.SignInRequired {
			if !ctx.SignedIn {
				ctx.SetCookie("redirect_to", url.QueryEscape(settings.AppSubURL+ctx.Req.RequestURI), 0, settings.AppSubURL)
				ctx.Redirect(settings.AppSubURL + "/user/login")
				return
			} else if ctx.User.IsLockedOut {
				ctx.Data["Title"] = "Your account is locked out"
				ctx.HTML(200, "user/auth/locked")
				return
			}
		}

		if !options.SignOutRequired && !ctx.SignedIn {
			success, err := AutoSignIn(ctx)
			if err != nil {
				ctx.Handle(500, "AutoSignIn", err)
				return
			} else if success {
				ctx.Redirect(settings.AppSubURL + ctx.Req.RequestURI)
				return
			}
		}

		if options.AdminRequired {
			if !ctx.User.IsAdmin {
				ctx.Error(403)
				return
			}
			ctx.Data["PageIsAdmin"] = true
		}
	}
}
