package gnetlog

import (
	"github.com/amsalt/log"
	"github.com/amsalt/log/adaptor/logrus"
)

// package gnetlog defines the log library used in gnet.
// it configures the log library.

// Init initial logger.
func Init() {
	defaultLogger()
}

func defaultLogger() {
	logger := logrus.NewBuilder(nil).Build()
	log.SetLogger(logger)
}
