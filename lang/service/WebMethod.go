package service

import ()

const (
	GET    byte = 1
	POST   byte = 2
	PUT    byte = 3
	DELETE byte = 4
)

var WebMethodName = map[string]byte{"GET": 1, "POST": 2, "PUT": 3, "DELETE": 4}
var WebMethodValue = map[byte]string{1: "GET", 2: "POST", 3: "PUT", 4: "DELETE"}
