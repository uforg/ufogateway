package db

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type Route struct {
	ID                   string    `db:"id" json:"id"`
	Project              string    `db:"project" json:"project"`
	Name                 string    `db:"name" json:"name"`
	Active               bool      `db:"active" json:"active"`
	Endpoint             string    `db:"endpoint" json:"endpoint"`
	OriginURL            string    `db:"origin_url" json:"origin_url"`
	StoreHits            bool      `db:"store_hits" json:"store_hits"`
	StoreReqHeaders      bool      `db:"store_req_headers" json:"store_req_headers"`
	StoreReqBody         bool      `db:"store_req_body" json:"store_req_body"`
	StoreReqBodyMaxBytes int       `db:"store_req_body_max_bytes" json:"store_req_body_max_bytes"`
	StoreResHeaders      bool      `db:"store_res_headers" json:"store_res_headers"`
	StoreResBody         bool      `db:"store_res_body" json:"store_res_body"`
	StoreResBodyMaxBytes int       `db:"store_res_body_max_bytes" json:"store_res_body_max_bytes"`
	RetentionDays        int       `db:"retention_days" json:"retention_days"`
	RetentionHits        int       `db:"retention_hits" json:"retention_hits"`
	TLSClientCert        string    `db:"tls_client_cert" json:"tls_client_cert"`
	TLSClientKey         string    `db:"tls_client_key" json:"tls_client_key"`
	TLSCaCert            string    `db:"tls_ca_cert" json:"tls_ca_cert"`
	TLSSkipCertVerify    bool      `db:"tls_skip_cert_verify" json:"tls_skip_cert_verify"`
	Created              time.Time `db:"created" json:"created"`
	Updated              time.Time `db:"updated" json:"updated"`
}

func NewRouteFromRecord(r *core.Record) Route {
	return Route{
		ID:                   r.Id,
		Project:              r.GetString("project"),
		Name:                 r.GetString("name"),
		Active:               r.GetBool("active"),
		Endpoint:             r.GetString("endpoint"),
		OriginURL:            r.GetString("origin_url"),
		StoreHits:            r.GetBool("store_hits"),
		StoreReqHeaders:      r.GetBool("store_req_headers"),
		StoreReqBody:         r.GetBool("store_req_body"),
		StoreReqBodyMaxBytes: r.GetInt("store_req_body_max_bytes"),
		StoreResHeaders:      r.GetBool("store_res_headers"),
		StoreResBody:         r.GetBool("store_res_body"),
		StoreResBodyMaxBytes: r.GetInt("store_res_body_max_bytes"),
		RetentionDays:        r.GetInt("retention_days"),
		RetentionHits:        r.GetInt("retention_hits"),
		TLSClientCert:        r.GetString("tls_client_cert"),
		TLSClientKey:         r.GetString("tls_client_key"),
		TLSCaCert:            r.GetString("tls_ca_cert"),
		TLSSkipCertVerify:    r.GetBool("tls_skip_cert_verify"),
		Created:              r.GetDateTime("created").Time(),
		Updated:              r.GetDateTime("updated").Time(),
	}
}

func (db *DB) GetRoutesFromDB() ([]Route, error) {
	records, err := db.app.FindRecordsByFilter(
		"routes",
		"active = true",
		"created",
		99999,
		0,
	)
	if err != nil {
		return nil, err
	}

	routes := []Route{}
	for _, record := range records {
		routes = append(routes, NewRouteFromRecord(record))
	}

	return routes, nil
}

func (db *DB) GetRoutesFromDBCached() ([]Route, error) {
	key := "db.GetRoutesFromDBCached"

	cachedRoutes, found := db.cacheInstance.Get(key)
	if found {
		return cachedRoutes.([]Route), nil
	}

	dbRoutes, err := db.GetRoutesFromDB()
	if err != nil {
		return nil, err
	}

	db.cacheInstance.Set(key, dbRoutes, 5*time.Second)
	return dbRoutes, nil
}

func (db *DB) GetRouteByID(routeID string) (Route, error) {
	record, err := db.app.FindRecordById("routes", routeID)
	if err != nil {
		return Route{}, err
	}

	return NewRouteFromRecord(record), nil
}

func (db *DB) GetRouteByIDCached(routeID string) (Route, error) {
	key := "db.GetRouteByIDCached." + routeID

	cachedRoute, found := db.cacheInstance.Get(key)
	if found {
		return cachedRoute.(Route), nil
	}

	dbRoute, err := db.GetRouteByID(routeID)
	if err != nil {
		return Route{}, err
	}

	db.cacheInstance.Set(key, dbRoute, 5*time.Second)
	return dbRoute, nil
}
