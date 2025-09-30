package serializer

import (
	"web_backend/internal/app/ds"
	"time"
)
type ApplicationJSON struct {
	ID				uint		`json:"application_id"`
	Status			string		`json:"status"`
	Created_At		time.Time	`json:"created_at"`
	Creator_Login	string		`json:"creator_login"`
	Moderator_Login	*string		`json:"moderator_login"`
	Forming_Date	*time.Time	`json:"form_date"`
	Finish_Date		*time.Time	`json:"finish_date"`
	Amperage		float64		`json:"amperage"`
}

func ApplicationToJSON(app ds.Application, creator_login string, moderator_login string) ApplicationJSON {
	var form_date, finish_date *time.Time
	if app.Forming_Date.Valid {
		form_date = &app.Forming_Date.Time
	}
	if app.Finish_Date.Valid {
		finish_date = &app.Finish_Date.Time
	}
	var m_login *string
	if moderator_login != "" {
		m_login = &moderator_login
	}
 	
	return ApplicationJSON{
		ID: 				app.Application_ID,
		Status: 			app.Status,
		Created_At: 		app.Created_At,
		Creator_Login: 		creator_login,
		Moderator_Login: 	m_login,
		Forming_Date: 		form_date,
		Finish_Date: 		finish_date,
		Amperage: 			app.Amperage,
	}
}

func ApplicationFromJSON(apps ApplicationJSON) ds.Application {
	if apps.Amperage == 0 {
		return ds.Application{}
	}
	return ds.Application{
		Amperage: apps.Amperage,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}