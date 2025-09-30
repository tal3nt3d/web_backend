package serializer

import "web_backend/internal/app/ds"

type ApplicationDeviceJSON struct {
	AppDev_ID 		uint		`json:"app_dev_id"`
	Application_ID 	uint		`json:"application_id"`
	Device_ID      	uint		`json:"device_id"`
	Amount 			int			`json:"amount"`
	Notes 			string  	`json:"notes"`
	Amperage 		float64		`json:"amperage"`
}

func ApplicationDeviceToJSON (app_dev ds.ApplicationDevices) ApplicationDeviceJSON {
	return ApplicationDeviceJSON{
		AppDev_ID: 		app_dev.AppDev_ID,
		Application_ID: app_dev.Application_ID,
		Device_ID: 		app_dev.Device_ID,
		Amount: 		app_dev.Amount,
		Notes: 			app_dev.Notes,
		Amperage: 		app_dev.Amperage,
	}
}

func ApplicationDeviceFromJSON (applicationDeviceJSON ApplicationDeviceJSON) ds.ApplicationDevices {
	return ds.ApplicationDevices{
		Amount: 		applicationDeviceJSON.Amount,
		Notes: 			applicationDeviceJSON.Notes,
		Amperage: 		applicationDeviceJSON.Amperage,
	}
}