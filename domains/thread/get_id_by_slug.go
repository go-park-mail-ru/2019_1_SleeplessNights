package thread

import (
	"github.com/jackc/pgx"
)

func getIdBySlug(slug string, tx *pgx.Tx)(id int64, err error) {
	row := tx.QueryRow(`SELECT * FROM func_thread_get_id_by_slug($1)`, slug)
	err = row.Scan(&id)
	return
}
