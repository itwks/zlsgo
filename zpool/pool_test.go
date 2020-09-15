/*
go test -race ./zpool -v
*/
package zpool_test

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zpool"
)

var curMem uint64

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

func TestPool(t *testing.T) {
	tt := zlsgo.NewTest(t)

	count := 100000
	workerNum := 1000
	g := sync.WaitGroup{}
	now := time.Now()
	mem := runtime.MemStats{}

	g = sync.WaitGroup{}
	now = time.Now()
	for i := 0; i < count; i++ {
		ii := i
		g.Add(1)
		go func() {
			_ = ii
			time.Sleep(time.Millisecond)
			g.Done()
		}()
	}
	g.Wait()

	mem = runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("NoPool memory:%dMB multiple:%v time:%v \n", curMem, float64(count)/float64(workerNum), time.Since(now))

	p := zpool.New(workerNum)
	g = sync.WaitGroup{}
	now = time.Now()
	for i := 0; i < count; i++ {
		ii := i
		g.Add(1)
		err := p.Do(func() {
			_ = ii
			time.Sleep(time.Millisecond)
			g.Done()
		})
		tt.EqualNil(err)
	}
	g.Wait()
	p.Close()
	mem = runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("Pool   memory:%dMB multiple:%v time:%v \n", curMem, float64(count)/float64(workerNum), time.Since(now))

	p = zpool.New(workerNum)
	g = sync.WaitGroup{}
	_ = p.PreInit()
	now = time.Now()
	for i := 0; i < count; i++ {
		ii := i
		g.Add(1)
		err := p.Do(func() {
			_ = ii
			time.Sleep(time.Millisecond)
			g.Done()
		})
		tt.EqualNil(err)
	}
	g.Wait()
	tt.EqualExit(uint(workerNum), p.Cap())
	p.Close()
	mem = runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("Pool   memory:%dMB multiple:%v time:%v \n", curMem, float64(count)/float64(workerNum), time.Since(now))

	err := p.PreInit()
	tt.EqualExit(true, err != nil)
	t.Log(err)

	p.Close()
	c := p.IsClosed()
	tt.EqualExit(true, c)

	err = p.Do(func() {})
	tt.EqualExit(true, err != nil)
	t.Log(err)
}

func TestPoolCap(t *testing.T) {
	tt := zlsgo.NewTest(t)
	p := zpool.New(0)
	g := sync.WaitGroup{}
	g.Add(1)
	err := p.Do(func() {
		time.Sleep(time.Second / 100)
		g.Done()
	})
	tt.EqualNil(err)
	tt.EqualExit(uint(1), p.Cap())
	g.Wait()

	maxsize := 10
	p = zpool.New(1, maxsize)
	g = sync.WaitGroup{}
	g.Add(1)
	err = p.Do(func() {
		time.Sleep(time.Second / 100)
		g.Done()
	})
	tt.EqualNil(err)
	tt.EqualExit(uint(1), p.Cap())
	g.Wait()
	p.Pause()
	tt.EqualExit(uint(0), p.Cap())

	newSize := 5
	go func() {
		tt.EqualExit(uint(0), p.Cap())
		time.Sleep(time.Second)
		p.Continue(newSize)
		tt.EqualExit(uint(newSize), p.Cap())
	}()

	restarSum := 6
	g.Add(restarSum)

	for i := 0; i < restarSum; i++ {
		ii := i
		go func() {
			now := time.Now()
			err := p.Do(func() {
				time.Sleep(time.Second / 100)
				t.Log("重新启动的任务", ii, time.Since(now))
				g.Done()
			})
			tt.EqualNil(err)
		}()
	}
	g.Wait()
	tt.EqualExit(uint(newSize), p.Cap())

	p.Continue(1000)
	tt.EqualExit(uint(maxsize), p.Cap())

	p.Close()
	p.Continue(1000)
}
