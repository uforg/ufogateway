package main

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/uforg/ufogateway/internal/gateway"
	_ "github.com/uforg/ufogateway/internal/migrations"
)

func main() {
	if err := start(); err != nil {
		panic(err)
	}
}

type routeProvider struct{}

func (routeProvider) Routes() ([]gateway.Route, error) {
	r := []gateway.Route{}
	return r, nil
}

type logStorer struct{}

func (logStorer) StoreRequestLog(reqLog gateway.RequestLog) {
}

func (logStorer) StoreResponseLog(reqLog gateway.ResponseLog) {
}

func start() error {
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	rp := routeProvider{}
	ls := logStorer{}

	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate:  isGoRun,
		TemplateLang: migratecmd.TemplateLangGo,
		Dir:          "./internal/migrations",
	})

	gat := gateway.NewGateway(rp, ls)
	wrappedGat := apis.WrapStdHandler(gat)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.Any("/", wrappedGat)
		return se.Next()
	})

	return app.Start()
}
