package ds

type AmperageApplicationDevices struct {
	AppDev_ID 					uint		`gorm:"primaryKey"`
	Amperage_Application_ID 	uint		`gorm:"not null;uniqueIndex:idx_amperage_application_device"`
	Device_ID      				uint		`gorm:"not null;uniqueIndex:idx_amperage_application_device"`
	Amount 						int			`gorm:"type:integer"`
	Notes 						string  	`gorm:"type:varchar(255)"`
	Amperage 					float64		`gorm:"type:numeric(10,3)"`

	Device 					Device 				`gorm:"foreignKey:Device_ID;references:Device_ID"`
	AmperageApplication 	AmperageApplication `gorm:"foreignKey:Amperage_Application_ID;references:Amperage_Application_ID"`
}

