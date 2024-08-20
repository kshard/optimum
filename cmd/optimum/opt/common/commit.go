//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package common

import (
	"context"
	"fmt"
	"time"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/schollz/progressbar/v3"
)

func AboutCommit(kind, extension string) string {
	return fmt.Sprintf(`
Batch writing to "%s" data structure requires commit after successful dataset
upload before dataset is available to reads.
%s
`, kind, extension)
}

// List all data structures of given type
func Commit(api *optimum.Client, id curie.IRI) (err error) {
	receipt, err := api.Commit(context.Background(), id)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(
			fmt.Sprintf("%s (vsn %s) | %s ...", curie.Reference(id), receipt.Version, "COMMITTING"),
		),
	)

	return spinner(bar, func() error {
		for {
			time.Sleep(IDLE_TIME)

			status, err := api.Status(context.Background(), receipt.Job)
			if err != nil {
				return err
			}

			bar.Describe(fmt.Sprintf("%s (vsn %s) | %s ...", curie.Reference(id), receipt.Version, status.Status))
			if status.Status == "SUCCEEDED" || status.Status == "FAILED" {
				return nil
			}
		}
	})
}
