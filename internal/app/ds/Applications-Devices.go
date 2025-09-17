package ds

type ApplicationDevices struct {
	AppDev_ID 		uint		`gorm:"primaryKey"`
	Application_ID 	uint		`gorm:"not null;uniqueIndex:idx_application_device"`
	Device_ID      	uint		`gorm:"not null;uniqueIndex:idx_application_device"`
	Amount 			int			`gorm:"type:integer"`
	Notes 			string  	`gorm:"type:varchar(255)"`
	Amperage 		float64		`gorm:"type:numeric(10,3)"`

	Device 			Device 		`gorm:"foreignKey:Device_ID;references:Device_ID"`
	Application 	Application `gorm:"foreignKey:Application_ID;references:Application_ID"`
}

