package adaptor

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&readTimeout=5s&writeTimeout=5s",
			config.Username, config.Password, config.Host, config.Port, config.DbName)
	case PSQL:
		driver = "pgx"
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?connect_timeout=5",
			config.Username, config.Password, config.Host, config.Port, config.DbName)
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("[adapter][ConnectDB][Open] error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("[adapter][ConnectDB][Ping] error: %w", err)
	}

	return db, nil
}
