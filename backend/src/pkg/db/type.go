package db

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint64, nil
}

func (n *NullUint64) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		n.Uint64 = 0
	case string:
		tmp, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		n.Uint64 = uint64(tmp)
		n.Valid = true
	case []byte:
		tmp, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return err
		}
		n.Uint64 = uint64(tmp)
		n.Valid = true
	case int32:
		n.Uint64 = uint64(value)
		n.Valid = true
	case int64:
		n.Uint64 = uint64(value)
		n.Valid = true
	case uint64:
		n.Uint64 = uint64(value)
		n.Valid = true
	default:
		return fmt.Errorf("incompatible type for NullUint64")
	}

	return nil
}
