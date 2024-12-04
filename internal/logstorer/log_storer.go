package logstorer

import (
	"io"

	"github.com/pocketbase/pocketbase"
	"github.com/uforg/ufogateway/internal/db"
	"github.com/uforg/ufogateway/internal/gateway"
)

type LogStorer struct {
	app *pocketbase.PocketBase
	db  *db.DB
}

func NewLogStorer(
	app *pocketbase.PocketBase,
	db *db.DB,
) *LogStorer {
	return &LogStorer{
		app: app,
		db:  db,
	}
}

func (ls *LogStorer) StoreRequestLog(reqLog gateway.RequestLog) {
	route, err := ls.db.GetRouteByIDCached(reqLog.RouteID)
	if err != nil {
		ls.app.Logger().Error(
			"failed to get route by id",
			"id", reqLog.RouteID,
			"fn", "StoreRequestLog",
			"error", err,
		)
		return
	}

	if !route.Active || !route.StoreHits {
		return
	}

	err = ls.db.CreateRequest(
		reqLog.RequestID,
		reqLog.RouteID,
		reqLog.Timestamp,
		reqLog.RequestIP,
		reqLog.RequestMethod,
		reqLog.RequestGatewayURL,
		reqLog.RequestOriginURL,
	)
	if err != nil {
		ls.app.Logger().Error(
			"failed to create request",
			"fn", "StoreRequestLog",
			"route_id", reqLog.RouteID,
			"request_id", reqLog.RequestID,
			"error", err,
		)
		return
	}

	if route.StoreReqHeaders {
		err = ls.db.StoreRequestReqHeaders(reqLog.RequestID, reqLog.RequestHeaders)
		if err != nil {
			ls.app.Logger().Error(
				"failed to store request request headers",
				"fn", "StoreRequestLog",
				"route_id", reqLog.RouteID,
				"request_id", reqLog.RequestID,
				"error", err,
			)
		}
	}

	if route.StoreReqBody {
		func() {
			bodyBytes, err := io.ReadAll(reqLog.RequestBody)
			if err != nil {
				ls.app.Logger().Error(
					"failed to read request body",
					"fn", "StoreRequestLog",
					"route_id", reqLog.RouteID,
					"request_id", reqLog.RequestID,
					"error", err,
				)
				return
			}

			if route.StoreReqBodyMaxBytes > 0 && len(bodyBytes) > route.StoreReqBodyMaxBytes {
				return
			}

			err = ls.db.StoreRequestReqBody(reqLog.RequestID, string(bodyBytes))
			if err != nil {
				ls.app.Logger().Error(
					"failed to store request request body",
					"fn", "StoreRequestLog",
					"route_id", reqLog.RouteID,
					"request_id", reqLog.RequestID,
					"error", err,
				)
			}
		}()
	}
}

func (ls *LogStorer) StoreResponseLog(reqLog gateway.ResponseLog) {
	route, err := ls.db.GetRouteByIDCached(reqLog.RouteID)
	if err != nil {
		ls.app.Logger().Error(
			"failed to get route by id",
			"id", reqLog.RouteID,
			"fn", "StoreResponseLog",
			"error", err,
		)
		return
	}

	if !route.Active || !route.StoreHits {
		return
	}

	err = ls.db.StoreRequestResTimestamp(
		reqLog.RequestID,
		reqLog.Timestamp,
		reqLog.Duration,
	)
	if err != nil {
		ls.app.Logger().Error(
			"failed to store request response timestamp",
			"fn", "StoreResponseLog",
			"route_id", reqLog.RouteID,
			"request_id", reqLog.RequestID,
			"error", err,
		)
	}

	if route.StoreResHeaders {
		err = ls.db.StoreRequestResHeaders(reqLog.RequestID, reqLog.ResponseHeaders)
		if err != nil {
			ls.app.Logger().Error(
				"failed to store request response headers",
				"fn", "StoreResponseLog",
				"route_id", reqLog.RouteID,
				"request_id", reqLog.RequestID,
				"error", err,
			)
		}
	}

	if route.StoreResBody {
		func() {
			bodyBytes, err := io.ReadAll(reqLog.ResponseBody)
			if err != nil {
				ls.app.Logger().Error(
					"failed to read response body",
					"fn", "StoreResponseLog",
					"route_id", reqLog.RouteID,
					"request_id", reqLog.RequestID,
					"error", err,
				)
				return
			}

			if route.StoreResBodyMaxBytes > 0 && len(bodyBytes) > route.StoreResBodyMaxBytes {
				return
			}

			err = ls.db.StoreRequestResBody(reqLog.RequestID, string(bodyBytes))
			if err != nil {
				ls.app.Logger().Error(
					"failed to store request response body",
					"fn", "StoreResponseLog",
					"route_id", reqLog.RouteID,
					"request_id", reqLog.RequestID,
					"error", err,
				)
			}
		}()
	}
}
