package store

/*
Store репозитории данных
*/
type Store interface {
	Auth() AuthRepository         // интерфейс для базы авторизации
	Sessions() SessionsRepository //  интерфейс для базы сессий
}
