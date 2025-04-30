package adaptor

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/michaelyusak/go-helper/entity"
)

func ConnectPostgres(config entity.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host, config.Port, config.DbName))
	if err != nil {
		return nil, fmt.Errorf("[adaptor][ConnectPostgres][Open] error: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("[adaptor][ConnectPostgres][Ping] error: %w", err)
	}

	return db, nil
}
