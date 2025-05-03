package adaptor

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/michaelyusak/go-helper/entity"
)

type DBType string

const (
	MYSQL DBType = "mysql"
	PSQL  DBType = "postgres"
)

func ConnectDB(dbType DBType, config entity.DBConfig) (*sql.DB, error) {
	var driver, dsn string

	switch dbType {
	case MYSQL:
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.Username, config.Password, config.Host, config.Port, config.DbName)
	case PSQL:
		driver = "pgx"
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			config.Username, config.Password, config.Host, config.Port, config.DbName)
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("[adapter][ConnectDB][Open] error: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("[adapter][ConnectDB][Ping] error: %w", err)
	}

	return db, nil
}
