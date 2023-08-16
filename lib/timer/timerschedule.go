package timer

import (
	"gnet/lib/vector"
	"sync"
	"time"
)

type TimerSchedule struct {
	timers      *vector.Vector
	addedCache  *vector.Vector
	deleteCache *vector.Vector
	mutex       sync.Mutex
	bStart      bool
	tick        *time.Ticker
}

func NewTimerSchedule() *TimerSchedule {
	ts := &TimerSchedule{}
	ts.timers = vector.New()
	ts.addedCache = vector.New()
	ts.deleteCache = vector.New()
	return ts
}

// Start the TimerSchedule
func (ts *TimerSchedule) Start() {
	if ts.bStart {
	} else {
		ts.bStart = true
		ts.tick = time.NewTicker(time.Duration(100) * time.Millisecond) //millisecond = 毫秒 = 千分之一秒
		go func() {
			for {
				<-ts.tick.C
				ts.Update(10)
			}
		}()
	}
}

// Stop the TimerSchedule
func (ts *TimerSchedule) Stop() bool {
	if ts.bStart {
		return false
	}
	ts.bStart = true
	ts.tick = time.NewTicker(time.Duration(100) * time.Millisecond)
	go func() {
		for {
			<-ts.tick.C
			ts.Update(10)
		}
	}()
	return true
}

// Update update all timers
func (ts *TimerSchedule) Update(dt int) {
	ts.mutex.Lock()
	ts.timers.AppendVec(ts.addedCache)
	ts.addedCache.Clear()
	ts.mutex.Unlock()
	for i := 0; i < ts.timers.Len(); i++ {
		t := ts.timers.At(i).(*Timer)
		t.update(dt)
		if t.isComplete {
			ts.UnSchedule(t)
		}
	}
	ts.mutex.Lock()
	for i := 0; i < ts.deleteCache.Len(); i++ {
		t := ts.deleteCache.At(i)
		for i := 0; i < ts.timers.Len(); i++ {
			if ts.timers.At(i) == t {
				ts.timers.Delete(i)
				break
			}
		}
	}
	ts.deleteCache.Clear()
	ts.mutex.Unlock()
}

// Schedule start a timer with interval(100=1s) and repeat.
// callback will be triggerd each interval, and timer will delete after trigger repeat times
// if interval is small than schedule's interval
// it may trigger multitimes at a update.
func (ts *TimerSchedule) Schedule(interval, repeat int, cb TimerCallback) *Timer {
	t := NewTimer(interval, repeat, cb)
	ts.mutex.Lock()
	ts.addedCache.Put(t)
	ts.mutex.Unlock()
	return t
}

func (ts *TimerSchedule) UnSchedule(t *Timer) {
	ts.mutex.Lock()
	ts.deleteCache.Put(t)
	ts.mutex.Unlock()
	t.cancel()
}
