package main

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/uforg/ufogateway/internal/cache"
	"github.com/uforg/ufogateway/internal/db"
	"github.com/uforg/ufogateway/internal/gateway"
	"github.com/uforg/ufogateway/internal/logstorer"
	_ "github.com/uforg/ufogateway/internal/migrations"
	"github.com/uforg/ufogateway/internal/routeprovider"
)

func main() {
	if err := start(); err != nil {
		panic(err)
	}
}

func start() error {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate:  isGoRun,
		TemplateLang: migratecmd.TemplateLangGo,
		Dir:          "./internal/migrations",
	})

	cacheInstance := cache.NewCacheInstance()
	db := db.NewDB(app, cacheInstance)

	routeProvider := routeprovider.NewRouteProvider(app, db)
	logStorer := logstorer.NewLogStorer(app, db)

	gat := gateway.NewGateway(routeProvider, logStorer)
	wrappedGat := apis.WrapStdHandler(gat)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.Any("/", wrappedGat)
		return se.Next()
	})

	return app.Start()
}
