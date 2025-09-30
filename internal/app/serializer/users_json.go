package serializer

import "web_backend/internal/app/ds"

type UserJSON struct {
	ID 			uint	`json:"id"` 
	Login 		string	`json:"login"`
	Password 	string	`json:"password"`
	IsModerator bool 	`json:"is_moderator"`
}

func UserToJSON(user ds.Users) UserJSON {
	return UserJSON{
		ID: 			uint(user.User_ID),
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