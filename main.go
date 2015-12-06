package main

import (
	"os"
	"runtime"
	"strings"
	"net/http"
	"github.com/codegangsta/cli"
	"gopkg.in/macaron.v1"
)

var (
	appVer = "1.0"
	rootPath = "/Users/Nicholas/Desktop/"
	appURL = "localhost:8080"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "Holderchem Source Guide"
	app.Version = appVer
	app.Commands = []cli.Command{
		cli.Command{
			Name: "web",
			Usage: "Starts HC Source Guide web server",
			Action: runWeb,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port, p",
					Value: "8881",
					Usage: "Override port number",
				},
				cli.StringFlag{
					Name: "host, h",
					Value: "localhost",
					Usage: "Override host address",
				},
			},
		},
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}

func runWeb(ctx *cli.Context) {
	if ctx.IsSet("port") {
		appURL = strings.Replace(appURL, "8080", ctx.String("port"), 1)
	}
	if ctx.IsSet("host") {
		appURL = strings.Replace(appURL, "localhost", ctx.String("host"), 1)
	}
	
	m := macaron.New()
	m.Use(macaron.Static("public", macaron.StaticOptions{
		Prefix: "public",
		SkipLogging: true,
		// // Expires defines which user-defined function to use for producing a HTTP Expires Header. Default is nil.
		// // https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
		// Expires: func() string { 
		// 	return time.Now().Add(24 * 60 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		// },
	}))
	
	http.ListenAndServe(appURL, m)
}