package serializer

import "web_backend/internal/app/ds"

type DeviceJSON struct {
	Device_ID 	uint	`json:"device_id"`
	IsDelete 	bool	`json:"is_delete"`
	Title	 	string  `json:"title"`
	Description string  `json:"description"`
	Dev_Power 	float64	`json:"dev_power"`
	Photo      	string  `json:"photo"`
}

func DeviceToJSON(device ds.Device) DeviceJSON {
	return DeviceJSON{
		Device_ID: 		device.Device_ID,
		IsDelete: 		device.IsDelete,
		Title: 			device.Title,
		Description: 	device.Description,
		Dev_Power: 		device.Dev_Power,
		Photo: 			device.Photo,	
	}
}

func DeviceFromJSON(deviceJSON DeviceJSON) ds.Device {
	return ds.Device{
		IsDelete: 		deviceJSON.IsDelete,
		Title: 			deviceJSON.Title,
		Description: 	deviceJSON.Description,
		Dev_Power: 		deviceJSON.Dev_Power,
		Photo: 			deviceJSON.Photo,
	}
}