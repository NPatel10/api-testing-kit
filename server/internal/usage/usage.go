package usage

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID           string          `json:"id"`
	UserID       *string         `json:"userId,omitempty"`
	SessionID    *string         `json:"sessionId,omitempty"`
	RequestRunID *string         `json:"requestRunId,omitempty"`
	Bucket       string          `json:"bucket"`
	EventKey     string          `json:"eventKey"`
	Quantity     int32           `json:"quantity"`
	Dimensions   json.RawMessage `json:"dimensions,omitempty"`
	OccurredAt   time.Time       `json:"occurredAt"`
}

type RecentFilter struct {
	UserID         *string
	SessionID      *string
	RequestRunID   *string
	Bucket         string
	EventKey       string
	OccurredAfter  *time.Time
	OccurredBefore *time.Time
	Limit          int32
}

type SummaryRow struct {
	Bucket        string    `json:"bucket"`
	EventKey      string    `json:"eventKey"`
	EventCount    int64     `json:"eventCount"`
	TotalQuantity int64     `json:"totalQuantity"`
	FirstOccurred time.Time `json:"firstOccurredAt"`
	LastOccurred  time.Time `json:"lastOccurredAt"`
}

type SummaryFilter struct {
	UserID         *string
	SessionID      *string
	RequestRunID   *string
	Bucket         string
	EventKey       string
	OccurredAfter  *time.Time
	OccurredBefore *time.Time
	Limit          int32
}
