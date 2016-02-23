package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	// Seconds field of the earliest valid Timestamp.
	// This is time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	minValidSeconds = -62135596800
	// Seconds field just after the latest valid Timestamp.
	// This is time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	maxValidSeconds = 253402300800
)

func (t *Timestamp) Validate() error {
	if t.Seconds < minValidSeconds {
		return fmt.Errorf("timestamp: %v before 0001-01-01", t)
	}

	if t.Seconds >= maxValidSeconds {
		return fmt.Errorf("timestamp: %v after 10000-01-01", t)
	}

	if t.Nanos < 0 || t.Nanos >= 1e9 {
		return fmt.Errorf("timestamp: %v: nanos not in range [0, 1e9)", t)
	}

	return nil
}

func (t *Timestamp) Scan(src interface{}) error {
	switch value := src.(type) {
	case int:
		t.Seconds = int64(value)
	case int32:
		t.Seconds = int64(value)
	case int64:
		t.Seconds = value
	case time.Time:
		nanos := value.UnixNano()
		seconds := nanos / 1e9
		t.Seconds = seconds
		t.Nanos = int32(nanos % 1e9)
	}

	return t.Validate()
}

func (t *Timestamp) Value() (driver.Value, error) {
	return time.Unix(t.Seconds, int64(t.Nanos)).UTC(), t.Validate()
}

func (t *Timestamp) Millis() int64 {
	millis := t.Seconds * 1000
	return millis + int64(t.Nanos/1e6)
}
