package util

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

// Convert uuid.UUID (Google) to pgtype.UUID (pgx)
func ConvertToPgtypeUUID(id uuid.UUID) pgtype.UUID {
	var pgUUID pgtype.UUID
	_ = pgUUID.Scan(id.String())
	return pgUUID
}

func ToPgTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

func FromPgUUID(pgUUID pgtype.UUID) (uuid.UUID, error) {
	if !pgUUID.Valid {
		return uuid.UUID{}, fmt.Errorf("UUID is not valid")
	}
	return uuid.FromBytes(pgUUID.Bytes[:])
}
