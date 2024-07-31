//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/optimum
//

package opt

import (
	"context"
	"fmt"
	"time"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "tbd",
	Long: `
tbd
`,
	SilenceUsage: true,
	RunE:         commit,
}

func commit(cmd *cobra.Command, args []string) (err error) {
	cli, err := stack()
	if err != nil {
		return err
	}

	api := optimum.New(cli, host)

	receipt, err := api.Commit(context.Background(), curie.IRI(cask))
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(
			fmt.Sprintf("%s | %s ...", cask, "COMMITTING"),
		),
	)

	return spinner(bar, func() error {
		for {
			time.Sleep(IDLE_TIME)

			status, err := api.Status(context.Background(), receipt.Job)
			if err != nil {
				return err
			}

			bar.Describe(fmt.Sprintf("%s | %s ...", cask, status.Status))
			if status.Status == "SUCCEEDED" || status.Status == "FAILED" {
				return nil
			}
		}
	})
}
