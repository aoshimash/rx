package postgres

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// encodeCursor encodes a UUID to a base64 string for cursor-based pagination
func encodeCursor(id uuid.UUID) string {
	return base64.URLEncoding.EncodeToString(id[:])
}

// decodeCursor decodes a base64 string to a UUID
func decodeCursor(cursor string) (uuid.UUID, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromBytes(data)
}
