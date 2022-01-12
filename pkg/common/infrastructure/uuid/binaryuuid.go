package uuid

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type BinaryUUID uuid.UUID

// Value returns []byte uuid
func (id BinaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(id).MarshalBinary()
}

func (id *BinaryUUID) Scan(src interface{}) error {
	var result uuid.UUID
	err := result.Scan(src)
	if err != nil {
		return err
	}
	*id = BinaryUUID(result)
	return nil
}
