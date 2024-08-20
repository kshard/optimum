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
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/schollz/progressbar/v3"
)

func AboutCreate(kind, extension string) string {
	return fmt.Sprintf(`
Creates new instance of "%s" data structure. Omitting the configuration
parameters causes usage of default params.
%s
`, kind, extension)
}

func Create(api *optimum.Client, id curie.IRI, fopts string) (err error) {
	opts := map[string]any{}

	if fopts != "" {
		b, err := os.ReadFile(fopts)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &opts); err != nil {
			return err
		}
	}

	receipt, err := api.Create(context.Background(), id, opts)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(
			fmt.Sprintf("%s (vsn %s) | %s ... opts: %+v", curie.Reference(id), receipt.Version, "CREATING", opts),
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
