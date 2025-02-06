package timer

import (
	"errors"
	"gnet/lib/logzap"
)

type TimerCbFunc func(uint64)

type Timer struct {
	interval  uint64        // 时间间隔(毫秒ms)
	elapsed   uint64        // 时间流逝(毫秒ms)
	repeat    int        	// 重复次数 <0永久重复
	repeated  int           // 已重复次数
	forever   bool          // 是否永久重复
	completed bool          // 是否已完成
	cb        TimerCbFunc   // 回调函数
}

func NewTimer(interval uint64, repeat int, cb TimerCbFunc) *Timer {
	if interval <= 0 {
		panic("NewTimer: interval is negative or zero.")
	}
	t := &Timer{}
	t.interval = interval
	t.elapsed = 0
	t.repeat = repeat
	t.repeated = 0
	t.forever = (t.repeat < 0)
	t.completed = false
	t.cb = cb
	
	return t
}

func (t *Timer) update(du uint64) {
	if t.completed {
		return
	}
	t.elapsed += du
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

// 重置倒计时
func (t *Timer) Reset() error {
	if t.completed {
		return errors.New("timer Reset err: is completed")
	}
	t.elapsed = 0
	t.repeated = 0
	
	return nil
}

// 取消倒计时
func (t *Timer) Cancel() {
	t.completed = true
}

// 是否已完成
func (t *Timer) Completed() bool {
	return t.completed
}
