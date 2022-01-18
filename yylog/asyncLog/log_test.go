package asyncLog

import (
	"testing"
	"time"
)

func TestNewLogFile(t *testing.T) {
	n := len(asyncLog.files)

	lf1 := NewLogFile("test1.log")
	if len(asyncLog.files) != n+1 {
		t.Fatalf("NewLogFile num != 1")
	}

	_ = NewLogFile("test1.log")
	if len(asyncLog.files) != n+1 {
		t.Fatalf("NewLogFile num != 1 repeat")
	}

	lf2 := NewLogFile("test2.log")
	if len(asyncLog.files) != n+2 {
		t.Fatalf("NewLogFile num != 2")
	}
	lf2.SetRotate(RotateDate)

	_ = lf1.Write("lf1: hello world1")
	_ = lf1.Write("lf1: hello world2")
	_ = lf2.Write("lf2: hello world1")
	_ = lf2.Write("lf2: hello world2")
	_ = lf1.Write("lf1: hello world3")
	_ = lf1.Write("lf1: hello world4")

	time.Sleep(time.Second * 2)

	_ = lf1.Write("lf1: ---hello world1")
	_ = lf1.Write("lf1: ---hello world2")
	_ = lf2.Write("lf2: ---hello world1")
	_ = lf2.Write("lf2: ---hello world2")
	_ = lf1.Write("lf1: ---hello world3")
	_ = lf1.Write("lf1: ---hello world4")

	time.Sleep(time.Second * 2)
}

func TestJson(t *testing.T) {
	lf1 := NewLogFile("test-json.log")
	lf1.SetFlags(NoFlag)
	lf1.SetUseCache(false)

	var hello = struct {
		Hello string
		World int
	}{
		Hello: "test",
		World: 12,
	}
	_ = lf1.WriteJson(hello)
	_ = lf1.WriteJson(hello)
	_ = lf1.WriteJson(hello)

	time.Sleep(time.Second * 2)
}

func TestProbability(t *testing.T) {
	lf := NewLogFile("test-probability.log")
	lf.SetProbability(0.5)
	for i := 0; i < 20; i++ {
		lf.Write("probability")
	}

	time.Sleep(time.Second)
}

func BenchmarkWrite(b *testing.B) {
	lf := NewLogFile("bench-test.log")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = lf.Write("hello world")
		}
	})
}

func BenchmarkWriteNoCache(b *testing.B) {
	lf := NewLogFile("bench-nocache-test.log")
	lf.SetUseCache(false)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = lf.Write("hello world")
		}
	})
}
