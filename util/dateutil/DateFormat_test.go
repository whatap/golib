package dateutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateFormat(t *testing.T) {
	df := NewDateFormat("y-m-d")
	fmt.Println(df.Format())

	assert.NotNil(t, df)

}

func TestDateFormatParse(t *testing.T) {
	df := NewDateFormat("y-m-d")
	fmt.Println(df.Parse("1976-02-18"))

	assert.NotNil(t, df)

}
