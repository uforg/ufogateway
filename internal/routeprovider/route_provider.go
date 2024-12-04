package routeprovider

import (
	"github.com/pocketbase/pocketbase"
	"github.com/uforg/ufogateway/internal/db"
	"github.com/uforg/ufogateway/internal/gateway"
)

type RouteProvider struct {
	app *pocketbase.PocketBase
	db  *db.DB
}

func NewRouteProvider(
	app *pocketbase.PocketBase,
	db *db.DB,
) *RouteProvider {
	return &RouteProvider{
		app: app,
		db:  db,
	}
}

func (rp *RouteProvider) Routes() ([]gateway.Route, error) {
	dbRoutes, err := rp.db.GetRoutesFromDBCached()
	if err != nil {
		return nil, err
	}

	routes := []gateway.Route{}
	for _, route := range dbRoutes {
		routes = append(routes, gateway.Route{
			ID:        route.ID,
			Endpoint:  route.Endpoint,
			OriginURL: route.OriginURL,
		})
	}

	return routes, nil
}
