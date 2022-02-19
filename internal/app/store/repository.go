package store

import "github.com/katelinlis/AuthBackend/internal/app/model"

//AuthRepository ...
type AuthRepository interface {
	Create(*model.Auth) error //Создание пользователя
	GetUserByUsername(*model.Auth) error
}

type SessionsRepository interface {
	Create(*model.Session) error //Создание пользователя
	FindByUserid(int) (model.Session, error)
}
