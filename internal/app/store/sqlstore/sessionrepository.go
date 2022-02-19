package sqlstore

import (
	//"database/sql"

	"time"

	"github.com/katelinlis/AuthBackend/internal/app/model"
)

type sessionsRepository struct {
	store *Store
}

func (r *sessionsRepository) Create(rt *model.Session) error {

	var err2 = r.store.db.QueryRow(
		"INSERT INTO wall (author,text,timestamp) VALUES ($1,$2,$3) RETURNING id",
		rt.RefrashToken,
		rt.UserID,
		rt.Useragent,
		time.Now().Unix(),
	).Scan(&rt.ID)

	return err2

}

func (r *sessionsRepository) FindByUserid(int) (model.Session, error) {

	return model.Session{}, nil
}
