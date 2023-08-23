package timer

import (
	"errors"
	"gnet/lib/logzap"
)

type TimerCallback func(int)

type Timer struct {
	cb        TimerCallback
	interval  int  //interval time of milloseconds per trigger
	elapsed   int  //time elapsed
	repeat    int  //repeat times, <= 0 forever
	repeated  int  //allready repeated times
	completed bool //is timer completed
	forever   bool //is timer forever
}

func NewTimer(interval, repeat int, cb TimerCallback) *Timer {
	if interval < 0 {
		interval = 0
	}
	if repeat < 0 {
		repeat = 0
	}
	t := &Timer{}
	t.interval = interval
	t.cb = cb
	t.repeat = repeat
	t.forever = (t.repeat <= 0)

	return t
}

func (t *Timer) update(dt int) {
	if t.completed {
		return
	}

	t.elapsed += dt
	if t.elapsed < t.interval {
		return
	}

	for t.elapsed >= t.interval {
		t.elapsed -= t.interval
		t.repeated += 1

		t.trigger()

		if !t.forever {
			if t.repeated >= t.repeat {
				t.completed = true
				return
			}
		}
	}
}

func (t *Timer) trigger() {
	defer func() {
		if err := recover(); err != nil {
			logzap.Errorw("timer trigger recover err")
		}
	}()
	t.cb(t.interval)
}

// Reset timer's time elapsed and repeated times.
func (t *Timer) Reset() error {
	if t.completed {
		return errors.New("timer Reset err: is completed")
	}
	t.elapsed = 0
	t.repeated = 0
	return nil
}

// Cancel timer
func (t *Timer) cancel() {
	t.completed = true
}
