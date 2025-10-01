package serializer

import "web_backend/internal/app/ds"

type AmperageApplicationDeviceJSON struct {
	AppDev_ID 				uint		`json:"app_dev_id"`
	Amperage_Application_ID uint		`json:"amperage_application_id"`
	Device_ID      			uint		`json:"device_id"`
	Amount 					int			`json:"amount"`
	Notes 					string  	`json:"notes"`
	Amperage 				float64		`json:"amperage"`
}

func AmperageApplicationDeviceToJSON (app_dev ds.AmperageApplicationDevices) AmperageApplicationDeviceJSON {
	return AmperageApplicationDeviceJSON{
		AppDev_ID: 					app_dev.AppDev_ID,
		Amperage_Application_ID: 	app_dev.Amperage_Application_ID,
		Device_ID: 					app_dev.Device_ID,
		Amount: 					app_dev.Amount,
		Notes: 						app_dev.Notes,
		Amperage: 					app_dev.Amperage,
	}
}

func AmperageApplicationDeviceFromJSON (amperage_applicationDeviceJSON AmperageApplicationDeviceJSON) ds.AmperageApplicationDevices {
	return ds.AmperageApplicationDevices{
		Amount: 		amperage_applicationDeviceJSON.Amount,
		Notes: 			amperage_applicationDeviceJSON.Notes,
		Amperage: 		amperage_applicationDeviceJSON.Amperage,
	}
}