package domain

import (
	"encoding/binary"

	"github.com/google/uuid"
)

const (
	MinCheckInterval = 60 // seconds
)

func IntToUUID(num uint32) uuid.UUID {
	var b [16]byte
	binary.BigEndian.PutUint32(b[:4], num)
	return uuid.Must(uuid.FromBytes(b[:]))
}
