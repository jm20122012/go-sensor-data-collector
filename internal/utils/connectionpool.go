package utils

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AcquirePoolConn(pool *pgxpool.Pool) *pgxpool.Conn {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		slog.Error("Unable to acquire connection", "error", err)
	}
	return conn
}
