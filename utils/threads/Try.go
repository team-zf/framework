package threads

import (
	"errors"
	"fmt"
)

func Try(fn func(), catch func(err error), finally ...func()) {
	defer func() {
		if len(finally) > 0 {
			finally[0]()
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			if catch != nil {
				catch(errors.New(fmt.Sprintf("%v", err)))
			}
		}
	}()
	fn()
}

func GoTry(fn func(), catch func(err error), finally ...func()) {
	go Try(fn, catch, finally...)
}
