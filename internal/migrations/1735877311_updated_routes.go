package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3090596648")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(15, []byte(`{
			"autogeneratePattern": "",
			"hidden": true,
			"id": "text45819390",
			"max": 0,
			"min": 0,
			"name": "tls_client_cert",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(16, []byte(`{
			"autogeneratePattern": "",
			"hidden": true,
			"id": "text1075645601",
			"max": 0,
			"min": 0,
			"name": "tls_client_key",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(17, []byte(`{
			"autogeneratePattern": "",
			"hidden": true,
			"id": "text1171539955",
			"max": 0,
			"min": 0,
			"name": "tls_ca_cert",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(18, []byte(`{
			"hidden": true,
			"id": "bool3690812205",
			"name": "tls_skip_cert_verify",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "bool"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3090596648")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text45819390")

		// remove field
		collection.Fields.RemoveById("text1075645601")

		// remove field
		collection.Fields.RemoveById("text1171539955")

		// remove field
		collection.Fields.RemoveById("bool3690812205")

		return app.Save(collection)
	})
}
