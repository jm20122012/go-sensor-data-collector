package sensordb

import (
	"context"
	"fmt"
	"log"
	"os"
	"sensor-data-collection-service/internal/sensordb/sqlc"
	"sensor-data-collection-service/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbWrapper struct {
	Conn *pgxpool.Conn
	DB   SensorDataDB
}

func (d *DbWrapper) Release() {
	d.Conn.Release()
}

func NewDbWrapper(pool *pgxpool.Pool) *DbWrapper {
	conn := utils.AcquirePoolConn(pool)
	db := NewSensorDataDB(conn)
	return &DbWrapper{
		Conn: conn,
		DB:   *db,
	}
}

type SensorDataDB struct {
	*sqlc.Queries
}

func NewSensorDataDB(conn *pgxpool.Conn) *SensorDataDB {
	return &SensorDataDB{
		Queries: sqlc.New(conn),
	}
}

func BuildPgConnectionString(user string, password string, host string, port string, database string) string {
	// connectionString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, database)
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)

	return connectionString
}

func CreateConnectionPool(connectionString string) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
