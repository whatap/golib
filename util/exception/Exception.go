package exception

import (
	"fmt"
)

type CustomException struct {
	ErrorClassName string
	ErrorMessage   string
	ErrorStack     string
	ErrorEsc       string
}

func NewCustomException(t, msg, stack, esc string) *CustomException {
	p := new(CustomException)
	p.ErrorClassName = t
	p.ErrorMessage = msg
	p.ErrorStack = stack
	p.ErrorEsc = esc
	return p
}

func (cex *CustomException) Error() string {
	return fmt.Sprintf("name:%s\n,message:%s\n,esc:%s\n,stack:%s", cex.ErrorClassName, cex.ErrorMessage, cex.ErrorEsc, cex.ErrorStack)
}
