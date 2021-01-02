package utils

func QueueRun(funs ...func() bool) {
	for i := 0; i < len(funs); i++ {
		if !funs[i]() {
			return
		}
	}
}
