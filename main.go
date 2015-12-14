package main

import (
	"crypto/tls"
	"fmt"
	gotmpl "html/template"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"gopkg.in/macaron.v1"

	"github.com/credli/hcsg/auth"
	"github.com/credli/hcsg/middleware"
	"github.com/credli/hcsg/models"
	"github.com/credli/hcsg/routers"
	"github.com/credli/hcsg/settings"
	"github.com/credli/hcsg/template"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "Holderchem Source Guide"
	app.Version = settings.AppVer + "\n(" + settings.BuildTime + ")\nCommit: " + settings.BuildGitHash
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "web",
			Usage:  "Starts HC Source Guide web server",
			Action: runWeb,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Value: "8881",
					Usage: "Override port number",
				},
				cli.StringFlag{
					Name:  "host",
					Value: "localhost",
					Usage: "Override host address",
				},
				cli.StringFlag{
					Name:  "config",
					Value: "custom/app.ini",
					Usage: "Override default configuration file path",
				},
			},
		},
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}

func runWeb(ctx *cli.Context) {
	routers.GlobalInit()

	if ctx.IsSet("config") {
		settings.CustomConf = ctx.String("config")
	}
	appURL := fmt.Sprintf("%s:%s", settings.HttpAddr, settings.HttpPort)
	// override settings
	if ctx.IsSet("port") {
		appURL = strings.Replace(appURL, settings.HttpPort, ctx.String("port"), 1)
		settings.HttpPort = ctx.String("port")
	}
	if ctx.IsSet("host") {
		appURL = strings.Replace(appURL, settings.HttpAddr, ctx.String("host"), 1)
		settings.HttpAddr = ctx.String("host")
	}

	m := newMacaron()

	reqSignIn := middleware.Toggle(&middleware.ToggleOptions{SignInRequired: true})
	//ignSignIn := middleware.Toggle(&middleware.ToggleOptions{SignInRequired: false})
	//ignSignInAndCsrf := middleware.Toggle(&middleware.ToggleOptions{DisableCsrf: true, SignInRequired: false})
	adminReq := middleware.Toggle(&middleware.ToggleOptions{SignInRequired: true, AdminRequired: true})

	//bind := binding.Bind
	bindIgnErr := binding.BindIgnErr

	m.Get("/", routers.Home)
	m.Group("/user", func() {
		m.Get("/login", routers.Login)
		m.Post("/login", bindIgnErr(auth.LoginForm{}), routers.LoginPost)
		m.Post("/logout", reqSignIn, routers.LogoutPost)
	})
	m.Get("/robots.txt", func(ctx *middleware.Context) {
		if settings.HasRobotsTxt {
			ctx.ServeFileContent(path.Join(settings.CustomPath, "robots.txt"))
		} else {
			ctx.Error(404)
		}
	})

	m.Get("/catalogs", adminReq, routers.CatalogList)

	m.NotFound(routers.NotFound)

	var err error
	log.Printf("Listening on %v://%s:%s%s\n", settings.Protocol, settings.HttpAddr, settings.HttpPort, settings.AppSubURL)
	switch settings.Protocol {
	case settings.HTTP:
		err = http.ListenAndServe(appURL, m)
	case settings.HTTPS:
		server := &http.Server{
			Addr: appURL,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS10,
			},
			Handler: m,
		}
		err = server.ListenAndServeTLS(settings.CertFile, settings.KeyFile)
	default:
		log.Panicf("Invalid protocol: %s", settings.Protocol)
	}

	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}

func newMacaron() *macaron.Macaron {
	m := macaron.New()
	if !settings.DisableRouterLog {
		m.Use(macaron.Logger())
	}
	m.Use(macaron.Recovery())
	m.Use(macaron.Static(settings.StaticRootPath, macaron.StaticOptions{
		Prefix:      "public",
		SkipLogging: true,
		// // Expires defines which user-defined function to use for producing a HTTP Expires Header. Default is nil.
		// // https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
		// Expires: func() string {
		// 	return time.Now().Add(24 * 60 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		// },
	}))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  path.Join(settings.StaticRootPath, "tmpl"),
		Funcs:      []gotmpl.FuncMap{template.Funcs},
		IndentJSON: macaron.Env != macaron.PROD,
	}))
	m.Use(session.Sessioner(settings.SessionConfig))
	m.Use(csrf.Csrfer(csrf.Options{
		Secret:     settings.SecretKey,
		SetCookie:  true,
		Header:     "X-Csrf-Token",
		CookiePath: settings.AppSubURL,
	}))
	m.Use(toolbox.Toolboxer(m, toolbox.Options{
		HealthCheckFuncs: []*toolbox.HealthCheckFuncDesc{
			&toolbox.HealthCheckFuncDesc{
				Desc: "Database connection",
				Func: models.Ping,
			},
		},
	}))
	m.Use(middleware.Contexter())
	return m
}
