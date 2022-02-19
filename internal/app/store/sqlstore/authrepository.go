package sqlstore

import (
	//"database/sql"

	"time"

	"github.com/katelinlis/AuthBackend/internal/app/model"
)

type authRepository struct {
	store *Store
}

func (r *authRepository) Create(a *model.Auth) error {
	if err := a.BeforeCreate(); err != nil {
		return err
	}

	if err := a.Validate(); err != nil {
		return err
	}

	var err2 = r.store.db.QueryRow(
		"INSERT INTO users (username,password) VALUES ($1,$2) RETURNING id",
		a.Login,
		a.EncryptedPassword,
		time.Now().Unix(),
	).Scan(&a.ID)

	return err2
}

func (r *authRepository) GetUserByUsername(auth *model.Auth) error {
	if err := auth.Validate(); err != nil {
		return err
	}
	err := r.store.db.QueryRow("Select id,password from users where username = $1 ORDER BY id DESC", auth.Login).Scan(&auth.ID, &auth.EncryptedPassword)

	return err
}
