package list

import (
	"fmt"
)

type LinkedListEntity struct {
	prev  *LinkedListEntity
	next  *LinkedListEntity
	Value interface{}
}

func (e *LinkedListEntity) ToString() string {
	if e.Value == nil {
		return ""
	}
	return fmt.Sprint(e.Value)
}
