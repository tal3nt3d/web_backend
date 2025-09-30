package ds

import(
	"time"
	"database/sql"
)

type Application struct {
	Application_ID 	uint		`gorm:"primaryKey;autoIncrement; not null"`
	Status 			string		`gorm:"type:varchar(15); not null"`
	Created_At      time.Time	`gorm:"not null"`
	Creator_ID	 	uint	  	`gorm:"type:integer(15); not null"`
	Moderator_ID  	uint  		`gorm:"type:integer(15); default: null"`
	Forming_Date 	sql.NullTime`gorm:"default:null"`
	Finish_Date 	sql.NullTime`gorm:"default:null"`
	Amperage		float64		`gorm:"type:numeric(3,1)"`
	
	Creator 		Users		`gorm:"foreignKey:Creator_ID"`
	Moderator 		Users		`gorm:"foreignKey:Moderator_ID"`
}

