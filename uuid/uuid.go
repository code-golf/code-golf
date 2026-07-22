package uuid

import (
	"database/sql/driver"
	"uuid"
)

type UUID struct{ uuid.UUID }

func Parse(s string) (UUID, error) {
	u, err := uuid.Parse(s)
	return UUID{u}, err
}

func (u *UUID) Scan(src any) error {
	return u.UnmarshalText(src.([]byte))
}

func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}
