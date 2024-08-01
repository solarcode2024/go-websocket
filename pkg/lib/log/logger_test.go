package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLog(t *testing.T) {
	Log(logrus.DebugLevel, "hello, world!")
}
