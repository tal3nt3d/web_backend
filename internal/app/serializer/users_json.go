package serializer

import (
	"web_backend/internal/app/ds"

	"github.com/google/uuid"
)

type UserJSON struct {
	ID 			uuid.UUID	`json:"id"` 
	Login 		string		`json:"login"`
	Password 	string		`json:"password"`
	IsModerator bool 		`json:"is_moderator"`
}

func UserToJSON(user ds.Users) UserJSON {
	return UserJSON{
		ID: 			uuid.UUID(user.User_ID),
		Login:			user.Login,
		Password: 		user.Password,
		IsModerator: 	user.IsModerator,
	}
}

func UserFromJSON(userJSON UserJSON) ds.Users {
	return ds.Users{
		Login: 			userJSON.Login,
		Password: 		userJSON.Password,
		IsModerator: 	userJSON.IsModerator,
	}
}