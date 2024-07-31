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
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fogfish/curie"
	"github.com/kshard/optimum"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&fopts, "json", "j", "", "json config file")
}

var (
	fopts string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create data collection",
	Long: `
tbd
`,
	SilenceUsage: true,
	RunE:         create,
}

func create(cmd *cobra.Command, args []string) (err error) {
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

	cli, err := stack()
	if err != nil {
		return err
	}

	api := optimum.New(cli, host)

	receipt, err := api.Create(context.Background(), curie.IRI(cask), opts)
	if err != nil {
		return err
	}

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(
			fmt.Sprintf("%s | %s ... opts: %+v", cask, "CREATING", opts),
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
