package helpers

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func Int4ToInt(v pgtype.Int4) int {
	return int(v.Int32)
}

func IntToInt4(v int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(v),
		Valid: true,
	}
}

func TimeToTimestamp(v time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  v,
		Valid: true,
	}
}

func TimestampToTime(v pgtype.Timestamp) time.Time {
	return v.Time
}
