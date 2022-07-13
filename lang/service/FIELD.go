package service

import (

)

type FIELD struct {
	Id byte
	Value string
}

func NewFIELD() *FIELD{
	p := new(FIELD)
	return p
}
