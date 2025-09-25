package adaptor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/michaelyusak/go-helper/adaptor"
	"github.com/michaelyusak/go-helper/entity"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestConnectDbPostgres(t *testing.T) {
	ctx := context.Background()

	dbConf := entity.DBConfig{
		Username: "postgres",
		Password: "password",
		DbName:   "test_db",
	}

	psqlC, err := testcontainers.Run(
		ctx,
		"postgres:latest",
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     dbConf.Username,
			"POSTGRES_PASSWORD": dbConf.Password,
			"POSTGRES_DB":       dbConf.DbName,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("ready to accept connections"),
		),
	)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	testcontainers.CleanupContainer(t, psqlC)

	fmt.Println("psql testcontainers is running")

	host, err := psqlC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mappedPort, err := psqlC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatal(err)
	}

	dbConf.Host = host
	dbConf.Port = mappedPort.Port()

	t.Run("Test DB Down", func(t *testing.T) {
		dbConfCopy := dbConf
		dbConfCopy.Port = "5434"

		_, err := adaptor.ConnectDB(adaptor.PSQL, dbConfCopy)

		assert.Error(t, err)
	})

	t.Run("Test DB Up", func(t *testing.T) {
		db, err := adaptor.ConnectDB(adaptor.PSQL, dbConf)

		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("Test DB Running", func(t *testing.T) {
		db, _ := adaptor.ConnectDB(adaptor.PSQL, dbConf)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.PingContext(ctx)

		assert.NoError(t, err)
	})
}

func TestConnectDbMysql(t *testing.T) {
	ctx := context.Background()

	dbConf := entity.DBConfig{
		Username: "mysql",
		Password: "password",
		DbName:   "test_db",
	}

	mysqlC, err := testcontainers.Run(
		ctx,
		"mysql:latest",
		testcontainers.WithExposedPorts("3306/tcp"),
		testcontainers.WithEnv(map[string]string{
			"MYSQL_DATABASE":      dbConf.DbName,
			"MYSQL_USER":          dbConf.Username,
			"MYSQL_PASSWORD":      dbConf.Password,
			"MYSQL_ROOT_PASSWORD": dbConf.Password,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("3306/tcp"),
			wait.ForLog("ready for connections"),
		),
	)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	testcontainers.CleanupContainer(t, mysqlC)

	fmt.Println("mysql testcontainers is running")

	host, err := mysqlC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mappedPort, err := mysqlC.MappedPort(ctx, "3306/tcp")
	if err != nil {
		t.Fatal(err)
	}

	dbConf.Host = host
	dbConf.Port = mappedPort.Port()

	t.Run("Test DB Down", func(t *testing.T) {
		dbConfCopy := dbConf
		dbConfCopy.Port = "5434"

		_, err := adaptor.ConnectDB(adaptor.MYSQL, dbConfCopy)

		assert.Error(t, err)
	})

	t.Run("Test DB Up", func(t *testing.T) {
		db, err := adaptor.ConnectDB(adaptor.MYSQL, dbConf)

		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("Test DB Running", func(t *testing.T) {
		db, _ := adaptor.ConnectDB(adaptor.MYSQL, dbConf)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := db.PingContext(ctx)

		assert.NoError(t, err)
	})
}
