package uuid

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
)

var errInvalid = errors.New("invalid uuid")

type UUID [16]byte

func Parse(s string) (UUID, error) {
	var u UUID
	err := u.UnmarshalText([]byte(s))
	return u, err
}

func (u *UUID) Scan(src any) error {
	return u.UnmarshalText(src.([]byte))
}

func (u UUID) String() string {
	dst := []byte("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	hex.Encode(dst[0:8], u[0:4])
	hex.Encode(dst[9:13], u[4:6])
	hex.Encode(dst[14:18], u[6:8])
	hex.Encode(dst[19:23], u[8:10])
	hex.Encode(dst[24:36], u[10:16])
	return string(dst)
}

func (u *UUID) UnmarshalText(b []byte) error {
	if len(b) != 36 {
		return errInvalid
	}
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return errInvalid
	}
	if _, err := hex.Decode(u[0:4], b[0:8]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(u[4:6], b[9:13]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(u[6:8], b[14:18]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(u[8:10], b[19:23]); err != nil {
		return errInvalid
	}
	if _, err := hex.Decode(u[10:16], b[24:36]); err != nil {
		return errInvalid
	}
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return u.String(), nil
}
