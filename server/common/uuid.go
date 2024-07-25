package common

import (
	"github.com/google/uuid"
)

func GentNewUuid() *uuid.UUID {
	newUUID := uuid.New()
	return &newUUID
}
