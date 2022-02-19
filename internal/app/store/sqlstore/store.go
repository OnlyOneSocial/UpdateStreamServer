package sqlstore

import (
	"database/sql"

	"github.com/katelinlis/AuthBackend/internal/app/store"

	//"githab.com/katelinlis/AuthBackend/internal/app/model"
	_ "github.com/lib/pq" //db import
)

//Store ...
type Store struct {
	db                  *sql.DB
	authRepository      *authRepository
	sessionshRepository *sessionsRepository
}

//New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

//Auth ...
func (s *Store) Auth() store.AuthRepository {
	if s.authRepository != nil {
		return s.authRepository
	}

	s.authRepository = &authRepository{
		store: s,
	}

	return s.authRepository
}

//Sessions ...
func (s *Store) Sessions() store.SessionsRepository {
	if s.sessionshRepository != nil {
		return s.sessionshRepository
	}

	s.sessionshRepository = &sessionsRepository{
		store: s,
	}

	return s.sessionshRepository
}
