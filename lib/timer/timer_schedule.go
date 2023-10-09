package timer

import (
	"gnet/lib/logzap"
	"gnet/lib/vector"
	"go.uber.org/zap"
	"sync"
	"time"
)

// 计时器
type TimerSchedule struct {
	timers   *vector.Vector
	addCache *vector.Vector
	delCache *vector.Vector
	started  bool         // 是否已经启动
	ticker   *time.Ticker //
	tick     int          // 间隔(毫秒ms)
	mutex    sync.Mutex
}

func NewTimerSchedule() *TimerSchedule {
	ts := &TimerSchedule{}
	ts.timers = vector.New()
	ts.addCache = vector.New()
	ts.delCache = vector.New()
	ts.started = false
	ts.ticker = nil
	ts.tick = 1000 // 默认tick为1秒=1000ms
	return ts
}

// 设置间隔(毫秒ms)
func (ts *TimerSchedule) SetTick(tick int) {
	if !ts.started && tick > 0 {
		ts.tick = tick
	} else {
		logzap.Error("timerschedule set tick error", zap.Int("tick", tick))
	}
}

// 启动计时器
func (ts *TimerSchedule) Start() {
	if !ts.started {
		ts.started = true
		ts.ticker = time.NewTicker(time.Duration(ts.tick) * time.Millisecond) // 1000毫秒
		go func() {
			for {
				if ts.ticker != nil {
					<-ts.ticker.C
					ts.Update(ts.tick)
				} else {
					break
				}
			}
		}()
	}
}

// 停止计时器计时器
func (ts *TimerSchedule) Stop() {
	if ts.ticker != nil {
		ts.ticker.Stop()
		ts.ticker = nil
	}
	ts.started = false
}

// Update all timers
func (ts *TimerSchedule) Update(d int) {
	ts.mutex.Lock()
	if ts.addCache.Len() > 0 {
		ts.timers.AppendVec(ts.addCache)
		ts.addCache.Clear()
	}
	ts.mutex.Unlock()
	for i := 0; i < ts.timers.Len(); i++ {
		t := ts.timers.At(i).(*Timer)
		t.update(d)
		if t.Completed() {
			ts.UnSchedule(t)
		}
	}
	ts.mutex.Lock()
	for i := 0; i < ts.delCache.Len(); i++ {
		t := ts.delCache.At(i)
		for i := 0; i < ts.timers.Len(); i++ {
			if ts.timers.At(i) == t {
				ts.timers.Delete(i)
				break
			}
		}
	}
	ts.delCache.Clear()
	ts.mutex.Unlock()
}

// 启动一个倒计时
// @interval 时间间隔(毫秒ms)
// @repeat 重复次数 -1=永久重复
func (ts *TimerSchedule) Schedule(interval, repeat int, cb TimerCallback) *Timer {
	if repeat == 0 {
		repeat = 1
	}
	t := NewTimer(interval, repeat, cb)
	ts.mutex.Lock()
	ts.addCache.Put(t)
	ts.mutex.Unlock()
	return t
}

// 停止一个倒计时
func (ts *TimerSchedule) UnSchedule(t *Timer) {
	ts.mutex.Lock()
	ts.delCache.Put(t)
	ts.mutex.Unlock()
	t.Cancel()
}
