package log_test

import (
	"gnet/lib/log"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log.Init("test", "game", log.DEBUG_LEVEL, log.DEBUG_LEVEL, 10000, 1000)
	for {
		time.Sleep(time.Second)
		log.Debug("hahaha %v, %v", 2, 3)
	}
}
