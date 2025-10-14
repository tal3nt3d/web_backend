package repository

import (
	"errors"
	"os"
	"time"
	"web_backend/internal/app/ds"
	"web_backend/internal/app/serializer"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (r *Repository) GetUserByID(id uuid.UUID) (ds.Users, error) {
	user := ds.Users{}
	sub := r.db.Where("user_id = ?", id).Find(&user)
	if sub.Error != nil {
		return ds.Users{}, sub.Error
	}
	if sub.RowsAffected == 0 {
		return ds.Users{}, ErrNotFound
	}
	err := sub.First(&user).Error
	if err != nil {
		return ds.Users{}, err
	}
	return user, nil
}

func (r *Repository) GetUserByLogin(login string) (ds.Users, error) {
	user := ds.Users{}
	sub := r.db.Where("login = ?", login).Find(&user)
	if sub.Error != nil {
		return ds.Users{}, sub.Error
	}
	if sub.RowsAffected == 0 {
		return ds.Users{}, ErrNotFound
	}
	err := sub.First(&user).Error
	if err != nil {
		return ds.Users{}, err
	}
	return user, nil
}

func (r *Repository) CreateUser(userJSON serializer.UserJSON) (ds.Users, error) {
	user := serializer.UserFromJSON(userJSON)
	if user.Login == "" {
		return ds.Users{}, errors.New("login is empty")
	}
	if user.Password == "" {
		return ds.Users{}, errors.New("password is empty")
	}
	if _, err := r.GetUserByLogin(user.Login); err == nil {
		return ds.Users{}, errors.New("user already exists")
	}

	user.User_ID = uuid.New()

	sub := r.db.Create(&user)
	if sub.Error != nil {
		return ds.Users{}, sub.Error
	}
	return user, nil
}

func (r *Repository) SignIn(userJSON serializer.UserJSON) (string, error) {
	user, err := r.GetUserByLogin(userJSON.Login)
	if err != nil {
		return "", err
	}

	token, err := GenerateToken(user.User_ID, user.IsModerator)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateToken(id uuid.UUID, isModerator bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = id.String()
	claims["is_moderator"] = isModerator
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *Repository) EditInfo(login string, userJSON serializer.UserJSON) (ds.Users, error) {
	currUser, err := r.GetUserByLogin(login)
	if err != nil {
		return ds.Users{}, err
	}

	if userJSON.Login != "" {
		currUser.Login = userJSON.Login
	}

	if userJSON.Password != "" {
		currUser.Password = userJSON.Password
	}

	if userJSON.IsModerator && !currUser.IsModerator {
		userJSON.IsModerator = false
	}
	currUser.IsModerator = userJSON.IsModerator

	err = r.db.Save(&currUser).Error
	if err != nil {
		return ds.Users{}, err
	}
	return currUser, nil
}