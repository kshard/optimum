//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package common

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

const IDLE_TIME = 20 * time.Second

func spinner(bar *progressbar.ProgressBar, f func() error) error {
	ch := make(chan bool)

	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				bar.Add(1)
				time.Sleep(40 * time.Millisecond)
			}
		}
	}()

	err := f()

	ch <- false
	bar.Finish()

	return err
}
