package threads

/**
 * 新开协程的类有回调用的
 */
type ThreadRun struct {
	Chanresult chan func()
}

func (e *ThreadRun) Go(fn func(), callback func()) {
	GoTry(fn, nil, func() {
		e.Chanresult <- callback
		close(e.Chanresult)
	})
}

func NewGoRun(fn func(), callback ...func()) *ThreadRun {
	result := new(ThreadRun)
	result.Chanresult = make(chan func(), 1)
	if len(callback) == 0 {
		result.Go(fn, nil)
	} else {
		result.Go(fn, callback[0])
	}
	return result
}

func NewGo() *ThreadRun {
	result := new(ThreadRun)
	result.Chanresult = make(chan func(), 1)
	return result

}
