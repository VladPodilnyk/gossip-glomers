package main

import (
	"math"
	"time"
)

type RetryableFunc = func(effect func() error, errHandler func(err error)) error

func exponentialBackoff(baseDelayMs time.Duration, attempts int) RetryableFunc {
	return func(effect func() error, errHandler func(err error)) error {
		var err error
		for i := 1; i <= attempts; i++ {
			err = effect()
			if err == nil {
				return nil
			}
			delay := time.Duration(math.Pow(2, float64(i))) * baseDelayMs
			time.Sleep(delay)
		}
		errHandler(err)
		return err
	}
}
