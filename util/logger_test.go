package util

import (
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	log := Log{}

	name := "test.log"
	log.Start(name, os.O_RDWR|os.O_CREATE, 0666)
	defer log.Stop()
	log.SetFlags(0)
	log.Printf("show on screen and write to log file: %s\n", name)
	log.LPrintf("only write to log file: %s\n", name)
}
