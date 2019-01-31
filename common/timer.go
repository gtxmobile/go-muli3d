package common
import (
	"time"
)
type MyTimer struct {
	start_time_ time.Time
}
func New_timer() *MyTimer{
	timer := MyTimer{}
	timer.Restart()
	return &timer
}

func (timer *MyTimer)Restart(){
	timer.start_time_ = timer.Current_time()
}

func (timer *MyTimer)Elapsed() float64{
	sec := timer.Current_time().Sub(timer.start_time_)
	return sec.Seconds()
}

func (*MyTimer)Current_time() time.Time{
	return time.Now()
}
type Fps_counter struct {
	timer_ MyTimer
	elapsed_frame_ uint32
	interval_,fps_,	elapsed_seconds_ float64
}
func New_fps_counter(interval float64) *Fps_counter{
	return &Fps_counter{}
}

func (fpc *Fps_counter)on_frame(fps *float64) bool{
	fpc.elapsed_seconds_ += float64(fpc.timer_.Elapsed())
	fpc.elapsed_frame_++
	fpc.timer_.Restart()
	if fpc.elapsed_seconds_ >= fpc.interval_{
		*fps = float64(fpc.elapsed_frame_) / fpc.elapsed_seconds_
		fpc.elapsed_seconds_ = 0
		fpc.elapsed_frame_ = 0
		return true
	}
	return false
}