package service

const (
	GET     byte = 1
	POST    byte = 2
	PUT     byte = 3
	DELETE  byte = 4
	PATCH   byte = 5
	OPTIONS byte = 6
	HEAD    byte = 7
	TRACE   byte = 8
)

var WebMethodName = map[string]byte{"GET": 1, "POST": 2, "PUT": 3, "DELETE": 4, "PATCH": 5, "OPTIONS": 6, "HEAD": 7, "TRACE": 8}
var WebMethodValue = map[byte]string{1: "GET", 2: "POST", 3: "PUT", 4: "DELETE", 5: "PATCH", 6: "OPTIONS", 7: "HEAD", 8: "TRACE"}
