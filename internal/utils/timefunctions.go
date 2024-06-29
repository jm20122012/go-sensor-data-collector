package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetUTCTimestamp() pgtype.Timestamptz {
	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Create a Timestamptz and set its value
	var timestamptz pgtype.Timestamptz
	timestamptz.Time = currentTime
	timestamptz.Valid = true

	return timestamptz
}
