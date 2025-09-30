package ds

type Device struct {
	Device_ID 	uint	`gorm:"primaryKey;autoIncrement"`
	IsDelete 	bool	`gorm:"type:boolean not null;default:false"`
	Photo      	string  `gorm:"type:varchar(100)"`
	Title	 	string  `gorm:"type:varchar(255) not null"`
	Description string  `gorm:"type:varchar(255) not null"`
	Dev_Power 	float64	`gorm:"type:numeric(10,3)"`
}
