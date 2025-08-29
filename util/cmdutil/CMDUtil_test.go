package cmdutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteShellCommand(t *testing.T) {

	result, err := ExecuteShellCommand("ls -alh")
	fmt.Println("result", result)
	assert.Nil(t, err)

}
