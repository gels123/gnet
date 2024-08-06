package logzap

import "testing"

func TestLogzap(t *testing.T) {
	SetSource("gwlog_test")
	//SetOutput([]string{"stderr", "gwlog_test.log"})
	//SetLevel(debugLevel)

	if lv := ParseLevel("debug"); lv != debugLevel {
		t.Fail()
	}
	if lv := ParseLevel("info"); lv != infoLevel {
		t.Fail()
	}
	if lv := ParseLevel("warn"); lv != warnLevel {
		t.Fail()
	}
	if lv := ParseLevel("error"); lv != errorLevel {
		t.Fail()
	}
	if lv := ParseLevel("panic"); lv != panicLevel {
		t.Fail()
	}
	if lv := ParseLevel("fatal"); lv != fatalLevel {
		t.Fail()
	}

	Debugf("this is a debug %d", 1)
	//SetLevel(infoLevel)
	Debugf("SHOULD NOT SEE THIS!")
	Infof("this is an info %d", 2)
	Warnf("this is a warning %d", 3)
	TraceError("this is a trace error %d", 4)
	func() {
		defer func() {
			_ = recover()
		}()
		Panicf("this is a panic %d", 4)
	}()

	func() {
		defer func() {
			_ = recover()
		}()
		//Fatalf("this is a fatal %d", 5)
	}()
}
