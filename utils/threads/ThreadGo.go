package threads

import (
	"context"
	"errors"
	"sync"
)

/**
 * 子协程管理计数，可以等子协程都完成
 * 用它来管理所有开的协程，需要等这些线程都跑完
 */
type ThreadGo struct {
	Wg  sync.WaitGroup
	Ctx context.Context
	Cal context.CancelFunc
}

// 在新前线程上跑
func (e *ThreadGo) Go(f func(ctx context.Context)) {
	if f == nil {
		panic(errors.New("Go func is nil."))
	}
	e.Wg.Add(1)
	GoTry(func() {
		f(e.Ctx)
	}, nil, func() {
		e.Wg.Done()
	})
}

// 返回可关掉的子协程
func (e *ThreadGo) SubGo(f func(ctx context.Context)) context.CancelFunc {
	if f == nil {
		panic(errors.New("Go func is nil."))
	}
	ctx, cal := context.WithCancel(e.Ctx)
	e.Wg.Add(1)
	GoTry(func() {
		f(ctx)
	}, nil, func() {
		e.Wg.Done()
	})
	return cal
}

// 在新协程上跑
func (e *ThreadGo) GoTry(f func(ctx context.Context), catch func(error), finally ...func()) {
	if f == nil {
		panic(errors.New("Go func is nil."))
	}
	e.Wg.Add(1)
	GoTry(
		func() {
			f(e.Ctx)
		},
		catch,
		func() {
			defer e.Wg.Done()
			if len(finally) > 0 {
				finally[0]()
			}
		},
	)
}

// 在当前协程上跑
func (e *ThreadGo) Try(f func(ctx context.Context), catch func(error), finally ...func()) {
	if f == nil {
		panic(errors.New("Go func is nil."))
	}
	e.Wg.Add(1)
	Try(
		func() {
			f(e.Ctx)
		},
		catch,
		func() {
			defer e.Wg.Done()
			if len(finally) > 0 {
				finally[0]()
			}
		},
	)
}

func (e *ThreadGo) CloseWait() {
	e.Cal()
	e.Wg.Wait()
}

func NewThreadGo() *ThreadGo {
	reuslt := new(ThreadGo)
	reuslt.Ctx, reuslt.Cal = context.WithCancel(context.Background())
	return reuslt

}
func NewThreadGoBySub(ctx context.Context) *ThreadGo {
	reuslt := new(ThreadGo)
	reuslt.Ctx, reuslt.Cal = context.WithCancel(ctx)
	return reuslt
}

func NewThreadGoByGo(thgo *ThreadGo) *ThreadGo {
	result := new(ThreadGo)
	result.Ctx, result.Cal = context.WithCancel(thgo.Ctx)
	return result
}
