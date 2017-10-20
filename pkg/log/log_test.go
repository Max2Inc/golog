package log

import (
	"testing"
	"time"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"math/rand"
)
type NodeConfig struct {
	ID        string `json:"_id"`
	Rev       string `json:"_rev"`
	MASServer string `json:"mas_server"`
}

func generateLogs (wg *sync.WaitGroup, logCount int) {
	defer wg.Done()
	r := rand.Intn(1000)

	err := fmt.Errorf("What happened here ?")
	nodeConfig := NodeConfig {
		ID:         "identity",
		Rev:        "1.3",
		MASServer: "mas.internal.vsys.io",
	}
	entry := logrus.WithFields(logrus.Fields{
		"VMesh": 99,
		"node": nodeConfig,
	})

	for ii := 0; ii < logCount; ii++ {

		time.Sleep(time.Duration(r) * time.Microsecond)

		logger.Debugf("Some debug")
		logger.Infof("Some information")
		logger.WithError(err).Errorf("Some error")

		entry.WithField("time", time.Now()).Debugf("Some debug")
		entry.Infof("Some information")
		entry.WithError(err).Errorf("Some error")
	}
}



// Test multiple concurrent logs is safe across go-routines
func Test(t *testing.T) {

	goRoutinesStart := runtime.NumGoroutine()
	r := rand.Intn(10)

	goRoutines := 5000
	var wg sync.WaitGroup
	wg.Add(goRoutines)
	for ii := 0; ii < goRoutines; ii++ {

		time.Sleep(time.Duration(r) * time.Microsecond)

		go generateLogs(&wg, 50)
	}

	wg.Wait()

	if runtime.NumGoroutine() != goRoutinesStart {
		t.Errorf("Number of goRoutines did return to start. %d != %d", goRoutinesStart, runtime.NumGoroutine())
		t.Fail()
	}
}
