package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()
	assert.NotNil(t, logger)

	logger.Error("Error!")
	logger.Errorf("Error!")

	logger.Warnf("Warn!")
	logger.Warn("Warn!")

	logger.Infof("Warn!")
	logger.Info("Warn!")
	logger.Infoln("Warn!")

	logger.Printf("[Print]", "Print!")
	logger.Println("[Print]", "Println!")

	logger.Debugf("Debug!")
	logger.Debug("Debug!")

}
