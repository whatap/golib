package uuidutil

import (
	"github.com/google/uuid"
)

func Generate() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return ""
	} else {
		return uuid.String()
	}
}

func ToLong(uuid string) int64 {
	h := int64(0)
	if uuid != "" {
		ln := len(uuid)
		for i := 0; i < ln; i++ {
			h = 31*h + int64(uuid[i])
		}
	}
	return h
}
