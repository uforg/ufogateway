package db

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
)

const requestsCollectionName = "requests"

func (db *DB) getRequestsCollection() (*core.Collection, error) {
	return db.app.FindCollectionByNameOrId(requestsCollectionName)
}

func (db *DB) CreateRequest(
	requestID string,
	routeID string,
	reqTimestamp time.Time,
	reqIP string,
	reqMethod string,
	reqGatewayURL string,
	reqOriginURL string,
) error {
	collection, err := db.getRequestsCollection()
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	record.Id = requestID
	record.Set("route", routeID)
	record.Set("req_timestamp", reqTimestamp)
	record.Set("req_ip", reqIP)
	record.Set("req_method", reqMethod)
	record.Set("req_gateway_url", reqGatewayURL)
	record.Set("req_origin_url", reqOriginURL)

	return db.app.Save(record)
}

func (db *DB) GetRequestRecordByID(requestID string) (*core.Record, error) {
	return db.app.FindRecordById(requestsCollectionName, requestID)
}

func (db *DB) StoreRequestReqHeaders(requestID string, reqHeaders map[string][]string) error {
	record, err := db.GetRequestRecordByID(requestID)
	if err != nil {
		return err
	}

	record.Set("req_headers", reqHeaders)

	return db.app.Save(record)
}

func (db *DB) StoreRequestReqBody(requestID string, reqBody string) error {
	record, err := db.GetRequestRecordByID(requestID)
	if err != nil {
		return err
	}

	record.Set("req_body", reqBody)

	return db.app.Save(record)
}

func (db *DB) StoreRequestResTimestamp(
	requestID string,
	resTimestamp time.Time,
	resDuration time.Duration,
) error {
	record, err := db.GetRequestRecordByID(requestID)
	if err != nil {
		return err
	}

	record.Set("res_timestamp", resTimestamp)
	record.Set("res_duration", resDuration.String())

	return db.app.Save(record)
}

func (db *DB) StoreRequestResHeaders(requestID string, resHeaders map[string][]string) error {
	record, err := db.GetRequestRecordByID(requestID)
	if err != nil {
		return err
	}

	record.Set("res_headers", resHeaders)

	return db.app.Save(record)
}

func (db *DB) StoreRequestResBody(requestID string, resBody string) error {
	record, err := db.GetRequestRecordByID(requestID)
	if err != nil {
		return err
	}

	record.Set("res_body", resBody)

	return db.app.Save(record)
}

func (db *DB) DeleteExpiredRequests() (int64, error) {
	res1, err := db.app.DB().
		NewQuery(`
			WITH requests_to_delete AS (
				SELECT requests.id
				FROM requests
				INNER JOIN routes ON routes.id = requests.route
				WHERE 
					routes.retention_days > 0
					AND requests.created < date('now', '-' || routes.retention_days || ' days')
			)
			DELETE FROM requests
			WHERE id IN (SELECT id FROM requests_to_delete);
		`).
		Execute()
	if err != nil {
		return 0, err
	}
	res1AffectedRows, err := res1.RowsAffected()
	if err != nil {
		return 0, err
	}

	res2, err := db.app.DB().
		NewQuery(`
			WITH ranked_requests AS (
				SELECT 
					requests.id,
					ROW_NUMBER() OVER (
						PARTITION BY requests.route 
						ORDER BY requests.created DESC
					) as rn,
					routes.retention_hits
				FROM requests
				INNER JOIN routes ON routes.id = requests.route
				WHERE routes.retention_hits > 0
			)
			DELETE FROM requests
			WHERE id IN (
				SELECT id 
				FROM ranked_requests 
				WHERE rn > retention_hits
			);
		`).
		Execute()
	if err != nil {
		return 0, err
	}
	res2AffectedRows, err := res2.RowsAffected()
	if err != nil {
		return 0, err
	}

	return res1AffectedRows + res2AffectedRows, nil
}
